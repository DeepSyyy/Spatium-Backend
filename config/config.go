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

// LoadEnv loads environment variables safely for both local and production
func LoadEnv() {
	// Try to load .env (only works locally)
	if err := godotenv.Load(); err == nil {
		log.Println("✅ Loaded .env file for local environment")
	} else {
		log.Println("🌐 Using environment variables from system (Railway/Production)")
	}

	// Allow app to boot even if some vars missing in Railway
	dbHost := getEnv("DB_HOST", "")
	dbUser := getEnv("DB_USER", "")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	dbPort := getEnv("DB_PORT", "5432")

	if dbHost == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Println("⚠️  Some database environment variables are missing — skipping DB connect check.")
	}

	// Parse JWT durations (fallback defaults)
	jwtExpire := parseDurationSafe("JWT_EXPIRED", "1h")
	jwtRefreshToken := parseDurationSafe("REFRESH_TOKEN_EXPIRED", "24h")

	AppConfig = &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DBHost:           dbHost,
		DBPort:           dbPort,
		DBUser:           dbUser,
		DBPass:           dbPass,
		DBName:           dbName,
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTResfreshToken: jwtRefreshToken.String(),
		JWTExpire:        jwtExpire.String(),
		OpenAIAPIKey:     getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}

func parseDurationSafe(key, fallback string) time.Duration {
	val := getEnv(key, fallback)
	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("⚠️  Invalid duration for %s: %s, using default %s", key, val, fallback)
		dur, _ = time.ParseDuration(fallback)
	}
	return dur
}

// ConnectDB tries to connect to the DB but won't crash the app if it fails
func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Println("⚠️  Continuing without DB connection (won’t crash Railway).")
		return
	}

	DB = db
	log.Println("✅ Database connected successfully.")
}
