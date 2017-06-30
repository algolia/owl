package owl

import "os"

var (
	config    Configuration
	useMetric bool
	useSentry bool
	useSlack  bool
)

func Init(c Configuration) (err error) {
	if err = checkConfiguration(c); err != nil {
		return
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

	initLogger()
	initMetric()
	initSlack()

	return
}

func Stop() {
	stopLogger()
}
