package registry

import (
	"common/config"
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
)

var ErrNotFound = errors.New("no service addresses found")

var Client *api.Client

func SetConsulAndRegisterServiceWithConsul(httpServerAvailable bool, grpcServerAvailable bool) string {
	err := NewConsulClient()
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = commonconfig.ConsulAddress

	id := commonconfig.AppName + "-" + uuid.New().String()

	if httpServerAvailable {
		registration := &api.AgentServiceRegistration{
			ID:      id,
			Name:    commonconfig.AppName,
			Address: commonconfig.ServiceHost,
			Port:    commonconfig.Port,
			Tags: []string{
				"rest",
				"traefik.http.routers." + commonconfig.AppName + "-http.rule=Host(`" + commonconfig.AppName + ".localhost`)",
				"traefik.http.routers." + commonconfig.AppName + "-http.entrypoints=web",
			},
			Check: &api.AgentServiceCheck{
				CheckID:                        id,
				TTL:                            "5s",
				DeregisterCriticalServiceAfter: "10s",
			},
		}
		errService := Client.Agent().ServiceRegister(registration)
		if errService != nil {
			log.Fatalf("Failed to register service: %v", errService)
			return ""
		}
	}

	if grpcServerAvailable {
		grpcRegistration := &api.AgentServiceRegistration{
			ID:      id + "-grpc",
			Name:    commonconfig.AppName + "-grpc",
			Address: commonconfig.ServiceHost,
			Port:    commonconfig.GrpcPort,
			Tags:    []string{"grpc", "traefik.enable=false"},
			Check: &api.AgentServiceCheck{
				CheckID:                        id + "-grpc",
				TTL:                            "5s",
				DeregisterCriticalServiceAfter: "10s",
			},
		}
		errGrpcService := Client.Agent().ServiceRegister(grpcRegistration)
		if errGrpcService != nil {
			log.Fatalf("Failed to register service: %v", errGrpcService)
			return ""
		}

	}

	return id
}

func ReportHealthyState(instanceID string, serviceName string, httpServerAvailable bool, grpcServerAvailable bool) {
	if httpServerAvailable {
		if err := Client.Agent().PassTTL(instanceID, serviceName); err != nil {
			log.Fatalf("Failed to report healthy state: %v", err)
		}
	}

	if grpcServerAvailable {
		if err := Client.Agent().PassTTL(instanceID+"-grpc", serviceName+"-grpc"); err != nil {
			log.Fatalf("Failed to report healthy state: %v", err)
		}
	}
}

func NewConsulClient() error {
	config := api.DefaultConfig()
	config.Address = commonconfig.ConsulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	Client = client
	return nil
}

func DeregisterService(id string) {
	_ = Client.Agent().ServiceDeregister(id)
	_ = Client.Agent().ServiceDeregister(id + "-grpc")
}
