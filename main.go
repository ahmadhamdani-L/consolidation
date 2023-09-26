package main

import (
	db "mcash-finance-console-core/internal/database"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/http"
	"mcash-finance-console-core/internal/middleware"
	"mcash-finance-console-core/pkg/logger"
	"mcash-finance-console-core/pkg/util/env"
	"os"
	"strings"

	"mcash-finance-console-core/pkg/redis"

	"github.com/labstack/echo/v4"
)

func init() {
	if selectedEnv := strings.ToUpper(os.Getenv("ENV")); selectedEnv == "LOCAL" {
		env.NewEnv().Load(selectedEnv)
		logger.Log().Info().Msg("Choosen environment: " + selectedEnv)
	}
}

// @title MCASH-CONSOLE
// @version 0.0.1
// @description This is a doc for mcash-console.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:3030
// @schemes http
// @BasePath /
func main() {
	db.Init()
	// migration.Init()
	// elasticsearch.Init()

	e := echo.New()
	middleware.Init(e)

	f := factory.NewFactory()
	http.Init(e, f)
	redis.Init()
	e.Logger.Fatal(e.Start(":3030"))
}
