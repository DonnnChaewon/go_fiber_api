package main

import "github.com/gofiber/fiber/v2"

func main() {
	db := InitializeDB()

	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	AuthHandlers(app.Group("/auth"), db)

	Download(app.Group("/download"), db)

	protected := app.Use(AuthMiddleware(db))

	MovieHandlers(protected.Group("/movie"), db)

	app.Listen(":3000")
}