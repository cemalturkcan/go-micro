package middlewares

import (
	"common/config"
	"common/middlewares/authentication"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RegisterMiddlewares(s *fiber.App) {
	s.Use(recover.New())
	s.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
	}))
	if commonconfig.LoggerEnabled {
		s.Use(logger.New())
	}

	s.Use(compress.New())

	authentication.RegisterAuthenticationMiddleware(s)
}

func RegisterFinalMiddlewares(s *fiber.App) {

	s.Static("/", "./public")
	s.Use(func(c *fiber.Ctx) error {
		log.Info("404: ", c.Path())
		log.Info(c.Method(), c.Path(), c.Response().StatusCode())
		return c.Redirect("/")
	})
}
