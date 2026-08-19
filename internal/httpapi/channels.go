package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/fajrlabs/getnotified/internal/notify"
	"github.com/fajrlabs/getnotified/internal/store"
)

// A test must answer while someone is waiting on it.
const testTimeout = 15 * time.Second

func (a *API) listChannels(w http.ResponseWriter, r *http.Request) error {
	out, err := store.ListChannels(r.Context(), a.Pool, a.OrgID)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, out)
}

func (a *API) createChannel(w http.ResponseWriter, r *http.Request) error {
	in, err := decode[store.Channel](w, r)
	if err != nil {
		return err
	}
	if in.Name == "" {
		return Invalid("name", "Give the channel a name, such as “Ops Slack”.")
	}
	if _, err := notify.For(in.Type); err != nil {
		return Invalid("type", "Choose one of Slack, email, webhook, SMS, WhatsApp, or iMessage.")
	}
	if err := notify.ValidateConfig(in.Type, in.Config); err != nil {
		return Invalid("config", err.Error())
	}

	ch, err := store.CreateChannel(r.Context(), a.Pool, a.OrgID, in.Name, in.Type, in.Config)
	if err != nil {
		return err
	}
	return respond(w, http.StatusCreated, ch)
}

func (a *API) deleteChannel(w http.ResponseWriter, r *http.Request) error {
	deleted, err := store.DeleteChannel(r.Context(), a.Pool, chi.URLParam(r, "id"))
	if err != nil {
		return notFound(err, ErrChannelNotFound)
	}
	if !deleted {
		return ErrChannelNotFound
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// testChannel delivers a test message straight away, rather than through the
// queue, so whoever pressed the button learns the result now.
func (a *API) testChannel(w http.ResponseWriter, r *http.Request) error {
	ch, err := store.GetChannel(r.Context(), a.Pool, chi.URLParam(r, "id"))
	if err != nil {
		return notFound(err, ErrChannelNotFound)
	}

	sender, err := notify.For(ch.Type)
	if err != nil {
		return Invalid("type", "That kind of channel is not supported.")
	}
	if err := notify.ValidateConfig(ch.Type, ch.Config); err != nil {
		return Invalid("config", err.Error())
	}

	ctx, cancel := context.WithTimeout(r.Context(), testTimeout)
	defer cancel()

	if err := sender.Send(ctx, notify.TestEvent(time.Now()), ch); err != nil {
		slog.Info("channel test failed", "channel", ch.ID, "type", ch.Type, "err", err)
		return TestFailed(notify.Explain(err))
	}
	return respond(w, http.StatusOK, map[string]bool{"delivered": true})
}

func (a *API) getMonitorChannels(w http.ResponseWriter, r *http.Request) error {
	ids, err := store.MonitorChannelIDs(r.Context(), a.Pool, chi.URLParam(r, "id"))
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	if ids == nil {
		ids = []string{}
	}
	return respond(w, http.StatusOK, ids)
}

func (a *API) setMonitorChannels(w http.ResponseWriter, r *http.Request) error {
	body, err := decode[struct {
		ChannelIDs []string `json:"channel_ids"`
	}](w, r)
	if err != nil {
		return err
	}
	if err := store.SetMonitorChannels(r.Context(), a.Pool, chi.URLParam(r, "id"), body.ChannelIDs); err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
