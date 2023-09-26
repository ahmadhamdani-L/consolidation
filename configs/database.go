package configs

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DbConfig struct {
	host         string
	port         string
	name         string
	user         string
	pass         string
	sslmode      string
	tz           string
	maxOpenConns int
	maxIdleConns int
	connLifetime time.Duration
}

var (
	db     *DbConfig
	dbOnce sync.Once
)

func (dbc *DbConfig) Host() string {
	return dbc.host
}

func (dbc *DbConfig) Port() string {
	return dbc.port
}

func (dbc *DbConfig) Name() string {
	return dbc.name
}

func (dbc *DbConfig) User() string {
	return dbc.user
}

func (dbc *DbConfig) Password() string {
	return dbc.pass
}

func (dbc *DbConfig) SslMode() string {
	return dbc.sslmode
}

func (dbc *DbConfig) Timezone() string {
	return dbc.tz
}

func (dbc *DbConfig) MaxOpenConnections() int {
	return dbc.maxOpenConns
}

func (dbc *DbConfig) MaxIdleConnections() int {
	return dbc.maxIdleConns
}

func (dbc *DbConfig) ConnectionLifetime() time.Duration {
	return dbc.connLifetime
}

func DB() *DbConfig {
	dbOnce.Do(func() {
		listSSLMode := map[string]bool{
			"disable":     true,
			"allow":       true,
			"prefer":      true,
			"require":     true,
			"verify-ca":   true,
			"verify-full": true,
		}
		givenSSLMode := strings.ToLower(strings.TrimSpace(PriorityString(fang.GetString("db.sslmode"), os.Getenv("DB_SSLMODE"), "disable")))
		if _, ok := listSSLMode[givenSSLMode]; !ok {
			givenSSLMode = "disable"
		}

		strMaxOpenConns := PriorityString(fang.GetString("db.max_open_connections"), os.Getenv("DB_MAX_OPEN_CONNECTIONS"), "20")
		mOpenConns, err := strconv.Atoi(strMaxOpenConns)
		if err != nil {
			mOpenConns = 0
		}

		strMaxIdleConns := PriorityString(fang.GetString("db.max_idle_connections"), os.Getenv("DB_MAX_IDLE_CONNECTIONS"), "18")
		mIdleConns, err := strconv.Atoi(strMaxIdleConns)
		if err != nil {
			mIdleConns = 18
		}

		strConnLifetime := PriorityString(fang.GetString("db.connection_lifetime"), os.Getenv("DB_CONNECTION_LIFETIME"))
		ltConn, err := time.ParseDuration(strConnLifetime)
		if err != nil {
			ltConn = 0 * time.Second
		}

		db = &DbConfig{
			host:         PriorityString(fang.GetString("db.host"), os.Getenv("DB_HOST"), "localhost"),
			port:         PriorityString(fang.GetString("db.port"), os.Getenv("DB_PORT"), "5432"),
			name:         PriorityString(fang.GetString("db.name"), os.Getenv("DB_NAME"), "db_name"),
			user:         PriorityString(fang.GetString("db.user"), os.Getenv("DB_USER"), "postgres"),
			pass:         PriorityString(fang.GetString("db.pass"), os.Getenv("DB_PASS"), "postgres"),
			sslmode:      givenSSLMode,
			tz:           PriorityString(fang.GetString("db.tz"), os.Getenv("DB_TZ"), "UTC"),
			maxOpenConns: mOpenConns,
			maxIdleConns: mIdleConns,
			connLifetime: ltConn,
		}
	})

	return db
}
