package dkron

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/serf/serf"
)

type RPC struct {
	agent *AgentCommand
}

func (r *RPC) ExecutionDone(execution Execution, reply *serf.NodeResponse) error {
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
		}).Error(err)

		return err
	}

	// Send notification
	Notification(r.agent.config, &execution, exg).Send()

	reply.From = r.agent.config.NodeName
	reply.Payload = []byte("saved")

	return nil
}

func listenRPC(a *AgentCommand) {
	r := RPC{
		agent: a,
	}

	rpc.Register(r)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func callExecutionDone(execution *Execution) error {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("error dialing:", err)
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	// Synchronous call
	var reply serf.NodeResponse
	err = client.Call("RPC.ExecutionDone", execution, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Debug("rpc: from: %s", reply.From)
	return nil
}
