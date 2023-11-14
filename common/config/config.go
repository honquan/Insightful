package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Config struct {
	MySQL    *MySQL
	Redis    *Redis
	Mongo    *Mongo
	Postgres *Postgres
	Worker   *Worker
}

type MySQL struct {
	Host         string `envconfig:"MYSQL_HOST"`
	Port         int64  `envconfig:"MYSQL_PORT"`
	UserName     string `envconfig:"MYSQL_USER_NAME"`
	Password     string `envconfig:"MYSQL_PASSWORD"`
	Database     string `envconfig:"MYSQL_DATABASE"`
	MaxIdleConns int    `envconfig:"MYSQL_MAX_IDLE_CONNS"`
	MaxOpenConns int    `envconfig:"MYSQL_MAX_OPEN_CONNS"`
}

type Redis struct {
	Address string `envconfig:"REDIS_ADDRESS" default:"127.0.0.1:6379"`

	RedisHost     string `envconfig:"REDIS_HOST" default:"0.0.0.0"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisDatabase int    `envconfig:"REDIS_DATABASE" default:"0"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
}

type Mongo struct {
	MongoURI          string `envconfig:"MONGO_DB_URI" default:"mongodb://0.0.0.0:27017"`
	MongoDatabaseName string `envconfig:"MONGO_DB_NAME" default:"insightful"`
	MongoDebug        bool   `envconfig:"MONGO_DB_DEBUG" default:"false"`
}

type Worker struct {
	MaxWorker int `envconfig:"MAX_WORKER" default:"3000"`
	MaxQueue  int `envconfig:"MAX_QUEUE" default:"3000"`
}

type Postgres struct {
	Host         string `envconfig:"POSTGRES_HOST" required:"true" default:"127.0.0.1"`
	Port         int64  `envconfig:"POSTGRES_PORT" required:"true" default:"5432"`
	UserName     string `envconfig:"POSTGRES_USER_NAME" required:"true" default:"postgres"`
	Password     string `envconfig:"POSTGRES_PASSWORD" default:""`
	Database     string `envconfig:"POSTGRES_DATABASE" required:"true" default:"insightful"`
	MaxIdleConns int    `envconfig:"POSTGRES_MAX_IDLE_CONNS" required:"true" default:"2"`
	MaxOpenConns int    `envconfig:"POSTGRES_MAX_OPEN_CONNS" required:"true" default:"4"`
}

func NewConfig() *Config {
	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot get env")
		}
	}

	var cnf Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Printf("Error when process get config %v", err)
		panic("Cannot envconfig Process")
	}

	return &cnf
}
