package probes

import (
	"context"
	"github.com/soider/go-metrics-collector/internal/message"
	"regexp"
)

// Function type wrapper to use plain functions as agent.prober interface
type Function func(ctx context.Context, from, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error)

// Probe interface implementation
func (pf Function) Probe(ctx context.Context, from, uri string, searchFor *regexp.Regexp) (message.ProbeResultMessage, error) {
	return pf(ctx, from, uri, searchFor)
}
