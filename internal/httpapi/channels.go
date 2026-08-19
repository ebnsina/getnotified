package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fajrlabs/getnotified/internal/notify"
	"github.com/fajrlabs/getnotified/internal/store"
)

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
