package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gwh.com/project-api/pkg/model/user"
	common "gwh.com/project-common"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/user/login"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) getCaptcha(c *gin.Context) {
	resp := &common.Result{}
	phone := c.PostForm("mobile")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	captchaResp, err := LoginServiceClient.GetCaptcha(ctx, &login.CaptchaRequest{Mobile: phone})
	if err != nil {
		code, msg := errs.PraseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, resp.Success(captchaResp.Code))
}

func (h *HandlerUser) register(c *gin.Context) {
	//接收参数
	resp := &common.Result{}
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//校验参数
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//调用user grpc服务 获取响应
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	registerRequest := &login.RegisterRequest{}
	err := copier.Copy(registerRequest, req)
	if err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	_, err = LoginServiceClient.Register(ctx, registerRequest)
	if err != nil {
		code, msg := errs.PraseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	//返回结果
	c.JSON(http.StatusOK, resp.Success(nil))
}
