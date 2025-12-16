package model

import (
	"gwh.com/project-common/errs"
)

var (
	RedisError              = errs.NewError(999, "redis error")
	DBError                 = errs.NewError(998, "db error")
	NoLegalPhone            = errs.NewError(10102001, "手机号不合法")
	CaptchaError            = errs.NewError(10102002, "验证码错误")
	CaptchaNotExist         = errs.NewError(10102003, "验证码不存在或已过期")
	EmailExist              = errs.NewError(10102004, "邮箱已存在")
	AccountExist            = errs.NewError(10102005, "账号已存在")
	PhoneExist              = errs.NewError(10102006, "手机号已存在")
	AccountAndPasswordError = errs.NewError(10102007, "账号密码不正确")
)
