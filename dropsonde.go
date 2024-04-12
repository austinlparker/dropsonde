package main

import (
	"encoding/json"
	"github.com/austinlparker/dropsonde/cmd"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type responseMsg struct {
	ResourceSpans   json.RawMessage `json:"resourceSpans,omitempty"`
	ResourceMetrics json.RawMessage `json:"resourceMetrics,omitempty"`
	ResourceLogs    json.RawMessage `json:"resourceLogs,omitempty"`
}

func (m *model) assignMessage(msg tapMessage) {
	if msg.data != nil {
		var response responseMsg
		err := json.Unmarshal(msg.data, &response)
		if err != nil {
			tea.Println("Error unmarshaling message", err)
			return
		}
		if response.ResourceSpans != nil {
			unmarshaler := &ptrace.JSONUnmarshaler{}
			span, err := unmarshaler.UnmarshalTraces(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.traces = append(m.traces, span)
		}
		if response.ResourceMetrics != nil {
			unmarshaler := &pmetric.JSONUnmarshaler{}
			metrics, err := unmarshaler.UnmarshalMetrics(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.metrics = append(m.metrics, metrics)
		}
		if response.ResourceLogs != nil {
			unmarshaler := &plog.JSONUnmarshaler{}
			log, err := unmarshaler.UnmarshalLogs(msg.data)
			if err != nil {
				tea.Println("Error unmarshaling message", err)
				return
			}
			m.logs = append(m.logs, log)
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForMessages(m.wsConn, m.channel),
		waitForMessages(m.channel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := tmsg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	case tapMessage:
		m.assignMessage(msg.(tapMessage))
		return m, waitForMessages(m.channel)
	}
	return m, nil
}

func (m model) View() string {

}

func (m model) renderMetrics() string {
	metricString := strings.Builder{}
	for i := 0; i < len(m.metrics); i++ {
		metric := m.metrics[i]
		for j := 0; j < metric.ResourceMetrics().Len(); j++ {
			resourceMetric := metric.ResourceMetrics().At(j)
			scopeMetrics := resourceMetric.ScopeMetrics()
			for k := 0; k < scopeMetrics.Len(); k++ {
				scopeMetric := scopeMetrics.At(k)
				name := scopeMetric.Scope().Name()
				metricString.WriteString(name)
				for v := 0; v < scopeMetrics.Len(); v++ {
					metrics := scopeMetric.Metrics()
					for l := 0; l < metrics.Len(); l++ {
						value := metrics.At(l)
						metricString.WriteString(value.Name())
						metricString.WriteString("\n")
					}
				}
			}
		}
	}
	return metricString.String()
}

func (m model) renderMetricsBetter() string {
	s := strings.Builder{}
	for _, metric := range m.metrics {
		resources := metric.ResourceMetrics()
		for i := 0; i < resources.Len(); i++ {
			resourceMetric := resources.At(i)
			scopeMetrics := resourceMetric.ScopeMetrics()
			for j := 0; j < scopeMetrics.Len(); j++ {
				sm := scopeMetrics.At(j)
				s.WriteString(sm.Scope().Name())
				s.WriteString(" | ")
				for k := 0; k < sm.Metrics().Len(); k++ {
					m := sm.Metrics().At(k)
					s.WriteString(m.Name())
					s.WriteString(": ")
					s.WriteString(m.Description())
					s.WriteString("\n")
				}
			}
		}
	}
	return s.String()
}

func (m model) renderTraces() string {
	traceString := strings.Builder{}
	for i := 0; i < len(m.traces); i++ {
		rs := m.traces[i].ResourceSpans()
		for j := 0; j < rs.Len(); j++ {
			ss := rs.At(j).ScopeSpans()
			for k := 0; k < ss.Len(); k++ {
				spanList := ss.At(k)
				spans := spanList.Spans()
				for v := 0; v < spans.Len(); v++ {
					span := spans.At(v)
					traceString.WriteString(span.Name())
					traceString.WriteString(" ")
					traceString.WriteString(span.TraceID().String())
					traceString.WriteString("\n")
				}
			}
		}
	}
	return traceString.String()
}

func main() {
	cmd.Execute()
}
