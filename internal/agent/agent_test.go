package agent

import (
	"context"
	"fmt"
	"github.com/soider/go-metrics-collector/internal/agent/probes"
	"github.com/soider/go-metrics-collector/internal/message"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

type dummyProber struct {
	callLog  map[string][]time.Time
	latency  map[string]time.Duration
	found    map[string]bool
	httpCode map[string]int
}

func (dp dummyProber) Probe(ctx context.Context, from string, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error) {
	dp.callLog[from] = append(dp.callLog[from], time.Now())
	return message.ProbeResultMessage{
		AgentName:    from,
		ResponseTime: dp.latency[uri],
		HTTPCode:     dp.httpCode[uri],
		ContentFound: dp.found[uri],
	}, nil

}

func TestMonitoringAgent(t *testing.T) {
	t.Run("smoke_test_monitoring_agent_run_loop_saves_probe_results_and_has_no_error", func(t *testing.T) {
		targets := Targets{
			Target{
				Name:     "test_target_1",
				URI:      "http://uri1/",
				Interval: 300 * time.Millisecond,
				RegExp:   "",
				Selector: nil,
			},
			Target{
				Name:     "test_target_2",
				URI:      "http://uri2/",
				Interval: 250 * time.Millisecond,
				RegExp:   ".*body.*",
				Selector: nil,
			},
		}
		resCh := make(chan message.ProbeResultMessage, 100)
		stopCh := make(chan struct{})
		prober := dummyProber{
			callLog: map[string][]time.Time{},
			latency: map[string]time.Duration{
				"http://uri1/": time.Millisecond * 300,
				"http://uri2/": time.Millisecond * 200,
			},
			found: map[string]bool{
				"http://uri1/": false,
				"http://uri2/": true,
			},
			httpCode: map[string]int{
				"http://uri1/": 301,
				"http://uri2/": 400,
			},
		}
		agent := NewMonitoringAgent(targets, resCh, prober)
		go func() {
			time.Sleep(time.Second)
			close(stopCh)
		}()

		err := agent.Run(stopCh)
		close(resCh)

		assert.NoError(t, err, "Must have no error")
		assert.True(t, len(resCh) > 0, "Must have probe results")

		t.Log("check call intervals")
		for i := 1; i < len(prober.callLog["test_target_1"]); i++ {
			runInterval := prober.callLog["test_target_1"][i].Sub(prober.callLog["test_target_1"][i-1])
			assert.True(t, runInterval > 290) // 10 ms should be more than enough to deal with scheduler glitches
		}

		for i := 1; i < len(prober.callLog["test_target_2"]); i++ {
			runInterval := prober.callLog["test_target_2"][i].Sub(prober.callLog["test_target_2"][i-1])
			assert.True(t, runInterval > 240) // 10 ms should be more than enough to deal with scheduler glitches
		}
	})

	t.Run("smoke_test_monitoring_agent_run_loop_stops_after_N_errors", func(t *testing.T) {
		targets := Targets{
			Target{
				Name:     "test_target_1",
				URI:      "http://uri1/",
				Interval: time.Millisecond,
				RegExp:   "",
				Selector: nil,
			},
		}
		resCh := make(chan message.ProbeResultMessage, 100)
		stopCh := make(chan struct{})
		prober := probes.Function(func(ctx context.Context, from, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error) {
			return message.ProbeResultMessage{}, fmt.Errorf("error from the mocked prober for uri `%s`, Agent `%s`", uri, from)
		})
		agent := NewMonitoringAgent(targets, resCh, prober)
		agent.failureThreshold = 3
		go func() {
			time.Sleep(time.Millisecond * 10)
			close(stopCh)
		}()

		err := agent.Run(stopCh)
		close(resCh)

		assert.Error(t, err, "Must have error")
		assert.Equal(t, err.(ErrTooManyFailures).Agent, "test_target_1")
		assert.True(t, len(resCh) == 0, "Must have no probe results")
	})
}
