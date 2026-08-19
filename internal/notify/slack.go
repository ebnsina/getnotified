package notify

import (
	"context"
	"fmt"

	"github.com/fajrlabs/getnotified/internal/store"
)

type slack struct{}

func (slack) Send(ctx context.Context, e Event, c store.Channel) error {
	emoji := ":red_circle:"
	if e.Kind == "up" {
		emoji = ":large_green_circle:"
	}
	return postJSON(ctx, cfg(c, "webhook_url"), map[string]any{
		"text": fmt.Sprintf("%s *%s*\n```%s```", emoji, e.Subject(), e.Body()),
	}, nil)
}
