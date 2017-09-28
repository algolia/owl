package owl

import (
	"os"
	"testing"
)

func TestMetricWavefront_extractWavefrontUrl(t *testing.T) {
	os.Setenv("OWL_TEST_WAVEFRONT_URL", "127.0.0.1")
	os.Setenv("OWL_TEST_WAVEFRONT_PORT", "43")

	for _, c := range []struct {
		InputUrl  string
		OutputUrl string
	}{
		// Valid ones
		{"1.2.3.4:42", "1.2.3.4:42"},
		{"$OWL_TEST_WAVEFRONT_URL:$OWL_TEST_WAVEFRONT_PORT", "127.0.0.1:43"},
		// Invalid ones
		{"", ""},
		{"1.2.3.4", ""},
		{":", ""},
		{"1.2.3.4:42:", ""},
		{"$OWL_TEST_WAVEFRONT_URL:42", ""},
		{"1.2.3.4:$OWL_TEST_WAVEFRONT_PORT", ""},
	} {
		url := extractWavefrontUrl(c.InputUrl)
		if url != c.OutputUrl {
			t.Errorf("Invalid Wavefront URL extraction: expected \"%s\" for \"%s\" but got \"%s\" instead\n", c.OutputUrl, c.InputUrl, url)
		}
	}
}
