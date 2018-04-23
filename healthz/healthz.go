package healthz

import (
	"fmt"
	"sync"
)

type (
	// Meta is used to provide mixed additional data.
	Meta map[string]string

	// ComponentHealthCheck is a callback that will return
	// an error if the dependency is not available.
	ComponentHealthCheck func() *Error

	// Component represents a crucial part an application
	// that effects the services state.
	Component struct {
		Name string

		// If there is an error returned from the component
		// we will use that to report the service as unhealthy.
		Check ComponentHealthCheck

		// For example when using a MySQL connection you may
		// want to add the hostname and user you connected as.
		Metadata Meta `json:"metadata"`

		// Required is the difference of the response status being
		// DOWN or DEGRADED.
		Required bool
	}

	// Error represents an issue with a component.
	Error struct {
		Description string `json:"description"`
		Metadata    Meta   `json:"metadata"`
	}

	Healthz struct {
		sync.Mutex

		metadata     Meta
		started      bool
		shuttingDown bool
		components   []*Component
	}
)

func New(m Meta, c ...*Component) *Healthz {
	return &Healthz{
		metadata:   m,
		components: c,
	}
}

func (h *Healthz) Started() {
	h.started = true
}

func (h *Healthz) Close() {
	h.shuttingDown = true
}

func (h *Healthz) Add(c *Component) {
	h.Lock()
	defer h.Unlock()

	h.components = append(h.components, c)
}

func NewComponent(n string, r bool, m Meta, c ComponentHealthCheck) *Component {
	return &Component{
		Name:     n,
		Check:    c,
		Metadata: m,
		Required: r,
	}
}

// NewError will take the error message and metadata
// and return a Error which can be used in the healthz
// response.
func NewError(msg string, meta Meta) *Error {
	return &Error{
		Description: msg,
		Metadata:    meta,
	}
}

// WrapError wraps an error object with a related message.
// However if err is nil it will return nil to make it easier
// to have something like:
//
// ```
// return healthz.WrapError(
//	pingDatabase(),
//  "Unable to connect to database.",
//  healthz.Meta{},
// )
// ````
func WrapError(err error, msg string, meta Meta) *Error {
	if err == nil {
		return nil
	}

	return NewError(fmt.Sprintf("%s (err: %v)", msg, err), meta)
}
