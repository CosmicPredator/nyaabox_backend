package main

import (
	"cosmic/nyaabox/internal/config"
	"cosmic/nyaabox/internal/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var version = "devbuild"

func main() {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
	})
	log.Infof("Nyaabox backend version %s", version)
	
	err := config.LoadEnv()
	if err != nil {
		os.Exit(-1)
	}
	
	routeManager := http.NewRouteHandler(app)
	routeManager.RegisterRoutes()
	
	port, err := config.GetEnv("NYAABOX_PORT")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	
	err = app.Listen(":" + port)
	
	if err != nil {
		log.Error(err.Error())
	}
}