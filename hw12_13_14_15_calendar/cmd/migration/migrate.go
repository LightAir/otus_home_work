package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if cfg.DB.Type != "sql" {
		log.Fatal("no sql database type selected")
	}

	s := cfg.DB.SQL

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", s.User, s.Password, s.Host, s.Port, s.Name)

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
		log.Println(err)
		os.Exit(0)
	}

	log.Println("successfully migrations")
}
