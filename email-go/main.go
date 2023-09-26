package main

import (
	"mail-go/kafka/config"
	"mail-go/kafka/consumer"
	"mail-go/kafka/controller"
	"mail-go/kafka/producer"
	"mail-go/kafka/routes"
	"os"
	"github.com/labstack/echo/v4"
)

var userController *controller.UserController



func main() {
	var PORT = os.Getenv("3031")
	
	e := echo.New()

	config.CORSConfig(e)

	routes.GetUserApiRoutes(e, userController)

	e.Logger.Fatal(e.Start(":" + PORT))


}

func init() {
	p := config.InitKafkaProducer()
	producer := producer.NewProducer(p)
	userController = controller.NewUserController(producer)
	c := config.InitKafkaConsumer(config.UserConsumerGroup)
	consumer := consumer.NewConsumer(c)
	go consumer.Consume([]string{config.UserNotificationTopic})
}
