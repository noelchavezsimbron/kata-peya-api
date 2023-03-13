package integration

import (
	"context"
	"database/sql"

	"kata-peya/config"
	"kata-peya/internal/storage/mysql"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	createTableSQl = `
		CREATE  TABLE IF NOT EXISTS pets(
		    id int primary key ,
		    name varchar(50) NOT NULL ,
		    vaccines varchar(100) ,
		    age_months int NOT NULL 
		);
		`
)

var (
	conf = config.Get()
)

func New() (*sql.DB, testcontainers.Container) {

	ctx := context.Background()
	mysqlC, err := startMySQLContainer(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	var (
		host, _ = mysqlC.Host(ctx)
		p, _    = mysqlC.MappedPort(ctx, "3306/tcp")
		port    = p.Int()
		user    = "root"
	)

	db := mysql.New(mysql.Config{
		Host:     host,
		Port:     port,
		Database: conf.DB.Database,
		User:     user,
		Password: conf.DB.Password,
	})

	if _, err := db.Exec(createTableSQl); err != nil {
		log.Fatal(err)

	}

	return db, mysqlC
}

func startMySQLContainer(ctx context.Context) (testcontainers.Container, error) {

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8",
			ExposedPorts: []string{"3306/tcp", "33060/tcp"},
			Env: map[string]string{
				"MYSQL_DATABASE":      conf.DB.Database,
				"MYSQL_ROOT_PASSWORD": conf.DB.Password,
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("port: 3306  MySQL Community Server - GPL"),
				wait.ForListeningPort("3306/tcp"),
			),
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}
