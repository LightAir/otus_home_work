package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type Log struct{}

func (l Log) Info(...interface{})  {}
func (l Log) Error(...interface{}) {}
func (l Log) LogGRPCRequest(ctx context.Context, info *grpc.UnaryServerInfo, d time.Duration, statusCode string) {
}

func prepareServer() *GRPCServer {
	cfg := &config.Config{}
	logg := &Log{}

	calendar := app.New(memorystorage.New(), cfg)
	return NewGRPCServer(logg, calendar, cfg)
}

func TestServer(t *testing.T) {
	t.Run("create test", func(t *testing.T) {
		s := prepareServer()

		tests := []struct {
			event  *Event
			result int32
			err    string
		}{
			{
				event: &Event{
					Id:            "e13b0eb3-4b87-41e6-bf29-e2e47b9dcdd8",
					Title:         "",
					DatetimeStart: "",
					DatetimeEnd:   "",
					Description:   "",
					UserId:        "e13b0eb3-4b87-41e6-bf29-e2e47b9dcdd8",
					WhenToNotify:  "",
				},
				result: 0,
				err:    `bad start date. parsing time "" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "2006"`,
			},
			{
				event: &Event{
					Id:            "4d9faa10-edf9-47a4-8d75-fa37bdc597c6",
					Title:         "",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "",
					Description:   "",
					UserId:        "4d9faa10-edf9-47a4-8d75-fa37bdc597c6",
					WhenToNotify:  "",
				},
				result: 0,
				err:    `bad end date. parsing time "" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "2006"`,
			},
			{
				event: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "872e211d-4f73-4564-816d-adcfd77a2450",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
				result: 1,
				err:    "",
			},
			{
				event: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "872e211d-4f73-4564-816d-adcfd77a2450",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
				result: 0,
				err:    "event already exist",
			},
			{
				event: &Event{
					Id:            "a43ec1d4-d805-4051-a0b4-79f36e9cf456",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "12",
					WhenToNotify:  "",
				},
				result: 0,
				err:    "bad userId. invalid UUID length: 2",
			},
			{
				event: &Event{
					Id:            "",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "9a66db8d-5714-4276-b860-852d888c95a9",
					WhenToNotify:  "",
				},
				result: 0,
				err:    "bad Id. invalid UUID length: 0",
			},
			{
				event: &Event{
					Id:            "9a66db8d-5714-4276-b860-852d888c95a9",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "9a66db8d-5714-4276-b860-852d888c95a9",
					WhenToNotify:  "",
				},
				result: 0,
				err:    "bad when date. parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\"",
			},
		}

		for _, test := range tests {
			resp, err := s.Create(context.Background(), test.event)
			if test.err == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}

			require.Equal(t, resp.Result, test.result)
		}
	})

	t.Run("update test", func(t *testing.T) {
		s := prepareServer()

		tests := []struct {
			event    *Event
			result   int32
			err      string
			oldEvent *Event
		}{
			{
				event: &Event{
					Id:            "e13b0eb3-4b87-41e6-bf29-e2e47b9dcdd8",
					Title:         "",
					DatetimeStart: "",
					DatetimeEnd:   "",
					Description:   "",
					UserId:        "e13b0eb3-4b87-41e6-bf29-e2e47b9dcdd8",
					WhenToNotify:  "",
				},
				result: 0,
				err:    `bad start date. parsing time "" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "2006"`,
			},
			{
				event: &Event{
					Id:            "4d9faa10-edf9-47a4-8d75-fa37bdc597c6",
					Title:         "",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "",
					Description:   "",
					UserId:        "4d9faa10-edf9-47a4-8d75-fa37bdc597c6",
					WhenToNotify:  "",
				},
				result: 0,
				err:    `bad end date. parsing time "" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "2006"`,
			},
			{
				event: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "872e211d-4f73-4564-816d-adcfd77a2450",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
				result: 0,
				err:    "event not found",
			},
			{
				event: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "872e211d-4f73-4564-816d-adcfd77a2450",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
				result: 1,
				err:    "",
				oldEvent: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test old",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "9a66db8d-5714-4276-b860-852d888c95a9",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
			},
			{
				event: &Event{
					Id:            "a43ec1d4-d805-4051-a0b4-79f36e9cf456",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "12",
					WhenToNotify:  "",
				},
				result: 0,
				err:    "bad userId. invalid UUID length: 2",
			},
			{
				event: &Event{
					Id:            "",
					Title:         "test",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "9a66db8d-5714-4276-b860-852d888c95a9",
					WhenToNotify:  "",
				},
				result: 0,
				err:    "bad Id. invalid UUID length: 0",
			},
		}

		for _, test := range tests {
			if test.oldEvent != nil {
				_, err := s.Create(context.Background(), test.oldEvent)
				require.Nil(t, err)
			}

			resp, err := s.Update(context.Background(), test.event)
			if test.err == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}

			require.Equal(t, resp.Result, test.result)
		}
	})

	t.Run("delete test", func(t *testing.T) {
		s := prepareServer()

		tests := []struct {
			id       string
			result   int32
			err      string
			oldEvent *Event
		}{
			{
				id:     "e13b0eb3-4b87-41e6-bf29-e2e47b9dcdd8",
				result: 0,
				err:    `event not found`,
			},
			{
				id:     "872e211d-4f73-4564-816d-adcfd77a2450",
				result: 1,
				err:    "",
				oldEvent: &Event{
					Id:            "872e211d-4f73-4564-816d-adcfd77a2450",
					Title:         "test old",
					DatetimeStart: "2010-05-12T10:10:20Z",
					DatetimeEnd:   "2010-05-12T10:10:20Z",
					Description:   "",
					UserId:        "9a66db8d-5714-4276-b860-852d888c95a9",
					WhenToNotify:  "2010-05-12T10:10:20Z",
				},
			},
		}

		for _, test := range tests {
			if test.oldEvent != nil {
				_, err := s.Create(context.Background(), test.oldEvent)
				require.Nil(t, err)
			}

			req := &DeleteRequest{
				Id: test.id,
			}

			resp, err := s.Delete(context.Background(), req)
			if test.err == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}

			require.Equal(t, resp.Result, test.result)
		}
	})
}
