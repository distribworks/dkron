package dkron

import (
	"bytes"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/abronan/valkeyrie/store"
	"github.com/hashicorp/serf/serf"
)

const (
	QuerySchedulerRestart = "scheduler:restart"
	QueryRunJob           = "run:job"
	QueryRPCConfig        = "rpc:config"
)

type RunQueryParam struct {
	Execution *Execution `json:"execution"`
	RPCAddr   string     `json:"rpc_addr"`
}

// Send a serf run query to the cluster, this is used to ask a node or nodes
// to run a Job.
func (a *Agent) RunQuery(ex *Execution) {
	var params *serf.QueryParam

	job, err := a.Store.GetJob(ex.JobName)

	if err != nil {
		//Job can be removed and the QuerySchedulerRestart not yet received.
		//In this case, the job will not be found in the store.
		if err == store.ErrKeyNotFound {
			log.Warning("agent: Job not found, cancelling this execution")
			return
		}
		log.WithError(err).Fatal("agent: Getting job error")
		return
	}

	// In the first execution attempt we build and filter the target nodes
	// but we use the existing node target in case of retry.
	if ex.Attempt <= 1 {
		filterNodes, filterTags, err := a.processFilteredNodes(job)
		if err != nil {
			log.WithFields(logrus.Fields{
				"job": job.Name,
				"err": err.Error(),
			}).Fatal("agent: Error processing filtered nodes")
		}
		log.Debug("agent: Filtered nodes to run: ", filterNodes)
		log.Debug("agent: Filtered tags to run: ", job.Tags)

		//serf match regexp but we want only match full tag
		serfFilterTags := make(map[string]string)
		for key, val := range filterTags {
			b := new(bytes.Buffer)
			b.WriteString("^")
			b.WriteString(val)
			b.WriteString("$")
			serfFilterTags[key] = b.String()
		}

		params = &serf.QueryParam{
			FilterNodes: filterNodes,
			FilterTags:  serfFilterTags,
			RequestAck:  true,
		}
	} else {
		params = &serf.QueryParam{
			FilterNodes: []string{ex.NodeName},
			RequestAck:  true,
		}
	}

	rqp := &RunQueryParam{
		Execution: ex,
		RPCAddr:   a.getRPCAddr(),
	}
	rqpJson, _ := json.Marshal(rqp)

	log.WithFields(logrus.Fields{
		"query":    QueryRunJob,
		"job_name": job.Name,
		"json":     string(rqpJson),
	}).Debug("agent: Sending query")

	qr, err := a.serf.Query(QueryRunJob, rqpJson, params)
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
func (a *Agent) schedulerRestartQuery(leaderName string) {
	params := &serf.QueryParam{
		FilterNodes: []string{leaderName},
		RequestAck:  true,
	}

	qr, err := a.serf.Query(QuerySchedulerRestart, []byte(""), params)
	if err != nil {
		log.WithError(err).Fatal("agent: Error sending the scheduler reload query")
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
func (a *Agent) queryRPCConfig() ([]byte, error) {
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
