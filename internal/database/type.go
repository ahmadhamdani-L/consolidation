package database

import (
	"fmt"
	"mcash-finance-console-core/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db interface {
	Init() (*gorm.DB, error)
}

type db struct {
	Host string
	User string
	Pass string
	Port string
	Name string
}

type dbPostgreSQL struct {
	db
	SslMode string
	Tz      string
}

func (c *dbPostgreSQL) Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", c.Host, c.User, c.Pass, c.Name, c.Port, c.SslMode, c.Tz)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxOpenConns(configs.DB().MaxOpenConnections())
	sqlDb.SetMaxIdleConns(configs.DB().MaxIdleConnections())
	sqlDb.SetConnMaxLifetime(configs.DB().ConnectionLifetime())

	return db, nil
}
