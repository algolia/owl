package owl

import (
	"errors"
	"os"
)

var (
	config    Configuration
	useMetric bool
	useSentry bool
	useSlack  bool
)

func Init(c Configuration) error {
	if c.AppName == "" {
		return errors.New("owl: `AppName` configuration field cannot be empty")
	}
	config = c

	if os.Getenv("OWL_USE_METRIC") != "" {
		useMetric = true
	}

	if os.Getenv("OWL_USE_SENTRY") != "" {
		useSentry = true
	}

	if os.Getenv("OWL_USE_SLACK") != "" {
		useSlack = true
	}

	if err := initLogger(); err != nil {
		return err
	}
	if err := initMetric(); err != nil {
		return err
	}
	if err := initSlack(); err != nil {
		return err
	}

	return nil
}

func Stop() {
	stopLogger()

	config = Configuration{}
	useMetric = false
	useSentry = false
	useSlack = false
}
