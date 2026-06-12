// Package config 从 .env 或环境变量读取配置，提供默认值。
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用全部配置项。
type Config struct {
	Port      string // HTTP 监听端口，默认 ":8080"
	Mode      string // Gin 运行模式：debug / release / test
	TZ        string // 时区，默认 Asia/Shanghai
	LogDir    string // 日志输出目录，默认 storage/logs
	JWTSecret string // JWT 签名密钥

	RedisHost     string // Redis 地址
	RedisPort     string // Redis 端口
	RedisPassword string // Redis 密码
	RedisDB       int    // Redis 数据库编号

	DBHost     string // MySQL 地址
	DBPort     string // MySQL 端口
	DBUser     string // MySQL 用户名
	DBPassword string // MySQL 密码
	DBName     string // MySQL 数据库名
}

// Load 读取配置并返回 Config 实例。.env 不存在时不报错，使用系统环境变量或默认值。
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用环境变量或默认值")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		Mode:      getEnv("GIN_MODE", "debug"),
		TZ:        getEnv("TZ", "Asia/Shanghai"),
		LogDir:    getEnv("LOG_DIR", "storage/logs"),
		JWTSecret: getEnv("JWT_SECRET", "simon-stu-secret-key"),

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

// getEnv 读取环境变量，不存在时返回 fallback。
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt 读取环境变量并转为 int，转换失败时打印警告并返回 fallback。
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("警告: %s 不是有效的数字，使用默认值 %d", key, fallback)
	}
	return fallback
}
