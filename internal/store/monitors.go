package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Monitor struct {
	ID                  string    `db:"id"                   json:"id"`
	OrgID               string    `db:"org_id"               json:"org_id"`
	Name                string    `db:"name"                 json:"name"`
	Type                string    `db:"type"                 json:"type"`
	Target              string    `db:"target"               json:"target"`
	IntervalSeconds     int32     `db:"interval_seconds"     json:"interval_seconds"`
	TimeoutSeconds      int32     `db:"timeout_seconds"      json:"timeout_seconds"`
	ExpectedStatus      []int32   `db:"expected_status"      json:"expected_status"`
	SSLWarnDays         int32     `db:"ssl_warn_days"        json:"ssl_warn_days"`
	FailureThreshold    int32     `db:"failure_threshold"    json:"failure_threshold"`
	Tags                []string  `db:"tags"                 json:"tags"`
	Paused              bool      `db:"paused"               json:"paused"`
	Status              string    `db:"status"               json:"status"`
	ConsecutiveFailures int32     `db:"consecutive_failures" json:"consecutive_failures"`
	NextCheckAt         time.Time `db:"next_check_at"        json:"next_check_at"`
	CreatedAt           time.Time `db:"created_at"           json:"created_at"`
}

// Summary is what the dashboard and the public status page both render, so
// uptime is computed once in SQL rather than per caller.
type Summary struct {
	ID        string   `db:"id"         json:"id"`
	Name      string   `db:"name"       json:"name"`
	Type      string   `db:"type"       json:"type"`
	Target    string   `db:"target"     json:"target"`
	Status    string   `db:"status"     json:"status"`
	Paused    bool     `db:"paused"     json:"paused"`
	Tags      []string `db:"tags"       json:"tags"`
	Up24h     *float64 `db:"up_24h"     json:"up_24h"`
	Up7d      *float64 `db:"up_7d"      json:"up_7d"`
	Up30d     *float64 `db:"up_30d"     json:"up_30d"`
	LatencyMS *int32   `db:"latency_ms" json:"latency_ms"`
	Recent    []bool   `db:"recent"     json:"recent"`
}

// Detail is a monitor together with how it has been doing, so the page for one
// monitor needs a single request.
type Detail struct {
	Monitor
	Up24h     *float64 `db:"up_24h"     json:"up_24h"`
	Up7d      *float64 `db:"up_7d"      json:"up_7d"`
	Up30d     *float64 `db:"up_30d"     json:"up_30d"`
	LatencyMS *int32   `db:"latency_ms" json:"latency_ms"`
}

// MonitorInput carries pointers so one type serves both create (defaults filled
// in by SQL) and patch (absent fields keep their stored value).
type MonitorInput struct {
	Name             *string  `json:"name"`
	Type             *string  `json:"type"`
	Target           *string  `json:"target"`
	IntervalSeconds  *int32   `json:"interval_seconds"`
	TimeoutSeconds   *int32   `json:"timeout_seconds"`
	ExpectedStatus   []int32  `json:"expected_status"`
	SSLWarnDays      *int32   `json:"ssl_warn_days"`
	FailureThreshold *int32   `json:"failure_threshold"`
	Tags             []string `json:"tags"`
	Paused           *bool    `json:"paused"`
}

const monitorCols = `id, org_id, name, type, target, interval_seconds, timeout_seconds,
	expected_status, ssl_warn_days, failure_threshold, tags, paused, status,
	consecutive_failures, next_check_at, created_at`

func GetMonitor(ctx context.Context, db Queryer, id string) (Monitor, error) {
	rows, err := db.Query(ctx, `select `+monitorCols+` from monitors where id = $1`, id)
	if err != nil {
		return Monitor{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Monitor])
}

// LockMonitor reads a monitor for update. The failure count is read, decided
// on, and written, so the row has to be held for the whole decision.
func LockMonitor(ctx context.Context, db Queryer, id string) (Monitor, error) {
	rows, err := db.Query(ctx, `select `+monitorCols+` from monitors where id = $1 for update`, id)
	if err != nil {
		return Monitor{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Monitor])
}

