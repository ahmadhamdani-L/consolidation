package main

import (
	"os"
	"strings"
	db "worker/internal/database"
	"worker/internal/factory"
	kafkaconsumer "worker/internal/kafka/consumer"
	"worker/pkg/logger"
	"worker/pkg/util/env"
)

func init() {
	if selectedEnv := strings.ToUpper(os.Getenv("ENV")); selectedEnv == "LOCAL" {
		env.NewEnv().Load(selectedEnv)
		logger.Log().Info().Msg("Choosen environment: " + selectedEnv)
	}
}

// @title codeid-boiler
// @version 0.0.1
// @description This is a doc for codeid-boiler.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:3030
// @BasePath /
func main() {

	db.Init()

	f := factory.NewFactory()
	kafkaconsumer.NewConsumer(f)
}
