package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Incident struct {
	ID         string     `db:"id"          json:"id"`
	MonitorID  string     `db:"monitor_id"  json:"monitor_id"`
	StartedAt  time.Time  `db:"started_at"  json:"started_at"`
	ResolvedAt *time.Time `db:"resolved_at" json:"resolved_at"`
	Cause      *string    `db:"cause"       json:"cause"`
}

const incidentCols = `id, monitor_id, started_at, resolved_at, cause`

func GetIncident(ctx context.Context, db Queryer, id string) (Incident, error) {
	rows, err := db.Query(ctx, `select `+incidentCols+` from incidents where id = $1`, id)
	if err != nil {
		return Incident{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Incident])
}

// ListIncidents returns the org's incidents, narrowed to one monitor when
// monitorID is non-empty.
func ListIncidents(ctx context.Context, db Queryer, orgID, monitorID string, limit int) ([]Incident, error) {
	var monitor *string
	if monitorID != "" {
		monitor = &monitorID
	}
	rows, err := db.Query(ctx, `
		select i.id, i.monitor_id, i.started_at, i.resolved_at, i.cause
		  from incidents i
		  join monitors m on m.id = i.monitor_id
		 where m.org_id = $1 and ($2::uuid is null or i.monitor_id = $2)
		 order by i.started_at desc limit $3`, orgID, monitor, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Incident])
}

// OpenIncident returns an empty id when one is already open for the monitor.
func OpenIncident(ctx context.Context, db Queryer, monitorID, cause string) (string, error) {
	var id string
	err := db.QueryRow(ctx, `
		insert into incidents (monitor_id, cause) values ($1, $2)
		on conflict do nothing returning id`, monitorID, cause).Scan(&id)
	return id, err
}

// ResolveIncident returns an empty id when nothing was open.
func ResolveIncident(ctx context.Context, db Queryer, monitorID string) (string, error) {
	var id string
	err := db.QueryRow(ctx, `
		update incidents set resolved_at = now()
		 where monitor_id = $1 and resolved_at is null returning id`, monitorID).Scan(&id)
	return id, err
}
