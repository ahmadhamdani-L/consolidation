package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"notification/internal/factory"
	"notification/internal/model"
	"notification/internal/service"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewConsumer(f *factory.Factory) {
	conf := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}
	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topics := []string{"NOTIFICATION"}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Println("Failed to subscribe to topics: ", err)
		os.Exit(1)
	}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Process messages
	service := service.NewService(f)
	log.Println("Kafka consumer started")
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Println("err: ", err)
				break
			}
			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n", *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			data := model.JsonData{}
			json.Unmarshal([]byte(ev.Value), &data)
			if strings.Trim(*ev.TopicPartition.Topic, " ") == "NOTIFICATION" {
				// service.SendExportNotif(data)
				service.SendNotif(data)
			}
		
		}
	}

	c.Close()
}
