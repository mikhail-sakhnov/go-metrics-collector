package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/pkg/env"
	"github.com/soider/go-metrics-collector/internal/pkg/message"
	"io"
	"log"
)

// NewKafkaWriterLoop builds new loop function and communication channel from the given kafkaWriter
func NewKafkaWriterLoop(kafkaWriter *kafka.Writer, failsThreshold int) (func(ctx context.Context) error, chan message.ProbeResultMessage) {
	resCh := make(chan message.ProbeResultMessage, env.GetEnvInt("KAFKA_IN_MEMORY_QUEUE_SIZE", 100))
	loop := func(ctx context.Context) error {
		fails := 0
		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-resCh:
				err := handleMessage(ctx, kafkaWriter, msg)
				switch errors.Unwrap(err) {
				case nil:
				case io.EOF:
					return fmt.Errorf("kafka connection error: %w", err)
				case context.Canceled:
					return nil
				default:
					fails++
					log.Print("error sending kafka message, re-send to the dead letter queue", err)
					if fails > failsThreshold {
						return fmt.Errorf("writer loop fails for %d times", fails)
					}
				}
			}
		}
		return nil
	}
	return loop, resCh
}

func handleMessage(ctx context.Context, client *kafka.Writer, msg message.ProbeResultMessage) error {
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't serialize message `%v`: %w", msg, err)
	}
	return client.WriteMessages(ctx, kafka.Message{
		Key:   uuid.NewV4().Bytes(),
		Value: rawMsg,
	})
}
