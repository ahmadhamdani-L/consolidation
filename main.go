package main

import (
	"fmt"
	"os"
	db "worker/database"
	"worker/internal/factory"
	kafkaconsumer "worker/internal/kafka/consumer"
	"worker/pkg/util/env"

	"github.com/sirupsen/logrus"
)

func init() {
	ENV := os.Getenv("ENV")
	env := env.NewEnv()
	env.Load(ENV)

	logrus.Info("Choosen environment " + ENV)
}

// @title worker
// @version 0.0.1
// @description This is a doc for worker.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:3030
// @BasePath /
func main() {

	db.Init()

	f := factory.NewFactory()
	fmt.Println("OKE")
	kafkaconsumer.NewConsumer(f)
}
