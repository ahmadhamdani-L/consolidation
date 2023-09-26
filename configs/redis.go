package configs

import (
	"os"
	"strconv"
	"sync"
)

type RedisConfig struct {
	host     string
	port     string
	pass     string
	db       int
	poolSize int
}

var (
	rdb     *RedisConfig
	rdbOnce sync.Once
)

func (rdc *RedisConfig) Host() string {
	return rdc.host
}

func (rdc *RedisConfig) Port() string {
	return rdc.port
}

func (rdc *RedisConfig) Password() string {
	return rdc.pass
}

func (rdc *RedisConfig) Db() int {
	return rdc.db
}

func (rdc *RedisConfig) PoolSize() int {
	return rdc.poolSize
}

func Redis() *RedisConfig {
	rdbOnce.Do(func() {
		strEnv := PriorityString(fang.GetString("redis.db"), os.Getenv("REDIS_DB"), "0")
		intDbIdx, err := strconv.Atoi(strEnv)
		if err != nil {
			intDbIdx = 0
		}

		strEnv = PriorityString(fang.GetString("redis.pool_size"), os.Getenv("REDIS_POOL_SIZE"), "10")
		intPoolSize, err := strconv.Atoi(strEnv)
		if err != nil {
			intPoolSize = 10
		}

		rdb = &RedisConfig{
			host:     PriorityString(fang.GetString("redis.host"), os.Getenv("REDIS_HOST"), "localhost"),
			port:     PriorityString(fang.GetString("redis.port"), os.Getenv("REDIS_PORT"), "6379"),
			db:       intDbIdx,
			pass:     PriorityString(fang.GetString("redis.password"), os.Getenv("REDIS_PASSWORD"), ""),
			poolSize: intPoolSize,
		}
	})
	return rdb
}
