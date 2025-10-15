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
	// Only load .env file if it exists (for local development)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  Failed to load .env file:", err)
		} else {
			log.Println("✅ Loaded .env file for local environment")
		}
	} else {
		log.Println("🌐 Using environment variables from system (Railway/Production)")
	}

	// --- Validate critical variables ---
	dbHost := getEnv("DB_HOST", "")
	if dbHost == "" {
		log.Fatalf("❌ DB_HOST environment variable is not set. Cannot start the application.")
	}

	dbUser := getEnv("DB_USER", "")
	if dbUser == "" {
		log.Fatalf("❌ DB_USER environment variable is not set. Cannot start the application.")
	}

	dbPass := getEnv("DB_PASSWORD", "")
	if dbPass == "" {
		log.Fatalf("❌ DB_PASSWORD environment variable is not set. Cannot start the application.")
	}

	dbName := getEnv("DB_NAME", "")
	if dbName == "" {
		log.Fatalf("❌ DB_NAME environment variable is not set. Cannot start the application.")
	}

	// --- Parse time durations ---
	jwtExpire, err := time.ParseDuration(getEnv("JWT_EXPIRED", "1h"))
	if err != nil {
		log.Fatalf("❌ Failed to parse JWT_EXPIRED duration: %v", err)
	}

	jwtRefreshToken, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRED", "24h"))
	if err != nil {
		log.Fatalf("❌ Failed to parse REFRESH_TOKEN_EXPIRED duration: %v", err)
	}

	// --- Populate AppConfig ---
	AppConfig = &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DBHost:           dbHost,
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           dbUser,
		DBPass:           dbPass,
		DBName:           dbName,
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTResfreshToken: jwtRefreshToken.String(),
		JWTExpire:        jwtExpire.String(),
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
