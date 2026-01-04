package user

import (
	"api/grpc"
	"common/errs"
	"common/httputil"
	userservice "grpc/user/user"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  httputil.ResponseData
// @Failure      400  {object}  httputil.ResponseData
// @Failure      500  {object}  httputil.ResponseData
// @Router       /api/test [post]
func (*TestHandler) Test(ctx *gin.Context) {
	_, err := grpc.UserServiceClient.Test(ctx, &userservice.Req{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		httputil.ErrorJsonResponse(ctx, "", code, msg, nil)
		return
	}
}
