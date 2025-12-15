package model

import (
	"gwh.com/project-common/errs"
)

var (
	NoLegalPhone = errs.NewError(2001, "手机号不合法")
)
