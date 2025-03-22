package app

import (
	"common/config"
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

type Config struct {
	// This one will be register before default middlewares
	RegisterMiddlewaresBefore func(app *fiber.App)
	// This one will be register after default middlewares
	RegisterMiddlewaresAfter func(app *fiber.App)

	RegisterRoutes func(app *fiber.App)

	// After all routes are registered and before default final middlewares
	RegisterFinalMiddlewaresBefore func(app *fiber.App)

	// After default final middlewares
	RegisterFinalMiddlewaresAfter func(app *fiber.App)

	// Grpc routes
	RegisterGrpcRoutes func(server *grpc.Server)

	ConnectDatabase bool
	ConnectKeystore bool
}

func Load(config Config) {
	//
	id := registry.SetConsulAndRegisterServiceWithConsul(config.RegisterRoutes != nil, config.RegisterGrpcRoutes != nil)
	registry.RegisterGrpcResolver()
	//

	if config.ConnectDatabase {
		database.Connect()
		database.MigrateDb()
	}

	if config.ConnectKeystore {
		keystore.Connect()
	}

	var wg sync.WaitGroup

	if config.RegisterRoutes == nil || config.RegisterGrpcRoutes == nil {
		wg.Add(1)
	} else {
		wg.Add(2)
	}

	//

	if config.RegisterRoutes != nil {
		go func() {
			defer wg.Done()
			server.NewWebServer(config.RegisterMiddlewaresBefore, config.RegisterMiddlewaresAfter, config.RegisterRoutes, config.RegisterFinalMiddlewaresBefore, config.RegisterFinalMiddlewaresAfter)
		}()
	}

	//
	if config.RegisterGrpcRoutes != nil {
		go func() {
			defer wg.Done()
			server.NewGrpcServer(config.RegisterGrpcRoutes)
		}()
	}

	//

	//
	go func() {
		for {
			registry.ReportHealthyState(id, commonconfig.AppName, config.RegisterRoutes != nil, config.RegisterGrpcRoutes != nil)
			time.Sleep(1 * time.Second)
		}
	}()
	//

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		if config.ConnectDatabase {
			database.Close()
		}
		if config.ConnectKeystore {
			keystore.Close()
		}

		registry.DeregisterService(id)
	}()
	wg.Wait()
}
