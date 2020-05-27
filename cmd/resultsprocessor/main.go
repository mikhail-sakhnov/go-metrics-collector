package main

import (
	"context"
	"github.com/soider/go-metrics-collector/internal/resultsprocessor"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "brokers", Aliases: []string{"b"}, Required: true},
			&cli.StringFlag{Name: "resultsTopic", Aliases: []string{"r"}, Required: true},
			&cli.StringFlag{Name: "certFile", Required: true},
			&cli.StringFlag{Name: "keyFile", Required: true},
			&cli.StringFlag{Name: "caFile", Required: true},
			&cli.StringFlag{Name: "databaseUri", Aliases: []string{"db"}, Required: true},
		},
		Name:  "metrics-collector-results-processor",
		Usage: "starts results processor",
		Action: func(c *cli.Context) error {
			return run(
				c.StringSlice("brokers"),
				c.String("resultsTopic"),
				c.String("certFile"),
				c.String("keyFile"),
				c.String("caFile"),
				c.String("databaseUri"),
			)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(brokers []string, topic string, certFile string, keyFile string, caFile string, databaseURI string) error {
	db := resultsprocessor.MustBuildDBClient(databaseURI)
	reader := resultsprocessor.MustBuildKafkaReadClient(brokers, topic, certFile, keyFile, caFile)
	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)
	go func() {
		defer cancel()
		select {
		case <-signalCh:
			log.Print("Caught interruption signal, stopping agent")
		}
	}()

	return resultsprocessor.Loop(ctx, db, reader)
}
