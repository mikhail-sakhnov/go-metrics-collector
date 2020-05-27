package e2e

import (
	"context"
	"github.com/soider/go-metrics-collector/internal/agent"
	"github.com/soider/go-metrics-collector/internal/resultsprocessor"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestWorkflowSmokeTest(t *testing.T) {

	db := resultsprocessor.MustBuildDBClient(os.Getenv("POSTGRES_URI"))
	var results []resultsprocessor.ProbeResultRow
	err := db.Model(&results).Select()
	require.NoError(t, err)
	require.Len(t, results, 0, "Must have no records before e2e tests")

	ctx, cancel := context.WithCancel(context.Background())

	go agent.Loop(ctx,
		"/app/targets.example.yaml",
		[]string{},
		[]string{os.Getenv("KAFKA_BROKERS")},
		os.Getenv("KAFKA_TOPIC"),
		"/app/service.cert",
		"/app/service.key",
		"/app/ca.pem",
		5)
	go resultsprocessor.Loop(ctx,
		db,
		resultsprocessor.MustBuildKafkaReadClient(
			[]string{os.Getenv("KAFKA_BROKERS")},
			os.Getenv("KAFKA_TOPIC"),
			"/app/service.cert",
			"/app/service.key",
			"/app/ca.pem",
		))
	time.Sleep(time.Second * 30)
	cancel()
	err = db.Model(&results).Select()
	require.NoError(t, err)
	require.Greater(t, len(results), 0, "Must have results in the database after e2e tests")
}
