package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/seed"
	"github.com/DeepSyyy/Spatium-Backend/server"
)

func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	config.ConnectDB()

	log.Println("📦 Environment PORT:", os.Getenv("PORT"))
	log.Println("📦 Config APP_PORT:", config.AppConfig.AppPort)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed-user":
			fmt.Println("🌱 Running database seeder...")
			seed.SeedUsers(config.DB)
			fmt.Println("✅ Seeding complete.")
			return
		case "seed-mood-tag":
			fmt.Println("🌱 Running mood tag seeder...")
			seed.SeedMoodTags(config.DB)
			fmt.Println("✅ Seeding complete.")
			return
		case "seed-reaction-type":
			fmt.Println("🌱 Running reaction type seeder...")
			seed.SeedReactionTypes(config.DB)
			fmt.Println("✅ Seeding complete.")
		case "serve":
			fmt.Println("🚀 Starting server...")
			server.Start()
			return
		default:
			fmt.Println("❌ Unknown command. Use:")
			fmt.Println("   go run main.go serve  → start API server")
			fmt.Println("   go run main.go seed   → seed sample data")
			return
		}
	}

	fmt.Println("ℹ️ No command provided. Defaulting to 'serve' mode.")
	server.Start()
}
