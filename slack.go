package owl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func initSlack() error {
	if config.Slack == nil {
		useSlack = false
		return nil
	}

	if useSlack && config.Slack.Token == "" {
		return errors.New("owl: cannot use Slack logging with `Token` field empty")
	}

	return nil
}

func Slack(channel string, format string, args ...interface{}) {
	if !useSlack {
		return
	}

	message := fmt.Sprintf(format, args...)

	slackUrl := fmt.Sprintf(
		"https://slack.com/api/chat.postMessage?token=%s&channel=%s&as_user=true&text=%s",
		config.Slack.Token,
		channel,
		url.QueryEscape(message),
	)

	res, err := http.Get(slackUrl)
	if err == nil {
		res.Body.Close()
	}
}
