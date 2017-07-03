![owl logo](https://raw.githubusercontent.com/algolia/owl/master/owl.png?token=AC-ItxuzISrqgsEX_VHR3ghHWUaMnKJ5ks5ZY2-fwA%3D%3D)

## Features

- Write (colored) logs on standard/error outputs
- Also write logs as [JSON lines](http://jsonlines.org/) on disk
- Automatically send logged errors to [Sentry](https://sentry.io/)
- Send metrics to Wavefront via StatsD proxy (running on Algolia API servers)
- Send metrics to Wavefront via Wavefront proxy (running on Kubernetes cluster)
- Send messages to Slack channel(s)

## Installation for your own project

### Local development

As you should already have access to the Github Algolia organization, you
should be able to access https://github.com/algolia/owl. If so, it means you
can simply clone it in your Go source tree by using:

```
git clone git@github.com:algolia/owl.git $GOPATH/src/github.com/algolia/owl
```

### Travis

As `owl` is a private Gitub repository, TravisCI will not be able to download
this project directly as a dependency of your own project. The following
section will guide you through the generation and upload of a dedicated SSH key
into both your Travis build and `owl` repo.


First, you need to generate a public/private key pair that will be used to let
your project access the `github.com/algolia/owl`:

```
ssh-keygen -t rsa -b 4096 -f travis_key_to_owl -P ''
```

Then add your freshly generated public key (the `travis_key_to_owl.pub`
file) to the list of authorized deploy keys of the `owl` repository
[here](https://github.com/algolia/owl/settings/keys). Click on `Add deploy
key`, fill the `Title` field with the name of your project and paste the
content of the public key in the `Key` field. Make sure the `Allow write access
checkbox` is left unchecked.

Now go in your own project repository directory, and add your private key to
the Travis configuration file using:

```
cd /path/to/your/own/project
travis encrypt-file /path/to/travis_key_to_owl --add
```

Once it's done, your project's `.travis.yaml` file will include a new line in
the `before_install` section which is there to decrypt your private SSH key. To
let Git make use of it from your future builds, you need to make the SSH agent
load it. To do it, just add those three lines after the `openssl` command of
the `before_install` section of your `.travis.yaml` file:

```
before_install:
- openssl ...
- chmod 600 travis_access_to_owl
- eval "$(ssh-agent)"
- ssh-add travis_access_to_owl
```

Finally, simply commit and push your changes in Git:

```
git add travis_key_to_owl.enc .travis.yaml
git ci -m 'chore: Let Travis' ssh-agent to access github.com/algolia/owl'
```

You can now safely delete your `travis_access_to_owl` (private key) and
`travis_key_to_owl.pub` (public key) files as they are correctly stored and
should not be shared by anything else than Github and Travis.

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
