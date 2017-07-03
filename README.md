![owl logo](https://raw.githubusercontent.com/algolia/owl/master/owl.png?token=AC-ItxuzISrqgsEX_VHR3ghHWUaMnKJ5ks5ZY2-fwA%3D%3D)

## Features

- Write (colored) logs on standard/error outputs
- Also write logs as [JSON lines](http://jsonlines.org/) on disk
- Automatically send logged errors to [Sentry](https://sentry.io/)
- Send metrics to Wavefront via StatsD proxy (running on Algolia API servers)
- Send metrics to Wavefront via Wavefront proxy (running on Kubernetes cluster)
- Send messages to Slack channel(s)

## Installation

## How to use?

Logging and monitoring should be easy to use in any project to avoid friction.
With this goal in mind, `owl` was written in a way it should be really
straightforward to use out of the box:

1. Generate a minimal `owl.Configuration` object (either in your code Go code
   or unmarshal it from a JSON file.
2. Pass it to the `owl.Init` function at the very beginning of your `main` and
   check its `error` return value to make sure the configuration was valid.
3. Put a `defer owl.Stop()` right after to make sure everything will clean up
   by itself when your program will terminate.

To prevent you from logging/sending metrics from your code when testing without
changing the your own code or the `owl` configuration, some features are only
enabled if specific `OWL_USE_*` environment variables are set.

## Examples

### Log messages and errors

```go
owlConfig := owl.Configuration{
	AppName: "my-project",
	Logger: &owl.LoggerConfiguration{
		DisplayLogs: true,
		LogFilePath: "logs.json",
		UseColors:   true,
	},
}

err := owl.Init(owlConfig)
defer owl.Stop()
if err != nil {
	fmt.Println(err)
}

owl.Info("this is an info")
owl.Warning("this is a warning %s", "anthony")
err = owl.Error("this is an error")
owl.Info("this error is displayed on the standard output: %s", err)
```

### Log errors to Sentry

To let you enable/disable Sentry easily without changing the configuration, the
`OWL_USE_SENTRY` environment variable must to be set in order to send errors to
Sentry.

You should find the Sentry DSN token [here](https://docs.sentry.io/clients/go/)
if you're logged to Sentry website with your Algolia Google account.

```go
owlConfig := owl.Configuration{
	AppName: "my-project",
	Logger: &owl.LoggerConfiguration{
		SentryDsn: "https://*****@sentry.io/124548",
	},
}

err := owl.Init(owlConfig)
defer owl.Stop()
if err != nil {
	fmt.Println(err)
}

owl.Info("this is an info message: it won't be logged to Sentry")
owl.Error("this is an error: it will be recorded by Sentry")
```

### Log metrics to Wavefront

To let you enable/disable Metric logging easily without changing the
configuration, the `OWL_USE_METRIC` environment variable must to be set in
order to metrics to Wavefront.

The following code shows how to configure `owl` and send few metrics to
Wavefront when your program runs on Algolia API servers, sending the metrics
via the StatsD proxy running. If you'd like to run it on a Kubernetes cluster
instead, replace the `StatsdUrl` configuration parameter with `WavefrontUrl`
and set it with the Wavefront proxy URL instead (usually `localhost:2878`). In
this case, you can use the exact same `Metric*` functions but also pass extra
maps of tags with `*WithTags` function variants.

```go
owlConfig := owl.Configuration{
	AppName: "my-project",
	Metric: &owl.MetricConfiguration{
		StatsdUrl: "127.0.0.1:8125",
	},
}

err := owl.Init(owlConfig)
defer owl.Stop()
if err != nil {
	fmt.Println(err)
}

t := owl.NewMetricTimer()

owl.MetricIncByOne("my_counter")
owl.MetricGauge("temporature", int64(30))

t.Stop("time_spent.main_function")
```

### Send message to a Slack channel

To let you enable/disable Slack logging easily without changing the
configuration, the `OWL_USE_SLACK` environment variable must to be set in
order to send messages to Slack.

You should generate a Slack token from [here](https://api.slack.com/tokens) if
you're logged to Slack with your Algolia Google account.

```go
owlConfig := owl.Configuration{
	AppName: "my-project",
	Slack: &owl.SlackConfiguration{
		Token: "your-slack-token",
	},
}

err := owl.Init(owlConfig)
defer owl.Stop()
if err != nil {
	fmt.Println(err)
}

owl.Slack("z-notif-operations", "This message will be seen in Slack")
owl.Slack("z-notif-operations", "And you can still use %s %s", "format", "strings")
```
