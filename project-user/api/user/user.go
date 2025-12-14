package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	common "gwh.com/project-common"
	"gwh.com/project-user/pkg/dao"
	"gwh.com/project-user/pkg/model"
	"gwh.com/project-user/pkg/repo"
)

type HandlerUser struct {
	cache repo.Cache
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{cache: dao.Rc}
}

func (h *HandlerUser) getCaptcha(c *gin.Context)  {
	resp := &common.Result{}
	//获取参数
	phone := c.PostForm("mobile")
	//校验参数
	if !common.VerifyMobile(phone){
		c.JSON(http.StatusOK, resp.Fail(model.NoLegalPhone, "手机号不合法"))
		return
	}
	//生成验证码
	code := "123456"
	//调用短信平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功")
		//存储验证码 redis 15min
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := h.cache.Put(c, "PHONE_"+phone, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("将手机号和验证码存入redis出错:", zap.Error(err))
		}
	}()
	c.JSON(http.StatusOK, resp.Success(code))
}