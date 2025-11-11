package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	MongoURI   string
	MongoDB    string
	CORSOrigin string
}

func MustLoadConfig() Config {
	_ = godotenv.Load(".env")

	cfg := Config{
		Port:       getEnv("PORT"),
		MongoURI:   getEnv("MONGODB_URI"),
		MongoDB:    getEnv("MONGODB_DB"),
		CORSOrigin: getEnv("CORS_ORIGIN"),
	}
	return cfg
}

func getEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("environment variable %s not set", k)
	}
	return v
}
