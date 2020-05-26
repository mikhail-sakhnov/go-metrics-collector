package message

import "time"

// ProbeResultMessage represents single probes result
type ProbeResultMessage struct {
	AgentName    string        `json:"agent_name"`
	ResponseTime time.Duration `json:"response_time"`
	HTTPCode     int           `json:"http_code"`
	ContentFound bool          `json:"content_found"`
}
