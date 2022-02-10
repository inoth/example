package main

import (
	"context"
	svc "grpc_server/service"
	pb "message/message"
	"time"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	gSrv "github.com/asim/go-micro/plugins/server/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/sirupsen/logrus"
)

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8500"),
	)
	// server.DefaultServer = gSrv.NewServer(
	// 	server.Registry(consulReg),
	// 	server.Name("go.micro.srv.TestUser"),
	// )

	// go-micro 踩坑，方案A注册方式某些情况下会导致错误
	// OK: HTTP status code 200; transport: missing content-type field

	// service := micro.NewService(
	// 	micro.Address(":9978"),
	// 	micro.Version("latest"),
	// 	micro.Name("go.micro.srv.TestUser"),
	// 	micro.Registry(consulReg),
	// 	micro.RegisterInterval(time.Second*300),
	// 	micro.WrapHandler(ExceptionHandle),
	// )

	// 目前来说正确的做法
	service := micro.NewService(micro.Server(gSrv.NewServer(
		server.Address(":9978"),
		server.Version("leatest"),
		server.Name("go.micro.srv.TestUser"),
		server.Registry(consulReg),
		server.RegisterInterval(time.Second*30),
		server.WrapHandler(ExceptionHandle),
	)))

	service.Init()

	pb.RegisterUserHandler(service.Server(), new(svc.UserService))
	// micro.RegisterSubscriber("go.micro.srv.TestUser", service.Server(), new(svc.UserService))
	service.Run()
}

func ExceptionHandle(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		defer func() {
			err := recover()
			if err != nil {
				logrus.Error(err)
			}
		}()
		return fn(ctx, req, rsp)
	}
}
