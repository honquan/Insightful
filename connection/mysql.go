package connection

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"insightful/common/config"
	"time"
)

func InitMysql(conf *config.Config) (*gorm.DB, error) {
	mysqlUsername := conf.MySQL.UserName
	mysqlPassword := conf.MySQL.Password
	mysqlHost := conf.MySQL.Host
	mysqlPort := conf.MySQL.Port
	mysqlDBName := conf.MySQL.Database

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
	sqlDB.SetMaxIdleConns(conf.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MySQL.MaxOpenConns)

	return conn, err
}
