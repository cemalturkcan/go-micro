package server

import (
	"common/commonconfig"
	"common/exitcode"
	"common/middlewares"
	"common/rest"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
)

func NewWebServer(
	RegisterMiddlewaresBefore func(app *fiber.App),
	RegisterMiddlewaresAfter func(app *fiber.App),
	RegisterRoutes func(app *fiber.App) fiber.Router,
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

	group := RegisterRoutes(app)

	group.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	RegisterFinalMiddlewaresBefore(app)
	middlewares.RegisterFinalMiddlewares(app)
	RegisterFinalMiddlewaresAfter(app)

	err := app.Listen(fmt.Sprintf(":%d", commonconfig.Port))
	if err != nil {
		os.Exit(exitcode.ServerStartError)
	}

	return app
}

func NewGrpcServer(RegisterGrpcRoutes func(server *grpc.Server)) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", commonconfig.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", commonconfig.GrpcPort, err)
	}
	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	RegisterGrpcRoutes(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
	return grpcServer
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	log.Error("Error: ", err)
	code, message := rest.Error(err)
	return rest.ErrorRes(c, code, message)
}
