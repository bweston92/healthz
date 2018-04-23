package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bweston92/healthz/check"
)

var (
	endpoint  = flag.String("url", "http://localhost:8000/healthz", "The absolute path to healthz endpoint.")
	outFormat = flag.String("format", "text", "The output format: `text` or `json`.")
	timeoutS  = flag.Int("timeout", 5, "The timeout in seconds.")
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	flag.Parse()

	c := &check.Checker{
		URL:     *endpoint,
		Timeout: time.Duration(*timeoutS) * time.Second,
		Client:  http.DefaultClient,
		Output: &check.Output{
			Format: *outFormat,
			Dest:   os.Stdout,
		},
	}

	if !c.IsAlive() {
		os.Exit(1)
	}
}
