package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/soider/go-metrics-collector/internal/pkg/message"
	"golang.org/x/sync/errgroup"
	"log"
	"regexp"
	"time"
)

// MonitoringAgent handles set of monitoring targets, in terms of application MonitoringAgent maps 1 to 1 with one running Agent
type MonitoringAgent struct {
	targets          Targets
	failureThreshold int

	prober    prober
	resultsCh chan message.ProbeResultMessage
}

type prober interface {
	Probe(ctx context.Context, from string, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error)
}

// NewMonitoringAgent constructor
func NewMonitoringAgent(t Targets, resCh chan message.ProbeResultMessage, prober prober, failureThreshold int) *MonitoringAgent {
	return &MonitoringAgent{
		targets:          t,
		resultsCh:        resCh,
		failureThreshold: failureThreshold,
		prober:           prober,
	}
}

// Run runs main loop for the coordinator
func (ma MonitoringAgent) Run(ctx context.Context) error {
	log.Print("Start Agent", ma.targets)
	var errGr errgroup.Group
	for _, t := range ma.targets {
		t := t // capture local var
		log.Print("Scheduling routine", t.Name)
		var pattern *regexp.Regexp
		if t.RegExp != "" {
			pattern = regexp.MustCompile(t.RegExp)
		}
		errGr.Go(func() error {
			ticker := time.NewTicker(t.Interval)
			fails := 0
			for {
				select {
				case <-ticker.C:
					log.Println("Doing job", t.Name, t.URI, t.RegExp)
					if t.Timeout > 0 {
						ctx, _ = context.WithTimeout(ctx, t.Timeout)
					}
					msg, err := ma.prober.Probe(ctx, t.Name, t.URI, pattern)
					switch errors.Unwrap(err) {
					case nil:
					case context.Canceled:
						return nil
					default:
						log.Printf("Agent `%s` probes failed %v for `%d` times", t.Name, err, fails)
						continue
					}

					if err := ma.sendResult(msg); err != nil {
						log.Printf("failed to send probe result: %s", err)
						fails++
						if fails > ma.failureThreshold {
							return ErrTooManyFailures{t.Name}
						}
					}
				case <-ctx.Done():
					return nil
				}
			}
		})
	}
	return errGr.Wait()
}

func (ma MonitoringAgent) sendResult(msg message.ProbeResultMessage) error {
	select {
	case ma.resultsCh <- msg:
	default:
		return fmt.Errorf("can't save probes result, in memory buffer full")
	}
	return nil
}
