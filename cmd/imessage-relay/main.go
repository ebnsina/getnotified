// Command imessage-relay sends iMessages on behalf of GetNotified. Apple has
// no public API, so this runs on a Mac and drives the Messages app.
package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const sendTimeout = 20 * time.Second

// script takes the recipient and the message as arguments rather than building
// AppleScript from them, so nothing in a message can be executed.
const script = `on run argv
	set theRecipient to item 1 of argv
	set theText to item 2 of argv
	tell application "Messages"
		set targetService to 1st account whose service type = iMessage
		set targetBuddy to participant theRecipient of targetService
		send theText to targetBuddy
	end tell
end run`

type config struct {
	addr string
	key  string
}

func loadConfig() (config, error) {
	var missing []string
	get := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	c := config{addr: get("RELAY_ADDR"), key: get("IMESSAGE_RELAY_KEY")}
	if len(missing) > 0 {
		return c, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return c, nil
}

func main() {
	if runtime.GOOS != "darwin" {
		slog.Error("the iMessage relay only runs on macOS", "os", runtime.GOOS)
		os.Exit(1)
	}
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /send", handleSend(cfg.key))

	srv := &http.Server{Addr: cfg.addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	slog.Info("iMessage relay listening", "addr", cfg.addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func handleSend(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authorised(r, key) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var body struct {
			To   string `json:"to"`
			Text string `json:"text"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if body.To == "" || body.Text == "" {
			http.Error(w, "both to and text are required", http.StatusBadRequest)
			return
		}

		if err := send(r.Context(), body.To, body.Text); err != nil {
			slog.Error("send failed", "to", body.To, "err", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func authorised(r *http.Request, key string) bool {
	got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	return subtle.ConstantTimeCompare([]byte(got), []byte(key)) == 1
}

func send(ctx context.Context, to, text string) error {
	ctx, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "osascript", "-", to, text)
	cmd.Stdin = strings.NewReader(script)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Messages refused it: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
