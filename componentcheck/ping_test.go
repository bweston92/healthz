package componentcheck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNewPingsHealthCheckStub struct {
	err error
}

func (s *testNewPingsHealthCheckStub) Ping() error {
	return s.err
}

func TestNewPingsHealthCheck(t *testing.T) {
	stub := &testNewPingsHealthCheckStub{}
	sut := NewPingsHealthCheck(stub)
	assert.Nil(t, sut())

	stub.err = errors.New("this is a test")
	assert.NotNil(t, sut())
}
