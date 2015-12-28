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

func (r *RPCServer) ExecutionDone(execution Execution, reply *serf.NodeResponse) error {
	log.WithFields(logrus.Fields{
		"group": execution.Group,
		"job":   execution.JobName,
	}).Debug("rpc: Received execution done")

	// Save job status
	job, err := r.agent.store.GetJob(execution.JobName)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Error(ErrExecutionDoneForDeletedJob)
			return ErrExecutionDoneForDeletedJob
		}
		log.Fatal("rpc:", err)
		return err
	}

	// Save the new execution to store
	if _, err := r.agent.store.SetExecution(&execution); err != nil {
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
		log.Fatal("rpc:", err)
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

func (r *RPCClient) callExecutionDone(execution *Execution) error {
	client, err := rpc.DialHTTP("tcp", r.ServerAddr)
	if err != nil {
		log.WithFields(logrus.Fields{
			"err":         err,
			"server_addr": r.ServerAddr,
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
		}).Error("rpc: Error calling ExecutionDone")
		return err
	}
	log.Debug("rpc: from: %s", reply.From)

	return nil
}
