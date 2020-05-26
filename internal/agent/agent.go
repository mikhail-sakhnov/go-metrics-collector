package agent

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

// MonitoringAgent handles set of monitoring targets, in terms of application MonitoringAgent maps 1 to 1 with one running agent
type MonitoringAgent struct {
	targets Targets
}

// NewMonitoringAgent constructor
func NewMonitoringAgent(t Targets) *MonitoringAgent {
	return &MonitoringAgent{
		targets: t,
	}
}

// Run runs main loop for the coordinator
func (ma MonitoringAgent) Run(stopCh chan struct{}) error {
	log.Print("Start agent", ma.targets)
	var errGr errgroup.Group
	for _, t := range ma.targets {
		t := t // capture local var
		log.Print("Scheduling routine", t.URI)
		errGr.Go(func() error {
			ticker := time.NewTicker(t.Interval)
			for {
				select {
				case <-ticker.C:
					log.Println("Doing job", t.URI, t.RegExp)
				case <-stopCh:
					log.Print("Stopping", t.URI)
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
