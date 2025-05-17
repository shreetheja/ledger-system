package kafka

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	Producer *kafka.Producer
	Consumer *kafka.Consumer
)

// InitKafka sets up the Kafka producer and consumer
func InitKafka(broker, groupID string) error {
	var err error

	// Initialize Producer
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"client.id":         "ledger-producer",
	})
	if err != nil {
		return err
	}

	// Initialize Consumer
	Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		Producer.Close()
		return err
	}
	log.Println("✅ Connected to Kafka:")
	return nil
}

// CloseKafka closes both the producer and consumer connections
func CloseKafka() {
	if Producer != nil {
		Producer.Close()
	}
	if Consumer != nil {
		Consumer.Close()
	}
}

func CreateTopics(broker string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return err
	}
	defer adminClient.Close()

	// Prepare topic specifications
	var topicSpecs []kafka.TopicSpecification
	for _, topic := range Topics {
		topicSpecs = append(topicSpecs, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := adminClient.CreateTopics(ctx, topicSpecs, kafka.SetAdminOperationTimeout(30*time.Second))
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error.Code() == kafka.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists, skipping creation", result.Topic)
		} else if result.Error.Code() != kafka.ErrNoError {
			return result.Error
		} else {
			log.Printf("Topic %s created successfully", result.Topic)
		}
	}

	return nil
}
