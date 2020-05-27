package resultsprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/pkg/message"
	"log"
)

// Loop reads messages from kafka and saves them to the postgres
func Loop(ctx context.Context, db *pg.DB, reader *kafka.Reader) error {
	var msg message.ProbeResultMessage
	for {
		raw, err := reader.ReadMessage(ctx)
		switch errors.Unwrap(err) {
		case nil:
		case context.Canceled:
			return nil
		default:
			return fmt.Errorf("can't read data from kafka: %w", err)
		}

		if err := json.Unmarshal(raw.Value, &msg); err != nil {
			handleBadMessage(raw.Value)
		}
		log.Printf("Processing incoming message: %v", msg)
		if err := db.Insert(&ProbeResultRow{
			AgentName:    msg.AgentName,
			ContentFound: msg.ContentFound,
			ResponseTime: msg.ResponseTime,
			HTTPCode:     msg.HTTPCode,
		}); err != nil {
			return fmt.Errorf("can't save result: %w", err)
		}
	}
	return nil
}

func handleBadMessage(message []byte) {
	log.Printf("[BAD FORMAT]: %s", message)
}
