package configs

import (
	"os"
	"sync"
)

type JwtConfig struct {
	secretKey string
}

var (
	jwt     *JwtConfig
	jwtOnce sync.Once
)

func (jc *JwtConfig) SecretKey() string {
	return jc.secretKey
}

func Jwt() *JwtConfig {
	jwtOnce.Do(func() {
		jwt = &JwtConfig{
			secretKey: PriorityString(fang.GetString("jwt.key"), os.Getenv("JWT_KEY")),
		}
	})
	return jwt
}
