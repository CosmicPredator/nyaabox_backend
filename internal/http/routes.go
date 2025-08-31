package http

import (
	"cosmic/nyaabox/internal/config"
	"cosmic/nyaabox/internal/db"
	"cosmic/nyaabox/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Routes struct {
	app *fiber.App
	dbHandle *db.DbHandler
}

func (r *Routes) apiIndexHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(newHealthResponse("Ok"))
}

func (r *Routes) uploadHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	
	files := form.File["documents"]
	var fileNameArr []string = make([]string, 0)
	
	for _, file := range files {
		fileName, err := storage.SaveToDisk(ctx, *r.dbHandle, file)
		if err != nil { return err }
		fileNameArr = append(fileNameArr, fileName)
	}
	return ctx.Status(fiber.StatusOK).JSON(newUploadSuccessResponse(fileNameArr))
}

func (r *Routes) handleBadApiAttempt(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendFile("surprise.jpg")
}

func (r *Routes) RegisterRoutes() error {
	savePath, err := config.GetEnv("NB_SAVE_PATH")
	if err != nil {
		return err
	}

	r.app.Static("/file", savePath, fiber.Static{
		ByteRange: true,
		Compress: true,
	})
	
	r.app.Get("/api", r.apiIndexHandler)
	log.Debug("registered route: '/api'")
	
	r.app.Post("/api/uploadx", r.uploadHandler)
	log.Debug("registered route: '/api/upload'")
	
	r.app.Get("/api/*", r.handleBadApiAttempt)

	return nil
}

func NewRouteHandler(app *fiber.App, dbHandle *db.DbHandler) *Routes {
	return &Routes{
		app: app,
		dbHandle: dbHandle,
	}
}