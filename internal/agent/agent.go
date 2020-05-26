package agent

import "log"

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
func (MonitoringAgent) Run(stopCh chan struct{}) error {
	log.Print("Start agent")
	select {
	case <-stopCh:
	}
	log.Print("Stop agent")
	return nil
}
