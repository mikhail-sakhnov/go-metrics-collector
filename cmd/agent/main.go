package main

import (
	"github.com/soider/go-metrics-collector/internal/agent"
	"github.com/soider/go-metrics-collector/internal/agent/probes"
	"github.com/soider/go-metrics-collector/internal/agent/writer"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "selector", Aliases: []string{"s"}},
			&cli.StringSliceFlag{Name: "brokers", Aliases: []string{"b"}, Required: true},
			&cli.StringFlag{Name: "targets", Aliases: []string{"t"}, DefaultText: "targets.yaml"},
			&cli.StringFlag{Name: "resultsTopic", Aliases: []string{"r"}, Required: true},
			&cli.StringFlag{Name: "certFile", Required: true},
			&cli.StringFlag{Name: "keyFile", Required: true},
			&cli.StringFlag{Name: "caFile", Required: true},
		},
		Name:  "metrics-collector-agent",
		Usage: "starts agent",
		Action: func(c *cli.Context) error {
			return run(c.String("targets"),
				c.StringSlice("selector"),
				c.StringSlice("brokers"),
				c.String("resultsTopic"),
				c.String("certFile"),
				c.String("keyFile"),
				c.String("caFile"),
			)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Main application entry point
func run(targetsFilePath string, rawSelector []string, brokers []string, topic, certFile, keyFile, caFile string) error {
	targets := agent.MustReadTargets(
		targetsFilePath,
		agent.MustParseSelector(rawSelector),
	)
	var gr errgroup.Group
	stopCh := make(chan struct{})
	loopFn, resCh := writer.NewKafkaWriterLoop(writer.MustBuildKafkaWriteClient(
		brokers,
		topic,
		certFile,
		keyFile,
		caFile,
	), stopCh)

	runningAgent := agent.NewMonitoringAgent(targets, resCh, probes.HTTPProbe)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go func() {
		<-signalCh
		log.Print("Caught interruption signal, stopping agent")
		close(stopCh)
	}()


	gr.Go(func() error {
		return loopFn()
	})

	gr.Go(func() error {
		return runningAgent.Run(stopCh)
	})
	return gr.Wait()
}
