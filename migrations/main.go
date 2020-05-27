package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"github.com/go-pg/migrations/v7"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "action", Aliases: []string{"a"}, Required: true, Usage: "Available actions: up, down, reset, version, set_version. See https://github.com/go-pg/migrations for the details"},
			&cli.StringFlag{Name: "databaseUri", Aliases: []string{"db"}, Required: true},
		},
		Name:  "metrics-collector-migrations",
		Usage: "manages db",
		Action: func(c *cli.Context) error {
			return migrate(c.String("databaseUri"), c.String("action"))
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func migrate(dbURI string, action string) error {
	opt, err := pg.ParseURL(dbURI)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)
	_, _, err = migrations.Run(db, action)
	return err
}
