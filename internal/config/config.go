package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     string
	MongoURI string
	S3       S3Config
	SNS      SNSConfig
	SQS      SQSConfig
}

type S3Config struct {
	ContractBucket                string
	Region                        string
	PresignedURLExpirationMinutes int
}

type SNSConfig struct {
	TopicArn string
}

type SQSConfig struct {
	QueueURL string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", "mongodb://admin:password@localhost:27017"),
		S3: S3Config{
			ContractBucket:                getEnv("AWS_S3_BUCKET", "distributed-jobs-platform-contracts"),
			Region:                        getEnv("AWS_REGION", "us-east-1"),
			PresignedURLExpirationMinutes: getEnvInt("AWS_S3_PRESIGNED_EXPIRATION_MINUTES", 15),
		},
		SNS: SNSConfig{
			TopicArn: getEnv("SNS_TOPIC_ARN", ""),
		},
		SQS: SQSConfig{
			QueueURL: getEnv("SQS_QUEUE_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
