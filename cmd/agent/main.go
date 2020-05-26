package main

import (
	"github.com/soider/go-metrics-collector/internal/agent"
	"github.com/urfave/cli/v2"
	"log"
	"os"
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
			return agent.Main(c.String("targets"), c.StringSlice("selector"))
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
