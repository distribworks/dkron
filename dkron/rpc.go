package dkron

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/abronan/valkeyrie/store"
	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/serf/serf"
)

var (
	ErrExecutionDoneForDeletedJob = errors.New("rpc: Received execution done for a deleted job.")
)

type RPCServer struct {
	agent *Agent
}

func (rpcs *RPCServer) GetJob(jobName string, job *Job) error {
	defer metrics.MeasureSince([]string{"rpc", "get_job"}, time.Now())
	log.WithFields(logrus.Fields{
		"job": jobName,
	}).Debug("rpc: Received GetJob")

	j, err := rpcs.agent.Store.GetJob(jobName)
	if err != nil {
		return err
	}

	// Copy the data structure
	job.Name = j.Name
	job.Shell = j.Shell
	job.EnvironmentVariables = j.EnvironmentVariables
	job.Command = j.Command
	job.Executor = j.Executor
	job.ExecutorConfig = j.ExecutorConfig

	return nil
}

func (rpcs *RPCServer) ExecutionDone(execution Execution, reply *serf.NodeResponse) error {
	defer metrics.MeasureSince([]string{"rpc", "execution_done"}, time.Now())
	log.WithFields(logrus.Fields{
		"group": execution.Group,
		"job":   execution.JobName,
	}).Debug("rpc: Received execution done")

	// Load the job from the store
	job, err := rpcs.agent.Store.GetJob(execution.JobName)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Warning(ErrExecutionDoneForDeletedJob)
			return ErrExecutionDoneForDeletedJob
		}
		log.Fatal("rpc:", err)
		return err
	}
	// Lock the job while editing
	if err = job.Lock(); err != nil {
		log.Fatal("rpc:", err)
	}

	// Get the defined output types for the job, and call them
	origExec := execution
	for k, v := range job.Processors {
		log.WithField("plugin", k).Debug("rpc: Processing execution with plugin")
		if processor, ok := rpcs.agent.ProcessorPlugins[k]; ok {
			v["reporting_node"] = rpcs.agent.config.NodeName
			e := processor.Process(&ExecutionProcessorArgs{Execution: origExec, Config: v})
			execution = e
		}
	}

	// Save the execution to store
	if _, err := rpcs.agent.Store.SetExecution(&execution); err != nil {
		return err
	}

	if execution.Success {
		job.LastSuccess = execution.FinishedAt
		job.SuccessCount++
	} else {
		job.LastError = execution.FinishedAt
		job.ErrorCount++
	}

	if err := rpcs.agent.Store.SetJob(job, nil); err != nil {
		log.Fatal("rpc:", err)
	}

	// Release the lock
	if err = job.Unlock(); err != nil {
		log.Fatal("rpc:", err)
	}

	reply.From = rpcs.agent.config.NodeName
	reply.Payload = []byte("saved")

	// If the execution failed, retry it until retries limit (default: don't retry)
	if !execution.Success && execution.Attempt < job.Retries+1 {
		execution.Attempt++

		// Keep all execution properties intact except the last output
		// as it could exceed serf query limits.
		execution.Output = []byte{}

		log.WithFields(logrus.Fields{
			"attempt":   execution.Attempt,
			"execution": execution,
		}).Debug("Retrying execution")

		rpcs.agent.RunQuery(&execution)
		return nil
	}

	exg, err := rpcs.agent.Store.GetExecutionGroup(&execution)
	if err != nil {
		log.WithError(err).WithField("group", execution.Group).Error("rpc: Error getting execution group.")
		return err
	}

	// Send notification
	Notification(rpcs.agent.config, &execution, exg, job).Send()

	// Jobs that have dependent jobs are a bit more expensive because we need to call the Status() method for every execution.
	// Check first if there's dependent jobs and then check for the job status to begin executiong dependent jobs on success.
	if len(job.DependentJobs) > 0 && job.Status() == Success {
		for _, djn := range job.DependentJobs {
			dj, err := rpcs.agent.Store.GetJob(djn)
			if err != nil {
				return err
			}
			dj.Run()
		}
	}

	return nil
}

var workaroundRPCHTTPMux = 0

func listenRPC(a *Agent) {
	r := &RPCServer{
		agent: a,
	}

	bindIp, err := a.GetBindIP()
	if err != nil {
		log.Fatal("rpc:", err)
	}
	RPCAddr := fmt.Sprintf("%s:%d", bindIp, a.config.RPCPort)
	log.WithFields(logrus.Fields{
		"rpc_addr": RPCAddr,
	}).Debug("rpc: Registering RPC server")

	rpc.Register(r)

	// ===== workaround ==========
	// This is needed mainly for testing
	// see: https://github.com/golang/go/issues/13395
	oldMux := http.DefaultServeMux
	if workaroundRPCHTTPMux > 0 {
		mux := http.NewServeMux()
		http.DefaultServeMux = mux
	}
	workaroundRPCHTTPMux = workaroundRPCHTTPMux + 1
	// ===========================

	rpc.HandleHTTP()

	// workaround
	http.DefaultServeMux = oldMux

	l, e := net.Listen("tcp", RPCAddr)
	if e != nil {
		log.Fatal("rpc:", e)
	}
	go http.Serve(l, nil)
}

type RPCClient struct {
	//Addres of the server to call
	ServerAddr string
}

func (rpcc *RPCClient) callExecutionDone(execution *Execution) error {
	defer metrics.MeasureSince([]string{"rpc", "call_execution_done"}, time.Now())
	client, err := rpc.DialHTTP("tcp", rpcc.ServerAddr)
	if err != nil {
		log.WithFields(logrus.Fields{
			"err":         err,
			"server_addr": rpcc.ServerAddr,
		}).Error("rpc: error dialing.")
		return err
	}
	defer client.Close()

	// Synchronous call
	var reply serf.NodeResponse
	err = client.Call("RPCServer.ExecutionDone", execution, &reply)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Warning("rpc: Error calling ExecutionDone")
		return err
	}
	log.Debug("rpc: from: ", reply.From)

	return nil
}

func (rpcc *RPCClient) GetJob(jobName string) (*Job, error) {
	defer metrics.MeasureSince([]string{"rpc", "call_get_job"}, time.Now())
	client, err := rpc.DialHTTP("tcp", rpcc.ServerAddr)
	if err != nil {
		log.WithFields(logrus.Fields{
			"err":         err,
			"server_addr": rpcc.ServerAddr,
		}).Error("rpc: error dialing.")
		return nil, err
	}
	defer client.Close()

	// Synchronous call
	var job Job
	err = client.Call("RPCServer.GetJob", jobName, &job)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Warning("rpc: Error calling GetJob")
		return nil, err
	}

	return &job, nil
}
