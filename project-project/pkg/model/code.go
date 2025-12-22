package model

import (
	"gwh.com/project-common/errs"
)

var (
	RedisError              = errs.NewError(999, "redis error")
	DBError                 = errs.NewError(998, "db error")
	ParamsError             = errs.NewError(401, "参数错误")
	NoLegalPhone            = errs.NewError(10102001, "手机号不合法")
	CaptchaError            = errs.NewError(10102002, "验证码错误")
	CaptchaNotExist         = errs.NewError(10102003, "验证码不存在或已过期")
	EmailExist              = errs.NewError(10102004, "邮箱已存在")
	AccountExist            = errs.NewError(10102005, "账号已存在")
	PhoneExist              = errs.NewError(10102006, "手机号已存在")
	AccountAndPasswordError = errs.NewError(10102007, "账号密码不正确")
	TaskNameNotNull         = errs.NewError(20102001, "任务标题不能为空")
	TaskStagesNotNull       = errs.NewError(20102002, "任务步骤不存在")
	ProjectAlreadyDeleted   = errs.NewError(20102003, "项目已经删除了")
)
