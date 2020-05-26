package agent

import (
	"github.com/soider/go-metrics-collector/internal/agent/probe"
	"github.com/soider/go-metrics-collector/internal/message"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Main(targetsFilePath string, rawSelector []string) error {
	targets := MustReadTargets(
		targetsFilePath,
		MustParseSelector(rawSelector),
	)
	resCh := make(chan message.MonitoringResultMessage, 10)
	runningAgent := NewMonitoringAgent(targets, resCh, probe.HTTPProbe)
	stopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)
	go func() {
		<-signalCh
		log.Print("Caught interruption signal, stopping agent")
		close(stopCh)
	}()
	return runningAgent.Run(stopCh)
}
