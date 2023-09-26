package kafkaproducer

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type service struct {
	kafka *kafka.Producer
	topic string
}

type JsonData struct {
	FileLoc   string
	CompanyID int
	UserID    int
	Timestamp *time.Time
	Data      string
	Name      string
	Filter    struct {
		Period   string
		Versions int
	}
}

func NewProducer(topic string) *service {
	kConfig := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}

	p, err := kafka.NewProducer(&kConfig)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
	return &service{
		kafka: p,
		topic: topic,
	}
}

func (s *service) SendMessage(key, data string) {
	go func() {
		for e := range s.kafka.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	s.kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(data),
	}, nil)

	// Wait for all messages to be delivered
	s.kafka.Flush(5 * 1000)
	s.kafka.Close()
}
