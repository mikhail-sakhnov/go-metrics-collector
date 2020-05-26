package probe

import (
	"context"
	"fmt"
	"github.com/soider/go-metrics-collector/internal/message"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// ProbeFunction type wrapper to use plain functions as agent.prober interface
type ProbeFunction func(ctx context.Context, from, uri string, searchFor *regexp.Regexp) (message.MonitoringResultMessage, error)

// Probe interface implementation
func (pf ProbeFunction) Probe(ctx context.Context, from, uri string, searchFor *regexp.Regexp) (message.MonitoringResultMessage, error) {
	return pf(ctx, from, uri, searchFor)
}

// HTTPProbe does http probe
var HTTPProbe = ProbeFunction(
	func(ctx context.Context, from string, uri string, searchFor *regexp.Regexp) (message.MonitoringResultMessage, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)

		if err != nil {
			return message.MonitoringResultMessage{}, fmt.Errorf("http probe request building failure for uri `%s`: %v", uri, err)
		}

		var duration time.Duration
		var resp *http.Response
		var httpErr error

		func() {
			start := time.Now()
			defer func() {
				duration = time.Since(start)
			}()
			resp, httpErr = http.DefaultClient.Do(req)
		}()

		if httpErr != nil {
			return message.MonitoringResultMessage{}, fmt.Errorf("http probe request execution failure for uri `%s`: %v", uri, err)
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return message.MonitoringResultMessage{}, fmt.Errorf("http probe read body failure for uri `%s`: %v", uri, err)
		}
		var found bool

		if searchFor != nil {
			found = searchFor.Find(data) != nil
		}

		return message.MonitoringResultMessage{
			ContentFound: found,
			ResponseTime: duration,
			HttpCode:     resp.StatusCode,
			AgentName:    from,
		}, nil
	})
