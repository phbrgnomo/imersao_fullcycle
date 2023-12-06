package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

// Consumer represents a Kafka consumer configuration
type Consumer struct {
	ConfigMap *ckafka.ConfigMap 
	Topics		[]string
}

//NewConsumer creates a new instance of the Consumer with the provided configuration
func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: 		configMap,
		Topics:        	topics,
	}
}

// Consume continuously reads messages from Kafka and sends them to the provided msgChan channel
func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	// Create a new Kafka consumer with the specified configuration
	consumer, err := ckafka.NewConsumer(c.ConfigMap)
	if err!= nil {
		panic(err)
	}

	// Subscribe the consumet tot he specified Kafka topics
	err = consumer.SubscribeTopics(c.Topics, nil)
	if err!= nil {
		panic(err)
	}

	// Continuously read messages from Kafka and send them to the msgChan channel
	for {
		msg, err := consumer.ReadMessage(-1)
		if err!= nil {
			panic(err)
		}
		msgChan <- msg
	}
}