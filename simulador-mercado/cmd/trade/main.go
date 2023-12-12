package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/dto"
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/entity"
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/infra/kafka"
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/transformer"
)

func main() {
	// Channels for communication between different parts of the program
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// Kafka message channel
	kafkaMsgChan := make(chan *ckafka.Message)
	// Kafka configuration
	configMap := &ckafka.ConfigMap{
		// the file /etc/hosts must have `127.0.0.1 kubernetes.docker.internal host.docker.internal`
		// To access kafka from ouside the containers, use the port 9094
		// To access kafka from other containers on the same network, use the port 9092

		// "bootstrap.servers": 	"host.docker.internal:9094",
		"bootstrap.servers": "kafka:9092"
		"group.id":				"myGroup",
		"auto.offset.reset":	"latest",
	}

	// Kafka producer
	producer := kafka.NewKafkaProducer(configMap)
	// Kafka consumer
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	// Concurrenty consume Kafka messages
	go kafka.Consume(kafkaMsgChan) // Thread 2

	// Create a book for trading
	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // Thread 3

	// Process Kafka messages and transform them into orders (Thread 4)
	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err!= nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	// Process orders and transform the into output, then publish to Kafka (Thread 5)
	for res := range ordersOut{
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "   ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}