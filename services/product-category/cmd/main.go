package main

import (
	"common/app"
	grpcutil "common/grpc"
	"common/grpc/helloworld"
	"common/middlewares/authentication"
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
)

func main() {
	app.Load(app.Config{
		RegisterMiddlewaresAfter: RegisterMiddlewaresAfter,
		RegisterRoutes:           RegisterRoutes,
		RegisterGrpcRoutes:       RegisterGrpcRoutes,
		ConnectKeystore:          true,
		ConnectDatabase:          true,
	})
}

func RegisterMiddlewaresAfter(app *fiber.App) {
	authentication.RegisterAuthenticationMiddleware(app, []string{"/"})
}
func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		conn, err := grpcutil.ServiceConnection("product-catalog")
		if err != nil {
			return c.SendString("Error connecting to grpc server")
		}

		helloClient := helloworld.NewGreeterClient(conn)
		log.Println("Calling grpc server")
		resp, err := helloClient.SayHello(context.Background(), &helloworld.HelloRequest{Name: "Product Category"})
		log.Println("Response from grpc server")
		if err != nil {
			log.Println(err)
			return c.SendString("Error calling grpc server")
		}

		return c.SendString(resp.GetMessage())
	})
}

func RegisterGrpcRoutes(s *grpc.Server) {
}
