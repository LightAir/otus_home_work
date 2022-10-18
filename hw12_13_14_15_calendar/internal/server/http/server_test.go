package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type Log struct{}

func (l Log) Info(...interface{})                                    {}
func (l Log) Error(...interface{})                                   {}
func (l Log) LogHTTPRequest(_ *http.Request, _ time.Duration, _ int) {}

type PeriodCases struct {
	period   string
	response string
	date     string
	eq       bool
}

func prepareServer() *Server {
	cfg := &config.Config{}
	logg := &Log{}

	calendar := app.New(memorystorage.New(), cfg)
	return NewServer(logg, calendar, cfg)
}

func checkResponse(t *testing.T, response *http.Response, message string) {
	t.Helper()
	date, err := ioutil.ReadAll(response.Body)
	require.Nil(t, err)

	require.Equal(t, message, string(date))
}

func getEventResponse(t *testing.T, event *EventRequest) (string, error) {
	t.Helper()

	dateStart, err := time.Parse(time.RFC3339, event.Start)
	if err != nil {
		return "", err
	}

	dateEnd, err := time.Parse(time.RFC3339, event.End)
	if err != nil {
		return "", err
	}

	dateWhen, err := time.Parse(time.RFC3339, event.When)
	if err != nil {
		return "", err
	}

	eventResponse := storage.Event{
		ID:            uuid.MustParse(event.ID),
		Title:         event.Title,
		DatetimeStart: dateStart,
		DatetimeEnd:   dateEnd,
		Description:   event.Desc,
		UserID:        uuid.MustParse(event.UserID),
		WhenToNotify:  dateWhen,
	}

	eventResponseByte, err := json.Marshal(eventResponse)
	if err != nil {
		return "", err
	}

	return string(eventResponseByte), nil
}

