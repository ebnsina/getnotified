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
	Kind     string         `json:"kind"` // "down" | "up"
	Monitor  store.Monitor  `json:"monitor"`
	Incident store.Incident `json:"incident"`
}

type Notifier interface {
	Send(ctx context.Context, e Event, c store.Channel) error
}

// required names the config key each channel type cannot work without.
var required = map[string]string{
	"slack":    "webhook_url",
	"webhook":  "url",
	"email":    "to",
	"sms":      "to",
	"whatsapp": "to",
	"imessage": "to",
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
	key, ok := required[kind]
	if !ok {
		return errors.New("That channel type is not supported.")
	}
	if s, _ := config[key].(string); strings.TrimSpace(s) == "" {
		return fmt.Errorf("This channel needs a %s before it can send anything.",
			strings.ReplaceAll(key, "_", " "))
	}
	return nil
}

// Subject and Body are deliberately plain. No countdowns, no urgency theatre.
func (e Event) Subject() string {
	if e.Kind == "up" {
		return fmt.Sprintf("%s is back up", e.Monitor.Name)
	}
	return fmt.Sprintf("%s is down", e.Monitor.Name)
}

func (e Event) Body() string {
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
