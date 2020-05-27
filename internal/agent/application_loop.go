package agent

import (
	"context"
	"github.com/soider/go-metrics-collector/internal/agent/probes"
	"golang.org/x/sync/errgroup"
)

// Loop main application loop
func Loop(ctx context.Context, targetsFilePath string, rawSelector []string, brokers []string, topic string, certFile, keyFile, caFile string, failureThreshold int) error {
	targets := MustReadTargets(
		targetsFilePath,
		MustParseSelector(rawSelector),
	)

	gr, ctx := errgroup.WithContext(ctx)

	loopFn, resCh := NewKafkaWriterLoop(MustBuildKafkaWriteClient(
		brokers,
		topic,
		certFile,
		keyFile,
		caFile,
	),
		failureThreshold,
	)

	runningAgent := NewMonitoringAgent(targets, resCh, probes.HTTPProbe, failureThreshold)
	gr.Go(func() error {
		return loopFn(ctx)
	})

	gr.Go(func() error {
		return runningAgent.Run(ctx)
	})

	return gr.Wait()

}
