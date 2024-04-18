package parsers

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"slices"
)

type Timeseries struct {
	Name        string
	Description string
	Scope       pcommon.InstrumentationScope
	Resources   pcommon.Resource
	Attributes  pcommon.Map
	DataPoints  pmetric.NumberDataPointSlice
}

func ParseMetricMessage(msg []byte) (pmetric.Metrics, error) {
	unmarshaler := &pmetric.JSONUnmarshaler{}
	metrics, err := unmarshaler.UnmarshalMetrics(msg)
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
