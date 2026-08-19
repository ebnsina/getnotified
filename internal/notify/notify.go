// Package notify turns an incident into a message on whatever channels the
// monitor is attached to.
package notify

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fajrlabs/getnotified/internal/store"
)

// Event is what a channel announces. Recoveries need the same context as
// outages, so this carries the event rather than a bare incident.
type Event struct {
	Kind     string         `json:"kind"` // "down" | "up" | "test"
	Monitor  store.Monitor  `json:"monitor"`
	Incident store.Incident `json:"incident"`
}

// TestEvent is what a channel sends when someone checks it works. It says so
// plainly, so nobody mistakes it for a real outage.
func TestEvent(now time.Time) Event {
	return Event{
		Kind:     "test",
		Monitor:  store.Monitor{Name: "Test message", Target: "no monitor — this is a check"},
		Incident: store.Incident{StartedAt: now},
	}
}

type Notifier interface {
	Send(ctx context.Context, e Event, c store.Channel) error
}

// required names the config key each channel cannot work without, and how to
// describe it to whoever is filling in the form.
var required = map[string]struct{ key, label string }{
	"slack":    {"webhook_url", "Slack address to post to"},
	"webhook":  {"url", "web address to send to"},
	"email":    {"to", "email address"},
	"sms":      {"to", "phone number"},
	"whatsapp": {"to", "phone number"},
	"imessage": {"to", "phone number or Apple ID"},
}

var notifiers = map[string]Notifier{
	"slack":    slack{},
	"webhook":  webhook{},
	"email":    email{},
	"sms":      twilio{},
	"whatsapp": twilio{},
	"imessage": imessage{},
}

func For(kind string) (Notifier, error) {
	n, ok := notifiers[kind]
	if !ok {
		return nil, fmt.Errorf("unknown channel type %q", kind)
	}
	return n, nil
}

// ValidateConfig reports a missing destination in words the sender can act on.
func ValidateConfig(kind string, config map[string]any) error {
	want, ok := required[kind]
	if !ok {
		return errors.New("That kind of channel is not supported.")
	}
	if s, _ := config[want.key].(string); strings.TrimSpace(s) == "" {
		return fmt.Errorf("Add a %s before this can send anything.", want.label)
	}
	return nil
}

// Subject and Body are deliberately plain. No countdowns, no urgency theatre.
func (e Event) Subject() string {
	switch e.Kind {
	case "test":
		return "Test message from GetNotified"
	case "up":
		return fmt.Sprintf("%s is back up", e.Monitor.Name)
	default:
		return fmt.Sprintf("%s is down", e.Monitor.Name)
	}
}

func (e Event) Body() string {
	if e.Kind == "test" {
		return "This is a test from GetNotified.\n\n" +
			"If you are reading it, this channel is set up correctly and real " +
			"notices will arrive here.\n"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\nTarget: %s\n", e.Subject(), e.Monitor.Target)
	if e.Incident.Cause != nil && *e.Incident.Cause != "" {
		fmt.Fprintf(&b, "Cause: %s\n", *e.Incident.Cause)
	}
	if e.Kind == "up" && e.Incident.ResolvedAt != nil {
		fmt.Fprintf(&b, "Downtime: %s\n",
			e.Incident.ResolvedAt.Sub(e.Incident.StartedAt).Round(time.Second))
	}
	fmt.Fprintf(&b, "Started: %s\n", e.Incident.StartedAt.Format(time.RFC1123))
	return b.String()
}

func cfg(c store.Channel, key string) string {
	s, _ := c.Config[key].(string)
	return strings.TrimSpace(s)
}
