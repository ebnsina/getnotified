package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fajrlabs/getnotified/internal/store"
)

func TestRunHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	m := store.Monitor{Type: "http", Target: srv.URL, TimeoutSeconds: 5, ExpectedStatus: []int32{200}}
	if r := Run(context.Background(), m); r.OK || r.StatusCode == nil || *r.StatusCode != 503 {
		t.Fatalf("503 against expected {200}: got %+v, want not-ok with code 503", r)
	}

	m.ExpectedStatus = []int32{200, 503}
	if r := Run(context.Background(), m); !r.OK {
		t.Fatalf("503 against expected {200,503}: got %+v, want ok", r)
	}
}

func TestHostPort(t *testing.T) {
	for target, want := range map[string]string{
		"https://example.com/health": "example.com:443",
		"http://example.com":         "example.com:80",
		"example.com:5432":           "example.com:5432",
		"example.com":                "example.com:80",
	} {
		if got := hostPort(target, "80"); got != want {
			t.Errorf("hostPort(%q) = %q, want %q", target, got, want)
		}
	}
}
