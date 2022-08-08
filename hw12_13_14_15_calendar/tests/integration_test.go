//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type httpTestSuite struct {
	suite.Suite
	host  string
	port  string
	event *internalhttp.EventRequest
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config-test.yaml", "Path to configuration file")
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, &httpTestSuite{})
}

func (s *httpTestSuite) SetupSuite() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	s.host = cfg.Server.Host
	s.port = cfg.Server.Port
}

func (s *httpTestSuite) SetupTest() {
	s.event = &internalhttp.EventRequest{
		ID:     "14670ec6-dbca-425b-a4c7-d13c269af380",
		Title:  "test",
		Start:  "2009-10-31T23:59:59Z",
		End:    "2010-01-01T08:00:00Z",
		Desc:   "new year",
		UserID: "9591d712-1b3e-4495-bb71-08c906273a09",
		When:   "2009-12-31T23:59:59Z",
	}
}

func (s *httpTestSuite) TearDownTest() {
	s.event = nil
}

func (s *httpTestSuite) req(method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%s/%s", s.host, s.port, path), body)
	s.NoError(err)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	return response
}

func (s *httpTestSuite) getBody(response *http.Response) string {
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	return strings.Trim(string(byteBody), "\n")
}

func (s *httpTestSuite) generatePath(period string) string {
	date, err := time.Parse(time.RFC3339, s.event.Start)
	s.NoError(err)

	return fmt.Sprintf("events/%s/%s", period, date.Format("2006/01/02"))
}

func (s *httpTestSuite) getEventResponse() (string, error) {
	dateStart, err := time.Parse(time.RFC3339, s.event.Start)
	if err != nil {
		return "", err
	}

	dateEnd, err := time.Parse(time.RFC3339, s.event.End)
	if err != nil {
		return "", err
	}

	dateWhen, err := time.Parse(time.RFC3339, s.event.When)
	if err != nil {
		return "", err
	}

	eventResponse := storage.Event{
		ID:            uuid.MustParse(s.event.ID),
		Title:         s.event.Title,
		DatetimeStart: dateStart,
		DatetimeEnd:   dateEnd,
		Description:   s.event.Desc,
		UserID:        uuid.MustParse(s.event.UserID),
		WhenToNotify:  dateWhen,
	}

	eventResponseByte, err := json.Marshal(eventResponse)
	if err != nil {
		return "", err
	}

	return string(eventResponseByte), nil
}

func (s *httpTestSuite) TestPing() {
	response := s.req(http.MethodGet, "", nil)
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"Pong"}`, s.getBody(response))
}

func (s *httpTestSuite) TestCreateEmptyBody() {
	response := s.req(http.MethodPost, "events/create", nil)
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":400,"Message":"failed to unmarshal request body"}`, s.getBody(response))
}

func (s *httpTestSuite) TestCreateBadBody() {
	reqStr := "{}"
	response := s.req(http.MethodPost, "events/create", strings.NewReader(reqStr))
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":500,"Message":"internal server error"}`, s.getBody(response))
}

func (s *httpTestSuite) TestFull() {
	reqBytes, err := json.Marshal(s.event)
	s.NoError(err)

	// create
	response := s.req(http.MethodPost, "events/create", bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"event was created"}`, s.getBody(response))
	response.Body.Close()

	// double
	response = s.req(http.MethodPost, "events/create", bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":500,"Message":"internal server error"}`, s.getBody(response))
	response.Body.Close()

	// day
	pathDay := s.generatePath("day")
	response = s.req(http.MethodGet, pathDay, bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)

	eventResponse, err := s.getEventResponse()
	s.NoError(err)

	s.Contains(s.getBody(response), eventResponse)
	response.Body.Close()

	// week
	pathWeek := s.generatePath("week")
	response = s.req(http.MethodGet, pathWeek, bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)

	eventResponse, err = s.getEventResponse()
	s.NoError(err)

	s.Contains(s.getBody(response), eventResponse)
	response.Body.Close()

	// month
	pathMonth := s.generatePath("month")
	response = s.req(http.MethodGet, pathMonth, bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)

	eventResponse, err = s.getEventResponse()
	s.NoError(err)

	s.Contains(s.getBody(response), eventResponse)
	response.Body.Close()

	// update
	event := &internalhttp.EventRequest{
		ID:     "14670ec6-dbca-425b-a4c7-d13c269af380",
		Title:  "test 2",
		Start:  "2009-10-31T23:59:59Z",
		End:    "2010-01-01T08:00:00Z",
		Desc:   "new year",
		UserID: "9591d712-1b3e-4495-bb71-08c906273a09",
		When:   "2009-12-31T23:59:59Z",
	}

	reqBytes, err = json.Marshal(event)
	s.NoError(err)

	response = s.req(http.MethodPatch, fmt.Sprintf("events/%s", s.event.ID), bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(
		`{"Status":200,"Message":"event 14670ec6-dbca-425b-a4c7-d13c269af380 was updated"}`,
		s.getBody(response))

	response.Body.Close()

	// check updated
	pathMonth = s.generatePath("month")
	response = s.req(http.MethodGet, pathMonth, bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)

	s.Contains(s.getBody(response), `test 2`)
	response.Body.Close()

	// delete
	response = s.req(http.MethodDelete, fmt.Sprintf("events/%s", s.event.ID), nil)

	s.Equal(http.StatusOK, response.StatusCode)

	expected := `{"Status":200,"Message":"event 14670ec6-dbca-425b-a4c7-d13c269af380 was deleted"}`
	s.Equal(expected, s.getBody(response))
	response.Body.Close()

	// check empty
	pathMonth = s.generatePath("month")
	response = s.req(http.MethodGet, pathMonth, bytes.NewReader(reqBytes))

	s.Equal(http.StatusOK, response.StatusCode)

	s.Equal(`{}`, s.getBody(response))
	response.Body.Close()
}
