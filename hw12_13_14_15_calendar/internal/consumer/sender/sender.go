package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/queue"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Sender struct {
	logger    Logger
	queue     queue.Queue
	queueName string
}

func NewSender(queue queue.Queue, logger Logger, queueName string) *Sender {
	return &Sender{
		logger:    logger,
		queue:     queue,
		queueName: queueName,
	}
}

func (s Sender) Send(body []byte) {
	data := storage.Event{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		s.logger.Errorf("unmarshal body error: %e", err)
	}

	if data.Description == "" {
		s.logger.Infof("Dear user. A notice to you \"%s\" on %s", data.Title, data.DatetimeStart)
	} else {
		s.logger.Infof("Dear user. A notice to you \"%s\": %s on %s", data.Title, data.Description, data.DatetimeStart)
	}
}

func (s Sender) Run(ctx context.Context) error {
	err := s.queue.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error connect to rmq: %w", err)
	}

	err = s.queue.Receive(s.queueName, s.Send)
	if err != nil {
		return fmt.Errorf("error of send: %w", err)
	}

	return nil
}
