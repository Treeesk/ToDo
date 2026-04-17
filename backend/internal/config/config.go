package config

// Пакет инициализирующий переменные окружения(Подключение к БД и к backend)

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("config: error loading .env file, use base vars")
	}
	return &Config{
		BaseURL:    loadEnv("BASE_URL", "localhost:8080"),
		DBHost:     loadEnv("DB_HOST", "localhost"),
		DBPort:     loadEnv("DB_PORT", "5432"),
		DBUser:     loadEnv("DB_USER", "youruser"),
		DBPassword: loadEnv("DB_PASSWORD", "yourpassword"),
		DBName:     loadEnv("DB_NAME", "yourdbname"),
	}
}

func loadEnv(key, defaultval string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultval
}
