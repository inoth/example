package main

import (
	svc "grpc_server/service"
	pb "message/message"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
)

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8500"),
	)
	// server.DefaultServer = gSrv.NewServer(
	// 	server.Registry(consulReg),
	// 	server.Name("go.micro.srv.TestUser"),
	// )
	service := micro.NewService(
		micro.Address(":9978"),
		micro.Version("latest"),
		micro.Name("go.micro.srv.TestUser"),
		micro.Registry(consulReg),
		// micro.WrapHandler(errWrapper),
	)
	service.Init()

	pb.RegisterUserHandler(service.Server(), new(svc.UserService))
	// micro.RegisterSubscriber("go.micro.srv.TestUser", service.Server(), new(svc.UserService))
	service.Run()
}
