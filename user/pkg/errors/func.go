package errors

import (
	"fmt"

	"common/errs"
)

func NewDBError(format string, a ...any) *errs.BError {
	err := fmt.Sprintf(format, a...)
	return errs.NewError(DBErrorCode, err)
}

func NewRedisError(format string, a ...any) *errs.BError {
	err := fmt.Sprintf(format, a...)
	return errs.NewError(RedisErrorCode, err)
}

func NewErrorWithCode(code errs.ErrorCode, format string, a ...any) *errs.BError {
	err := fmt.Sprintf(format, a...)
	return errs.NewError(code, err)
}
