package errs

import (
	"common/httputil"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) (httputil.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return httputil.BusinessCode(fromError.Code()), fromError.Message()
}

func ToBError(err error) *BError {
	fromError, _ := status.FromError(err)
	return NewError(ErrorCode(fromError.Code()), fromError.Message())
}
