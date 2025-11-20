package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	HTTPPort string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error; environment may already be set

	cfg := &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "yard_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		RedisAddr:     getEnv("REDIS_ADDR", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}

	// parse REDIS_DB
	// default 0
	cfg.RedisDB = 0
	if s := getEnv("REDIS_DB", "0"); s != "" {
		fmt.Sscanf(s, "%d", &cfg.RedisDB)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
