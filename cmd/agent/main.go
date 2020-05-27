package main

import (
	"context"
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
			&cli.StringSliceFlag{Name: "brokers", Aliases: []string{"b"}, Required: true},
			&cli.StringFlag{Name: "targets", Aliases: []string{"t"}, DefaultText: "targets.yaml"},
			&cli.StringFlag{Name: "resultsTopic", Aliases: []string{"r"}, Required: true},
			&cli.StringFlag{Name: "certFile", Required: true},
			&cli.StringFlag{Name: "keyFile", Required: true},
			&cli.StringFlag{Name: "caFile", Required: true},
			&cli.IntFlag{Name: "failureThreshold", DefaultText: "5"},
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
				c.Int("failureThreshold"),
			)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Main application entry point
func run(targetsFilePath string, rawSelector []string, brokers []string, topic, certFile, keyFile, caFile string, failureThreshold int) error {
	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)
	go func() {
		select {
		case <-signalCh:
			log.Print("Caught interruption signal, stopping agent")
		}
		cancel()
	}()

	return agent.Loop(ctx, targetsFilePath, rawSelector, brokers, topic, certFile, keyFile, caFile, failureThreshold)
}
