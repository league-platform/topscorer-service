package main

import (
    "github.com/gofiber/fiber/v2"
    "topscorer-service/handlers"
)

func main() {
    app := fiber.New()

    app.Post("/scorers", handlers.CreateScorer)
    app.Get("/scorers", handlers.GetScorers)

    app.Listen(":3000")
}
