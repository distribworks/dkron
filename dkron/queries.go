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

var responseStatusCheck bool

// RunQueryParam defines the struct used to send a Run query
// using serf.
type RunQueryParam struct {
	Execution        *Execution `json:"execution"`
	RPCAddr          string     `json:"rpc_addr"`
	ExpectedResponse time.Time
}

// QueryResponse hold status about Serf Query Responses for a job query
type QueryResponse struct {
	Node map[string]map[string]bool
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

	// In case the job is not a child job, compute the next execution time
	if job.ParentJob == "" {
		if e, ok := a.sched.GetEntry(jobName); ok {
			job.Next = e.Next
			if err := a.applySetJob(job.ToProto()); err != nil {
				return nil, fmt.Errorf("agent: RunQuery error storing job %s before running: %w", jobName, err)
			}
		} else {
			return nil, fmt.Errorf("agent: RunQuery error retrieving job: %s from scheduler", jobName)
		}
	}

	for i := 0; ; i++ {
		if i > 10 {
			log.WithFields(logrus.Fields{
				"query": QueryRunJob,
				"job":   job.Name,
			}).Info("RunQuery retry limit reached. Skipping execution")
			break
		}

		if i > 0 {
			time.Sleep(time.Duration(i) * time.Second)
		}

		// In the first execution attempt we build and filter the target nodes
		// but we use the existing node target in case of retry.
		if ex.Attempt <= 1 {

			filterNodes, filterTags, err := a.processFilteredNodes(job)
			if err != nil {
				return nil, fmt.Errorf("agent: RunQuery error processing filtered nodes: %w", err)
			}

			log.WithFields(logrus.Fields{
				"nodes":    filterNodes,
				"job_name": job.Name,
			}).Debug("agent: Filtered nodes to run")

			log.WithFields(logrus.Fields{
				"tags":     filterTags,
				"job_name": job.Name,
			}).Debug("agent: Filtered tags to run")

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
				Timeout:     10 * time.Second,
			}

		} else {
			params = &serf.QueryParam{
				FilterNodes: []string{ex.NodeName},
				RequestAck:  true,
				Timeout:     10 * time.Second,
			}
		}
		rqp := &RunQueryParam{
			Execution:        ex,
			RPCAddr:          a.getRPCAddr(),
			ExpectedResponse: time.Now().Add(params.Timeout),
		}

		rqpJSON, _ := json.Marshal(rqp)

		log.WithFields(logrus.Fields{
			"query":    QueryRunJob,
			"job_name": job.Name,
			"json":     string(rqpJSON),
		}).Debug("agent: Sending query")

		qrs := QueryResponse{Node: make(map[string]map[string]bool)}

		log.WithFields(logrus.Fields{
			"job_name": job.Name,
		}).Infof("agent: Sending query. Attempt: %v", i)

		qr, err := a.serf.Query(QueryRunJob, rqpJSON, params)
		if err != nil {
			return nil, fmt.Errorf("agent: RunQuery sending query error: %w", err)
		}
		defer qr.Close()

		ackCh := qr.AckCh()

		var responseCounter int
		for !qr.Finished() {
			select {
			case ack, ok := <-ackCh:

				if ok {
					log.WithFields(logrus.Fields{
						"from":  ack,
						"job":   job.Name,
					}).Info("agent: Received ack")
					responseCounter++
					n, o := qrs.Node[ack]
					if !o {
						n = make(map[string]bool)
						qrs.Node[ack] = n
					}
					n["ack"] = true
				}
			}

			if len(params.FilterNodes) == responseCounter {
				log.WithFields(logrus.Fields{
					"query": QueryRunJob,
					"job":   job.Name,
				}).Debug("All responses received. Not waiting for query timeout")
				qr.Close()
				break
			}
		}

		responseStatusCheck = true

		for _, v := range params.FilterNodes {
			if _, e := qrs.Node[v]; e {
				if _, b := qrs.Node[v]["ack"]; b {
					log.WithFields(logrus.Fields{
						"query": QueryRunJob,
						"job":   job.Name,
					}).Debugf("Ack check validated from expected node: %v", v)
				} else {
					log.WithFields(logrus.Fields{
						"query": QueryRunJob,
						"job":   job.Name,
					}).Debugf("No ack received from expected node: %v", v)
					responseStatusCheck = false
					break
				}

			} else {
				log.WithFields(logrus.Fields{
					"query": QueryRunJob,
					"job":   job.Name,
				}).Debugf("No responses from expected node: %v", v)
				responseStatusCheck = false
				break
			}
		}

		if responseStatusCheck == true {
			break
		}

		log.WithFields(logrus.Fields{
			"time":  time.Since(start),
			"query": QueryRunJob,
			"job":   job.Name,
		}).Debug("agent: Done waiting responses")
	}

	return job, nil
}
