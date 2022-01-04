package main

import (
	"context"
	"fmt"
	pb "message/message"

	gCli "github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
)

var UserSvc = "go.micro.srv.TestUser"

func main() {
	consulReg := consul.NewRegistry(
		registry.Addrs("Consul host"),
	)
	client.DefaultClient = gCli.NewClient(
		client.Registry(consulReg),
	)

	rsp := &pb.UserIdReply{}
	_ = call(UserSvc, "User.GetUserById", pb.UserIdRequest{}, &rsp)

	fmt.Printf("%v:%v", rsp.Uid, rsp.Name)
}

func call(svc string, action string, req interface{}, rsp interface{}) error {
	rq := client.NewRequest(svc, action, req, client.WithContentType("application/json"))
	if err := client.Call(context.TODO(), rq, rsp); err != nil {
		return err
	}
	return nil
}