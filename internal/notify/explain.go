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

	// Anything else was written to be read: the far end's own explanation, or
	// one of ours. Only socket wording is unfit to show.
	msg := err.Error()
	if strings.Contains(msg, "dial ") || strings.Contains(msg, "connect: ") {
		return "it could not be delivered."
	}
	return msg
}
