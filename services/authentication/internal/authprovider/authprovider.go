package authprovider

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserInformation struct {
	Id          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Permissions []string `json:"permissions"`
}

var (
	authProvider *gocloak.GoCloak
	path         = ""
	userName     = ""
	password     = ""
	realm        = ""
	clientId     = ""
	clientSecret = ""
	adminToken   string
)

func NewClient() {
	_ = godotenv.Load(".env")

	path = os.Getenv("KEYCLOAK_PATH")
	userName = os.Getenv("KEYCLOAK_USER")
	password = os.Getenv("KEYCLOAK_PASSWORD")
	realm = os.Getenv("KEYCLOAK_REALM")
	clientId = os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")

	authProvider = gocloak.NewClient(path)
	loginAdmin()
	go func() {
		for {
			<-time.After(4 * time.Minute)
			loginAdmin()
		}
	}()
}

func loginAdmin() {
	tokentemp, err := authProvider.LoginAdmin(context.Background(), userName, password, realm)
	if err != nil {
		log.Println("Error: ", err)
		log.Fatal(err)
	}
	log.Println("Admin Token: ", tokentemp.AccessToken)
	adminToken = tokentemp.AccessToken
}

func RegisterUser(user User) error {
	log.Println("Register User")
	log.Println("realm: ", realm)
	id, err := authProvider.CreateUser(context.Background(), adminToken, realm, gocloak.User{
		Username:  &user.Username,
		FirstName: &user.FirstName,
		LastName:  &user.LastName,
		Email:     &user.Email,
		Enabled:   gocloak.BoolP(true),
	})

	if err != nil {
		log.Println("Errr ", err)
		return err
	}

	err = authProvider.SetPassword(context.Background(), adminToken, id, realm, user.Password, false)

	return err
}

func Login(username, password string) (string, error) {
	jwt, err := authProvider.Login(context.Background(), clientId, clientSecret, realm, username, password)

	if err != nil {
		return "", err
	}
	return jwt.AccessToken, nil
}

func CheckAndGetUserInfo(token string) (*UserInformation, error) {
	res, err := authProvider.RetrospectToken(context.Background(), token, clientId, clientSecret, realm)

	if err != nil {
		return nil, err
	}
	if !*res.Active {
		return nil, errors.New("adminToken is not active")
	}

	user, err := authProvider.GetUserInfo(context.Background(), token, realm)

	if user == nil {
		return nil, errors.New("user not found")
	}

	roles, err := authProvider.GetRoleMappingByUserID(context.Background(), token, realm, *user.Sub)
	if err != nil {
		return nil, err
	}

	log.Println(res.Aud)

	return &UserInformation{
		Id:          *user.Sub,
		Username:    *user.PreferredUsername,
		Email:       *user.Email,
		FirstName:   *user.GivenName,
		LastName:    *user.FamilyName,
		Permissions: GetRealmMappings(roles),
	}, nil
}

func GetRealmMappings(mapping *gocloak.MappingsRepresentation) []string {
	var roles []string

	if mapping.RealmMappings == nil {
		return roles
	}

	for _, role := range *mapping.RealmMappings {
		roles = append(roles, *role.Name)
	}
	return roles
}
