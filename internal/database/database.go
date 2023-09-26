package database

import (
	"errors"
	"fmt"
	"strings"
	cfg "worker/configs"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	dbConnections map[string]*gorm.DB
)

func Init() {
	dbConfigurations := map[string]Db{
		"SAMPLE1": &dbPostgreSQL{
			db: db{
				Host: cfg.DB().Host(),
				User: cfg.DB().User(),
				Pass: cfg.DB().Password(),
				Port: cfg.DB().Port(),
				Name: cfg.DB().Name(),
			},
			SslMode: cfg.DB().SslMode(),
			Tz:      cfg.DB().Timezone(),
		},
	}

	dbConnections = make(map[string]*gorm.DB)
	for k, v := range dbConfigurations {
		db, err := v.Init()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to database %s", k))
		}
		dbConnections[k] = db
		logrus.Info(fmt.Sprintf("Successfully connected to database %s", k))
	}
}

func Connection(name string) (*gorm.DB, error) {
	if dbConnections[strings.ToUpper(name)] == nil {
		return nil, errors.New("Connection is undefined")
	}
	return dbConnections[strings.ToUpper(name)], nil
}
