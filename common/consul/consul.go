package consul

import (
	"common/commonconfig"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
)

func RegisterServiceWithConsul() (*api.Client, string) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = commonconfig.ConsulAddress

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	id := commonconfig.AppName + "-" + uuid.New().String()

	checkHTTP := fmt.Sprintf("http://%s:%d/%s/health", commonconfig.ServiceHost, commonconfig.Port, commonconfig.AppName)
	log.Printf("Check HTTP: %s", checkHTTP)

	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    commonconfig.AppName,
		Address: commonconfig.ServiceHost,
		Port:    commonconfig.Port,
		Tags: []string{
			// Traefik gRPC routing metadata
			"traefik.enable=true",
			"traefik.tcp.routers." + commonconfig.AppName + ".entrypoints=grpc",
			"traefik.tcp.routers." + commonconfig.AppName + ".service=" + commonconfig.AppName,
			"traefik.tcp.services." + commonconfig.AppName + ".loadbalancer.server.port=" + fmt.Sprintf("%d", commonconfig.GrpcPort),
		},
		Check: &api.AgentServiceCheck{
			HTTP:     checkHTTP,
			Interval: "10s",
			Notes:    "Consul check",
		},
		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				GRPC:     fmt.Sprintf("%s:%d", commonconfig.ServiceHost, commonconfig.GrpcPort),
				Interval: "10s",
				Notes:    "gRPC health check",
			},
		},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
		return nil, ""
	}

	fmt.Println("Service registered with Consul")

	return client, id
}

func DeregisterService(client *api.Client, id string) error {
	return client.Agent().ServiceDeregister(id)
}
