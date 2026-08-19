// Package httpapi is the product surface: everything the dashboard can do is
// done through these routes.
package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxBodyBytes = 1 << 20

type API struct {
	Pool   *pgxpool.Pool
	OrgID  string
	APIKey string
}

// handler lets every route return an error and leaves the response shape to one place.
type handler func(http.ResponseWriter, *http.Request) error

func (a *API) h(fn handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			write(w, r, err)
		}
	}
}

func (a *API) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.NotFound(a.h(func(http.ResponseWriter, *http.Request) error { return ErrNotFound }))
	r.Get("/healthz", a.h(a.health))
	r.Get("/status/{slug}", a.h(a.publicStatus))

	r.Route("/api", func(r chi.Router) {
		r.Use(a.auth)
		r.Get("/monitors", a.h(a.listMonitors))
		r.Post("/monitors", a.h(a.createMonitor))
		r.Get("/monitors/{id}", a.h(a.getMonitor))
		r.Patch("/monitors/{id}", a.h(a.patchMonitor))
		r.Delete("/monitors/{id}", a.h(a.deleteMonitor))
		r.Get("/monitors/{id}/checks", a.h(a.listChecks))
		r.Get("/monitors/{id}/incidents", a.h(a.listIncidents))
		r.Get("/monitors/{id}/channels", a.h(a.getMonitorChannels))
		r.Put("/monitors/{id}/channels", a.h(a.setMonitorChannels))

		r.Get("/channels", a.h(a.listChannels))
		r.Post("/channels", a.h(a.createChannel))
		r.Delete("/channels/{id}", a.h(a.deleteChannel))
		r.Post("/channels/{id}/test", a.h(a.testChannel))
	})
	return r
}

func (a *API) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+a.APIKey {
			write(w, r, ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) health(w http.ResponseWriter, r *http.Request) error {
	if err := a.Pool.Ping(r.Context()); err != nil {
		return err
	}
	return respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func respond(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func decode[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&v); err != nil {
		return v, ErrMalformedJSON
	}
	return v, nil
}
