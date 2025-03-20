package registry

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/resolver"
)

type registryResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	service    string
	ticker     *time.Ticker
	cancelFunc context.CancelFunc
}

func (r *registryResolver) start() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel
	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.resolve()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *registryResolver) resolve() {
	services, _, err := Client.Health().Service(r.service, "grpc", true, nil)
	if err != nil {
		return
	}
	addresses := make([]resolver.Address, len(services))
	for i, service := range services {
		addresses[i] = resolver.Address{
			Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
		}
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: addresses})
}

func (r *registryResolver) ResolveNow(o resolver.ResolveNowOptions) {
	//r.resolve()
}

func (r *registryResolver) Close() {
	r.ticker.Stop()
	r.cancelFunc()
}

type registryBuilder struct{}

func (b *registryBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	log.Println("building resolver" + target.Endpoint())
	r := &registryResolver{
		target:  target,
		cc:      cc,
		service: target.Endpoint(),
		ticker:  time.NewTicker(10 * time.Second),
	}
	r.start()
	r.resolve()
	return r, nil
}

func (b *registryBuilder) Scheme() string {
	return "registry"
}

func RegisterGrpcResolver() {
	resolver.Register(&registryBuilder{})
}
