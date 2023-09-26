package kafka

import (
	"bufio"
	"fmt"
	"mcash-finance-console-core/internal/model"
	"os"
	"strings"
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
		Request  string
	}
}

type NewToken struct {
	Token string
	UserID int
}

type JsonDataJurnal struct {
	TbID int
	DataJurnal []model.AdjustmentDetailEntity
}
type JsonDataImport struct {
	TrialBalance               	string
	AgingUtangPiutang          	string
	InvestasiTbk               	string
	InvestasiNonTbk            	string
	MutasiFA                   	string
	MutasiDta                  	string
	MutasiIa                   	string
	MutasiRua                  	string
	MutasiPersediaan           	string
	PembelianPenjualanBerelasi 	string
	EmployeeBenefit			   	string
	FNTrialBalance               string
	FNAgingUtangPiutang          string
	FNInvestasiTbk               string
	FNInvestasiNonTbk            string
	FNMutasiFA                   string
	FNMutasiDta                  string
	FNMutasiIa                   string
	FNMutasiRua                  string
	FNMutasiPersediaan           string
	FNPembelianPenjualanBerelasi string
	FNEmployeeBenefit			 string
	CompanyID                  int
	UserID                     int
	Version                    int
	ImportedWorkSheetID        int
	Period					   string
}

type JsonDataReUpload struct {
	TrialBalance                 string
	AgingUtangPiutang            string
	InvestasiTbk                 string
	InvestasiNonTbk              string
	MutasiFA                     string
	MutasiDta                    string
	MutasiIa                     string
	MutasiRua                    string
	MutasiPersediaan             string
	PembelianPenjualanBerelasi   string
	EmployeeBenefit			     string
	AllWorksheet				 string
	IDTrialBalance               int
	IDAgingUtangPiutang          int
	IDInvestasiTbk               int
	IDInvestasiNonTbk            int
	IDMutasiFA                   int
	IDMutasiDta                  int
	IDMutasiIa                   int
	IDMutasiRua                  int
	IDMutasiPersediaan           int
	IDPembelianPenjualanBerelasi int
	IDEmployeeBenefit			 int
	CompanyID                    int
	UserID                       int
	Version                      int
	ImportedWorkSheetID          int
	IDWorksheetDetailTrialBalance               int
	IDWorksheetDetailAgingUtangPiutang          int
	IDWorksheetDetailInvestasiTbk               int
	IDWorksheetDetailInvestasiNonTbk            int
	IDWorksheetDetailMutasiFA                   int
	IDWorksheetDetailMutasiDta   				int
	IDWorksheetDetailMutasiIa                   int
	IDWorksheetDetailMutasiRua                  int
	IDWorksheetDetailMutasiPersediaan           int
	IDWorksheetDetailPembelianPenjualanBerelasi int
	IDWorksheetDetailEmployeeBenefit            int
	FNTrialBalance               string
	FNAllWorksheet               string
}
func NewService(topic string) *service {
	kConfig := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset": os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}

	// configFile := "local-kafka.properties"
	// conf := ReadConfig(&kConfig)

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

func ReadConfig(kConfig string) kafka.ConfigMap {

	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(kConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m
}
