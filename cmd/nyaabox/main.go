package main

import (
	"context"
	"cosmic/nyaabox/internal/config"
	"cosmic/nyaabox/internal/cron"
	"cosmic/nyaabox/internal/db"
	"cosmic/nyaabox/internal/http"
	"cosmic/nyaabox/internal/storage"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var version = "devbuild"

func checkErr(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
		BodyLimit: 2 * 1024 * 1024 * 1024, // 2GB
		ErrorHandler: func (ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusNotFound).JSON(
				http.NewErrorResponse(err.Error()),
			)
		},
	})
	
	log.Infof("Nyaabox backend version %s", version)
	
	err := config.LoadEnv()
	checkErr(err)
	
	_, err = storage.PrepareSaveDir()
	checkErr(err)
	
	dbHandle, err := db.NewDbHandler()
	checkErr(err)
	defer dbHandle.Close()
	
	err = dbHandle.Init()
	checkErr(err)
	
	interval, err := config.GetEnv("NB_CRON_INTERVAL")
	if err != nil {
		interval = "30"
	}
	intervalInt, _ := strconv.Atoi(interval)
	
	cron := cron.NewCronJob(context.Background(), time.Duration(intervalInt)*time.Second, dbHandle)
	go cron.RunCronJob()
	
	routeManager := http.NewRouteHandler(app, dbHandle)
	routeManager.RegisterRoutes()

	port, err := config.GetEnv("NB_PORT")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	
	err = app.Listen(":" + port)
	
	if err != nil {
		log.Error(err.Error())
	}
}