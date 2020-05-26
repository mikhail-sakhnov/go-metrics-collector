package agent

import (
	"context"
	"fmt"
	"github.com/soider/go-metrics-collector/internal/message"
	"golang.org/x/sync/errgroup"
	"log"
	"regexp"
	"time"
)

// MonitoringAgent handles set of monitoring targets, in terms of application MonitoringAgent maps 1 to 1 with one running agent
type MonitoringAgent struct {
	targets          Targets
	resultsChan      chan message.MonitoringResultMessage
	failureThreshold int
	prober           prober
	writeTimeout     time.Duration
	readTimeout      time.Duration
}

type prober interface {
	Probe(ctx context.Context, from string, uri string, searchFor *regexp.Regexp) (message.MonitoringResultMessage, error)
}

// NewMonitoringAgent constructor
func NewMonitoringAgent(t Targets, resCh chan message.MonitoringResultMessage, prober prober) *MonitoringAgent {
	return &MonitoringAgent{
		targets:          t,
		resultsChan:      resCh,
		failureThreshold: 5, // TODO: make configurable
		writeTimeout:     time.Second * 1,
		readTimeout:      time.Second * 10,
		prober:           prober,
	}
}

// Run runs main loop for the coordinator
func (ma MonitoringAgent) Run(stopCh chan struct{}) error {
	log.Print("Start agent", ma.targets)
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
					log.Println("Doing job ", t.Name, t.URI, t.RegExp)
					msg, err := ma.prober.Probe(context.Background(), t.Name, t.URI, pattern)
					if err != nil {
						log.Printf("Agent `%s` probe failed %v for `%d` times", t.Name, err, fails)
						fails++
						if fails > ma.failureThreshold {
							return fmt.Errorf("agent `%s` has failed too many times", t.Name)
						}
					}
					ma.sendResult(msg)
				case <-stopCh:
					log.Print("Stopping ", t.URI)
					return nil
				}
			}
		})
	}
	err := errGr.Wait()
	fmt.Println(err)
	log.Print("Stop agent")
	return nil
}

func (ma MonitoringAgent) sendResult(msg message.MonitoringResultMessage) {
	select {
	case ma.resultsChan <- msg:
	case <-time.After(ma.writeTimeout):
		log.Printf("Can't save probe result, tried for %s, in memory buffer full", ma.writeTimeout)
	}
}
