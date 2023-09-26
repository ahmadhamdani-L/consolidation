package main

import (
	// "fmt"
	"log"
	db "notification/internal/database"
	"notification/internal/factory"
	"notification/internal/kafka"
	"notification/pkg/logger"
	"notification/pkg/util/env"
	"os"
	"strings"
	// "os"
	// "github.com/centrifugal/centrifuge-go"
)

func init() {
	if selectedEnv := strings.ToUpper(os.Getenv("ENV")); selectedEnv == "LOCAL" {
		env.NewEnv().Load(selectedEnv)
		logger.Log().Info().Msg("Choosen environment: " + selectedEnv)
	}
}

func main() {
	log.Println("Starting notification service...")
	db.Init()
	f := factory.NewFactory()
	kafka.NewConsumer(f)

}
