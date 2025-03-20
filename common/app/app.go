package app

import (
	"common/commonconfig"
	"common/database"
	"common/keystore"
	"common/registry"
	"common/server"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Load(
	RegisterMiddlewaresBefore func(app *fiber.App),
	RegisterMiddlewaresAfter func(app *fiber.App),
	RegisterRoutes func(app *fiber.App),
	RegisterFinalMiddlewaresBefore func(app *fiber.App),
	RegisterFinalMiddlewaresAfter func(app *fiber.App),
	RegisterGrpcRoutes func(server *grpc.Server),

) {
	//
	id := registry.SetConsulAndRegisterServiceWithConsul()
	registry.RegisterGrpcResolver()
	//

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

	go func() {
		for {
			registry.ReportHealthyState(id, commonconfig.AppName)
			time.Sleep(1 * time.Second)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		database.Close()
		keystore.Close()
		registry.DeregisterService(id)
	}()
	wg.Wait()
}
