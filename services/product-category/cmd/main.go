package main

import (
	"common/app"
	"common/commonconfig"
	"common/grpc/helloworld"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
)

func main() {
	app.Load(RegisterMiddlewaresBefore, RegisterMiddlewaresAfter, RegisterRoutes, RegisterFinalMiddlewaresBefore, RegisterFinalMiddlewaresAfter, RegisterGrpcRoutes)
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
		conn, err := grpc.NewClient("")

		if err != nil {
			return c.SendString("Error connecting to grpc server")
		}
		helloClient := helloworld.NewGreeterClient(conn)
		resp, err := helloClient.SayHello(c.Context(), &helloworld.HelloRequest{Name: "Product Category"})
		return c.SendString(resp.GetMessage())
	})

	return prefixGroup
}

func RegisterFinalMiddlewaresBefore(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresBefore")
}

func RegisterFinalMiddlewaresAfter(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresAfter")
}

func RegisterGrpcRoutes(s *grpc.Server) {
}
