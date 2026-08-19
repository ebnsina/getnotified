package notify

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"testing"
)

func TestExplain(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"nothing listening", &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			"nothing answered at that address."},
		{"unknown host", &net.DNSError{Err: "no such host", Name: "nope.example"},
			"that address could not be found."},
		{"gave up waiting", context.DeadlineExceeded, "it took too long to answer."},
		{"rejected by the far end", fmt.Errorf("https://hooks.slack.com/services/x returned 404"),
			"the address answered with 404."},
		{"our own wording survives", errors.New("Add an email address before this can send anything."),
			"Add an email address before this can send anything."},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Explain(c.err); got != c.want {
				t.Fatalf("Explain(%v) = %q, want %q", c.err, got, c.want)
			}
		})
	}

	// Whatever happens, nothing socket-shaped reaches the reader.
	for _, err := range []error{
		&net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connect: no route to host")},
		fmt.Errorf(`Post "http://x/y": dial tcp 1.2.3.4:80: i/o timeout`),
	} {
		if got := Explain(err); strings.Contains(got, "dial ") || strings.Contains(got, "tcp") {
			t.Errorf("Explain(%v) leaked socket wording: %q", err, got)
		}
	}
}
