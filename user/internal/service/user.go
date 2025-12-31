package service

import (
	"context"

	"grpc/user/user"
	"user/pkg/errors"
)

type UserService struct {
	user_service.UnimplementedUserServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Test(ctx context.Context, req *user_service.Req) (*user_service.Resp, error) {
	return nil, errors.NewDBError("错误")
}
