package main

import (
	svc "grpc_server/service"
	pb "message/message"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	gSrv "github.com/asim/go-micro/plugins/server/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
)

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("Consul Host"),
	)
	server.DefaultServer = gSrv.NewServer(
		server.Registry(consulReg),
		server.Name("go.micro.srv.TestUser"),
	)
	service := micro.NewService(
		micro.Address(":9978"),
		micro.Version("latest"),
		// micro.WrapHandler(errWrapper),
	)
	service.Init()

	pb.RegisterUserHandler(service.Server(), new(svc.UserService))
	// micro.RegisterSubscriber("go.micro.srv.TestUser", service.Server(), new(svc.UserService))
	service.Run()
}
