package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func New(c Config) *sql.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=60s&readTimeout=60s", c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := otelsql.Open("mysql", dns,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithDBName(c.Database))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("successfully MySQL Connection")
	return db
}

func _New(c Config) *sql.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=60s&readTimeout=60s", c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("successfully MySQL Connection")
	return db
}
