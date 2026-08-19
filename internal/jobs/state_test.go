package jobs

import (
	"testing"

	"github.com/fajrlabs/getnotified/internal/store"
)

func TestNextState_FlappingProtection(t *testing.T) {
	m := store.Monitor{Status: "up", FailureThreshold: 2}

	// One failure is not an incident.
	status, failures, transition := NextState(m, false)
	if status != "up" || failures != 1 || transition != "" {
		t.Fatalf("first failure: got %q %d %q, want up 1 <none>", status, failures, transition)
	}

	// The second consecutive one is.
	m.ConsecutiveFailures = failures
	status, failures, transition = NextState(m, false)
	if status != "down" || failures != 2 || transition != "down" {
		t.Fatalf("second failure: got %q %d %q, want down 2 down", status, failures, transition)
	}

	// Staying down does not re-fire.
	m.Status, m.ConsecutiveFailures = status, failures
	if _, _, transition = NextState(m, false); transition != "" {
		t.Fatalf("still down: got %q, want no transition", transition)
	}

	// Recovery fires once, then goes quiet.
	status, failures, transition = NextState(m, true)
	if status != "up" || failures != 0 || transition != "up" {
		t.Fatalf("recovery: got %q %d %q, want up 0 up", status, failures, transition)
	}
	m.Status, m.ConsecutiveFailures = status, failures
	if _, _, transition = NextState(m, true); transition != "" {
		t.Fatalf("still up: got %q, want no transition", transition)
	}

	// A single blip after recovery is absorbed.
	if _, _, transition = NextState(m, false); transition != "" {
		t.Fatalf("blip: got %q, want no transition", transition)
	}
}
