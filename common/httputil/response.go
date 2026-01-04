package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BusinessCode int

type ResponseData struct {
	Status    int          `json:"status"`
	Code      BusinessCode `json:"code"`
	Message   string       `json:"message"`
	RequestID string       `json:"requestId"`
	Data      any          `json:"data"`
}

const (
	CodeOK = 0
)

func SuccessJsonResponse(ctx *gin.Context, requestID string, data any) {
	if data == nil {
		data = gin.H{}
	}
	resp := &ResponseData{
		Status:    http.StatusOK,
		Code:      CodeOK,
		Message:   "Success",
		RequestID: requestID,
		Data:      data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func ErrorJsonResponse(ctx *gin.Context, requestID string, code BusinessCode, msg string, data any) {
	if data == nil {
		data = gin.H{}
	}
	resp := &ResponseData{
		Status:    http.StatusBadRequest,
		Code:      code,
		Message:   msg,
		RequestID: requestID,
		Data:      data,
	}
	ctx.JSON(http.StatusOK, resp)
}
