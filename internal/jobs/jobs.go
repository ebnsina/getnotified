// Package jobs runs the check schedule and the notification fan-out on River.
package jobs

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"github.com/fajrlabs/getnotified/internal/notify"
	"github.com/fajrlabs/getnotified/internal/probe"
	"github.com/fajrlabs/getnotified/internal/store"
)

// TickInterval bounds scheduling jitter, not the check interval itself.
const TickInterval = 5 * time.Second

// ScheduleArgs sweeps the whole fleet on a fixed tick. One periodic job instead
// of one per monitor means monitor CRUD needs no queue coordination.
type ScheduleArgs struct{}

func (ScheduleArgs) Kind() string { return "schedule" }

type CheckArgs struct {
	MonitorID string `json:"monitor_id"`
}

func (CheckArgs) Kind() string { return "check" }

type NotifyArgs struct {
	ChannelID  string `json:"channel_id"`
	IncidentID string `json:"incident_id"`
	Event      string `json:"event"` // "down" | "up"
}

func (NotifyArgs) Kind() string { return "notify" }

type ScheduleWorker struct {
	river.WorkerDefaults[ScheduleArgs]
	Pool *pgxpool.Pool
}

func (w *ScheduleWorker) Work(ctx context.Context, _ *river.Job[ScheduleArgs]) error {
	ids, err := store.ClaimDue(ctx, w.Pool)
	if err != nil || len(ids) == 0 {
		return err
	}

	params := make([]river.InsertManyParams, len(ids))
	for i, id := range ids {
		params[i] = river.InsertManyParams{Args: CheckArgs{MonitorID: id}}
	}
	_, err = river.ClientFromContext[pgx.Tx](ctx).InsertMany(ctx, params)
	return err
}

type CheckWorker struct {
	river.WorkerDefaults[CheckArgs]
	Pool *pgxpool.Pool
}

func (w *CheckWorker) Work(ctx context.Context, job *river.Job[CheckArgs]) error {
	m, err := store.GetMonitor(ctx, w.Pool, job.Args.MonitorID)
	if errors.Is(err, pgx.ErrNoRows) {
		return river.JobCancel(err)
	}
	if err != nil || m.Paused {
		return err
	}

	res := probe.Run(ctx, m)
	var errMsg *string
	if res.Error != "" {
		errMsg = &res.Error
	}
	if err := store.RecordCheck(ctx, w.Pool, m.ID, res.OK, res.LatencyMS, res.StatusCode, errMsg); err != nil {
		return err
	}
	return w.transition(ctx, m, res)
}

// transition persists the new state and, on a real flip, opens or closes the
// incident and fans out notifications.
func (w *CheckWorker) transition(ctx context.Context, m store.Monitor, res probe.Result) error {
	status, failures, kind := NextState(m, res.OK)

	var incidentID string
	err := pgx.BeginFunc(ctx, w.Pool, func(tx pgx.Tx) error {
		if err := store.SetMonitorState(ctx, tx, m.ID, status, failures); err != nil {
			return err
		}
		var err error
		switch kind {
		case "down":
			incidentID, err = store.OpenIncident(ctx, tx, m.ID, res.Error)
		case "up":
			incidentID, err = store.ResolveIncident(ctx, tx, m.ID)
		default:
			return nil
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // a racing check already opened or closed it
		}
		return err
	})
	if err != nil || incidentID == "" {
		return err
	}
	return w.fanOut(ctx, m.ID, incidentID, kind)
}

// fanOut gives every channel its own job, so one slow channel never holds up
// the rest.
func (w *CheckWorker) fanOut(ctx context.Context, monitorID, incidentID, kind string) error {
	deliveries, err := store.Deliveries(ctx, w.Pool, monitorID)
	if err != nil {
		return err
	}
	if len(deliveries) == 0 {
		slog.Info("incident with no channels attached", "monitor", monitorID, "event", kind)
		return nil
	}

	params := make([]river.InsertManyParams, len(deliveries))
	for i, d := range deliveries {
		params[i] = river.InsertManyParams{
			Args:       NotifyArgs{ChannelID: d.ChannelID, IncidentID: incidentID, Event: kind},
			InsertOpts: &river.InsertOpts{MaxAttempts: int(d.MaxAttempts)},
		}
	}
	_, err = river.ClientFromContext[pgx.Tx](ctx).InsertMany(ctx, params)
	return err
}

type NotifyWorker struct {
	river.WorkerDefaults[NotifyArgs]
	Pool *pgxpool.Pool
}

func (w *NotifyWorker) Work(ctx context.Context, job *river.Job[NotifyArgs]) error {
	ch, err := store.GetChannel(ctx, w.Pool, job.Args.ChannelID)
	if err != nil {
		return cancelIfGone(err)
	}
	inc, err := store.GetIncident(ctx, w.Pool, job.Args.IncidentID)
	if err != nil {
		return cancelIfGone(err)
	}
	m, err := store.GetMonitor(ctx, w.Pool, inc.MonitorID)
	if err != nil {
		return cancelIfGone(err)
	}

	n, err := notify.For(ch.Type)
	if err != nil {
		return river.JobCancel(err) // no retry will ever fix an unknown type
	}
	return n.Send(ctx, notify.Event{Kind: job.Args.Event, Monitor: m, Incident: inc}, ch)
}

func cancelIfGone(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return river.JobCancel(err)
	}
	return err
}

func Workers(pool *pgxpool.Pool) *river.Workers {
	w := river.NewWorkers()
	river.AddWorker(w, &ScheduleWorker{Pool: pool})
	river.AddWorker(w, &CheckWorker{Pool: pool})
	river.AddWorker(w, &NotifyWorker{Pool: pool})
	return w
}

func PeriodicJobs() []*river.PeriodicJob {
	return []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(TickInterval),
			func() (river.JobArgs, *river.InsertOpts) { return ScheduleArgs{}, nil },
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	}
}
