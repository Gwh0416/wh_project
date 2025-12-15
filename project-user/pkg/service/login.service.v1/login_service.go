package login_service_v1

import (
	"context"
	"time"

	"go.uber.org/zap"
	common "gwh.com/project-common"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/user/login"
	"gwh.com/project-user/internal/dao"
	"gwh.com/project-user/internal/repo"
	"gwh.com/project-user/pkg/model"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache repo.Cache
}

func NewLoginService() *LoginService {
	return &LoginService{cache: dao.Rc}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, req *login.CaptchaRequest) (*login.CaptchaResponse, error) {
	//获取参数
	phone := req.Mobile
	//校验参数
	if !common.VerifyMobile(phone) {
		return nil, errs.GrpcError(model.NoLegalPhone)
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
		err := ls.cache.Put(c, "PHONE_"+phone, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("将手机号和验证码存入redis出错:", zap.Error(err))
		}
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, req *login.RegisterRequest) (*login.RegisterResponse, error) {
	return nil, nil
}
