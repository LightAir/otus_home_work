package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Elements map[uuid.UUID]storage.Event

var (
	errEventAlreadyExist = errors.New("event already exist")
	errEventNotFound     = errors.New("event not found")
)

type Storage struct {
	items Elements
	mu    sync.RWMutex
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.items[event.ID]; isExist {
		return errEventAlreadyExist
	}

	s.items[event.ID] = event

	return nil
}

func (s *Storage) ChangeEvent(id uuid.UUID, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.items[id]; !isExist {
		return errEventNotFound
	}

	s.items[id] = event

	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.items[id]; isExist {
		delete(s.items, id)
	} else {
		return errEventNotFound
	}

	return nil
}

func (s *Storage) ListEventsByRange(p storage.DateRange) (map[uuid.UUID]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[uuid.UUID]storage.Event)

	for _, item := range s.items {
		lessOrEqualStart := item.DatetimeStart.After(p.Start) || item.DatetimeStart.Equal(p.Start)
		moreOrEqualStart := item.DatetimeStart.Before(p.End) || item.DatetimeStart.Equal(p.End)

		lessOrEqualEnd := item.DatetimeEnd.Before(p.Start) || item.DatetimeEnd.Equal(p.Start)
		moreOrEqualEnd := item.DatetimeEnd.After(p.End) || item.DatetimeEnd.Equal(p.End)
		if (lessOrEqualStart && moreOrEqualStart) || (lessOrEqualEnd && moreOrEqualEnd) {
			result[item.ID] = item
		}
	}

	return result, nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = nil

	return nil
}

func New() *Storage {
	return &Storage{
		items: make(map[uuid.UUID]storage.Event),
	}
}
