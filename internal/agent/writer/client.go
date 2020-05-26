package writer

import (
	"github.com/segmentio/kafka-go"
	"github.com/soider/go-metrics-collector/internal/pkg/tls"
	"log"
	"time"
)

// MustBuildKafkaWriteClient builds kafka write client or panics
func MustBuildKafkaWriteClient(brokers []string, topic, certFile, keyFile, caFile string) *kafka.Writer {
	cli, err := buildKafkaWriteClient(brokers, topic, certFile, keyFile, caFile)
	mustNotErr(err)
	return cli
}

func mustNotErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
		Brokers: brokers,
		Topic:   topic,
		Dialer:  dialer,
	})

	return w, nil
}
