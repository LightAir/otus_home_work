package app

import (
	"context"
	"fmt"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage interface {
	AddEvent(event storage.Event) error
	ChangeEvent(id uuid.UUID, event storage.Event) error
	RemoveEvent(id uuid.UUID) error
	ListEventsByRange(p storage.DateRange) (map[uuid.UUID]storage.Event, error)
	Connect(ctx context.Context) error
	Close() error
}

type App struct {
	config  config.Config
	storage Storage
}

func New(storage Storage, cfg *config.Config) *App {
	return &App{
		config:  *cfg,
		storage: storage,
	}
}

func buildEvent(id, title, start, end, desc, userID, when string) (*storage.Event, error) {
	dateStart, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, fmt.Errorf("bad start date. %w", err)
	}

	dateEnd, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, fmt.Errorf("bad end date. %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad Id. %w", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("bad userId. %w", err)
	}

	dateWhen, err := time.Parse(time.RFC3339, when)
	if err != nil {
		return nil, fmt.Errorf("bad when date. %w", err)
	}

	event := &storage.Event{
		ID:            parsedID,
		Title:         title,
		DatetimeStart: dateStart,
		DatetimeEnd:   dateEnd,
		Description:   desc,
		UserID:        parsedUserID,
		WhenToNotify:  dateWhen,
	}

	return event, nil
}

func (a *App) CreateEvent(id, title, start, end, desc, userID, when string) error {
	event, err := buildEvent(id, title, start, end, desc, userID, when)
	if err != nil {
		return err
	}

	return a.storage.AddEvent(*event)
}

func (a *App) UpdateEvent(id, title, start, end, desc, userID, when string) error {
	event, err := buildEvent(id, title, start, end, desc, userID, when)
	if err != nil {
		return err
	}

	return a.storage.ChangeEvent(event.ID, *event)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.RemoveEvent(uuid.MustParse(id))
}

func (a *App) FindEventsByPeriod(start, end time.Time) (map[uuid.UUID]storage.Event, error) {
	dateRange := storage.DateRange{
		Start: start,
		End:   end,
	}

	return a.storage.ListEventsByRange(dateRange)
}
