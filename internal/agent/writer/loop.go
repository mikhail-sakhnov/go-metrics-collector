package writer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/message"
	"io"
	"log"
)

const failsThreshold = 2

// NewKafkaWriterLoop builds new loop function and communication channel from the given kafkaWriter
func NewKafkaWriterLoop(kafkaWriter *kafka.Writer, stopCh chan struct{}) (func() error, chan message.ProbeResultMessage) {
	resCh := make(chan message.ProbeResultMessage)
	loop := func() error {
		fails := 0
		for {
			select {
			case <-stopCh:
				return nil
			case msg := <-resCh:
				err := handleMessage(kafkaWriter, msg)
				if err != nil {
					if errors.Is(err, io.EOF) { // connection closed
						return fmt.Errorf("kafka connection error: %w", err)
					}
					fails++
					if fails > failsThreshold {
						return fmt.Errorf("writer loop fails for %d times", fails)
					}
					log.Print("error sending kafka message, re-send to the dead letter queue", err)
				}
			}
		}
		return nil
	}
	return loop, resCh
}

func handleMessage(client *kafka.Writer, msg message.ProbeResultMessage) error {
	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't serialize message `%v`: %w", msg, err)
	}
	return client.WriteMessages(context.Background(), kafka.Message{
		Key:   uuid.NewV4().Bytes(),
		Value: rawMsg,
	})
}
