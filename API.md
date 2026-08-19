# GetNotified API

The API is the product. Everything the dashboard does, it does through these
routes — there is no private back channel.

Base URL is whatever `PORT` the server is bound to. All request and response
bodies are JSON.

## Authentication

Every route under `/api` needs a bearer token:

```
Authorization: Bearer $API_KEY
```

`/status/{slug}` and `/healthz` are public and take no credentials.

## Errors

Every failure returns the same envelope, with the same HTTP status:

```json
{ "error": { "code": "monitor_not_found", "message": "That monitor no longer exists." } }
```

`message` is written for a person to read and can be shown to them as-is — the
API is the single source of truth for error copy. Database and network faults
are logged server-side and reported as `internal_error`; internal detail is
never returned.

| Code | Status | Meaning |
| --- | --- | --- |
| `unauthorized` | 401 | Missing or wrong bearer token |
| `monitor_not_found` | 404 | No such monitor |
| `channel_not_found` | 404 | No such notification channel |
| `status_page_not_found` | 404 | No org with that slug |
| `not_found` | 404 | No such route |
| `malformed_json` | 400 | Body was not valid JSON |
| `invalid_<field>` | 400 | That field failed validation; `message` names the fix |
| `internal_error` | 500 | Something failed on our side |

## Monitors

### `GET /api/monitors`

Returns every monitor with uptime already computed.

```json
[
  {
    "id": "0d1e…",
    "name": "Marketing site",
    "type": "http",
    "target": "https://example.com",
    "status": "up",
    "paused": false,
    "tags": ["web"],
    "up_24h": 1.0,
    "up_7d": 0.9993,
    "up_30d": 0.9991,
    "latency_ms": 128
  }
]
```

Uptime fields are ratios between 0 and 1, or `null` when nothing has been
checked in that window.

### `POST /api/monitors`

| Field | Type | Default | Notes |
| --- | --- | --- | --- |
| `name` | string | — | Required |
| `target` | string | — | Required; a full URL for `http`, `host[:port]` otherwise |
| `type` | string | `http` | `http`, `tcp`, or `ssl_expiry` |
| `interval_seconds` | int | `60` | Minimum 10 |
| `timeout_seconds` | int | `10` | 1–120 |
| `expected_status` | int[] | `[200]` | HTTP only |
| `ssl_warn_days` | int | `14` | `ssl_expiry` only |
| `failure_threshold` | int | `2` | Consecutive failures before an incident opens |
| `tags` | string[] | `[]` | |
| `paused` | bool | `false` | |

Returns `201` and the full monitor.

### `GET /api/monitors/{id}`

Returns the full monitor, including scheduling state.

### `PATCH /api/monitors/{id}`

Same fields as create; omitted fields keep their stored value. Changing
`interval_seconds` also pulls the next check forward if it was scheduled
further out. Returns `200` and the updated monitor.

### `DELETE /api/monitors/{id}`

Returns `204`. Check history and incidents go with it.

### `GET /api/monitors/{id}/checks?limit=100`

Most recent first. `limit` defaults to 100, caps at 1000.

### `GET /api/monitors/{id}/incidents`

Most recent first, up to 100. `resolved_at` is `null` while an incident is open.

### `GET /api/monitors/{id}/channels`

Returns an array of channel ids currently attached.

### `PUT /api/monitors/{id}/channels`

```json
{ "channel_ids": ["…", "…"] }
```

Replaces the whole set. Returns `204`.

## Notification channels

### `GET /api/channels`

### `POST /api/channels`

```json
{ "name": "Ops Slack", "type": "slack", "config": { "webhook_url": "https://hooks.slack.com/…" } }
```

| Type | Required config | Optional config |
| --- | --- | --- |
| `slack` | `webhook_url` | `max_attempts` |
| `webhook` | `url` | `secret`, `max_attempts` |
| `email` | `to` | `max_attempts` |
| `sms` | `to` | `max_attempts` |
| `whatsapp` | `to` | `max_attempts` |
| `imessage` | `to` | `max_attempts` |

`max_attempts` sets the retry budget for that channel (default 5). Credentials
— SMTP, Twilio, the iMessage relay — come from server environment variables,
never from `config`.

Returns `201` and the channel.

### `DELETE /api/channels/{id}`

Returns `204`.

### `POST /api/channels/{id}/test`

Sends a test message through the channel immediately, rather than through the
queue, and waits for the result. Returns `200` with `{"delivered": true}`, or
`channel_test_failed` with the reason in plain words — the address refused it,
nothing answered, it timed out.

The message says it is a test, so nobody mistakes it for a real outage. A
`webhook` channel receives the usual payload with `"kind": "test"`.

## Public status

### `GET /status/{slug}`

No authentication. Cached for 30 seconds.

```json
{
  "org": { "id": "…", "name": "Fajr Labs", "slug": "default" },
  "monitors": [ … ],
  "incidents": [ … ],
  "overall": "operational",
  "as_of": "2026-08-19T10:20:30Z"
}
```

`overall` is `operational` or `degraded`.

### `GET /healthz`

Returns `{"status":"ok"}` once the database answers.

## Webhook payload

A `webhook` channel receives this on every transition:

```json
{
  "kind": "down",
  "monitor": { … },
  "incident": { "id": "…", "started_at": "…", "resolved_at": null, "cause": "unexpected status 503" }
}
```

`kind` is `down`, `up`, or `test`. If the channel has a `secret`, it arrives as the
`X-GetNotified-Secret` header. Any non-2xx response is retried with exponential
backoff up to `max_attempts`.
