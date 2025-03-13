package app

import (
	"common/commonconfig"
	"common/consul"
	"common/database"
	"common/exitcode"
	"common/keystore"
	"common/server"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

func Load(
	RegisterMiddlewaresBefore func(app *fiber.App),
	RegisterMiddlewaresAfter func(app *fiber.App),
	RegisterRoutes func(app *fiber.App),
	RegisterFinalMiddlewaresBefore func(app *fiber.App),
	RegisterFinalMiddlewaresAfter func(app *fiber.App),
) {
	client, id := consul.RegisterServiceWithConsul()

	if client == nil {
		log.Fatalf("Failed to register service")
		return
	}

	database.Connect()
	database.MigrateDb()
	keystore.Connect()
	app := server.New(RegisterMiddlewaresBefore, RegisterMiddlewaresAfter, RegisterRoutes, RegisterFinalMiddlewaresBefore, RegisterFinalMiddlewaresAfter)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		database.Close()
		keystore.Close()
		_ = consul.DeregisterService(client, id)
		_ = app.Shutdown()
	}()

	err := app.Listen(fmt.Sprintf(":%d", commonconfig.Port))
	if err != nil {
		os.Exit(exitcode.ServerStartError)
	}
}
