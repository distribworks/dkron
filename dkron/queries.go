package dkron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
)

const (
	// QueryRunJob define a run job query string
	QueryRunJob = "run:job"
	// QueryExecutionDone define the execution done query string
	QueryExecutionDone = "execution:done"
)

// RunQueryParam defines the struct used to send a Run query
// using serf.
type RunQueryParam struct {
	Execution *Execution `json:"execution"`
	RPCAddr   string     `json:"rpc_addr"`
}

// RunQuery sends a serf run query to the cluster, this is used to ask a node or nodes
// to run a Job. Returns a job with it's new status and next schedule.
func (a *Agent) RunQuery(jobName string, ex *Execution) (*Job, error) {
	start := time.Now()
	var params *serf.QueryParam

	job, err := a.Store.GetJob(jobName, nil)
	if err != nil {
		return nil, fmt.Errorf("agent: RunQuery error retrieving job: %s from store: %w", jobName, err)
	}

	// In case the job is not a child job
	if job.ParentJob == "" {
		if e, ok := a.sched.GetEntry(jobName); ok {
			job.Next = e.Next
		} else {
			return nil, fmt.Errorf("agent: RunQuery error retrieving job: %s from scheduler", jobName)
		}
	}
	job.Status = StatusRunning

	// TODO change
	// if err := a.applySetJob(job.ToProto()); err != nil {
	// 	return nil, fmt.Errorf("agent: RunQuery error storing job %s before running: %w", jobName, err)
	// }
	if err := a.Store.SetJob(job, false); err != nil {
		return nil, fmt.Errorf("agent: RunQuery error storing job %s before running: %w", jobName, err)
	}

	// In the first execution attempt we build and filter the target nodes
	// but we use the existing node target in case of retry.
	if ex.Attempt <= 1 {
		filterNodes, filterTags, err := a.processFilteredNodes(job)
		if err != nil {
			return nil, fmt.Errorf("agent: RunQuery error processing filtered nodes: %w", err)
		}
		log.Debug("agent: Filtered nodes to run: ", filterNodes)
		log.Debug("agent: Filtered tags to run: ", filterTags)

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
	rqpJSON, _ := json.Marshal(rqp)

	log.WithFields(logrus.Fields{
		"query":    QueryRunJob,
		"job_name": job.Name,
	}).Info("agent: Sending query")

	log.WithFields(logrus.Fields{
		"query":    QueryRunJob,
		"job_name": job.Name,
		"json":     string(rqpJSON),
	}).Debug("agent: Sending query")

	qr, err := a.serf.Query(QueryRunJob, rqpJSON, params)
	if err != nil {
		return nil, fmt.Errorf("agent: RunQuery sending query error: %w", err)
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
		"time":  time.Since(start),
		"query": QueryRunJob,
	}).Debug("agent: Done receiving acks and responses")

	return job, nil
}

// Broadcast a ExecutionDone to the cluster.
func (a *Agent) executionDoneQuery(nodes []string, group string) map[string]string {
	params := &serf.QueryParam{
		FilterNodes: nodes,
		RequestAck:  true,
	}

	log.WithFields(logrus.Fields{
		"query":   QueryExecutionDone,
		"members": nodes,
	}).Info("agent: Sending query")

	qr, err := a.serf.Query(QueryExecutionDone, []byte(group), params)
	if err != nil {
		log.WithError(err).Fatal("agent: Error sending the execution done query")
	}
	defer qr.Close()

	statuses := make(map[string]string)
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

				statuses[resp.From] = string(resp.Payload)
			}
		}
	}
	log.WithField("query", QueryExecutionDone).Debug("agent: Done receiving acks and responses")

	// In case the query finishes by deadline without receiving a response from the node
	// set the execution as finished, maybe the node is gone by now.
	return statuses
}
