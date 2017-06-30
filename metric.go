package owl

import "time"

func initMetric() {
	if config.Metric == nil {
		useMetric = false
		return
	}

	if config.Metric.StatsdUrl != "" {
		initMetricStatsd()
	}

	if config.Metric.WavefrontUrl != "" {
		initMetricWavefront()
	}
}

type MetricTimer struct {
	start time.Time
}

func NewMetricTimer() *MetricTimer {
	return &MetricTimer{
		start: time.Now(),
	}
}

func (t *MetricTimer) Stop(stat string) {
	t.StopWithTags(stat, nil)
}

func (t *MetricTimer) StopWithTags(stat string, tags map[string]string) {
	MetricDurationWithTags(stat, time.Now().Sub(t.start), tags)
}

func MetricDuration(stat string, delta time.Duration) {
	MetricDurationWithTags(stat, delta, nil)
}

func MetricDurationWithTags(stat string, delta time.Duration, tags map[string]string) {
	MetricGaugeWithTags(stat, int64(delta.Nanoseconds()/1000000), tags)
}

func MetricIncByOne(stat string) {
	MetricIncByOneWithTags(stat, nil)
}

func MetricIncByOneWithTags(stat string, tags map[string]string) {
	MetricIncWithTags(stat, 1, tags)
}

func MetricInc(stat string, value int64) {
	MetricIncWithTags(stat, value, nil)
}

func MetricIncWithTags(stat string, value int64, tags map[string]string) {
	if !useMetric {
		return
	}
	metricIncWavefront(stat, value, tags)
	metricIncStatsd(stat, value)
}

func MetricGauge(stat string, value int64) {
	MetricGaugeWithTags(stat, value, nil)
}

func MetricGaugeWithTags(stat string, value int64, tags map[string]string) {
	if !useMetric {
		return
	}
	metricGaugeWavefront(stat, value, tags)
	metricGaugeStatsd(stat, value)
}