func CreateMonitor(ctx context.Context, db Queryer, orgID string, in MonitorInput) (Monitor, error) {
	rows, err := db.Query(ctx, `
		insert into monitors (org_id, name, type, target, interval_seconds, timeout_seconds,
		                      expected_status, ssl_warn_days, failure_threshold, tags, paused)
		values ($1, $2, coalesce($3::text,'http'), $4, coalesce($5::int,60), coalesce($6::int,10),
		        coalesce($7::int[],'{200}'), coalesce($8::int,14), coalesce($9::int,2),
		        coalesce($10::text[],'{}'), coalesce($11::bool,false))
		returning `+monitorCols,
		orgID, in.Name, in.Type, in.Target, in.IntervalSeconds, in.TimeoutSeconds,
		in.ExpectedStatus, in.SSLWarnDays, in.FailureThreshold, in.Tags, in.Paused)
	if err != nil {
		return Monitor{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Monitor])
}

func UpdateMonitor(ctx context.Context, db Queryer, id string, in MonitorInput) (Monitor, error) {
	rows, err := db.Query(ctx, `
		update monitors set
		  name = coalesce($2, name),
		  type = coalesce($3, type),
		  target = coalesce($4, target),
		  interval_seconds = coalesce($5, interval_seconds),
		  timeout_seconds = coalesce($6, timeout_seconds),
		  expected_status = coalesce($7::int[], expected_status),
		  ssl_warn_days = coalesce($8, ssl_warn_days),
		  failure_threshold = coalesce($9, failure_threshold),
		  tags = coalesce($10::text[], tags),
		  paused = coalesce($11, paused),
		  next_check_at = least(next_check_at, now() + make_interval(secs => coalesce($5, interval_seconds)))
		where id = $1
		returning `+monitorCols,
		id, in.Name, in.Type, in.Target, in.IntervalSeconds, in.TimeoutSeconds,
		in.ExpectedStatus, in.SSLWarnDays, in.FailureThreshold, in.Tags, in.Paused)
	if err != nil {
		return Monitor{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Monitor])
}

func DeleteMonitor(ctx context.Context, db Queryer, id string) (bool, error) {
	tag, err := db.Exec(ctx, `delete from monitors where id = $1`, id)
	return tag.RowsAffected() > 0, err
}

// uptimeLateral computes uptime windows and the latest latency for the monitor
// in scope. It is shared so the list and the single view can never disagree.
const uptimeLateral = `
	left join lateral (
	  select avg(ok::int) filter (where checked_at > now() - interval '24 hours') as up_24h,
	         avg(ok::int) filter (where checked_at > now() - interval '7 days')   as up_7d,
	         avg(ok::int)                                                         as up_30d,
	         (array_agg(latency_ms order by checked_at desc))[1]                  as latency_ms
	    from checks
	   where monitor_id = m.id and checked_at > now() - interval '30 days'
	) u on true`

// RecentCheckCount is how many results the little strip on each row shows.
const RecentCheckCount = 24

func Summaries(ctx context.Context, db Queryer, orgID string) ([]Summary, error) {
	rows, err := db.Query(ctx, `
		select m.id, m.name, m.type, m.target, m.status, m.paused, m.tags,
		       u.up_24h, u.up_7d, u.up_30d, u.latency_ms,
		       coalesce(r.recent, '{}') as recent
		  from monitors m`+uptimeLateral+`
		  left join lateral (
		    select array_agg(ok order by checked_at) as recent
		      from (
		        select ok, checked_at from checks
		         where monitor_id = m.id order by checked_at desc limit $2
		      ) latest
		  ) r on true
		 where m.org_id = $1
		 order by m.name`, orgID, RecentCheckCount)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Summary])
}

func GetMonitorDetail(ctx context.Context, db Queryer, id string) (Detail, error) {
	rows, err := db.Query(ctx, `
		select m.id, m.org_id, m.name, m.type, m.target, m.interval_seconds, m.timeout_seconds,
		       m.expected_status, m.ssl_warn_days, m.failure_threshold, m.tags, m.paused, m.status,
		       m.consecutive_failures, m.next_check_at, m.created_at,
		       u.up_24h, u.up_7d, u.up_30d, u.latency_ms
		  from monitors m`+uptimeLateral+`
		 where m.id = $1`, id)
	if err != nil {
		return Detail{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Detail])
}

// ClaimDue reschedules every monitor whose check is due and returns their ids,
// atomically, so overlapping ticks cannot double-queue a check.
func ClaimDue(ctx context.Context, db Queryer) ([]string, error) {
	rows, err := db.Query(ctx, `
		update monitors
		   set next_check_at = now() + make_interval(secs => interval_seconds)
		 where not paused and next_check_at <= now()
		returning id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func SetMonitorState(ctx context.Context, db Queryer, id, status string, failures int32) error {
	_, err := db.Exec(ctx,
		`update monitors set status = $2, consecutive_failures = $3 where id = $1`,
		id, status, failures)
	return err
}
