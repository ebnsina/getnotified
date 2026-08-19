package httpapi

import (
	"testing"

	"github.com/fajrlabs/getnotified/internal/store"
)

func ptr[T any](v T) *T { return &v }

func TestValidateMonitor(t *testing.T) {
	valid := store.MonitorInput{
		Name:            ptr("Marketing site"),
		Target:          ptr("https://example.com"),
		IntervalSeconds: ptr(int32(60)),
		TimeoutSeconds:  ptr(int32(10)),
	}

	if err := validateMonitor(valid, true); err != nil {
		t.Fatalf("a well formed monitor was rejected: %v", err)
	}

	cases := []struct {
		name  string
		input store.MonitorInput
		field string
	}{
		{"no name", store.MonitorInput{Target: ptr("https://example.com")}, "invalid_name"},
		{"no target", store.MonitorInput{Name: ptr("x")}, "invalid_target"},
		{"target is not a web address",
			store.MonitorInput{Name: ptr("x"), Target: ptr("example.com")}, "invalid_target"},
		{"unknown check type",
			store.MonitorInput{Name: ptr("x"), Target: ptr("example.com"), Type: ptr("carrier-pigeon")}, "invalid_type"},
		{"checks too often",
			store.MonitorInput{Name: ptr("x"), Target: ptr("https://a.com"), IntervalSeconds: ptr(int32(5))}, "invalid_interval_seconds"},
		{"timeout outruns the interval",
			store.MonitorInput{Name: ptr("x"), Target: ptr("https://a.com"),
				IntervalSeconds: ptr(int32(30)), TimeoutSeconds: ptr(int32(45))}, "invalid_timeout_seconds"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateMonitor(c.input, true)
			if err == nil {
				t.Fatalf("expected %s, got no error", c.field)
			}
			apiErr, ok := err.(*Error)
			if !ok || apiErr.Code != c.field {
				t.Fatalf("got %v, want code %s", err, c.field)
			}
			if apiErr.Message == "" {
				t.Error("a validation failure must say what to do about it")
			}
		})
	}
}

// A bare TCP monitor takes a host, not a URL, so the web address rule must not
// apply to it.
func TestValidateMonitorNonHTTPTarget(t *testing.T) {
	in := store.MonitorInput{Name: ptr("db"), Target: ptr("example.com:5432"), Type: ptr("tcp")}
	if err := validateMonitor(in, true); err != nil {
		t.Fatalf("a host and port was rejected for a tcp monitor: %v", err)
	}
}
