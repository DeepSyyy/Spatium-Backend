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
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	JWTSecret        string
	JWTResfreshToken string
	JWTExpire        string
	OpenAIAPIKey     string
}

func LoadEnv() {
	// Hanya load .env kalau file-nya memang ada (untuk lokal)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  Failed to load .env file:", err)
		} else {
			log.Println("✅ Loaded .env file for local environment")
		}
	} else {
		log.Println("🌐 Using environment variables from system (Railway/Production)")
	}

	AppConfig = &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DBHost:           getEnv("DB_HOST", ""),
		DBPort:           getEnv("DB_PORT", ""),
		DBUser:           getEnv("DB_USER", ""),
		DBPass:           getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTResfreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
		JWTExpire:        getEnv("JWT_EXPIRED", "1h"),
		OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
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
