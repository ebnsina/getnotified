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

- Send a test message through any channel from its row in the dashboard, or
  `POST /api/channels/{id}/test`. It reports back what happened rather than
  queueing, so a channel can be checked before anyone relies on it.
- The macOS iMessage relay (`cmd/imessage-relay`), the one channel Apple gives
  no API for. Recipient and message reach AppleScript as arguments, so a
  message cannot execute anything.
- A failed delivery now carries the far end's own explanation instead of just a
  status code, which is the whole point of being able to test a channel. The
  iMessage relay recognises the three macOS refusals that matter — no Automation
  permission, no iMessage account signed in, and an unrecognised recipient — and
  says what to do about each.
- Delivery failures are explained in plain words — nothing answered, it timed
  out, the address refused it — instead of socket wording.

### Fixed

- Email over port 465 failed with a confusing handshake error. Go's SMTP client
  starts plain and upgrades with STARTTLS, so it cannot open a connection that
  is encrypted from the first byte. It now says to use 587 or 2525.

- The front page demo shifted the whole page each cycle. Failed bars are taller
  than healthy ones, so the row grew the moment one appeared and everything
  below it jumped. The row now reserves the tallest bar's height, and the
  caption reserves its line.

- Consecutive failures could be undercounted. The count was read in one
  statement and written in another, so two checks running at once could each
  write back the same figure and an incident would open late. The monitor row
  is now held for the whole decision.
- A check could outlive the gap between checks, which is what let two run at
  once. `timeout_seconds` may no longer exceed `interval_seconds`, enforced by
  a database constraint as well as by the API, so a partial update cannot slip
  past it either.
- Check results are now trimmed on a schedule. The table previously grew
  without limit.

### Changed

- The front page now says what the product does: what it checks, where it tells
  you, what you get, and an example of the API. It was atmosphere with no
  substance before.
- The dashboard is set in the sans. The serif stays on the front page and the
  status pages, where it is doing brand work rather than being read all day.
- Each monitor in the list shows its last 24 results as a strip, and the page
  for one monitor now shows its uptime and latency, which it never did.
- Status reads as a pill rather than a dot and a word, and panels share the
  radius of the controls inside them.
- Secondary buttons are outlined rather than ghosted, and the ghost variant is
  gone. The status page link uses an arrow, not a chevron.
- Every serif heading is italic, and the wordmark is set in the serif too.
- The front page sets its text in Inter Tight; the dashboard keeps Geist, which
  pairs with the mono it sits beside. The marketing surface and the working
  surface each get the face that suits the job.
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
