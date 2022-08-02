package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Logger interface {
	Error(msg ...interface{})
	Info(msg ...interface{})
	LogHTTPRequest(r *http.Request, d time.Duration, statusCode int)
}

type Server struct {
	host   string
	port   string
	logger Logger
	app    Application
	server *http.Server
}

type Application interface {
	CreateEvent(id, title, start, end, desc, userID, when string) error
	UpdateEvent(id, title, start, end, desc, userID, when string) error
	DeleteEvent(id string) error
	FindEventsByPeriod(start, end time.Time) (map[uuid.UUID]storage.Event, error)
}

type EventRequest struct {
	ID     string
	Title  string
	Start  string
	End    string
	Desc   string
	UserID string
	When   string
}

type TypicalResponse struct {
	Status  int
	Message string
}

const (
	day   = "day"
	week  = "week"
	month = "month"
)

func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	return &Server{
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		logger: logger,
		app:    app,
	}
}

func (s *Server) pingHandler(w http.ResponseWriter, _ *http.Request) {
	s.message(http.StatusOK, "Pong", w)
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to read request body", w)

		return
	}

	data := &EventRequest{}

	err = json.Unmarshal(body, data)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to unmarshal request body", w)

		return
	}

	err = s.app.CreateEvent(data.ID, data.Title, data.Start, data.End, data.Desc, data.UserID, data.When)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)

		return
	}

	s.message(http.StatusOK, "event was created", w)
}

func (s *Server) updateEventByGUID(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to read request body", w)

		return
	}

	data := &EventRequest{}

	err = json.Unmarshal(body, data)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to unmarshal request body", w)

		return
	}

	err = s.app.UpdateEvent(data.ID, data.Title, data.Start, data.End, data.Desc, data.UserID, data.When)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)

		return
	}

	s.message(http.StatusOK, fmt.Sprintf("event %s was updated", data.ID), w)
}

func (s *Server) deleteEventByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventID"]

	err := s.app.DeleteEvent(eventID)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)

		return
	}

	s.message(http.StatusOK, fmt.Sprintf("event %s was deleted", eventID), w)
}

func (s *Server) getEventsByPeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	period := vars["period"]
	date := fmt.Sprintf("%s-%s-%s", vars["YYYY"], vars["MM"], vars["DD"])

	dateStart, err := time.Parse("2006-01-02", date)
	if err != nil {
		s.message(http.StatusBadRequest, "date parse error", w)

		return
	}

	var dateEnd time.Time

	switch period {
	case day:
		dateEnd = dateStart.AddDate(0, 0, 1)
	case week:
		dateEnd = dateStart.AddDate(0, 0, 7)
	case month:
		dateEnd = dateStart.AddDate(0, 1, 0)
	default:
		s.message(http.StatusBadRequest, "the period does not exist", w)

		return
	}

	result, err := s.app.FindEventsByPeriod(dateStart, dateEnd)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)
		return
	}

	res, err := json.Marshal(result)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)
	}
}

func (s *Server) message(status int, message string, w http.ResponseWriter) {
	res, err := json.Marshal(TypicalResponse{
		Status:  status,
		Message: message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}

	w.WriteHeader(status)

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.host, s.port)

	r := mux.NewRouter()
	r.HandleFunc("/", s.pingHandler)
	r.HandleFunc("/events/{period}/{YYYY}/{MM}/{DD}", s.getEventsByPeriod).Methods("GET")
	r.HandleFunc("/events/create", s.createEventHandler).Methods("POST")
	r.HandleFunc("/events/{eventID}", s.updateEventByGUID).Methods("PATCH")
	r.HandleFunc("/events/{eventID}", s.deleteEventByGUID).Methods("DELETE")

	r.Use(s.loggingMiddleware)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.server = server

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	s.logger.Info("server starting on http://" + addr)

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
