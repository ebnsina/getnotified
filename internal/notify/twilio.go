package notify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fajrlabs/getnotified/internal/store"
)

// twilio serves both SMS and WhatsApp — the same API call, differing only in
// whether the numbers carry a "whatsapp:" prefix.
type twilio struct{}

func (twilio) Send(ctx context.Context, e Event, c store.Channel) error {
	sid, token := os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN")
	from, to := os.Getenv("TWILIO_FROM"), cfg(c, "to")
	if sid == "" || token == "" || from == "" {
		return fmt.Errorf("%s channel needs TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN and TWILIO_FROM", c.Type)
	}
	if c.Type == "whatsapp" {
		from = "whatsapp:" + strings.TrimPrefix(from, "whatsapp:")
		to = "whatsapp:" + strings.TrimPrefix(to, "whatsapp:")
	}

	form := url.Values{"From": {from}, "To": {to}, "Body": {e.Body()}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.twilio.com/2010-04-01/Accounts/"+sid+"/Messages.json",
		strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(sid, token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("twilio returned %d", resp.StatusCode)
	}
	return nil
}
