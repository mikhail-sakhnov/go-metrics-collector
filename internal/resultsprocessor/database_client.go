package resultsprocessor

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"log"
)

// MustBuildDBClient builds postgres client, panics on failure
func MustBuildDBClient(databaseURI string) *pg.DB {
	db, err := buildDBClient(databaseURI)
	mustNotErr(err)
	return db
}

func buildDBClient(databaseURI string) (*pg.DB, error) {
	opts, err := pg.ParseURL(databaseURI)
	if err != nil {
		return nil, fmt.Errorf("can't parse postgres connection uri `%s`: %w", databaseURI, err)
	}
	return pg.Connect(opts), nil
}

func mustNotErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
