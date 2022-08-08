package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/sql"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func initStorage(cfg *config.Config, dsn string) (app.Storage, error) {
	switch cfg.DB.Type {
	case "mem":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(cfg, dsn), nil
	}

	return nil, fmt.Errorf("unknown database type: %q", cfg.DB.Type)
}

func buildDsn(cfg *config.Config) string {
	s := cfg.DB.SQL
	return "postgres://" + s.User + ":" + s.Password + "@" + s.Host + ":" + s.Port + "/" + s.Name + "?sslmode=disable"
}

func migrations(cfg *config.Config, dsn string) {
	if cfg.DB.Type != "sql" {
		return
	}

	db, err := sql.Open(cfg.DB.SQL.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations/", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}
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

	dsn := buildDsn(cfg)

	migrations(cfg, dsn)

	logg := logger.New(cfg.Logger.Level)

	storage, err := initStorage(cfg, dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = storage.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
	}

	calendar := app.New(logg, storage, cfg)

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

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