func TestServer(t *testing.T) {
	event := &EventRequest{
		ID:     "14670ec6-dbca-425b-a4c7-d13c269af380",
		Title:  "test",
		Start:  "2009-12-31T23:59:59Z",
		End:    "2010-01-01T08:00:00Z",
		Desc:   "new year",
		UserID: "9591d712-1b3e-4495-bb71-08c906273a09",
		When:   "2009-12-31T23:59:59Z",
	}

	t.Run("createEventHandler bad json", func(t *testing.T) {
		s := prepareServer()

		serv := httptest.NewServer(http.HandlerFunc(s.createEventHandler))
		defer serv.Close()

		resp, err := http.Post(serv.URL, "application/json", strings.NewReader("{ badjson"))
		require.Nil(t, err)
		defer resp.Body.Close()

		checkResponse(t, resp, `{"Status":400,"Message":"failed to unmarshal request body"}`)
	})

	t.Run("createEventHandler event created", func(t *testing.T) {
		s := prepareServer()

		serv := httptest.NewServer(http.HandlerFunc(s.createEventHandler))
		defer serv.Close()

		res, err := json.Marshal(event)
		require.Nil(t, err)

		resp, err := http.Post(serv.URL, "application/json", bytes.NewReader(res))
		require.Nil(t, err)
		defer resp.Body.Close()

		checkResponse(t, resp, `{"Status":200,"Message":"event was created"}`)

		respTwo, err := http.Post(serv.URL, "application/json", bytes.NewReader(res))
		require.Nil(t, err)

		defer respTwo.Body.Close()
		checkResponse(t, respTwo, `{"Status":500,"Message":"internal server error"}`)
	})

	t.Run("updateEventByGUID bad json", func(t *testing.T) {
		s := prepareServer()

		serv := httptest.NewServer(http.HandlerFunc(s.updateEventByGUID))
		defer serv.Close()

		resp, err := http.Post(serv.URL, "application/json", strings.NewReader("{ badjson"))
		require.Nil(t, err)
		defer resp.Body.Close()

		checkResponse(t, resp, `{"Status":400,"Message":"failed to unmarshal request body"}`)
	})

	t.Run("updateEventByGUID event updated", func(t *testing.T) {
		s := prepareServer()

		router := mux.NewRouter()

		router.HandleFunc("/create", s.createEventHandler)
		router.HandleFunc("/update", s.updateEventByGUID)

		serv := httptest.NewServer(router)
		defer serv.Close()

		res, err := json.Marshal(event)
		require.Nil(t, err)

		resp, err := http.Post(serv.URL+"/create", "application/json", bytes.NewReader(res))
		require.Nil(t, err)
		defer resp.Body.Close()

		respUpdate, err := http.Post(serv.URL+"/update", "application/json", bytes.NewReader(res))
		require.Nil(t, err)
		defer respUpdate.Body.Close()

		checkResponse(t, respUpdate, `{"Status":200,"Message":"event 14670ec6-dbca-425b-a4c7-d13c269af380 was updated"}`)

		resBadDate, err := json.Marshal(EventRequest{
			ID:     "14670ec6-dbca-425b-a4c7-d13c269af380",
			Title:  "test",
			Start:  "2009-12-31",
			End:    "2010-01-01T08:00:00Z",
			Desc:   "new year",
			UserID: "9591d712-1b3e-4495-bb71-08c906273a09",
			When:   "never",
		})
		require.Nil(t, err)

		respBadDateUpdate, err := http.Post(serv.URL+"/update", "application/json", bytes.NewReader(resBadDate))
		require.Nil(t, err)
		defer respBadDateUpdate.Body.Close()

		checkResponse(t, respBadDateUpdate, `{"Status":500,"Message":"internal server error"}`)
	})

	t.Run("test deleteEventByGUID", func(t *testing.T) {
		s := prepareServer()

		router := mux.NewRouter()

		router.HandleFunc("/create", s.createEventHandler)
		router.HandleFunc("/delete/{eventID}", s.deleteEventByGUID)

		serv := httptest.NewServer(router)
		defer serv.Close()

		res, err := json.Marshal(event)
		require.Nil(t, err)

		resp, err := http.Post(serv.URL+"/create", "application/json", bytes.NewReader(res))
		require.Nil(t, err)
		defer resp.Body.Close()

		respDelete, err := http.Get(serv.URL + "/delete/26109d4b-1d69-4e32-a189-7ccab6c4230b")
		require.Nil(t, err)

		checkResponse(t, respDelete, `{"Status":500,"Message":"internal server error"}`)

		deleted, err := http.Get(serv.URL + "/delete/" + event.ID)
		require.Nil(t, err)

		checkResponse(t, deleted, `{"Status":200,"Message":"event `+event.ID+` was deleted"}`)
	})

	t.Run("test eventsByPeriod", func(t *testing.T) {
		s := prepareServer()

		router := mux.NewRouter()

		router.HandleFunc("/create", s.createEventHandler)
		router.HandleFunc("/{period}/{YYYY}/{MM}/{DD}", s.getEventsByPeriod)

		serv := httptest.NewServer(router)
		defer serv.Close()

		res, err := json.Marshal(event)
		require.Nil(t, err)

		resp, err := http.Post(serv.URL+"/create", "application/json", bytes.NewReader(res))
		require.Nil(t, err)
		resp.Body.Close()

		eventResponseStr, err := getEventResponse(t, event)
		require.Nil(t, err)

		cases := []*PeriodCases{
			{
				period:   "day",
				date:     "2009/12/31",
				response: eventResponseStr,
			},
			{
				period:   "week",
				date:     "2009/12/31",
				response: eventResponseStr,
			},
			{
				period:   "month",
				date:     "2009/12/31",
				response: eventResponseStr,
			},
			{
				period:   "bad",
				date:     "2009/12/31",
				eq:       true,
				response: `{"Status":400,"Message":"the period does not exist"}`,
			},
			{
				period:   "bad date",
				date:     "2009/12/331",
				eq:       true,
				response: `{"Status":400,"Message":"date parse error"}`,
			},
		}

		for _, oneCase := range cases {
			resp, err := http.Get(fmt.Sprintf("%s/%s/%s", serv.URL, oneCase.period, oneCase.date))
			require.Nil(t, err)

			defer resp.Body.Close()

			require.Nil(t, err)
			date, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			if oneCase.eq {
				require.Equal(t, string(date), oneCase.response)
			} else {
				require.Contains(t, string(date), oneCase.response)
			}
		}
	})
}
