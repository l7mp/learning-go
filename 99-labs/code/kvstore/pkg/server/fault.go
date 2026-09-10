package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// FaultEnvVar names the environment variable that enables fault injection. Fault injection
// is a testing facility: unless this variable is set to a non-empty value the key-value
// store behaves exactly as if none of the code in this file existed.
const FaultEnvVar = "KVSTORE_FAULT_INJECTION"

// FaultConfig describes the faults injected into a single API path.
type FaultConfig struct {
	// Path is the API path the faults apply to, e.g. "/api/put".
	Path string `json:"path"`
	// FailNext makes the next FailNext requests fail, then stops. Each failed request
	// decrements the counter, so a value of 3 fails exactly three requests.
	FailNext int `json:"failNext"`
	// AbortPercent fails the given percentage of requests, deterministically: every
	// n-th request fails, where n is 100/AbortPercent.
	AbortPercent int `json:"abortPercent"`
	// Status is the HTTP status returned for a failed request, 500 if unset.
	Status int `json:"status"`
	// Delay is slept before the request is served, in Go duration syntax, e.g. "2s".
	Delay string `json:"delay"`
}

// FaultStatus is the response of a GET on the fault API.
type FaultStatus struct {
	// Faults holds the fault configuration per path.
	Faults map[string]FaultConfig `json:"faults"`
	// Requests counts the requests served per path since startup, including the failed
	// ones. Tests use this to assert that a client actually retried.
	Requests map[string]int `json:"requests"`
	// Failures counts the requests failed per path since startup.
	Failures map[string]int `json:"failures"`
}

// faultInjector injects configurable faults into the API. It is inert unless enabled.
type faultInjector struct {
	mu       sync.Mutex
	enabled  bool
	faults   map[string]FaultConfig
	requests map[string]int
	failures map[string]int
}

func newFaultInjector() *faultInjector {
	return &faultInjector{
		enabled:  os.Getenv(FaultEnvVar) != "",
		faults:   make(map[string]FaultConfig),
		requests: make(map[string]int),
		failures: make(map[string]int),
	}
}

// wrap decorates an HTTP handler with fault injection for the given path. When fault
// injection is disabled the handler is returned unchanged.
func (f *faultInjector) wrap(path string, h http.HandlerFunc) http.HandlerFunc {
	if !f.enabled {
		return h
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status, delay := f.next(path)

		if delay > 0 {
			time.Sleep(delay)
		}

		if status != 0 {
			log.Printf("fault injection: failing %s with status %d", path, status)
			w.WriteHeader(status)
			return
		}

		h(w, r)
	}
}

// next accounts for one request on the path and reports the fault to apply: the HTTP
// status to fail with (zero to serve the request normally) and the delay to sleep first.
func (f *faultInjector) next(path string) (int, time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests[path]++

	c, ok := f.faults[path]
	if !ok {
		return 0, 0
	}

	delay, _ := time.ParseDuration(c.Delay)

	status := c.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	switch {
	case c.FailNext > 0:
		c.FailNext--
		f.faults[path] = c
		f.failures[path]++
		return status, delay
	case c.AbortPercent > 0 && f.requests[path]%(100/c.AbortPercent) == 0:
		f.failures[path]++
		return status, delay
	}

	return 0, delay
}

// registerAPI installs the fault control API. It is not installed at all unless fault
// injection is enabled, so a production key-value store does not serve it.
func (f *faultInjector) registerAPI(mux *http.ServeMux) {
	if !f.enabled {
		return
	}

	log.Printf("fault injection ENABLED (%s is set), serving the control API on /api/fault",
		FaultEnvVar)

	mux.HandleFunc("/api/fault", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			c := FaultConfig{}
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if c.Path == "" {
				http.Error(w, "no path in fault config", http.StatusBadRequest)
				return
			}
			if _, err := time.ParseDuration(c.Delay); c.Delay != "" && err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			f.mu.Lock()
			f.faults[c.Path] = c
			f.mu.Unlock()

			log.Printf("fault injection: %#v", c)
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(f.status())
		case http.MethodDelete:
			f.mu.Lock()
			f.faults = make(map[string]FaultConfig)
			f.requests = make(map[string]int)
			f.failures = make(map[string]int)
			f.mu.Unlock()

			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func (f *faultInjector) status() FaultStatus {
	f.mu.Lock()
	defer f.mu.Unlock()

	s := FaultStatus{
		Faults:   make(map[string]FaultConfig, len(f.faults)),
		Requests: make(map[string]int, len(f.requests)),
		Failures: make(map[string]int, len(f.failures)),
	}
	for k, v := range f.faults {
		s.Faults[k] = v
	}
	for k, v := range f.requests {
		s.Requests[k] = v
	}
	for k, v := range f.failures {
		s.Failures[k] = v
	}

	return s
}
