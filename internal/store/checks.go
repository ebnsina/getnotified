package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Check struct {
	ID         int64     `db:"id"          json:"id"`
	MonitorID  string    `db:"monitor_id"  json:"monitor_id"`
	CheckedAt  time.Time `db:"checked_at"  json:"checked_at"`
	OK         bool      `db:"ok"          json:"ok"`
	LatencyMS  int32     `db:"latency_ms"  json:"latency_ms"`
	StatusCode *int32    `db:"status_code" json:"status_code"`
	Error      *string   `db:"error"       json:"error"`
	Region     string    `db:"region"      json:"region"`
}

func RecordCheck(ctx context.Context, db Queryer, monitorID string, ok bool, latencyMS int32, code *int32, errMsg *string) error {
	_, err := db.Exec(ctx, `
		insert into checks (monitor_id, ok, latency_ms, status_code, error)
		values ($1, $2, $3, $4, $5)`, monitorID, ok, latencyMS, code, errMsg)
	return err
}

func ListChecks(ctx context.Context, db Queryer, monitorID string, limit int) ([]Check, error) {
	rows, err := db.Query(ctx, `
		select id, monitor_id, checked_at, ok, latency_ms, status_code, error, region
		  from checks where monitor_id = $1 order by checked_at desc limit $2`,
		monitorID, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Check])
}
