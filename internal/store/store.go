// Package store owns every SQL statement in GetNotified. Nothing above it
// writes queries.
package store

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

//go:embed grants.sql
var grantsSQL string

// Queryer is the slice of pgx that both *pgxpool.Pool and pgx.Tx satisfy.
type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// Beginner is satisfied by *pgxpool.Pool, for the few multi-statement writes.
type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Migrate applies the idempotent schema, including row-level security policies.
func Migrate(ctx context.Context, conn Queryer) error {
	_, err := conn.Exec(ctx, schemaSQL)
	return err
}

// Grants hands the app role access to every table, including River's own.
func Grants(ctx context.Context, conn Queryer) error {
	_, err := conn.Exec(ctx, grantsSQL)
	return err
}

type Org struct {
	ID   string `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

// EnsureOrg runs as the admin role: row-level security would otherwise reject
// the insert, since no org is in scope yet.
func EnsureOrg(ctx context.Context, conn Queryer, slug, name string) (string, error) {
	var id string
	err := conn.QueryRow(ctx, `
		insert into organizations (slug, name) values ($1, $2)
		on conflict (slug) do update set name = excluded.name
		returning id`, slug, name).Scan(&id)
	return id, err
}

func OrgBySlug(ctx context.Context, conn Queryer, slug string) (Org, error) {
	rows, err := conn.Query(ctx, `select id, name, slug from organizations where slug = $1`, slug)
	if err != nil {
		return Org{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Org])
}

// NewPool opens the app-role pool and pins every connection to one org, which
// is what the row-level security policies read.
func NewPool(ctx context.Context, url, orgID string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `select set_config('app.org_id', $1, false)`, orgID)
		return err
	}
	return pgxpool.NewWithConfig(ctx, cfg)
}
