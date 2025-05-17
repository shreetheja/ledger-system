package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvPostgresUser     = "POSTGRES_USER"
	EnvPostgresPassword = "POSTGRES_PASSWORD"
	EnvPostgresDB       = "POSTGRES_DB"
	EnvPostgresHost     = "POSTGRES_HOST"
	EnvPostgresPort     = "POSTGRES_PORT"

	EnvMongoHost = "MONGO_HOST"
	EnvMongoPort = "MONGO_PORT"
	EnvMongoDB   = "MONGO_DB"

	EnvKafkaBroker    = "KAFKA_BROKER"
	EnvKafkaTopic     = "KAFKA_TOPIC"
	EnvKafkaClusterID = "KAFKA_CLUSTER_ID"
	EnvPort           = "PORT"
)

// Global variables populated during init
var (
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string

	MongoHost string
	MongoPort string
	MongoDB   string

	KafkaBroker    string
	KafkaTopic     string
	KafkaClusterID string
	Port           string
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	PostgresUser = os.Getenv(EnvPostgresUser)
	PostgresPassword = os.Getenv(EnvPostgresPassword)
	PostgresDB = os.Getenv(EnvPostgresDB)
	PostgresHost = os.Getenv(EnvPostgresHost)
	PostgresPort = os.Getenv(EnvPostgresPort)

	MongoHost = os.Getenv(EnvMongoHost)
	MongoPort = os.Getenv(EnvMongoPort)
	MongoDB = os.Getenv(EnvMongoDB)

	KafkaBroker = os.Getenv(EnvKafkaBroker)
	KafkaTopic = os.Getenv(EnvKafkaTopic)
	KafkaClusterID = os.Getenv(EnvKafkaClusterID)
	Port = os.Getenv(EnvPort)
}
