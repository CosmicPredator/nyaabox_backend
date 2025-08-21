package http

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
)

type Middlewares struct {
	app *fiber.App
}

func (m *Middlewares) jwtMiddleware() {
	m.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{ Key: []byte("secret") },
	}))
}

func (m *Middlewares) Regsiter() {
	m.jwtMiddleware()
}

func NewMiddlewares(app *fiber.App) *Middlewares {
	return &Middlewares{ app: app }
}