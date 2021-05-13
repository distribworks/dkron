package dkron

import (
	"fmt"
	"log"
	"strings"
	"time"

	discover "github.com/hashicorp/go-discover"
	discoverk8s "github.com/hashicorp/go-discover/provider/k8s"
	"github.com/sirupsen/logrus"
)

func (a *Agent) retryJoinLAN() {
	r := &retryJoiner{
		cluster:     "LAN",
		addrs:       a.config.RetryJoinLAN,
		maxAttempts: a.config.RetryJoinMaxAttemptsLAN,
		interval:    a.config.RetryJoinIntervalLAN,
		join:        a.JoinLAN,
	}
	if err := r.retryJoin(a.logger); err != nil {
		a.retryJoinCh <- err
	}
}

// retryJoiner is used to handle retrying a join until it succeeds or all
// retries are exhausted.
type retryJoiner struct {
	// cluster is the name of the serf cluster, e.g. "LAN" or "WAN".
	cluster string

	// addrs is the list of servers or go-discover configurations
	// to join with.
	addrs []string

	// maxAttempts is the number of join attempts before giving up.
	maxAttempts int

	// interval is the time between two join attempts.
	interval time.Duration

	// join adds the discovered or configured servers to the given
	// serf cluster.
	join func([]string) (int, error)
}

func (r *retryJoiner) retryJoin(logger *logrus.Entry) error {
	if len(r.addrs) == 0 {
		return nil
	}

	// Copy the default providers, and then add the non-default
	providers := make(map[string]discover.Provider)
	for k, v := range discover.Providers {
		providers[k] = v
	}
	providers["k8s"] = &discoverk8s.Provider{}

	disco, err := discover.New(
		discover.WithUserAgent(UserAgent()),
		discover.WithProviders(providers),
	)
	if err != nil {
		return err
	}

	logger.Infof("agent: Retry join %s is supported for: %s", r.cluster, strings.Join(disco.Names(), " "))
	logger.WithField("cluster", r.cluster).Info("agent: Joining cluster...")
	attempt := 0
	for {
		var addrs []string
		var err error

		for _, addr := range r.addrs {
			switch {
			case strings.Contains(addr, "provider="):
				servers, err := disco.Addrs(addr, log.New(logger.Logger.Writer(), "", log.LstdFlags|log.Lshortfile))
				if err != nil {
					logger.WithError(err).WithField("cluster", r.cluster).Error("agent: Error Joining")
				} else {
					addrs = append(addrs, servers...)
					logger.Infof("agent: Discovered %s servers: %s", r.cluster, strings.Join(servers, " "))
				}

			default:
				ipAddr, err := ParseSingleIPTemplate(addr)
				if err != nil {
					logger.WithField("addr", addr).WithError(err).Error("agent: Error parsing retry-join ip template")
					continue
				}
				addrs = append(addrs, ipAddr)
			}
		}

		if len(addrs) > 0 {
			n, err := r.join(addrs)
			if err == nil {
				logger.Infof("agent: Join %s completed. Synced with %d initial agents", r.cluster, n)
				return nil
			}
		}

		if len(addrs) == 0 {
			err = fmt.Errorf("no servers to join")
		}

		attempt++
		if r.maxAttempts > 0 && attempt > r.maxAttempts {
			return fmt.Errorf("agent: max join %s retry exhausted, exiting", r.cluster)
		}

		logger.Warningf("agent: Join %s failed: %v, retrying in %v", r.cluster, err, r.interval)
		time.Sleep(r.interval)
	}
}
