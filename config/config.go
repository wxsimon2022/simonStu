package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Mode string

	LogDir string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用环境变量或默认值")
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		Mode: getEnv("GIN_MODE", "debug"),

		LogDir: getEnv("LOG_DIR", "storage/logs"),

		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "test"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("警告: %s 不是有效的数字，使用默认值 %d", key, fallback)
	}
	return fallback
}
