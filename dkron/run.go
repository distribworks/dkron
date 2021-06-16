package dkron

import (
	"fmt"
	"sync"

	"github.com/hashicorp/serf/serf"
)

// Run call the agents to run a job. Returns a job with its new status and next schedule.
func (a *Agent) Run(jobName string, ex *Execution) (*Job, error) {
	job, err := a.Store.GetJob(jobName, nil)
	if err != nil {
		return nil, fmt.Errorf("agent: Run error retrieving job: %s from store: %w", jobName, err)
	}

	// In case the job is not a child job, compute the next execution time
	if job.ParentJob == "" {
		if e, ok := a.sched.GetEntry(jobName); ok {
			job.Next = e.Next
			if err := a.applySetJob(job.ToProto()); err != nil {
				return nil, fmt.Errorf("agent: Run error storing job %s before running: %w", jobName, err)
			}
		} else {
			return nil, fmt.Errorf("agent: Run error retrieving job: %s from scheduler", jobName)
		}
	}

	// In the first execution attempt we build and filter the target nodes
	// but we use the existing node target in case of retry.
	var targetNodes []Node
	if ex.Attempt <= 1 {
		targetNodes = a.getTargetNodes(job.Tags, defaultSelector)
	} else {
		// In case of retrying, find the node or return with an error
		for _, m := range a.serf.Members() {
			if ex.NodeName == m.Name {
				if m.Status == serf.StatusAlive {
					targetNodes = []Node{m}
					break
				} else {
					return nil, fmt.Errorf("retry node is gone: %s for job %s", ex.NodeName, ex.JobName)
				}
			}
		}
	}

	// In case no nodes found, return reporting the error
	if len(targetNodes) < 1 {
		return nil, fmt.Errorf("no target nodes found to run job %s", ex.JobName)
	}
	a.logger.WithField("nodes", targetNodes).Debug("agent: Filtered nodes to run")

	var wg sync.WaitGroup
	for _, v := range targetNodes {
		// Determine node address
		addr, ok := v.Tags["rpc_addr"]
		if !ok {
			addr = v.Addr.String()
		}

		// Call here client GRPC AgentRun
		wg.Add(1)
		go func(node string, wg *sync.WaitGroup) {
			defer wg.Done()
			a.logger.WithFields(map[string]interface{}{
				"job_name": job.Name,
				"node":     node,
			}).Info("agent: Calling AgentRun")

			err := a.GRPCClient.AgentRun(node, job.ToProto(), ex.ToProto())
			if err != nil {
				a.logger.WithFields(map[string]interface{}{
					"job_name": job.Name,
					"node":     node,
				}).Error("agent: Error calling AgentRun")
			}
		}(addr, &wg)
	}

	wg.Wait()
	return job, nil
}
