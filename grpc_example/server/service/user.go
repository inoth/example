package service

import (
	"context"
	pb "message/message"
)

type UserService struct{}

// func (c *User) GetUserById(ctx context.Context, in *pb.UserIdRequest, opts ...client.CallOption) (*pb.UserIdReply, error) {
func (s *UserService) GetUserById(ctx context.Context, req *pb.UserIdRequest, rsp *pb.UserIdReply) error {
	*rsp = pb.UserIdReply{
		Uid:  "test",
		Name: "test",
	}
	return nil
}
