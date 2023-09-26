package kafkaconsumer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/app/consolidate"
	"worker-consol/internal/factory"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	topics   []string
	jsonData abstraction.JsonData
)

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

	topics = []string{"CONSOLIDATION"}

	err = c.SubscribeTopics(topics, nil)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run == true {
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
			data := abstraction.JsonData{}
			json.Unmarshal([]byte(ev.Value), &data)
			if strings.Trim(*ev.TopicPartition.Topic, " ") == "CONSOLIDATION" {
				// trialbalance.NewHandler(f).Action(string(ev.Key), data)
				consolidate.NewHandler(f).Action(string(ev.Key), data)
			}
		}
	}

	c.Close()
}
