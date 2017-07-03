package owl

import "github.com/cactus/go-statsd-client/statsd"

var (
	statsdClient     statsd.Statter
	statsdStatPrefix string
)

func initMetricStatsd() error {
	var err error

	if statsdClient, err = statsd.NewClient(config.Metric.StatsdUrl, config.AppName); err != nil {
		useMetric = false
		return Error("owl: cannot instantiate client to StatsD proxy (%s): %s", config.Metric.StatsdUrl, err)
	}

	statsdStatPrefix = config.AppName + "."

	return nil
}

func metricIncStatsd(stat string, value int64) {
	err := statsdClient.Inc(statsdStatPrefix+stat, value, 1.0)
	if err != nil {
		Warning("owl: cannot increment metric %s of %d in StatsD: %s", stat, value, err)
	}
}

func metricGaugeStatsd(stat string, value int64) {
	err := statsdClient.Gauge(statsdStatPrefix+stat, value, 1.0)
	if err != nil {
		Warning("owl: cannot gauge metric %s with %s in StatsD: %s", stat, value, err)
	}
}
