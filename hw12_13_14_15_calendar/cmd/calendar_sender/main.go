package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/consumer/sender"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	rmqqueue "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/queue/rmq"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config-sender.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rmq := rmqqueue.NewRmq(cfg)
	logg := logger.New(cfg.Logger.Level)

	sdr := sender.NewSender(rmq, *logg, cfg.QueueName)
	go func() {
		if err := sdr.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Sender started...")
	<-ctx.Done()
}
