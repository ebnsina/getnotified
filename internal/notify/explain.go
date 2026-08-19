package notify

import (
	"context"
	"errors"
	"net"
	"strings"
	"syscall"
)

// Explain turns a delivery failure into a sentence the sender can act on.
// Transport errors carry Go and socket wording that means nothing to them.
func Explain(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, context.DeadlineExceeded):
		return "it took too long to answer."
	case errors.Is(err, syscall.ECONNREFUSED):
		return "nothing answered at that address."
	}

	var dns *net.DNSError
	if errors.As(err, &dns) {
		return "that address could not be found."
	}
	var timeout net.Error
	if errors.As(err, &timeout) && timeout.Timeout() {
		return "it took too long to answer."
	}

	// Errors we wrote ourselves already read as plain sentences.
	if msg := err.Error(); !strings.Contains(msg, "dial ") && !strings.Contains(msg, "http://") &&
		!strings.Contains(msg, "https://") {
		return msg
	}
	if _, after, found := strings.Cut(err.Error(), " returned "); found {
		return "the address answered with " + after + "."
	}
	return "it could not be delivered."
}
