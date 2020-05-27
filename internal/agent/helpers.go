package agent

import (
	"log"
)

func mustNotErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
