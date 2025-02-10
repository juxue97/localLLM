package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	PublicHost           string
	Port                 string
	DBUser               string
	DBPassword           string
	DBAddress            string
	DBName               string
	JWTSecret            string
	JWTExpirationSeconds int64
	LLMModel             string
	LLMIp                string
	MongoUri             string
	MongoDatabase        string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		PublicHost:           getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                 getEnv("PORT", ":8080"),
		DBUser:               getEnv("DBUser", "root"),
		DBPassword:           getEnv("DBPassword", "hungchuno0o"),
		DBAddress:            fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:               getEnv("DBName", "chatbot"),
		JWTSecret:            getEnv("JWTSecret", "onaonaina"),
		JWTExpirationSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		LLMModel:             getEnv("LLMModel", "qwen2.5-coder:32b"),
		LLMIp:                getEnv("LLMIp", "http://192.168.50.164:11434"),
		// MongoUri:             getEnv("MongoUri", "mongodb://root:rootpass@192.168.50.164:27017/?authSource=admin"),
		MongoUri:      getEnv("MongoUri", "mongodb://root:rootpass@localllm.hopto.org:27017/?authSource=admin"),
		MongoDatabase: getEnv("MongoDB", "HR"),
		// LLMIp:         getEnv("LLMIp", "http://localllm.hopto.org:11434"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
