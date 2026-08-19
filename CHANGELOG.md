# Changelog

All notable changes to GetNotified are recorded here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Check engine on River: one periodic sweep claims every due monitor and
  queues a probe, so monitor changes need no queue coordination.
- HTTP, TCP and SSL-expiry probes.
- Incident state machine with flapping protection — an incident opens only
  after a configurable number of consecutive failures.
- Notification fan-out with one job and one retry budget per channel: Slack,
  email over SMTP, generic webhook, SMS and WhatsApp over Twilio, and iMessage
  through an optional relay.
- REST API covering everything the dashboard can do, documented in `API.md`.
- Single error envelope across the API, with messages written for people.
- Row-level security on every org-scoped table; the app role is confined to one
  org and fails closed when no org is in scope.
- SvelteKit dashboard: monitor list with uptime, monitor create and edit,
  incident timeline, notification channel management.
- Public status page at `/status/{slug}` — server rendered, no authentication,
  no third-party trackers.
- Password sign-in with a stateless signed-cookie session, built on Node's
  `scrypt` and HMAC.
- Client-side form validation with Valibot, mirroring the API's own rules.
- Mona Sans for text, Geist Mono for numbers, both self-hosted.
- Locale-aware formatting throughout via the `Intl` APIs.
- Friendly error pages for 404, 500 and anything else that goes wrong.
- Open Graph and canonical metadata on every page.
- Tests: the flapping rule and probes in Go, the form validators and formatters
  under Node's built-in test runner.
- Public landing page at `/`, with a live demonstration of the flapping rule.
  The dashboard now lives under `/app`.

### Changed

- Every serif heading is italic, matching the front page.
- Form controls are custom drawn and defined once: fields, selects, checkboxes
  and buttons share one radius, one focus treatment, and one set of states.
- One dark palette across the whole product, with Newsreader for headings,
  Geist for text, and Geist Mono for anything numeric.
- Interface copy rewritten in plain words — no implementation terms on screen.
- Icons render on the server. Hugeicons' own component builds its SVG in
  `onMount`, so icons were invisible until the page hydrated; a small local
  component renders the same data directly, and drops a dependency.
- `Intl.DurationFormat` drops zero units, so a zero-length duration rendered as
  an empty string and short ones carried stray gaps. Both are fixed.
