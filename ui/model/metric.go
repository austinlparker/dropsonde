package model

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"slices"
	"strconv"
	"strings"
)

type Timeseries struct {
	Name        string
	Description string
	Scope       pcommon.InstrumentationScope
	Resources   pcommon.Resource
	Attributes  pcommon.Map
	DataPoints  pmetric.NumberDataPointSlice
}

func TimeseriesToString(ts Timeseries) string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Name: %v", ts.Name))
	s.WriteString(fmt.Sprintf("\nDescription: %v", ts.Description))
	s.WriteString("\n\nScope Information")
	s.WriteString(fmt.Sprintf("\nName: %v", ts.Scope.Name()))
	s.WriteString(fmt.Sprintf("\nVersion: %v", ts.Scope.Version()))
	s.WriteString(fmt.Sprintf("\nAttributes: %v", ts.Scope.Attributes().AsRaw()))
	s.WriteString("\n\nResource Information")
	s.WriteString(fmt.Sprintf("\nAttributes: %v", ts.Resources.Attributes().AsRaw()))
	s.WriteString("\nAttributes\n")
	s.WriteString(fmt.Sprintf("\nAttributes: %v", ts.Attributes.AsRaw()))
	s.WriteString("\nData Points\n")
	for i := 0; i < ts.DataPoints.Len(); i++ {
		var dp string
		var t string
		t = ts.DataPoints.At(i).Timestamp().String()
		switch ts.DataPoints.At(i).ValueType() {
		case pmetric.NumberDataPointValueTypeInt:
			dp = strconv.FormatInt(ts.DataPoints.At(i).IntValue(), 10)
		case pmetric.NumberDataPointValueTypeDouble:
			dp = strconv.FormatFloat(ts.DataPoints.At(i).DoubleValue(), 'f', -1, 64)
		case pmetric.NumberDataPointValueTypeEmpty:
			dp = "Empty"
		}
		s.WriteString(fmt.Sprintf("\nTime: %v", t))
		s.WriteString(fmt.Sprintf("\nValue: %v\n", dp))
		s.WriteString(fmt.Sprintf("\nAttributes: %v\n", ts.DataPoints.At(i).Attributes().AsRaw()))
	}
	return s.String()
}

func ParseMetricMessage(msg WSMessage) (pmetric.Metrics, error) {
	unmarshaler := &pmetric.JSONUnmarshaler{}
	metrics, err := unmarshaler.UnmarshalMetrics(msg.data)
	if err != nil {
		return metrics, err
	}
	return metrics, nil
}

func newTimeseriesFromMetric(md pmetric.Metric, sm pcommon.InstrumentationScope, rs pcommon.Resource) Timeseries {
	var ts = Timeseries{
		Scope:      pcommon.NewInstrumentationScope(),
		Resources:  pcommon.NewResource(),
		Attributes: pcommon.NewMap(),
		DataPoints: pmetric.NewNumberDataPointSlice(),
	}
	ts.Name = md.Name()
	ts.Description = md.Description()
	sm.CopyTo(ts.Scope)
	rs.CopyTo(ts.Resources)

	switch md.Type() {
	case pmetric.MetricTypeGauge:

		md.Gauge().DataPoints().CopyTo(ts.DataPoints)
		for i := 0; i < md.Gauge().DataPoints().Len(); i++ {
			dp := ts.DataPoints.At(i)
			dp.Attributes().CopyTo(ts.Attributes)
		}
	case pmetric.MetricTypeSum:
		md.Sum().DataPoints().CopyTo(ts.DataPoints)
		for i := 0; i < md.Sum().DataPoints().Len(); i++ {
			dp := ts.DataPoints.At(i)
			dp.Attributes().CopyTo(ts.Attributes)
		}
	}

	return ts
}

func updateTimeseriesFromMetric(md pmetric.Metric, ts Timeseries) Timeseries {
	switch md.Type() {
	case pmetric.MetricTypeGauge:
		md.Gauge().DataPoints().MoveAndAppendTo(ts.DataPoints)
		for i := 0; i < md.Gauge().DataPoints().Len(); i++ {
			dp := ts.DataPoints.At(i)
			dp.Attributes().CopyTo(ts.Attributes)
		}
	case pmetric.MetricTypeSum:
		md.Sum().DataPoints().MoveAndAppendTo(ts.DataPoints)
		for i := 0; i < md.Sum().DataPoints().Len(); i++ {
			dp := ts.DataPoints.At(i)
			dp.Attributes().CopyTo(ts.Attributes)
		}
	}
	return ts
}

func AddOrUpdateMetricsToTimeseries(md *pmetric.Metrics, tl []Timeseries) []Timeseries {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			ilm := rm.ScopeMetrics().At(j)
			for k := 0; k < ilm.Metrics().Len(); k++ {
				metric := ilm.Metrics().At(k)
				name := metric.Name()
				idx := slices.IndexFunc(tl, func(ts Timeseries) bool {
					return ts.Name == name
				})
				if idx == -1 {
					timeseries := newTimeseriesFromMetric(metric, ilm.Scope(), rm.Resource())
					tl = append(tl, timeseries)
				} else {
					timeseries := updateTimeseriesFromMetric(metric, tl[idx])
					tl[idx] = timeseries
				}

			}
		}
	}
	return tl
}
