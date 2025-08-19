package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()
	
	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("'/' called")
		return c.JSON(map[string]int {
			"status": 200,
		})
	})
	
	err := app.Listen(":8787")
	if err != nil {
		log.Error(err.Error())
	}
}