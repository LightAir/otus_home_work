package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const datetimeFormat = "2006-Jan-02 15:04:05"

type Storage struct {
	dsn    string
	db     *sqlx.DB
	config config.Config
	ctx    context.Context
}

func (s *Storage) AddEvent(e storage.Event) error {
	query := `insert
				into events(id, title, datetime_start, datetime_end, description, user_id, when_to_notify)
				values($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.ExecContext(
		s.ctx,
		query,
		e.ID,
		e.Title,
		e.DatetimeStart,
		e.DatetimeEnd,
		e.Description,
		e.UserID,
		e.WhenToNotify)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ChangeEvent(id uuid.UUID, event storage.Event) error {
	query := `update
				events
			  set
				title = $1,
				datetime_start = $2,
				datetime_end = $3,
				description = $4,
				user_id = $5,
				when_to_notify = $6
			  where
			    id = $7`

	_, err := s.db.ExecContext(
		s.ctx,
		query,
		event.Title,
		event.DatetimeStart.Format(time.RFC3339),
		event.DatetimeEnd.Format(time.RFC3339),
		event.Description,
		event.UserID,
		event.WhenToNotify,
		id,
	)

	return err
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	query := "delete from events where id = $1"

	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListEventsByRange(p storage.DateRange) (map[uuid.UUID]storage.Event, error) {
	query := `select 
    			id,
    			title,
    			datetime_start as datetimeStart,
    			datetime_end as datetimeEnd,
    			description,
    			user_id as userId,
    			when_to_notify as whenToNotify,
    			is_notified as isNotified
			  from
			    events
			  where
			    (datetime_start >= $1 and datetime_start <= $2) or (datetime_end >= $1 and datetime_end <= $2)`

	rows, err := s.db.QueryxContext(s.ctx, query, p.Start.Format(time.RFC3339), p.End.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make(map[uuid.UUID]storage.Event)

	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events[event.ID] = event
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) RemoveOldEvents(datetime time.Time) error {
	query := `delete from events where datetime_end < $1`

	_, err := s.db.ExecContext(s.ctx, query, datetime.Format(datetimeFormat))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) SetIsNotified(id uuid.UUID) error {
	query := `update events set is_notified = true where id = $1`

	_, err := s.db.ExecContext(s.ctx, query, id)

	return err
}

func (s *Storage) ListEventsForNotification(datetime time.Time) (map[uuid.UUID]storage.Event, error) {
	query := `select 
    			id,
    			title,
    			datetime_start as datetimeStart,
    			datetime_end as datetimeEnd,
    			description,
    			user_id as userId,
    			is_notified as isNotified,
    			when_to_notify as whenToNotify
			  from
			    events
			  where
			    when_to_notify <= $1
				and is_notified = false`

	rows, err := s.db.QueryxContext(s.ctx, query, datetime.Format(datetimeFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make(map[uuid.UUID]storage.Event)

	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events[event.ID] = event
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func New(config *config.Config, dsn string) *Storage {
	return &Storage{
		dsn:    dsn,
		config: *config,
		db:     nil,
	}
}
