package dkron

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/serf/serf"
)

const (
	QuerySchedulerRestart = "scheduler:restart"
	QueryRunJob           = "run:job"
	QueryRPCConfig        = "rpc:config"
)

// Send a serf run query to the cluster, this is used to ask a node or nodes
// to run a Job.
func (a *AgentCommand) RunQuery(execution *Execution) {
	filterNodes, filterTags, err := a.processFilteredNodes(job)
	if err != nil {
		log.WithFields(logrus.Fields{
			"job": job.Name,
			"err": err.Error(),
		}).Fatal("agent: Error processing filtered nodes")
	}
	log.Debug("agent: Filtered nodes to run: ", filterNodes)
	log.Debug("agent: Filtered tags to run: ", job.Tags)

	params := &serf.QueryParam{
		FilterNodes: filterNodes,
		FilterTags:  filterTags,
		RequestAck:  true,
	}

	exJson, _ := json.Marshal(ex)
	log.WithFields(logrus.Fields{
		"query":    QueryRunJob,
		"job_name": ex.JobName,
		"json":     string(exJson),
	}).Debug("agent: Sending query")

	qr, err := a.serf.Query(QueryRunJob, exJson, params)
	if err != nil {
		log.WithField("query", QueryRunJob).WithError(err).Fatal("agent: Sending query error")
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query": QueryRunJob,
					"from":  ack,
				}).Debug("agent: Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query":    QueryRunJob,
					"from":     resp.From,
					"response": string(resp.Payload),
				}).Debug("agent: Received response")

				// Save execution to store
				a.setExecution(resp.Payload)
			}
		}
	}
	log.WithFields(logrus.Fields{
		"query": QueryRunJob,
	}).Debug("agent: Done receiving acks and responses")
}

// Broadcast a SchedulerRestartQuery to the cluster, only server members
// will attend to this. Forces a scheduler restart and reload all jobs.
func (a *AgentCommand) schedulerRestartQuery(leaderName string) {
	params := &serf.QueryParam{
		FilterNodes: []string{leaderName},
		RequestAck:  true,
	}

	qr, err := a.serf.Query(QuerySchedulerRestart, []byte(""), params)
	if err != nil {
		log.Fatal("agent: Error sending the scheduler reload query", err)
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"from": ack,
				}).Debug("agent: Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"from":    resp.From,
					"payload": string(resp.Payload),
				}).Debug("agent: Received response")
			}
		}
	}
	log.WithField("query", QuerySchedulerRestart).Debug("agent: Done receiving acks and responses")
}

// Broadcast a query to get the RPC config of one dkron_server, any that could
// attend later RPC calls.
func (a *AgentCommand) queryRPCConfig() ([]byte, error) {
	nodeName := a.selectServer().Name

	params := &serf.QueryParam{
		FilterNodes: []string{nodeName},
		FilterTags:  map[string]string{"dkron_server": "true"},
		RequestAck:  true,
	}

	qr, err := a.serf.Query(QueryRPCConfig, nil, params)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": QueryRPCConfig,
			"error": err,
		}).Fatal("proc: Error sending query")
		return nil, err
	}
	defer qr.Close()

	ackCh := qr.AckCh()
	respCh := qr.ResponseCh()

	var rpcAddr []byte
	for !qr.Finished() {
		select {
		case ack, ok := <-ackCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query": QueryRPCConfig,
					"from":  ack,
				}).Debug("proc: Received ack")
			}
		case resp, ok := <-respCh:
			if ok {
				log.WithFields(logrus.Fields{
					"query":   QueryRPCConfig,
					"from":    resp.From,
					"payload": string(resp.Payload),
				}).Debug("proc: Received response")

				rpcAddr = resp.Payload
			}
		}
	}
	log.WithFields(logrus.Fields{
		"query": QueryRPCConfig,
	}).Debug("proc: Done receiving acks and responses")

	return rpcAddr, nil
}
