package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/server/http"
	initstorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/init"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := initstorage.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = storage.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	logg := logger.New(cfg.Logger.Level)
	calendar := app.New(storage, cfg)
	server := internalhttp.NewServer(logg, calendar, cfg)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	GRPCServer := internalgrpc.NewGRPCServer(logg, calendar, cfg)

	go func() {
		if err := GRPCServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()
}
