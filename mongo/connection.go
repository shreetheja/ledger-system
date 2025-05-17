package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"ledger/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database
var LedgerCollection *mongo.Collection

// RecordTransaction inserts one or multiple ledger records in a MongoDB transaction.
func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := fmt.Sprintf("mongodb://%s:%s", config.MongoHost, config.MongoPort)
	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Ping to confirm connection
	if err := MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	MongoDB = MongoClient.Database(config.MongoDB)
	LedgerCollection = MongoDB.Collection("ledger_records")
	log.Println("✅ Connected to MongoDB:", config.MongoDB)
}
