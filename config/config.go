package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	AppConfig *Config
)

type Config struct {
	AppPort           string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	JWTSecret         string
	JWTExpiredMinutes string
	JWTResfreshToken  string
	JWTExpire         string
	OpenAIAPIKey      string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	AppConfig = &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASS", "password"),
		DBName:            getEnv("DB_NAME", "mydb"),
		JWTSecret:         getEnv("JWT_SECRET", "your_jwt_secret_key"),
		JWTExpiredMinutes: getEnv("JWT_EXPIRY_MINUTES", "60"),
		JWTResfreshToken:  getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
	}
}

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)
	if exist {
		return value
	}
	return fallback
}

func ConnectDB() {
	cfg := AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}
