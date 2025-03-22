package main

import (
	"common/app"
	"common/grpc/helloworld"
	"common/role"
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
)

func main() {
	app.Load(app.Config{
		RegisterMiddlewaresBefore:      RegisterMiddlewaresBefore,
		RegisterMiddlewaresAfter:       RegisterMiddlewaresAfter,
		RegisterRoutes:                 RegisterRoutes,
		RegisterFinalMiddlewaresBefore: RegisterFinalMiddlewaresBefore,
		RegisterFinalMiddlewaresAfter:  RegisterFinalMiddlewaresAfter,
		RegisterGrpcRoutes:             RegisterGrpcRoutes,
		ConnectKeystore:                true,
		ConnectDatabase:                true,
	})
}

func RegisterMiddlewaresBefore(app *fiber.App) {
	log.Println("RegisterMiddlewaresBefore")
}

func RegisterMiddlewaresAfter(app *fiber.App) {
	log.Println("RegisterMiddlewaresAfter")
}

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		if role.HasAny(c, []string{"manager"}) {
			return c.SendString("Hello, Manager!")
		}
		return c.SendString("Hello, World!")
	})
}
func RegisterFinalMiddlewaresBefore(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresBefore")
}

func RegisterFinalMiddlewaresAfter(app *fiber.App) {
	log.Println("RegisterFinalMiddlewaresAfter")
}

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func RegisterGrpcRoutes(s *grpc.Server) {
	helloworld.RegisterGreeterServer(s, &server{})
}
