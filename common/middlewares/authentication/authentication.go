package authentication

import (
	grpcutil "common/grpc"
	"common/grpc/authentication"
	"common/rest"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func RegisterAuthenticationMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {

		if filterJwtMiddleware(c) {
			return c.Next()
		}
		authHeader := c.Cookies("Authorization")

		if authHeader == "" {
			authHeader = strings.Replace(c.Get("Authorization"), "Bearer ", "", 1)
		}

		if authHeader == "" {
			return onUnAuthorized(c)
		}

		resp, err := CheckAndGetUserInfo(authHeader)
		if err != nil {
			return onUnAuthorized(c)
		}

		c.Locals("ID", resp.Id)
		c.Locals("Email", resp.Email)
		c.Locals("Permissions", resp.Permissions)
		return c.Next()
	})

}

func CheckAndGetUserInfo(token string) (*authentication.CheckAndGetUserInfoResponse, error) {
	conn, err := grpcutil.ServiceConnection("authentication")
	if err != nil {
		return nil, err
	}
	authClient := authentication.NewAuthenticationServiceClient(conn)

	resp, err := authClient.CheckAndGetUserInfo(context.Background(), &authentication.CheckAndGetUserInfoRequest{Token: token})

	log.Println("Response from grpc server")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func onUnAuthorized(c *fiber.Ctx) error {
	return rest.ErrorRes(c, rest.Unauthorized, rest.ErrorCode[rest.Unauthorized])
}

var specialAuthorizedRoutes = []string{
	"/login",
	"/register",
}

func filterJwtMiddleware(c *fiber.Ctx) bool {
	for _, route := range specialAuthorizedRoutes {
		if route == c.Path() {
			return true
		}
	}
	return false
}
