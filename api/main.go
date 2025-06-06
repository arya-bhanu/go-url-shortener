package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/arya-bhanu/go-url-shortener/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading env: %+v", err)
	}
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}

func setupRoutes(app *fiber.App) {
	app.Post("/api/v1", routes.ShortenUrlHandler)
}
