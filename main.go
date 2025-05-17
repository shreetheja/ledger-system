package main

import (
	"context"
	"fmt"
	"ledger/api"
	"ledger/config"
	"ledger/kafka"
	"ledger/mongo"
	"ledger/pg"
	"ledger/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.Initialize()

	// Initialize the PostgreSQL connection
	pg.InitPostgres()
	defer pg.DB.Close()

	// Initialize the Redis connection
	mongo.InitMongo()
	defer mongo.MongoClient.Disconnect(context.Background())

	// Initialize the Kafka producer
	kafka.InitKafka(config.KafkaBroker, config.KafkaClusterID)
	kafka.CreateTopics(config.KafkaBroker)
	defer kafka.Producer.Close()

	// setup...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on Ctrl+C
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutdown signal received")
		cancel()
		if kafka.Consumer != nil {
			kafka.Consumer.Close()
		}
	}()

	service.Initialize(ctx)

	err := service.CreateAccount("12", 10)
	log.Printf("err: %v\n", err)

	log.Print("Listening on : ", config.Port)
	http.DefaultClient.Timeout = time.Second * 10
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), api.InitialiseRoutes()))
}
