package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	PoolSize int
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database,
	)
}

func New(config Config) (*sql.DB, error) {
	c, err := otelsql.Open("postgres", config.DSN(),
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(config.Database),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create postgres sql.DB driver: %w", err)
	}

	c.SetMaxOpenConns(config.PoolSize)
	err = c.Ping()
	if err != nil {
		panic(err)
	}

	return c, nil
}
