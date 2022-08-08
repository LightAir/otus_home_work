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

	if _, isExist := s.items[event.ID]; isExist {
		return errEventAlreadyExist
	}

	s.items[event.ID] = event
	s.mu.Unlock()

	return nil
}

func (s *Storage) ChangeEvent(id uuid.UUID, event storage.Event) error {
	s.mu.Lock()

	if _, isExist := s.items[id]; isExist {
		s.items[id] = event
	} else {
		return errEventNotFound
	}

	s.mu.Unlock()

	return nil
}

func (s *Storage) RemoveEvent(id uuid.UUID) error {
	if _, isExist := s.items[id]; isExist {
		delete(s.items, id)
	} else {
		return errEventNotFound
	}

	return nil
}

func (s *Storage) ListEventsByUserID(userID uuid.UUID) (map[uuid.UUID]storage.Event, error) {
	result := make(map[uuid.UUID]storage.Event)

	for _, item := range s.items {
		if item.UserID == userID {
			result[item.ID] = item
		}
	}

	return result, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	s.items = make(map[uuid.UUID]storage.Event)

	return nil
}

func (s *Storage) Close() error {
	s.items = nil

	return nil
}

func New() *Storage {
	return &Storage{}
}
