package service

import (
	"context"
	pb "message/message"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct{}

// func (c *User) GetUserById(ctx context.Context, in *pb.UserIdRequest, opts ...client.CallOption) (*pb.UserIdReply, error) {
func (s *UserService) GetUserById(ctx context.Context, req *pb.UserIdRequest, rsp *pb.UserIdReply) error {
	*rsp = pb.UserIdReply{
		Uid: "test",
	}
	return nil
}

func (h *UserService) Test(ctx context.Context, req *pb.ApplicationRequest, rsp *pb.ApplicationDetailReply) error {
	res := &pb.ApplicationDetailReply{
		Id:      1,
		AppName: "qweqwe",
		Desc:    "asdasdas",
		AppType: "asdadsa",
	}
	for i := 0; i < 3; i++ {
		res.Versions = append(res.Versions, &pb.ApplicationVersion{
			Id:          int32(i + 1),
			Version:     "iasas",
			PackagePath: "asdadds",
			CreateTime:  timestamppb.New(time.Now()),
		})
	}
	*rsp = *res
	return nil
}
