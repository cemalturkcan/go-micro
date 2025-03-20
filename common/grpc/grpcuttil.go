package grpcutil

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clients = make(map[string]*grpc.ClientConn)

func ServiceConnection(serviceName string) (*grpc.ClientConn, error) {
	if conn, ok := clients[serviceName]; ok {
		return conn, nil
	}
	conn, err := grpc.NewClient(fmt.Sprintf("registry:///%s-grpc", serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return nil, err
	}
	clients[serviceName] = conn
	return conn, nil
}
