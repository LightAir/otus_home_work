package initstorage

import (
	"fmt"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/scheduler"
	memorystorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

func getDsn(s config.SQLDatabase) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", s.User, s.Password, s.Host, s.Port, s.Name)
}

func NewStorage(cfg *config.Config) (app.Storage, error) {
	switch cfg.DB.Type {
	case "mem":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(cfg, getDsn(cfg.DB.SQL)), nil
	}

	return nil, fmt.Errorf("unknown database type: %q", cfg.DB.Type)
}

func NewSchedulerStorage(cfg *config.Config) (scheduler.Storage, error) {
	switch cfg.DB.Type {
	case "mem":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(cfg, getDsn(cfg.DB.SQL)), nil
	}

	return nil, fmt.Errorf("unknown database type: %q", cfg.DB.Type)
}
