package notify

import (
	"context"

	"github.com/fajrlabs/getnotified/internal/store"
)

type webhook struct{}

func (webhook) Send(ctx context.Context, e Event, c store.Channel) error {
	hdr := map[string]string{}
	if s := cfg(c, "secret"); s != "" {
		hdr["X-GetNotified-Secret"] = s
	}
	return postJSON(ctx, cfg(c, "url"), e, hdr)
}
