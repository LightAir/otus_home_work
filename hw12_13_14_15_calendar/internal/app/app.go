package app

import (
	"context"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	config  config.Config
	storage Storage
	logger  Logger
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	AddEvent(event storage.Event) error
	ChangeEvent(id uuid.UUID, event storage.Event) error
	RemoveEvent(id uuid.UUID) error
	ListEventsByUserID(userID uuid.UUID) (map[uuid.UUID]storage.Event, error)
	Connect(ctx context.Context) error
	Close() error
}

func New(logger Logger, storage Storage, cfg *config.Config) *App {
	return &App{
		config:  *cfg,
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
