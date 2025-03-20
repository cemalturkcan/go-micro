package registry

import (
	"common/commonconfig"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
)

var ErrNotFound = errors.New("no service addresses found")

var Client *api.Client

func SetConsulAndRegisterServiceWithConsul() string {
	err := NewConsulClient()
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = commonconfig.ConsulAddress

	id := commonconfig.AppName + "-" + uuid.New().String()

	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    commonconfig.AppName,
		Address: commonconfig.ServiceHost,
		Port:    commonconfig.Port,
		Check: &api.AgentServiceCheck{
			CheckID:                        id,
			TTL:                            "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	grpcRegistration := &api.AgentServiceRegistration{
		ID:      id + "-grpc",
		Name:    commonconfig.AppName + "-grpc",
		Address: commonconfig.ServiceHost,
		Port:    commonconfig.GrpcPort,
		Tags:    []string{"grpc"},
		Check: &api.AgentServiceCheck{
			CheckID:                        id + "-grpc",
			TTL:                            "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	errService := Client.Agent().ServiceRegister(registration)
	errGrpcService := Client.Agent().ServiceRegister(grpcRegistration)

	if errService != nil {
		log.Fatalf("Failed to register service: %v", errService)
		return ""
	}

	if errGrpcService != nil {
		log.Fatalf("Failed to register service: %v", errGrpcService)
		return ""
	}

	return id
}

func ReportHealthyState(instanceID string, serviceName string) {
	if err := Client.Agent().PassTTL(instanceID, serviceName); err != nil {
		log.Fatalf("Failed to report healthy state: %v", err)
	}
	if err := Client.Agent().PassTTL(instanceID+"-grpc", serviceName+"-grpc"); err != nil {
		log.Fatalf("Failed to report healthy state: %v", err)
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

func ServiceAddresses(serviceName string) ([]string, error) {
	entries, _, err := Client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}
