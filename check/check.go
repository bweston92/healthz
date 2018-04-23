package check

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bweston92/healthz/healthz"
)

type (
	requestHandler interface {
		Do(*http.Request) (*http.Response, error)
	}

	componentStatus struct {
		Healthy  bool             `json:"healthy"`
		Metadata healthz.Meta     `json:"metadata"`
		Errors   []*healthz.Error `json:"errors"`
	}

	response struct {
		Status     string                      `json:"status"`
		Components map[string]*componentStatus `json:"components"`
		Meta       healthz.Meta                `json:"metadata"`
	}
)

type Checker struct {
	URL     string
	Timeout time.Duration
	Client  requestHandler
	Output  *Output
}

// IsAlive will return true or false depending on the state
// of the service the healthz URL is for.
func (c *Checker) IsAlive() bool {
	req, err := http.NewRequest("GET", c.URL, nil)

	if err != nil {
		logrus.WithError(err).Error("cannot continue, bad config")
		return false
	}

	res, err := c.Client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("unable to get healthz response")
		return false
	}
	defer res.Body.Close()

	if isStartingOrShuttingDown(res.StatusCode) {
		logrus.Warning("service is starting up or shutting down")
		return false
	}

	return isResponseValid(res, c.Output)
}

func isStartingOrShuttingDown(c int) bool {
	return c == http.StatusServiceUnavailable
}

func isResponseValid(res *http.Response, out *Output) bool {
	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		logrus.WithError(err).Error("unable to read healthz response fully")
		return false
	}

	r := &response{}
	if err := json.Unmarshal(b, r); err != nil {
		logrus.WithError(err).Error("server responded with invalid json")
		return false
	}

	if out != nil {
		out.write(r)
	}

	return res.StatusCode == http.StatusOK
}
