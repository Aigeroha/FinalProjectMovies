package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBConnStr  string
	JWTSecret  string
	RedisAddr  string
}

var AppConfig Config

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: .env файл не найден, используются системные переменные окружения")
	}

	dbConnStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	AppConfig = Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBConnStr:  dbConnStr,
		JWTSecret:  os.Getenv("JWT_SECRET"),
		RedisAddr:  os.Getenv("REDIS_ADDR"),
	}

	if AppConfig.ServerPort == "" {
		AppConfig.ServerPort = ":8080"
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("Критическая ошибка: Переменная JWT_SECRET не задана в файле .env")
	}

	if os.Getenv("REDIS_ADDR") == "" {
		log.Fatal("Критическая ошибка: Переменная REDIS_ADDR не задана в файле .env")
	}
}
