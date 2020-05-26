package agent

import (
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
	runningAgent := NewMonitoringAgent(targets)
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
