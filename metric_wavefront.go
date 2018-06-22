package owl

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront"
)

var wavefrontConfig *wavefront.WavefrontConfig

func initMetricWavefront() error {
	url := extractWavefrontUrl(config.Metric.WavefrontUrl)
	if url == "" {
		useMetric = false
		return Error("owl: invalid Wavefront proxy address (%s): should be URL:PORT or $URL:$PORT", config.Metric.WavefrontUrl)
	}

	addr, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		useMetric = false
		return Error("owl: cannot resolve Wavefront proxy address (%s): %s", config.Metric.WavefrontUrl, err)
	}

	var hostTags map[string]string
	if GitTag == "" {
		hostTags = nil
	} else {
		hostTags = map[string]string{"git_tag": GitTag}
	}

	wavefrontConfig = &wavefront.WavefrontConfig{
		Addr:         addr,
		Registry:     metrics.DefaultRegistry,
		DurationUnit: time.Nanosecond,
		Prefix:       config.AppName,
		HostTags:     hostTags,
		Percentiles:  []float64{0.5, 0.75, 0.95, 0.99, 0.999},
	}

	return nil
}

func extractWavefrontUrl(wavefrontUrl string) string {
	splits := strings.Split(wavefrontUrl, ":")
	if len(splits) != 2 {
		useMetric = false
		return ""
	}

	address, port := splits[0], splits[1]

	if address == "" || port == "" {
		return ""
	}

	if strings.HasPrefix(address, "$") && strings.HasPrefix(port, "$") {
		// Start with $ -> Load from the environment
		address = os.Getenv(address[1:])
		port = os.Getenv(port[1:])
		if address != "" && port != "" {
			return address + ":" + port
		}
	} else if !strings.HasPrefix(address, "$") && !strings.HasPrefix(port, "$") {
		// Start without $ -> Use the URL directly
		return wavefrontUrl
	}

	return ""
}

func metricIncWavefront(stat string, value int64, tags map[string]string) {
	if !useMetric || wavefrontConfig == nil {
		return
	}

	counter := wavefront.GetOrRegisterMetric(stat, metrics.NewCounter(), tags).(metrics.Counter)
	counter.Inc(value)
	wavefront.WavefrontSingleMetric(wavefrontConfig, stat, counter, tags)
}

func metricGaugeWavefront(stat string, value int64, tags map[string]string) {
	if !useMetric || wavefrontConfig == nil {
		return
	}

	gauge := wavefront.GetOrRegisterMetric(stat, metrics.NewGauge(), tags).(metrics.Gauge)
	gauge.Update(value)
	wavefront.WavefrontSingleMetric(wavefrontConfig, stat, gauge, tags)
}
