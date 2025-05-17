package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendCreateAccountMessage(msg CreateAccountMessage) error {
	msg.Timestamp = time.Now()
	return sendMessage(TopicCreateAccount, msg.UserID, msg)
}

func SendAddBalanceMessage(msg AddBalanceMessage) error {
	msg.Timestamp = time.Now()
	return sendMessage(TopicAddBalance, msg.UserID, msg)
}

func SendDeductBalanceMessage(msg DeductBalanceMessage) error {
	msg.Timestamp = time.Now()
	return sendMessage(TopicDeductBalance, msg.UserID, msg)
}

func sendMessage(topic, key string, msg interface{}) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	log.Printf("Message delivered to %v\n", m.TopicPartition)
	return nil
}
