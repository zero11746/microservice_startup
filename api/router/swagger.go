package router

import (
	_ "api/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Swagger struct {
}

func init() {
	ru := &Swagger{}
	Register(ru)
}

func (*Swagger) Route(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
