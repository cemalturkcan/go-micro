package app

import (
	"common/consul"
	"common/database"
	"common/keystore"
	"common/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"sync"
)

func Load(
	RegisterMiddlewaresBefore func(app *fiber.App),
	RegisterMiddlewaresAfter func(app *fiber.App),
	RegisterRoutes func(app *fiber.App) fiber.Router,
	RegisterFinalMiddlewaresBefore func(app *fiber.App),
	RegisterFinalMiddlewaresAfter func(app *fiber.App),
	RegisterGrpcRoutes func(server *grpc.Server),

) {
	client, id := consul.RegisterServiceWithConsul()

	if client == nil {
		log.Fatalf("Failed to register service")
		return
	}

	database.Connect()
	database.MigrateDb()
	keystore.Connect()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		server.NewWebServer(RegisterMiddlewaresBefore, RegisterMiddlewaresAfter, RegisterRoutes, RegisterFinalMiddlewaresBefore, RegisterFinalMiddlewaresAfter)
	}()

	go func() {
		defer wg.Done()
		server.NewGrpcServer(RegisterGrpcRoutes)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		database.Close()
		keystore.Close()
		_ = consul.DeregisterService(client, id)
	}()
	wg.Wait()
}
