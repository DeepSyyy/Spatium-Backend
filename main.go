package main

import (
	"fmt"
	"os"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/seed"
	"github.com/DeepSyyy/Spatium-Backend/server"
)

func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	config.ConnectDB()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed":
			fmt.Println("🌱 Running database seeder...")
			seed.SeedUsers(config.DB)
			seed.SeedMoodTags(config.DB)
			fmt.Println("✅ Seeding complete.")
			return
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
