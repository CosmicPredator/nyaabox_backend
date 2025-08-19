package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Routes struct {
	app *fiber.App
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
		log.Debugf("Processing: %s | %s | %s", file.Filename, file.Size, file.Header["Content-Type"][0])
		err := ctx.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
		if err != nil {
			return err
		}
		fileNameArr = append(fileNameArr, fmt.Sprintf("%s/%s", ctx.BaseURL(), file.Filename))
	}
	return ctx.Status(fiber.StatusOK).JSON(newUploadSuccessResponse(fileNameArr))
}

func (r *Routes) RegisterRoutes() {
	r.app.Get("/api", r.apiIndexHandler)
	log.Debug("registered route: '/api'")
	
	r.app.Post("/api/upload", r.uploadHandler)
	log.Debug("registered route: '/api/upload'")
}

func NewRouteHandler(app *fiber.App) *Routes {
	return &Routes{
		app: app,
	}
}