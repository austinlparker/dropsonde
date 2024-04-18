package parsers

import "go.opentelemetry.io/collector/pdata/ptrace"

func ParseTraceMessage(msg []byte) (ptrace.Traces, error) {
	unmarshaler := &ptrace.JSONUnmarshaler{}
	traces, err := unmarshaler.UnmarshalTraces(msg)
	if err != nil {
		return traces, err
	}
	return traces, nil
}

func AddToSpanList(traces ptrace.Traces, list []ptrace.Span) []ptrace.Span {
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		rs := traces.ResourceSpans().At(i)
		for j := 0; j < rs.ScopeSpans().Len(); j++ {
			ss := rs.ScopeSpans().At(j)
			for k := 0; k < ss.Spans().Len(); k++ {
				span := ss.Spans().At(k)
				list = append(list, span)
			}
		}
	}
	return list
}
