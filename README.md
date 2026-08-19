# GetNotified

Health-check monitoring that pings your endpoints and tells you — once, calmly —
when something breaks and when it comes back.

Part of [Fajr Labs](https://fajrlabs.com). `// halal software`: transparent
behaviour, no dark patterns, no third-party trackers on public status pages.

The API is the product. The dashboard is one client of it — see [`API.md`](API.md).

## Getting set up

Postgres needs two roles: an owner that runs migrations, and an app role that
the server queries as. Row-level security is enforced against the app role, so
they cannot be the same user.

```sh
createdb getnotified
psql -d getnotified -c "create role gn_app login password '…'"
psql -d getnotified -c "grant connect on database getnotified to gn_app"
```

Start the backend once. It creates the schema, runs River's migrations, creates
your org, and creates the `getnotified_app` role:

```sh
cp .env.example .env    # fill in every value
DATABASE_ADMIN_URL=… DATABASE_URL=… PORT=8080 API_KEY=… ORG_SLUG=default ORG_NAME=… \
  go run ./cmd/server
```

Then grant the app role its membership and restart:

```sh
psql -d getnotified -c "grant getnotified_app to gn_app"
```

Frontend:

```sh
node web/scripts/hash-password.js 'your password'   # -> AUTH_PASSWORD_HASH
cd web && npm install && npm run dev
```

Landing page on `:5173`, dashboard at `/app`, public status page at `/status/default`.

## How it works

**Check engine.** One River periodic job ticks every 5 seconds, claims every
monitor whose next check has come due in a single `UPDATE … RETURNING`, and
queues one check job each. Monitor changes need no queue coordination at all.

**Flapping protection.** A monitor opens an incident only after
`failure_threshold` consecutive failures, so one blip never notifies anyone. The
decision is `NextState` in `internal/jobs/state.go` — pure, and tested without a
database.

**Notification fan-out.** An incident is just a Postgres row. Each attached
channel gets its own job with its own retry budget, so a slow Slack never delays
email.

**Bounded growth.** A monitor checked every 30 seconds writes about a million
check results a year, so an hourly job trims anything past
`CHECK_RETENTION_DAYS` in batches. Incidents are small and are kept for good.

**Isolation.** Every org-scoped table has a row-level security policy bound to
the app role and keyed on `app.org_id`, which the connection pool pins at
startup. With no org in scope the policies return nothing, so it fails closed.

**Errors.** The API defines every failure, its code, and the sentence a person
should read. The dashboard renders that sentence and adds nothing.

## Layout

```
cmd/server/           wiring: config, bootstrap, River client, HTTP server
internal/store/       every SQL statement, the schema, and RLS policies
internal/probe/       http / tcp / ssl_expiry probes
internal/jobs/        River workers and the up/down state machine
internal/notify/      one file per notification channel
internal/httpapi/     routes, validation, and the error envelope
web/                  SvelteKit landing page, dashboard and status page
```

## The iMessage relay

Apple has no public iMessage API, so this one channel needs a Mac. Run
`cmd/imessage-relay` there and point the main server at it:

```sh
RELAY_ADDR=127.0.0.1:8123 IMESSAGE_RELAY_KEY=… go run ./cmd/imessage-relay
```

It exposes `POST /send` behind a bearer key and drives the Messages app. That
Mac needs two things: Messages signed in to iMessage, and Automation permission
for the terminal running the relay (System Settings, Privacy & Security,
Automation). The relay names whichever one is missing. The
recipient and message are passed to AppleScript as arguments, never
interpolated into it, so nothing in a message can be executed. Give the Mac
Automation permission for Messages the first time it runs.

Every other channel works without it.

## Not built yet

- **Kamal deploy config** — Dockerfiles for both services plus `deploy.yml`.
- **Multi-region checks** — `checks.region` exists and is written as `local`;
  the multi-VPS orchestration does not.
- **Billing, teams, RBAC** — deliberately out of scope for v1.

## Tests

```sh
go test ./...
cd web && npm test && npm run check && npm run build
```
