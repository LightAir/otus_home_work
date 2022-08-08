package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	rmqqueue "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/queue/rmq"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/scheduler"
	initstorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/init"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := initstorage.NewSchedulerStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	err = storage.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	logg := logger.New(cfg.Logger.Level)

	rmq := rmqqueue.NewRmq(cfg)
	sch := scheduler.NewScheduler(rmq, *logg, storage, cfg.QueueName)

	go func() {
		if err := sch.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Logger started...")
	<-ctx.Done()
}
