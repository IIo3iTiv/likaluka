package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config структура, обозначающая структуру .env файла
type Config struct {
	Mode string

	MinioEndpoint           string
	MinioBucketName         string
	MinioU                  string
	MinioP                  string
	MinioUseSSL             bool
	MinioFileTimeExpiration int

	PgHost string
	PgPort string
	PgDB   string
	PgU    string
	PgP    string
	PgMOC  string
	PgMCLT string
	PgMILT string
	PgSSL  string

	ServerPort string
}

var conf *Config

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig() {
	// Загружаем переменные окружения из файла .env
	if getEnv("MODE", "") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	// Устанавливаем конфигурационные параметры
	conf = &Config{
		Mode: getEnv("MODE", "pensil"),

		// Minio
		MinioEndpoint:           getEnv("MINIO_ENDPOINT", "pensil"),
		MinioBucketName:         getEnv("MINIO_BUCKET_NAME", "pensil"),
		MinioU:                  getEnv("MINIO_ROOT_USER", "pensil"),
		MinioP:                  getEnv("MINIO_ROOT_PASSWORD", "pensil"),
		MinioUseSSL:             getEnvAsBool("MINIO_USE_SSL", false),
		MinioFileTimeExpiration: getEnvAsInt("FILE_TIME_EXPIRATION", 0),

		// Postgre
		PgHost: getEnv("PG_HOST", "pensil"),
		PgPort: getEnv("PG_PORT", "pensil"),
		PgDB:   getEnv("PG_DB", "pensil"),
		PgU:    getEnv("PG_U", "pensil"),
		PgP:    getEnv("PG_P", "pensil"),
		PgMOC:  getEnv("PG_MOC", "pensil"),
		PgMCLT: getEnv("PG_MCLT", "pensil"),
		PgMILT: getEnv("PG_MILT", "pensil"),
		PgSSL:  getEnv("PG_SSL", "pensil"),

		// Server
		ServerPort: getEnv("PORT", "pensil"),
	}
}

// getEnv считывает значение переменной окружения или возвращает значение по умолчанию, если переменная не установлена
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt считывает значение переменной окружения как целое число или возвращает значение по умолчанию, если переменная не установлена или не может быть преобразована в целое число
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// getEnvAsBool считывает значение переменной окружения как булево или возвращает значение по умолчанию, если переменная не установлена или не может быть преобразована в булево
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
