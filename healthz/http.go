package healthz

import (
	"encoding/json"
	"net/http"
)

type (
	componentStatus struct {
		Healthy  bool     `json:"healthy"`
		Metadata Meta     `json:"metadata"`
		Errors   []*Error `json:"errors"`
	}

	response struct {
		Status     string                      `json:"status"`
		Components map[string]*componentStatus `json:"components"`
		Meta       Meta                        `json:"metadata"`
	}
)

func (h *Healthz) getResponse() *response {
	res := &response{
		Status:     StatusHealthy,
		Components: map[string]*componentStatus{},
		Meta:       h.metadata,
	}

	for _, c := range h.components {
		errs := []*Error{}

		if err := c.Check(); err != nil {
			errs = append(errs, err)

			if c.Required && res.Status != StatusUnresponsive {
				res.Status = StatusUnresponsive
			} else {
				res.Status = StatusDegraded
			}
		}

		res.Components[c.Name] = &componentStatus{
			Healthy:  len(errs) == 0,
			Metadata: c.Metadata,
			Errors:   errs,
		}
	}

	return res
}

func (h *Healthz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if !h.started {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "starting"}`))
		return
	}

	if h.shuttingDown {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "shutting_down"}`))
		return
	}

	res := h.getResponse()
	b, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "unhealthy"}`))
		return
	}

	if res.Status == StatusHealthy || res.Status == StatusDegraded {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}
