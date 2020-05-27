package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table probe_results...")
		_, err := db.Exec(`CREATE TABLE probe_results(
		agent_name VARCHAR(100) NOT NULL,
		response_time bigint NOT NULL,
		http_code smallint NOT NULL,
		content_found bool DEFAULT false,
		timestamp timestamp NOT NULL DEFAULT NOW(),
		PRIMARY KEY(agent_name, timestamp)
	)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table probe_results...")
		_, err := db.Exec(`DROP TABLE probe_results`)
		return err
	})
}
