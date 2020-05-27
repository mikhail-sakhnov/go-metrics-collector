package resultsprocessor

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/pkg/tls"
)

// MustBuildKafkaReadClient builds kafka read client, panics on failure
func MustBuildKafkaReadClient(brokers []string, topic, certFile, keyFile, caFile string) *kafka.Reader {
	cli, err := buildKafkaReadClient(brokers, topic, certFile, keyFile, caFile)
	mustNotErr(err)
	return cli
}

func buildKafkaReadClient(brokers []string, topic, certFile, keyFile, caFile string) (*kafka.Reader, error) {
	tlsConfig, err := tls.CreateTLSConfig(certFile, keyFile, caFile)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka read client TLS config: %w", err)
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,

		Dialer: &kafka.Dialer{
			TLS: tlsConfig,
		},
		Topic:    topic,
		GroupID:  "results-processor-cg",
		MinBytes: 100,
		MaxBytes: 500,
	})
	return reader, nil
}
