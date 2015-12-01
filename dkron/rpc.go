package dkron

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/serf/serf"
)

type RPCServer struct {
	agent *AgentCommand
}

func (r *RPCServer) ExecutionDone(execution Execution, reply *serf.NodeResponse) error {
	log.WithFields(logrus.Fields{
		"group": execution.Group,
		"job":   execution.JobName,
	}).Debug("rpc: Received execution done")

	// Save the new execution to store
	if _, err := r.agent.store.SetExecution(&execution); err != nil {
		log.Fatal(err)
		return err
	}

	// Save job status
	job, err := r.agent.store.GetJob(execution.JobName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if execution.Success {
		job.LastSuccess = execution.FinishedAt
		job.SuccessCount = job.SuccessCount + 1
	} else {
		job.LastError = execution.FinishedAt
		job.ErrorCount = job.ErrorCount + 1
	}

	if err := r.agent.store.SetJob(job); err != nil {
		log.Fatal(err)
	}

	exg, err := r.agent.store.GetExecutionGroup(&execution)
	if err != nil {
		log.WithFields(logrus.Fields{
			"group": execution.Group,
			"err":   err,
		}).Error("rpc: Error getting execution group.")

		return err
	}

	// Send notification
	Notification(r.agent.config, &execution, exg).Send()

	reply.From = r.agent.config.NodeName
	reply.Payload = []byte("saved")

	return nil
}

func listenRPC(a *AgentCommand) {
	r := &RPCServer{
		agent: a,
	}

	rpc.Register(r)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":3234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
