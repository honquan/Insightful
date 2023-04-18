package connection

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"insightful/src/apis/conf"
	"time"
)

func InitMysql() (*gorm.DB, error) {
	mysqlUsername := conf.EnvConfig.DBMysqlUsername
	mysqlPassword := conf.EnvConfig.DBMysqlPassword
	mysqlHost := conf.EnvConfig.DBMysqlHost
	mysqlPort := conf.EnvConfig.DBMysqlPort
	mysqlDBName := conf.EnvConfig.DBMysqlName

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}

	sqlDB, err := conn.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// set max idle and max open conns
	sqlDB.SetMaxIdleConns(conf.EnvConfig.DBMysqlMaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.EnvConfig.DBMysqlMaxOpenConns)

	return conn, err
}
