package errs

import (
	"common/httputils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) (httputils.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return httputils.BusinessCode(fromError.Code()), fromError.Message()
}

func ToBError(err error) *BError {
	fromError, _ := status.FromError(err)
	return NewError(ErrorCode(fromError.Code()), fromError.Message())
}
