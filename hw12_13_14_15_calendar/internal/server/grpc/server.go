//go:generate protoc -I ../../../api/ --go_out=. --go-grpc_out=. ../../../api/EventService.proto

package internalgrpc

import (
	"context"
	"net"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Period int

const (
	day Period = iota
	week
	month
)

type Logger interface {
	Error(msg ...interface{})
	Info(msg ...interface{})
	LogGRPCRequest(ctx context.Context, info *grpc.UnaryServerInfo, duration time.Duration, statusCode string)
}

type GRPCServer struct {
	host   string
	port   string
	logger Logger
	app    Application
	server *grpc.Server
}

type Application interface {
	CreateEvent(id, title, start, end, desc, userID, when string) error
	UpdateEvent(id, title, start, end, desc, userID, when string) error
	DeleteEvent(id string) error
	FindEventsByPeriod(start, end time.Time) (map[uuid.UUID]storage.Event, error)
}

func (srv *GRPCServer) Create(_ context.Context, event *Event) (*EventResponse, error) {
	err := srv.app.CreateEvent(
		event.Id,
		event.Title,
		event.DatetimeStart,
		event.DatetimeEnd,
		event.Description,
		event.UserId,
		event.WhenToNotify,
	)
	if err != nil {
		return &EventResponse{
			Result: 0,
		}, err
	}

	return &EventResponse{
		Result: 1,
	}, nil
}

func (srv *GRPCServer) Update(_ context.Context, event *Event) (*EventResponse, error) {
	err := srv.app.UpdateEvent(
		event.Id,
		event.Title,
		event.DatetimeStart,
		event.DatetimeEnd,
		event.Description,
		event.UserId,
		event.WhenToNotify,
	)
	if err != nil {
		return &EventResponse{
			Result: 0,
		}, err
	}

	return &EventResponse{
		Result: 1,
	}, nil
}

func (srv *GRPCServer) Delete(_ context.Context, request *DeleteRequest) (*EventResponse, error) {
	err := srv.app.DeleteEvent(request.Id)
	if err != nil {
		return &EventResponse{
			Result: 0,
		}, err
	}

	return &EventResponse{
		Result: 1,
	}, nil
}

func (srv *GRPCServer) EventListOfDay(_ context.Context, request *DateRequest) (*EventsResponse, error) {
	return srv.getEventListOfPeriod(request, day)
}

func (srv *GRPCServer) EventListOfWeek(_ context.Context, request *DateRequest) (*EventsResponse, error) {
	return srv.getEventListOfPeriod(request, week)
}

func (srv *GRPCServer) EventListOfMonth(_ context.Context, request *DateRequest) (*EventsResponse, error) {
	return srv.getEventListOfPeriod(request, month)
}

func (srv *GRPCServer) getEventListOfPeriod(request *DateRequest, period Period) (*EventsResponse, error) {
	dateStart, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, err
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
		return &EventsResponse{
			Events: nil,
		}, err
	}

	result, err := srv.app.FindEventsByPeriod(dateStart, dateEnd)
	if err != nil {
		return &EventsResponse{
			Events: nil,
		}, err
	}

	events := make([]*Event, 0)

	for _, item := range result {
		event := &Event{
			Id:            item.ID.String(),
			Title:         item.Title,
			DatetimeStart: item.DatetimeStart.Format(time.RFC3339),
			DatetimeEnd:   item.DatetimeEnd.Format(time.RFC3339),
			Description:   item.Description,
			UserId:        item.UserID.String(),
			WhenToNotify:  item.WhenToNotify.Format(time.RFC3339),
		}
		events = append(events, event)
	}

	return &EventsResponse{
		Events: events,
	}, nil
}

func (srv *GRPCServer) mustEmbedUnimplementedEventServiceServer() {
	srv.logger.Error("unimplemented event")
}

func NewGRPCServer(logger Logger, app Application, cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		host:   cfg.GRPCServer.Host,
		port:   cfg.GRPCServer.Port,
		logger: logger,
		app:    app,
	}
}

func (srv *GRPCServer) Start(ctx context.Context) error {
	s := grpc.NewServer(grpc.UnaryInterceptor(srv.RequestLogInterceptor))

	srv.server = s

	RegisterEventServiceServer(s, srv)

	addr := net.JoinHostPort(srv.host, srv.port)

	go func() {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Error(err)

			return
		}

		err = s.Serve(l)
		if err != nil {
			srv.logger.Error(err)
		}
	}()

	srv.logger.Info("GRPC server starting on tcp://" + addr)

	<-ctx.Done()
	return nil
}

func (srv *GRPCServer) Stop() {
	srv.server.GracefulStop()
}
