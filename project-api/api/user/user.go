package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gwh.com/project-api/api/rpc"
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
	captchaResp, err := rpc.LoginServiceClient.GetCaptcha(ctx, &login.CaptchaRequest{Mobile: phone})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
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
	_, err = rpc.LoginServiceClient.Register(ctx, registerRequest)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	//返回结果
	c.JSON(http.StatusOK, resp.Success(nil))
}

func (h *HandlerUser) login(c *gin.Context) {
	resp := &common.Result{}
	var req user.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用user grpc 完成登录
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	registerRequest := &login.LoginMessage{}
	err := copier.Copy(registerRequest, req)
	if err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	loginResp, err := rpc.LoginServiceClient.Login(ctx, registerRequest)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	result := &user.LoginResp{}
	err = copier.Copy(result, loginResp)
	if err != nil {
		c.JSON(http.StatusOK, resp.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	//返回结果
	c.JSON(http.StatusOK, resp.Success(result))
}

func (u *HandlerUser) myOrgList(c *gin.Context) {
	result := &common.Result{}
	memberId := c.GetInt64("memberId")
	list, err := rpc.LoginServiceClient.MyOrgList(context.Background(), &login.UserMessage{MemId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	if list.OrganizationList == nil {
		c.JSON(http.StatusOK, result.Success([]*user.OrganizationList{}))
		return
	}
	var orgs []*user.OrganizationList
	err = copier.Copy(&orgs, list.OrganizationList)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	c.JSON(http.StatusOK, result.Success(orgs))
}
