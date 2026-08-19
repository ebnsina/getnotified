package httpapi

import (
	"net/url"
	"slices"

	"github.com/fajrlabs/getnotified/internal/store"
)

var monitorTypes = []string{"http", "tcp", "ssl_expiry"}

// validateMonitor is the authoritative rule set; the dashboard mirrors it only
// to save a round trip.
func validateMonitor(in store.MonitorInput, creating bool) error {
	if creating && (in.Name == nil || *in.Name == "") {
		return Invalid("name", "Give the monitor a name so you can recognise it later.")
	}
	if creating && (in.Target == nil || *in.Target == "") {
		return Invalid("target", "Tell us what to check — a web address or a server name.")
	}
	if in.Type != nil && !slices.Contains(monitorTypes, *in.Type) {
		return Invalid("type", "Choose a web address, a port on a server, or a security certificate.")
	}

	kind := "http"
	if in.Type != nil {
		kind = *in.Type
	}
	if in.Target != nil && kind == "http" {
		u, err := url.Parse(*in.Target)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return Invalid("target", "That does not look like a web address. Try something like https://example.com.")
		}
	}
	if in.IntervalSeconds != nil && *in.IntervalSeconds < 10 {
		return Invalid("interval_seconds", "Check no more often than once every 10 seconds.")
	}
	if in.TimeoutSeconds != nil && (*in.TimeoutSeconds < 1 || *in.TimeoutSeconds > 120) {
		return Invalid("timeout_seconds", "Set a timeout between 1 and 120 seconds.")
	}
	// A check that can outlive its own schedule would overlap the next one.
	if in.TimeoutSeconds != nil && in.IntervalSeconds != nil &&
		*in.TimeoutSeconds > *in.IntervalSeconds {
		return Invalid("timeout_seconds",
			"Give up sooner than the gap between checks, so one check cannot run into the next.")
	}
	if in.FailureThreshold != nil && *in.FailureThreshold < 1 {
		return Invalid("failure_threshold", "It takes at least one failure to open an incident.")
	}
	if in.SSLWarnDays != nil && *in.SSLWarnDays < 1 {
		return Invalid("ssl_warn_days", "Warn at least one day before the certificate runs out.")
	}
	return nil
}
