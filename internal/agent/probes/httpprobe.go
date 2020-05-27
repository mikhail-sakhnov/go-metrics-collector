package probes

import (
	"context"
	"errors"
	"fmt"
	"github.com/soider/go-metrics-collector/internal/pkg/message"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// HTTPProbe does http probes
var HTTPProbe = Function(
	func(ctx context.Context, from string, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)

		if err != nil {
			return message.ProbeResultMessage{}, fmt.Errorf("http probes request building failure for uri `%s`: %w", uri, err)
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
		switch errors.Unwrap(httpErr) {
		case nil:
		case context.Canceled:
			return message.ProbeResultMessage{}, httpErr
		default:
			return message.ProbeResultMessage{}, fmt.Errorf("http probes request execution failure for uri `%s`: %w", uri, httpErr)
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return message.ProbeResultMessage{}, fmt.Errorf("http probes read body failure for uri `%s`: %w", uri, err)
		}
		var found bool

		if searchFor != nil {
			found = searchFor.Find(data) != nil
		}

		return message.ProbeResultMessage{
			ContentFound: found,
			ResponseTime: duration,
			HTTPCode:     resp.StatusCode,
			AgentName:    from,
		}, nil
	})
