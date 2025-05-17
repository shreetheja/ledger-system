package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerHandler struct {
	MessageChannel   chan *kafka.Message
	SubscribedTopics []string
}

var globalSubscribedChannels map[string][]chan *kafka.Message = make(map[string][]chan *kafka.Message)

// NewConsumerHandler creates a ConsumerHandler with initialized channels for topics
func NewConsumerHandler(topics []string) *ConsumerHandler {
	// Initialize the global subscribed channels map with buffer.
	var newConsumerChannel = make(chan *kafka.Message, 100)
	for _, topic := range topics {
		if _, ok := globalSubscribedChannels[topic]; !ok {
			globalSubscribedChannels[topic] = make([]chan *kafka.Message, 0)
		}
		globalSubscribedChannels[topic] = append(globalSubscribedChannels[topic], newConsumerChannel)
	}
	return &ConsumerHandler{
		MessageChannel:   newConsumerChannel,
		SubscribedTopics: topics,
	}
}

func (h *ConsumerHandler) StartConsuming(ctx context.Context) {
	if Consumer == nil {
		log.Fatal("Kafka consumer is not initialized. Call InitKafka first.")
	}

	if err := Consumer.SubscribeTopics(Topics, nil); err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}

	log.Println("Kafka consumer started. Waiting for messages...")

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Kafka consumer context cancelled, stopping message loop")
				return
			default:
				msg, err := Consumer.ReadMessage(-1)
				if err != nil {
					log.Printf("Consumer error: %v", err)
					continue
				}

				log.Printf("Received message: %s", msg.Value)
				var subscribedChannels = globalSubscribedChannels[*msg.TopicPartition.Topic]
				for _, channel := range subscribedChannels {
					select {
					case channel <- msg:
						log.Printf("Message sent to channel for topic %s", *msg.TopicPartition.Topic)
					default:
						log.Printf("Channel for topic %s is full, dropping message", *msg.TopicPartition.Topic)
					}
				}
			}
		}
	}(ctx)
}
