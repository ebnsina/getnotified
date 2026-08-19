package httpapi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fajrlabs/getnotified/internal/store"
)

const (
	defaultCheckLimit = 100
	maxCheckLimit     = 1000
	incidentLimit     = 100
)

func (a *API) listMonitors(w http.ResponseWriter, r *http.Request) error {
	out, err := store.Summaries(r.Context(), a.Pool, a.OrgID)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, out)
}

func (a *API) getMonitor(w http.ResponseWriter, r *http.Request) error {
	m, err := store.GetMonitor(r.Context(), a.Pool, chi.URLParam(r, "id"))
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	return respond(w, http.StatusOK, m)
}

func (a *API) createMonitor(w http.ResponseWriter, r *http.Request) error {
	in, err := decode[store.MonitorInput](w, r)
	if err != nil {
		return err
	}
	if err := validateMonitor(in, true); err != nil {
		return err
	}

	m, err := store.CreateMonitor(r.Context(), a.Pool, a.OrgID, in)
	if err != nil {
		return err
	}
	return respond(w, http.StatusCreated, m)
}

func (a *API) patchMonitor(w http.ResponseWriter, r *http.Request) error {
	in, err := decode[store.MonitorInput](w, r)
	if err != nil {
		return err
	}
	if err := validateMonitor(in, false); err != nil {
		return err
	}

	m, err := store.UpdateMonitor(r.Context(), a.Pool, chi.URLParam(r, "id"), in)
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	return respond(w, http.StatusOK, m)
}

func (a *API) deleteMonitor(w http.ResponseWriter, r *http.Request) error {
	deleted, err := store.DeleteMonitor(r.Context(), a.Pool, chi.URLParam(r, "id"))
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	if !deleted {
		return ErrMonitorNotFound
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (a *API) listChecks(w http.ResponseWriter, r *http.Request) error {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > maxCheckLimit {
		limit = defaultCheckLimit
	}
	out, err := store.ListChecks(r.Context(), a.Pool, chi.URLParam(r, "id"), limit)
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	return respond(w, http.StatusOK, out)
}

func (a *API) listIncidents(w http.ResponseWriter, r *http.Request) error {
	out, err := store.ListIncidents(r.Context(), a.Pool, a.OrgID, chi.URLParam(r, "id"), incidentLimit)
	if err != nil {
		return notFound(err, ErrMonitorNotFound)
	}
	return respond(w, http.StatusOK, out)
}
