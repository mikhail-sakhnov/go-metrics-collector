package agent

import (
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/pkg/env"
	"github.com/soider/go-metrics-collector/internal/pkg/tls"
	"time"
)

// MustBuildKafkaWriteClient builds kafka write client or panics
func MustBuildKafkaWriteClient(brokers []string, topic, certFile, keyFile, caFile string) *kafka.Writer {
	cli, err := buildKafkaWriteClient(brokers, topic, certFile, keyFile, caFile)
	mustNotErr(err)
	return cli
}

func buildKafkaWriteClient(brokers []string, topic, certFile, keyFile, caFile string) (*kafka.Writer, error) {
	tlsConfig, err := tls.CreateTLSConfig(
		certFile, keyFile, caFile,
	)
	if err != nil {
		return nil, err
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       tlsConfig,
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:       brokers,
		Topic:         topic,
		Dialer:        dialer,
		QueueCapacity: env.GetEnvInt("KAFKA_WRITE_QUEUE_CAPACITY", 10),
		BatchSize:     env.GetEnvInt("KAFKA_WRITE_BATCH_SIZE", 1),
		BatchTimeout:  time.Second * time.Duration(env.GetEnvInt("KAFKA_WRITE_BATCH_TIMEOUT_SECONDS", 1)),
		Async:         false,
	})

	return w, nil
}
