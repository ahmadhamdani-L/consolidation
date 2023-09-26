package kafkaconsumer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"worker/internal/abstraction"
	"worker/internal/app/imports"
	"worker/internal/app/reupload"
	"worker/internal/factory"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var topics []string = []string{
	"IMPORT",
	"NOTIFICATION",
	"REUPLOAD",
}

func NewConsumer(f *factory.Factory) {

	kConfig := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}

	c, err := kafka.NewConsumer(&kConfig)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n", *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			dataImport := abstraction.JsonDataImport{}
			json.Unmarshal([]byte(ev.Value), &dataImport)
			dataReUpload := abstraction.JsonDataReUpload{}
			json.Unmarshal([]byte(ev.Value), &dataReUpload)

			if strings.Trim(*ev.TopicPartition.Topic, " ") == "IMPORT" {
				imports.NewHandler(f).ActionImport(string(ev.Key), dataImport)
			}
			// if strings.Trim(*ev.TopicPartition.Topic, " ") == "NOTIFICATION" {
			// 	notification.Action(string(ev.Key), dataImport)
			// }
			if strings.Trim(*ev.TopicPartition.Topic, " ") == "REUPLOAD" {
				reupload.NewHandler(f).ActionImport(string(ev.Key), dataReUpload)
			}
		}
	}

	c.Close()
}
