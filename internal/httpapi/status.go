package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/fajrlabs/getnotified/internal/store"
)

// StatusPage is public and unauthenticated by design.
type StatusPage struct {
	Org       store.Org        `json:"org"`
	Monitors  []store.Summary  `json:"monitors"`
	Incidents []store.Incident `json:"incidents"`
	Overall   string           `json:"overall"`
	AsOf      time.Time        `json:"as_of"`
}

func (a *API) publicStatus(w http.ResponseWriter, r *http.Request) error {
	org, err := store.OrgBySlug(r.Context(), a.Pool, chi.URLParam(r, "slug"))
	if err != nil {
		return notFound(err, ErrStatusPageNotFound)
	}

	monitors, err := store.Summaries(r.Context(), a.Pool, org.ID)
	if err != nil {
		return err
	}
	incidents, err := store.ListIncidents(r.Context(), a.Pool, org.ID, "", incidentLimit)
	if err != nil {
		return err
	}

	overall := "operational"
	for _, m := range monitors {
		if !m.Paused && m.Status == "down" {
			overall = "degraded"
			break
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=30")
	return respond(w, http.StatusOK, StatusPage{org, monitors, incidents, overall, time.Now()})
}
