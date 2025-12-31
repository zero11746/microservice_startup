package user

import (
	"api/grpc"
	"common/errs"
	"common/httputils"
	userservice "grpc/user/user"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (*TestHandler) Test(ctx *gin.Context) {
	_, err := grpc.UserServiceClient.Test(ctx, &userservice.Req{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		httputils.ErrorJsonResponse(ctx, "", code, msg, nil)
		return
	}
}
