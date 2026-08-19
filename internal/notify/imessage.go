package notify

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fajrlabs/getnotified/internal/store"
)

// imessage posts to the optional macOS relay; without it configured, the
// channel simply fails and every other channel carries on.
type imessage struct{}

func (imessage) Send(ctx context.Context, e Event, c store.Channel) error {
	relay := os.Getenv("IMESSAGE_RELAY_URL")
	if relay == "" {
		return fmt.Errorf("imessage channel needs IMESSAGE_RELAY_URL")
	}
	return postJSON(ctx, strings.TrimSuffix(relay, "/")+"/send",
		map[string]string{"to": cfg(c, "to"), "text": e.Body()},
		map[string]string{"Authorization": "Bearer " + os.Getenv("IMESSAGE_RELAY_KEY")})
}
