package owl

import (
	"net"
	"time"

	"github.com/aseure/go-metrics-wavefront"
	"github.com/rcrowley/go-metrics"
)

var wavefrontConfig *wavefront.WavefrontConfig

func initMetricWavefront() {
	addr, err := net.ResolveTCPAddr("tcp", config.Metric.WavefrontUrl)
	if err != nil {
		useMetric = false
		Error("owl: cannot resolve Wavefront proxy address (%s): %s", config.Metric.WavefrontUrl, err)
		return
	}

	wavefrontConfig = &wavefront.WavefrontConfig{
		Addr:         addr,
		Registry:     metrics.DefaultRegistry,
		DurationUnit: time.Nanosecond,
		Prefix:       config.AppName,
		HostTags:     map[string]string{"git_tag": GitTag},
		Percentiles:  []float64{0.5, 0.75, 0.95, 0.99, 0.999},
	}
}

func metricIncWavefront(stat string, value int64, tags map[string]string) {
	mergedTags := mergeTags(tags, wavefrontConfig.HostTags)

	counter := wavefront.GetOrRegisterMetric(stat, metrics.NewCounter(), mergedTags).(metrics.Counter)
	counter.Inc(value)
	wavefront.WavefrontSingleMetric(wavefrontConfig, stat, counter, mergedTags)
}

func metricGaugeWavefront(stat string, value int64, tags map[string]string) {
	mergedTags := mergeTags(tags, wavefrontConfig.HostTags)

	gauge := wavefront.GetOrRegisterMetric(stat, metrics.NewGauge(), mergedTags).(metrics.Gauge)
	gauge.Update(value)
	wavefront.WavefrontSingleMetric(wavefrontConfig, stat, gauge, mergedTags)
}

func mergeTags(m1, m2 map[string]string) map[string]string {
	tags := make(map[string]string)

	if m1 != nil {
		for k, v := range m1 {
			tags[k] = v
		}
	}

	if m2 != nil {
		for k, v := range m2 {
			tags[k] = v
		}
	}

	return tags
}
