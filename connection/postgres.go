package connection

import (
	"database/sql"
	"fmt"
	"insightful/common/config"

	_ "github.com/lib/pq"
	"log"
	"time"
)

func InitPostgres(conf *config.Config) *sql.DB {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", conf.Postgres.UserName, conf.Postgres.Password, conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.Database)
	//dsn := fmt.Sprintf("user=%v password=%v dbname=postgres sslmode=disable", conf.Postgres.UserName, conf.Postgres.Password, conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.Database)

	poolConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer poolConn.Close()
	poolConn.SetMaxOpenConns(conf.Postgres.MaxOpenConns)
	poolConn.SetMaxIdleConns(conf.Postgres.MaxIdleConns)
	poolConn.SetConnMaxLifetime(2 * time.Minute)

	return poolConn
}
