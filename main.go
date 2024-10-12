package main

import (
	"golang-auth/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://go-auth-frontend.onrender.com/", // Frontend origin
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",                        // Allowed methods
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",                   // Allow Authorization header
		AllowCredentials: true,                                                            // If using credentials like cookies or authorization
	}))

	// Initialize database and store
	store := db.NewStore()

	// Define routes from routes.go
	SetupRoutes(app, store)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
