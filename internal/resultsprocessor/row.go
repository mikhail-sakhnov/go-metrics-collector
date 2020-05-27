package resultsprocessor

import "time"

// ProbeResultRow represent single tuple in the probe_results table
type ProbeResultRow struct {
	tableName    struct{}      `pg:"probe_results"`
	AgentName    string        `pg:"agent_name"`
	HTTPCode     int           `pg:"http_code"`
	ContentFound bool          `pg:"content_found"`
	Timestamp    time.Time     `pg:"timestamp"`
	ResponseTime time.Duration `pg:"response_time"`
}
