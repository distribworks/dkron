package dkron

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	"github.com/hashicorp/serf/serf"
)

var (
	ErrExecutionDoneForDeletedJob = errors.New("rpc: Received execution done for a deleted job.")
)

type RPCServer struct {
	agent *AgentCommand
}

func (rpcs *RPCServer) ExecutionDone(execution Execution, reply *serf.NodeResponse) error {
	log.WithFields(logrus.Fields{
		"group": execution.Group,
		"job":   execution.JobName,
	}).Debug("rpc: Received execution done")

	// Save job status
	job, err := rpcs.agent.store.GetJob(execution.JobName)
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

	// Save the new execution to store
	if _, err := rpcs.agent.store.SetExecution(&execution); err != nil {
		return err
	}

	if execution.Success {
		job.LastSuccess = execution.FinishedAt
		job.SuccessCount++
	} else {
		job.LastError = execution.FinishedAt
		job.ErrorCount++
	}

	if err := rpcs.agent.store.SetJob(job); err != nil {
		log.Fatal("rpc:", err)
	}

	// Release the lock
	if err = job.Unlock(); err != nil {
		log.Fatal("rpc:", err)
	}

	reply.From = rpcs.agent.config.NodeName
	reply.Payload = []byte("saved")

	// If the job failed, retry it until retries limit (default: don't retry)
	if !execution.Success && execution.Attempt < job.Retries+1 {
		execution.Attempt++

		log.WithFields(logrus.Fields{
			"attempt":   execution.Attempt,
			"execution": execution,
		}).Debug("Retrying execution")

		rpcs.agent.RunQuery(&execution)
		return nil
	}

	exg, err := rpcs.agent.store.GetExecutionGroup(&execution)
	if err != nil {
		log.WithError(err).WithField("group", execution.Group).Error("rpc: Error getting execution group.")
		return err
	}

	// Send notification
	Notification(rpcs.agent.config, &execution, exg).Send()

	// Run dependent jobs
	for _, djn := range job.DependentJobs {
		dj, err := rpcs.agent.store.GetJob(djn)
		if err != nil {
			return err
		}
		dj.Run()
	}

	return nil
}

var workaroundRPCHTTPMux = 0

func listenRPC(a *AgentCommand) {
	r := &RPCServer{
		agent: a,
	}

	log.WithFields(logrus.Fields{
		"rpc_addr": a.getRPCAddr(),
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

	l, e := net.Listen("tcp", a.getRPCAddr())
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
	log.Debug("rpc: from: %s", reply.From)

	return nil
}
