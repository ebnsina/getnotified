# GetNotified

Uptime monitoring that checks sites, ports and certificates, and says something
once when one breaks and once when it comes back. Go and Postgres on the back,
SvelteKit on the front. Single-tenant on purpose: one deploy per customer.

## Commands

```sh
go build ./... && go test ./...          # backend
cd web && npm test && npm run check      # frontend types and unit tests
cd web && npm run build                  # production build
```

Running it locally needs two Postgres roles, because row-level security is
enforced against the app role and an owner cannot be subject to it:

```sh
DATABASE_ADMIN_URL=postgres:///getnotified            # owner: migrations, org, grants
DATABASE_URL=postgres://gn_app:…@localhost/getnotified # app role, member of getnotified_app
PORT=8080 API_KEY=… ORG_SLUG=default ORG_NAME=… CHECK_RETENTION_DAYS=90 go run ./cmd/server
cd web && npm run dev            # :5173, reads web/.env
```

## Layout

```
cmd/server/          wiring: config, bootstrap, River client, HTTP server
cmd/imessage-relay/  optional macOS service; drives Messages.app
internal/store/      every SQL statement, the schema, RLS policies
internal/probe/      http / tcp / ssl_expiry probes
internal/jobs/       River workers and the up/down state machine
internal/notify/     one file per channel
internal/httpapi/    routes, validation, the error envelope
web/                 landing page, dashboard under /app, public status pages
```

## How it works

One River periodic job sweeps every due monitor in a single `UPDATE … RETURNING`
and queues a probe each, so monitor changes need no queue coordination. An
incident opens only after `failure_threshold` consecutive failures; the decision
is `NextState` in `internal/jobs/state.go`, kept pure and tested, and it runs
against a locked row so concurrent checks cannot lose a failure. Each attached
channel gets its own job and its own retry budget.

## Conventions

**Git.** Author as `ebnsina <ebnsina.me@gmail.com>`. No `Co-Authored-By` trailer.
Remote uses the `github-es` SSH host alias. `docs/` and `data/` are local only
and gitignored — no plans or roadmaps in the public repo.

**Errors.** The API is the single source of truth. Every failure returns
`{"error":{"code","message"}}` where the message is written for a person, and
the frontend shows it verbatim. Database and socket wording never reaches a
reader: `internal/notify/explain.go` translates transport failures, and
`internal/httpapi/errors.go` maps database faults to a generic message while
logging the real one. When a remote rejects a delivery, quote its explanation —
that is the whole point of being able to test a channel.

**Language.** No implementation terms in the interface. Not endpoints, webhook
headers, binaries, or status codes as nouns. "A web address", not "HTTP target".

**Comments.** One line, two at most. Say why, not what.

**Configuration.** Nothing has a default. A missing variable stops the server
with the name of what is missing. On the frontend the check is lazy, so a build
does not need the values. `PUBLIC_` is reserved by SvelteKit — use `SITE_ORIGIN`.

**Validation.** The API decides. Valibot mirrors the rules client-side only to
save a round trip. Rules that a partial update could slip past belong in a
database constraint as well, reported in the same words.

**Formatting.** `Intl` everywhere — numbers, dates, durations, relative time.
Locale comes from `Accept-Language` through the layout. Never hand-roll.

**Isolation.** Every org-scoped table has an RLS policy bound to
`getnotified_app` and keyed on `app.org_id`, which the pool pins at startup. It
fails closed when no org is in scope.

**Changelog.** Keep `CHANGELOG.md` current. `API.md` documents the whole API.

## Design

Dark throughout, one palette. The accent is the product's own up/down green
rather than an arbitrary brand colour.

- **Newsreader**, always italic, for display — front page and status pages only.
  The dashboard is a working surface and uses the sans.
- **Inter Tight** for the landing page, **Geist** for the app, **Geist Mono**
  for anything numeric, data-like, or a label.
- Buttons, inputs, selects and panels share a `0.75rem` radius. Focus is a light
  border plus a 2px primary ring at 2px offset. Secondary buttons are outlined,
  never ghosted.
- Icons come from `@hugeicons/core-free-icons` and render through
  `web/src/lib/Icon.svelte`. Do not use `@hugeicons/svelte` — it builds its SVG
  in `onMount`, so nothing server-renders.
- Status pages carry no third-party scripts, ever.

## State of the channels

Verified end to end: **webhook**, **email** (including `smtp.PlainAuth` against
a real server). Written but never run against a live provider: **Slack**,
**SMS**, **WhatsApp**. **iMessage** works only from a Mac that is signed in to
Messages with Automation permission granted.

## Not built

Kamal deploy config. Multi-tenancy — deliberately: the product is self-hosted,
so onboarding a customer means they deploy their own instance. Billing, teams
and RBAC. Multi-region checks, though `checks.region` exists for it.
