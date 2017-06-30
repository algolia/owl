package owl

type Configuration struct {
	AppName string               `json:"app_name"`
	Logger  *LoggerConfiguration `json:"logger"`
	Metric  *MetricConfiguration `json:"metric"`
	Slack   *SlackConfiguration  `json:"slack"`
}

type LoggerConfiguration struct {
	DisplayLogs bool   `json:"display_logs"`
	LogFilePath string `json:"log_file_path"`
	Logger      string `json:"logger"`
	SentryDsn   string `json:"sentry_dsn"`
	UseColors   bool   `json:"use_colors"`
}

type MetricConfiguration struct {
	StatsdUrl    string `json:"statsd_url"`
	WavefrontUrl string `json:"wavefront_url"`
}

type SlackConfiguration struct {
	Token string `json:"token"`
}

func checkConfiguration(c Configuration) error {
	// TODO
	return nil
}
