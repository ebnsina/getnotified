package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// bodyPeek is how much of a rejection we quote back. Enough to explain, not
// enough to paste a web page into a notification.
const bodyPeek = 400

var client = &http.Client{Timeout: 15 * time.Second}

// postJSON treats any non-2xx as an error so River retries the delivery.
func postJSON(ctx context.Context, url string, body any, hdr map[string]string) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range hdr {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		peek, _ := io.ReadAll(io.LimitReader(resp.Body, bodyPeek))
		if said := strings.TrimSpace(string(peek)); said != "" {
			return fmt.Errorf("it answered %d — %s", resp.StatusCode, said)
		}
		return fmt.Errorf("it answered %d", resp.StatusCode)
	}
	return nil
}
