package config

import "os"

type Config struct {
	PrimaryDBURL string
	ReplicaDBURL string
	RabbitMQURL  string
	Port         string
}

func LoadConfig() Config {
	return Config{
		PrimaryDBURL: getEnv("PRIMARY_DB_URL", "postgres://postgres:postgres@localhost:5433/order_db?sslmode=disable"),
		ReplicaDBURL: getEnv("REPLICA_DB_URL", "postgres://postgres:postgres@localhost:5434/order_db?sslmode=disable"),
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
