package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// invalidTextRepresentation is what Postgres reports for a malformed uuid in a
// path parameter — a bad address, not a server fault.
const invalidTextRepresentation = "22P02"

// Error is the only shape this API ever returns for a failure. Message is
// written for a person to read, so callers can surface it verbatim.
type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

// Every failure this API can report. Keeping them in one list is what makes
// the API the single source of truth for error copy.
var (
	ErrUnauthorized = &Error{http.StatusUnauthorized, "unauthorized",
		"Your session is not valid. Please sign in again."}
	ErrMonitorNotFound = &Error{http.StatusNotFound, "monitor_not_found",
		"That monitor no longer exists."}
	ErrChannelNotFound = &Error{http.StatusNotFound, "channel_not_found",
		"That notification channel no longer exists."}
	ErrStatusPageNotFound = &Error{http.StatusNotFound, "status_page_not_found",
		"There is no status page at this address."}
	ErrNotFound = &Error{http.StatusNotFound, "not_found",
		"We could not find what you asked for."}
	ErrMalformedJSON = &Error{http.StatusBadRequest, "malformed_json",
		"We could not read that request. Please check the format and try again."}
	ErrInternal = &Error{http.StatusInternalServerError, "internal_error",
		"Something went wrong on our side. The team has been notified."}
)

// Invalid builds a validation failure. The message is shown to the person who
// submitted the form, so it names the field in plain words.
func Invalid(field, message string) *Error {
	return &Error{http.StatusBadRequest, "invalid_" + field, message}
}

// write turns any error into the envelope. Database and network failures are
// logged in full and reported as a generic message, never leaked.
func write(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *Error
	switch {
	case errors.As(err, &apiErr):
	case errors.Is(err, pgx.ErrNoRows):
		apiErr = ErrNotFound
	default:
		slog.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
		apiErr = ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)
	json.NewEncoder(w).Encode(map[string]*Error{"error": apiErr})
}

// notFound maps a missing or unparseable row reference to the right message.
func notFound(err error, missing *Error) error {
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) ||
		(errors.As(err, &pgErr) && pgErr.Code == invalidTextRepresentation) {
		return missing
	}
	return err
}
