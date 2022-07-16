package internalhttp

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
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
}

type Application interface{}

func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	return &Server{
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		logger: logger,
	}
}

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, "Hello!\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.host, s.port)

	r := mux.NewRouter()
	r.HandleFunc("/", s.HelloHandler)
	r.Use(s.loggingMiddleware)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.logger.Info("server starting on http://" + addr)

	go func() {
		err := server.ListenAndServe()
		s.logger.Info("server started on " + addr)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
