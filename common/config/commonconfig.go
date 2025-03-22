package commonconfig

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	// DATABASE
	DbHost     = ""
	DbPort     = ""
	DbUsername = ""
	DbPassword = ""
	DbDatabase = ""
	DbSchema   = ""
	DbSslMode  = ""

	// KEYSTORE
	KeyStoreHost     = ""
	KeyStorePort     = ""
	KeyStoreDb       = 0
	KeyStorePassword = ""

	AppName          = ""
	Port             = 8080
	LoadBalancerPort = 8081
	GrpcPort         = 50051

	PreFork = false

	LoggerEnabled = false

	// Consul
	ConsulAddress = ""
	ServiceHost   = ""
)

const (
	Development = "development"
	Production  = "production"
)

var (
	Mode = Development
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
		log.Println(err)
	}

	log.Println("Loading .env file")

	//DB
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbDatabase = os.Getenv("DB_DATABASE")
	DbSchema = os.Getenv("DB_SCHEMA")
	DbSslMode = os.Getenv("DB_SSL_MODE")

	KeyStoreHost = os.Getenv("KEYSTORE_HOST")
	KeyStorePort = os.Getenv("KEYSTORE_PORT")
	tempKeyStoreDb, _ := strconv.Atoi(os.Getenv("KEYSTORE_DB"))
	KeyStoreDb = tempKeyStoreDb

	KeyStorePassword = os.Getenv("KEYSTORE_PASSWORD")

	AppName = os.Getenv("APP_NAME")

	PreFork = os.Getenv("PreFork") == "true"

	LoggerEnabled = os.Getenv("LOGGER_ENABLED") == "true"

	ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ServiceHost = os.Getenv("SERVICE_ADDRESS")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		Port = port
	}

	gPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err == nil {
		GrpcPort = gPort
	}

	lbPort, err := strconv.Atoi(os.Getenv("LOAD_BALANCER_PORT"))
	if err == nil {
		LoadBalancerPort = lbPort
	}

	if os.Getenv("MODE") == Production {
		Mode = Production
	}

}
