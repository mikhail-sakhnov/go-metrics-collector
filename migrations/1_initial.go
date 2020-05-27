package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table probe_results...")
		_, err := db.Exec(`CREATE TABLE probe_results(
		AgentName VARCHAR(100) NOT NULL,
		ResponseTime bigint NOT NULL,
		HTTPCode smallint NOT NULL,
		ContentFound bool DEFAULT false,
		TimeStamp timestamp NOT NULL DEFAULT NOW(),
		PRIMARY KEY(Agentname, TimeStamp)
	)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table probe_results...")
		_, err := db.Exec(`DROP TABLE probe_results`)
		return err
	})
}
