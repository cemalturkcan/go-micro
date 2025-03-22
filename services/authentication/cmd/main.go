package main

import (
	"authentication/internal/authprovider"
	"common/app"
	"common/grpc/authentication"
	"common/rest"
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"log"
)

func main() {
	app.Load(app.Config{
		RegisterGrpcRoutes: RegisterGrpcRoutes,
		RegisterRoutes:     RegisterRoutes,
		ConnectKeystore:    false,
		ConnectDatabase:    false,
	})
}

type server struct {
	authentication.UnimplementedAuthenticationServiceServer
}

func (s *server) CheckAndGetUserInfo(_ context.Context, in *authentication.CheckAndGetUserInfoRequest) (*authentication.CheckAndGetUserInfoResponse, error) {
	res, err := authprovider.CheckAndGetUserInfo(in.Token)
	if err != nil {
		return nil, err
	}
	return &authentication.CheckAndGetUserInfoResponse{
		Id:          res.Id,
		Username:    res.Username,
		Email:       res.Email,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Permissions: res.Permissions,
	}, nil

}

func RegisterGrpcRoutes(s *grpc.Server) {
	authprovider.NewClient()
	authentication.RegisterAuthenticationServiceServer(s, &server{})
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func RegisterRoutes(app *fiber.App) {
	app.Post("/login", Login)
	app.Post("/register", Register)
}

func Login(c *fiber.Ctx) error {
	log.Println("Login")
	var req LoginRequest
	err := rest.SetBodyAndValidate(c, &req)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	res, err := authprovider.Login(req.Username, req.Password)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	return rest.Res(c, nil, LoginResponse{Token: res})
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	err := rest.SetBodyAndValidate(c, &req)
	if err != nil {
		return rest.Res(c, err, nil)
	}

	err = authprovider.RegisterUser(authprovider.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})

	if err != nil {
		return rest.Res(c, err, nil)
	}

	res, err := authprovider.Login(req.Username, req.Password)
	if err != nil {
		return rest.Res(c, err, nil)
	}
	return rest.Res(c, nil, LoginResponse{Token: res})
}
