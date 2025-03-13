package server

import (
	"common/commonconfig"
	"common/middlewares"
	"common/rest"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func New(
	RegisterMiddlewaresBefore func(app *fiber.App),
	RegisterMiddlewaresAfter func(app *fiber.App),
	RegisterRoutes func(app *fiber.App),
	RegisterFinalMiddlewaresBefore func(app *fiber.App),
	RegisterFinalMiddlewaresAfter func(app *fiber.App),
) *fiber.App {

	app := fiber.New(fiber.Config{
		AppName:      commonconfig.AppName,
		ErrorHandler: ErrorHandler,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	RegisterMiddlewaresBefore(app)
	middlewares.RegisterMiddlewares(app)
	RegisterMiddlewaresAfter(app)

	RegisterRoutes(app)

	app.Get(fmt.Sprintf("/%s/health", commonconfig.AppName), func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	RegisterFinalMiddlewaresBefore(app)
	middlewares.RegisterFinalMiddlewares(app)
	RegisterFinalMiddlewaresAfter(app)

	return app
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	log.Error("Error: ", err)
	code, message := rest.Error(err)
	return rest.ErrorRes(c, code, message)
}
