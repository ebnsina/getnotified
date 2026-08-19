package jobs

import "github.com/fajrlabs/getnotified/internal/store"

// NextState is the whole up/down decision, kept pure so the flapping rule is
// testable without a database. transition is "", "down" or "up".
func NextState(m store.Monitor, ok bool) (status string, failures int32, transition string) {
	if ok {
		status = "up"
	} else {
		failures = m.ConsecutiveFailures + 1
		status = m.Status
		if failures >= m.FailureThreshold {
			status = "down"
		}
	}

	switch {
	case status == "down" && m.Status != "down":
		transition = "down"
	case status == "up" && m.Status == "down":
		transition = "up"
	}
	return status, failures, transition
}
