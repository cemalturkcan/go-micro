package main

import (
	"common/app"
	"common/commonconfig"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app.Load(RegisterMiddlewaresBefore, RegisterMiddlewaresAfter, RegisterRoutes, RegisterFinalMiddlewaresBefore, RegisterFinalMiddlewaresAfter)
}

func RegisterMiddlewaresBefore(app *fiber.App) {
	log.Println("RegisterMiddlewaresBefore")
}

func RegisterMiddlewaresAfter(app *fiber.App) {
	log.Println("RegisterMiddlewaresAfter")
}

func RegisterRoutes(app *fiber.App) fiber.Router {
	prefixGroup := app.Group(fmt.Sprintf("/%s", commonconfig.AppName))
	prefixGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return prefixGroup
}

func RegisterFinalMiddlewaresBefore(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresBefore")
}

func RegisterFinalMiddlewaresAfter(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresAfter")
}
