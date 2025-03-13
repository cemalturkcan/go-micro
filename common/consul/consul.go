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
	checkHTTP := fmt.Sprintf("http://%s:%d/%s/health", commonconfig.ServiceAddress, commonconfig.Port, commonconfig.AppName)
	log.Printf("Check HTTP: %s", checkHTTP)

	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    commonconfig.AppName,
		Address: commonconfig.ServiceAddress,
		Port:    commonconfig.Port,
		Check: &api.AgentServiceCheck{
			HTTP:     checkHTTP,
			Interval: "10s",
			Notes:    "Consul check",
		},
		Tags: []string{"traefik.enable=true"},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
		return nil, ""
	}

	fmt.Println("Service A registered with Consul")

	return client, id
}

func DeregisterService(client *api.Client, id string) error {
	return client.Agent().ServiceDeregister(id)
}
