package errors

import "common/errs"

const (
	RedisErrorCode    errs.ErrorCode = 999
	DBErrorCode       errs.ErrorCode = 998
	ParamsErrorCode   errs.ErrorCode = 401
	NoLegalMobileCode errs.ErrorCode = 10102001
)

var (
	RedisError    = errs.NewError(RedisErrorCode, "redis错误")
	DBError       = errs.NewError(DBErrorCode, "db错误")
	ParamsError   = errs.NewError(ParamsErrorCode, "参数错误")
	NoLegalMobile = errs.NewError(NoLegalMobileCode, "手机号不合法")
)

var UnknownError = errs.NewError(-1, "未知错误")

var Errors = map[errs.ErrorCode]error{
	RedisErrorCode:    RedisError,
	DBErrorCode:       DBError,
	ParamsErrorCode:   ParamsError,
	NoLegalMobileCode: NoLegalMobile,
}
