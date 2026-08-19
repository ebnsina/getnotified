// Package probe answers one question: is this target healthy right now.
package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/fajrlabs/getnotified/internal/store"
)

const defaultTimeout = 10 * time.Second

type Result struct {
	OK         bool
	LatencyMS  int32
	StatusCode *int32
	Error      string
}

// Run never returns an error: a failed probe is a Result, which is data.
func Run(ctx context.Context, m store.Monitor) Result {
	timeout := time.Duration(m.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	var code *int32
	var err error

	switch m.Type {
	case "tcp":
		err = probeTCP(ctx, m.Target)
	case "ssl_expiry":
		err = probeSSL(ctx, m.Target, int(m.SSLWarnDays))
	default:
		code, err = probeHTTP(ctx, m)
	}

	r := Result{OK: err == nil, LatencyMS: int32(time.Since(start).Milliseconds()), StatusCode: code}
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

func probeHTTP(ctx context.Context, m store.Monitor) (*int32, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.Target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GetNotified/1.0 (+https://getnotified.sh)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	code := int32(resp.StatusCode)
	expected := m.ExpectedStatus
	if len(expected) == 0 {
		expected = []int32{http.StatusOK}
	}
	if !slices.Contains(expected, code) {
		return &code, fmt.Errorf("unexpected status %d", code)
	}
	return &code, nil
}

func probeTCP(ctx context.Context, target string) error {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", hostPort(target, "80"))
	if err != nil {
		return err
	}
	return conn.Close()
}

func probeSSL(ctx context.Context, target string, warnDays int) error {
	conn, err := (&tls.Dialer{}).DialContext(ctx, "tcp", hostPort(target, "443"))
	if err != nil {
		return err
	}
	defer conn.Close()

	cert := conn.(*tls.Conn).ConnectionState().PeerCertificates[0]
	if left := time.Until(cert.NotAfter); left < time.Duration(warnDays)*24*time.Hour {
		return fmt.Errorf("certificate expires in %d days (%s)",
			int(left.Hours()/24), cert.NotAfter.Format(time.DateOnly))
	}
	return nil
}

// hostPort accepts a bare host, host:port, or a full URL.
func hostPort(target, defaultPort string) string {
	if u, err := url.Parse(target); err == nil && u.Host != "" {
		target = u.Host
		if u.Scheme == "https" {
			defaultPort = "443"
		}
	}
	if _, _, err := net.SplitHostPort(target); err == nil {
		return target
	}
	return net.JoinHostPort(target, defaultPort)
}
