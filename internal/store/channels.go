package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Channel struct {
	ID     string         `db:"id"     json:"id"`
	OrgID  string         `db:"org_id" json:"org_id"`
	Name   string         `db:"name"   json:"name"`
	Type   string         `db:"type"   json:"type"`
	Config map[string]any `db:"config" json:"config"`
}

const channelCols = `id, org_id, name, type, config`

func GetChannel(ctx context.Context, db Queryer, id string) (Channel, error) {
	rows, err := db.Query(ctx, `select `+channelCols+` from notification_channels where id = $1`, id)
	if err != nil {
		return Channel{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Channel])
}

func ListChannels(ctx context.Context, db Queryer, orgID string) ([]Channel, error) {
	rows, err := db.Query(ctx,
		`select `+channelCols+` from notification_channels where org_id = $1 order by name`, orgID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Channel])
}

func CreateChannel(ctx context.Context, db Queryer, orgID, name, kind string, config map[string]any) (Channel, error) {
	rows, err := db.Query(ctx, `
		insert into notification_channels (org_id, name, type, config)
		values ($1, $2, $3, $4) returning `+channelCols, orgID, name, kind, config)
	if err != nil {
		return Channel{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Channel])
}

func DeleteChannel(ctx context.Context, db Queryer, id string) (bool, error) {
	tag, err := db.Exec(ctx, `delete from notification_channels where id = $1`, id)
	return tag.RowsAffected() > 0, err
}

func MonitorChannelIDs(ctx context.Context, db Queryer, monitorID string) ([]string, error) {
	rows, err := db.Query(ctx,
		`select channel_id from monitor_channels where monitor_id = $1`, monitorID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

// Delivery is one channel to notify, with its own retry budget.
type Delivery struct {
	ChannelID   string `db:"channel_id"`
	MaxAttempts int32  `db:"max_attempts"`
}

func Deliveries(ctx context.Context, db Queryer, monitorID string) ([]Delivery, error) {
	rows, err := db.Query(ctx, `
		select c.id as channel_id, coalesce((c.config->>'max_attempts')::int, 5) as max_attempts
		  from monitor_channels mc
		  join notification_channels c on c.id = mc.channel_id
		 where mc.monitor_id = $1`, monitorID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Delivery])
}

func SetMonitorChannels(ctx context.Context, db Beginner, monitorID string, channelIDs []string) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx,
			`delete from monitor_channels where monitor_id = $1`, monitorID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			insert into monitor_channels (monitor_id, channel_id)
			select $1, c.id from notification_channels c where c.id = any($2::uuid[])`,
			monitorID, channelIDs)
		return err
	})
}
