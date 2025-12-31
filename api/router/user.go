package router

import (
	"api/grpc"
	"api/handler/user"

	"github.com/gin-gonic/gin"
)

type User struct {
}

func init() {
	ru := &User{}
	Register(ru)
}

func (*User) Route(r *gin.Engine) {
	grpc.InitRpcServiceClient()

	h := user.NewTestHandler()
	r.POST("/api/user/test", h.Test)
}
