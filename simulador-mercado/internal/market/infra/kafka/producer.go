package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

// Producer representes a Kafka message producer
type Producer struct{
	ConfigMap *ckafka.ConfigMap
}

// NewKafkaProducer creates a new instanca of the Kafka producer.
func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

// Publish sends a message to the specified Kafka topic with the given key
func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	// Create a new Kafka producer instance
	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err!= nil {
		return err
	}

	// Create a Kafka message with the provided data
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Key:             key,
		Value:           msg.([]byte),
	}

	// Produce the message to the Kafka topic
	err = producer.Produce(message, nil)
	if err!= nil {
		return err
	}
	return nil
}
