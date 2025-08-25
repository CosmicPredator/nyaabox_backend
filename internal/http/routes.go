package http

import (
	"cosmic/nyaabox/internal/db"
	"cosmic/nyaabox/internal/storage"
	"mime"
	"path/filepath"

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

func (r *Routes) getFileHandler(ctx *fiber.Ctx) error {
	fileId := ctx.Params("file_id")
	file, err := storage.GetFile(fileId, r.dbHandle)
	ext := filepath.Ext(file.FileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if err != nil {
		return err
	}
	ctx.Set("Content-Type", mimeType)
	ctx.Set("Access-Control-Allow-Methods", "GET, HEAD")
	ctx.Set("Cache-Control", "public, max-age=31536000")
	ctx.Set("Access-Control-Allow-Origin", "*")
	
	return ctx.Status(fiber.StatusOK).SendFile(file.FilePath, true)
}

func (r *Routes) handleBadApiAttempt(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendFile("surprise.jpg")
}

func (r *Routes) RegisterRoutes() {
	r.app.Get("/file/:file_id", r.getFileHandler)
	log.Debug("registered route: '/:file_id'")
	
	r.app.Get("/api", r.apiIndexHandler)
	log.Debug("registered route: '/api'")
	
	r.app.Post("/api/uploadx", r.uploadHandler)
	log.Debug("registered route: '/api/upload'")
	
	r.app.Get("/api/*", r.handleBadApiAttempt)
}

func NewRouteHandler(app *fiber.App, dbHandle *db.DbHandler) *Routes {
	return &Routes{
		app: app,
		dbHandle: dbHandle,
	}
}