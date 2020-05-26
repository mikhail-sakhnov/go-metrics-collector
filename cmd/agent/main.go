package main

import (
	"github.com/soider/go-metrics-collector/internal/agent"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "selector", Aliases: []string{"s"}},
			&cli.StringFlag{Name: "targets", Aliases: []string{"t"}, DefaultText: "targets.yaml"},
		},
		Name:  "metrics-collector-agent",
		Usage: "starts agent",
		Action: func(c *cli.Context) error {
			targets := agent.MustReadTargets(
				c.String("targets"),
				agent.MustParseSelector(c.StringSlice("selector")))
			runningAgent := agent.NewMonitoringAgent(targets)
			stopCh := make(chan struct{})
			signalCh := make(chan os.Signal, 1)
			signal.Notify(signalCh, syscall.SIGINT)
			go func() {
				<-signalCh
				log.Print("Caught interruption signal, stopping agent")
				stopCh <- struct{}{}
			}()
			return runningAgent.Run(stopCh)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
