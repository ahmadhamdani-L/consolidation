package main

import (
	"os"
	"strings"
	db "worker-consol/internal/database"
	"worker-consol/internal/factory"
	kafkaconsumer "worker-consol/internal/kafka/consumer"
	"worker-consol/pkg/logger"
	"worker-consol/pkg/util/env"
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
