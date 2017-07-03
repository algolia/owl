package owl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setEnvironmentVariables(t *testing.T) {
	for _, envVar := range []string{"OWL_USE_METRIC", "OWL_USE_SENTRY", "OWL_USE_SLACK"} {
		err := os.Setenv(envVar, "true")
		require.Nil(t, err, "should set environment variable without error")
	}
}

func TestOwlInit_OK(t *testing.T) {
	setEnvironmentVariables(t)

	validLoggerConfig := &LoggerConfiguration{
		DisplayLogs: true,
		LogFilePath: "logs.json",
		SentryDsn:   "https://some-token@sentry.io/some-digits",
		UseColors:   true,
	}
	validMetricConfiguration := &MetricConfiguration{
		WavefrontUrl: "localhost:2878",
		StatsdUrl:    "127.0.0.1:8125",
	}
	validSlackConfiguration := &SlackConfiguration{
		Token: "some-token",
	}

	validConfigs := []Configuration{
		{AppName: "test", Logger: nil, Metric: nil, Slack: nil},
		{AppName: "test", Logger: validLoggerConfig, Metric: nil, Slack: nil},
		{AppName: "test", Logger: nil, Metric: validMetricConfiguration, Slack: nil},
		{AppName: "test", Logger: nil, Metric: nil, Slack: validSlackConfiguration},
	}

	for _, c := range validConfigs {
		err := Init(c)
		require.Nil(t, err, "should init `owl` without error")
	}
}

func TestOwlInit_KO_EmptyAppName(t *testing.T) {
	setEnvironmentVariables(t)

	c := Configuration{AppName: "", Logger: nil, Metric: nil, Slack: nil}
	err := Init(c)
	require.NotNil(t, err, "should not init `owl` without error with empty `AppName`")
	Stop()
}

func TestOwlInit_KO_EmptySentryDsn(t *testing.T) {
	setEnvironmentVariables(t)

	loggerConfig := &LoggerConfiguration{
		DisplayLogs: false,
		LogFilePath: "logs.json",
		SentryDsn:   "",
		UseColors:   false,
	}

	c := Configuration{AppName: "test", Logger: loggerConfig, Metric: nil, Slack: nil}
	err := Init(c)
	require.NotNil(t, err, "should not init `owl` without error with empty `AppName`")
	Stop()
}

func TestOwlInit_KO_EmptyStatsdWavefrontUrl(t *testing.T) {
	setEnvironmentVariables(t)

	metricConfig := &MetricConfiguration{
		WavefrontUrl: "",
		StatsdUrl:    "",
	}

	c := Configuration{AppName: "test", Logger: nil, Metric: metricConfig, Slack: nil}
	err := Init(c)
	require.NotNil(t, err, "should not init `owl` without error with both empty `StatsdUrl` and `WavefrontUrl`")
	Stop()
}

func TestOwlInit_KO_EmptySlackToken(t *testing.T) {
	setEnvironmentVariables(t)

	slackConfig := &SlackConfiguration{
		Token: "",
	}

	c := Configuration{AppName: "test", Logger: nil, Metric: nil, Slack: slackConfig}
	err := Init(c)
	require.NotNil(t, err, "should not init `owl` without error with empty Slack `Token`")
	Stop()
}
