package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	ElasticUrl string

	RabbitUrl      string
	RabbitHost     string
	RabbitPort     string
	RabbitUser     string
	RabbitPass     string
	QueueCreated   string
	QueueCancelled string
}

func LoadConfig() *Config {
	// Load .env only in local development
	_ = godotenv.Load()

	cfg := &Config{
		Port: os.Getenv("PORT"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),

		RabbitUrl:      os.Getenv("RABBITMQ_URL"),
		RabbitHost:     os.Getenv("RABBITMQ_HOST"),
		RabbitPort:     os.Getenv("RABBITMQ_PORT"),
		RabbitUser:     os.Getenv("RABBITMQ_USER"),
		RabbitPass:     os.Getenv("RABBITMQ_PASS"),
		QueueCreated:   os.Getenv("ORDER_CREATED_QUEUE"),
		QueueCancelled: os.Getenv("ORDER_CANCELLED_QUEUE"),
		ElasticUrl:     os.Getenv("ELASTICSEARCH_URL"),
	}

	if cfg.Port == "" {
		log.Fatal("❌ APP_PORT must be set in env")
	}

	return cfg
}
