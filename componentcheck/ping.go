package componentcheck

import (
	"context"
	"time"

	"github.com/bweston92/healthz/healthz"
)

type Pings interface {
	Ping() error
}

type PingsWithTimeout interface {
	Ping(ctx context.Context) error
}

// NewPingsHealthCheck ...
func NewPingsHealthCheck(c Pings) healthz.ComponentHealthCheck {
	return func() *healthz.Error {
		return healthz.WrapError(c.Ping(), "unable to ping service", healthz.Meta{})
	}
}

// NewPingsWithTimeoutHealthCheck ...
func NewPingsWithTimeoutHealthCheck(c PingsWithTimeout, timeout time.Duration) healthz.ComponentHealthCheck {
	return func() *healthz.Error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return healthz.WrapError(c.Ping(ctx), "unable to ping service", healthz.Meta{
			"timeout": timeout.String(),
		})
	}
}
