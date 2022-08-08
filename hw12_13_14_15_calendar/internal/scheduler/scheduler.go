package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/queue"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Error(msg ...interface{})
	Info(msg ...interface{})
}

type Storage interface {
	RemoveOldEvents(datetime time.Time) error
	ListEventsForNotification(datetime time.Time) (map[uuid.UUID]storage.Event, error)
	SetIsNotified(id uuid.UUID) error
	Connect(ctx context.Context) error
	Close() error
}

type Scheduler struct {
	queue     queue.Queue
	logger    Logger
	storage   Storage
	queueName string
}

func NewScheduler(queue queue.Queue, logger logger.Logger, storage Storage, queueName string) *Scheduler {
	return &Scheduler{
		queue:     queue,
		logger:    logger,
		storage:   storage,
		queueName: queueName,
	}
}

func (sch *Scheduler) Run(ctx context.Context) error {
	if err := sch.queue.Connect(ctx); err != nil {
		return fmt.Errorf("error queue start: %w", err)
	}

	notyTimer := time.NewTicker(30 * time.Second)
	removeTimer := time.NewTicker(10 * time.Minute)

	sch.logger.Info("Calendar scheduler started...")

	for {
		select {
		case <-notyTimer.C:
			err := sch.Noty()
			if err != nil {
				sch.logger.Errorf("notification error: %e", err)
			}
		case <-removeTimer.C:
			if err := sch.Remove(); err != nil {
				sch.logger.Errorf("remove events error: %e", err)
			}
		case <-ctx.Done():
			notyTimer.Stop()
			removeTimer.Stop()

			return nil
		}
	}
}

func (sch *Scheduler) Noty() error {
	events, err := sch.storage.ListEventsForNotification(time.Now())
	if err != nil {
		return err
	}

	for _, event := range events {
		body, err := json.Marshal(event)
		if err != nil {
			sch.logger.Errorf("marshal error: %e", err)
		}

		if err := sch.queue.Sent(body, sch.queueName); err != nil {
			sch.logger.Errorf("sent error: %e", err)
		}

		if err := sch.storage.SetIsNotified(event.ID); err != nil {
			sch.logger.Errorf("set is notified error: %e", err)
		}
	}

	return nil
}

func (sch *Scheduler) Remove() error {
	t := time.Now().AddDate(-1, 0, 0)

	return sch.storage.RemoveOldEvents(t)
}
