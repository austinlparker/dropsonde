package model

import (
	"encoding/json"
)

type responseMsg struct {
	ResourceSpans   json.RawMessage `json:"resourceSpans,omitempty"`
	ResourceMetrics json.RawMessage `json:"resourceMetrics,omitempty"`
	ResourceLogs    json.RawMessage `json:"resourceLogs,omitempty"`
}
