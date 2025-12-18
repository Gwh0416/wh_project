package login_service_v1

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	common "gwh.com/project-common"
	"gwh.com/project-common/encrypts"
	"gwh.com/project-common/errs"
	"gwh.com/project-common/jwts"
	"gwh.com/project-common/tms"
	"gwh.com/project-grpc/user/login"
	"gwh.com/project-user/config"
	"gwh.com/project-user/internal/dao"
	"gwh.com/project-user/internal/data/member"
	"gwh.com/project-user/internal/data/organization"
	"gwh.com/project-user/internal/database"
	"gwh.com/project-user/internal/database/tran"
	"gwh.com/project-user/internal/repo"
	"gwh.com/project-user/pkg/model"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction
}

func NewLoginService() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		transaction:      dao.NewTransactionImpl(),
	}
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
		err := ls.cache.Put(c, model.PhoneRegisKey+phone, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("将手机号和验证码存入redis出错:", zap.Error(err))
		}
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, req *login.RegisterRequest) (*login.RegisterResponse, error) {
	//校验验证码
	cx := context.Background()
	redisCode, err := ls.cache.Get(cx, model.PhoneRegisKey+req.Mobile)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errs.GrpcError(model.CaptchaNotExist)
		}
		zap.L().Error("redis Get phone err:", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	////校验业务逻辑（邮箱是否被注册 账号是否被注册 手机号是否被注册）
	exist, err := ls.memberRepo.GetMemberByEmail(cx, req.Email)
	if err != nil {
		zap.L().Error("Register DB get err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	exist, err = ls.memberRepo.GetMemberByAccount(cx, req.Name)
	if err != nil {
		zap.L().Error("Register DB get err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	exist, err = ls.memberRepo.GetMemberByMobile(cx, req.Mobile)
	if err != nil {
		zap.L().Error("Register DB get err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.PhoneExist)
	}
	//执行业务 存入数据库
	pwd := encrypts.Md5(req.Password)
	mem := &member.Member{
		Account:       req.Name,
		Password:      pwd,
		Name:          req.Name,
		Mobile:        req.Mobile,
		Email:         req.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = ls.transaction.Action(func(conn database.DBConn) error {
		err = ls.memberRepo.SaveMember(conn, ctx, mem)
		if err != nil {
			zap.L().Error("Register DB SaveMember err:", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		//存入组织
		err = ls.organizationRepo.SaveOrganization(conn, cx, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})

	return &login.RegisterResponse{}, err
}

func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	cx := context.Background()
	//查询数据库
	pwd := encrypts.Md5(msg.Password)
	mem, err := ls.memberRepo.FindMember(cx, msg.Account, pwd)
	if err != nil {
		zap.L().Error("Login DB FindMember err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountAndPasswordError)
	}
	memMessage := &login.MemberMessage{}
	err = copier.Copy(&memMessage, mem)
	memMessage.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)
	memMessage.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	memMessage.CreateTime = tms.FormatByMill(mem.CreateTime)

	//根据用户id查询组织
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(cx, mem.Id)
	if err != nil {
		zap.L().Error("Login DB FindOrganizationByMemId err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
		org.OwnerCode = memMessage.Code
		org.CreateTime = tms.FormatByMill(organization.ToMap(orgs)[org.Id].CreateTime)
	}
	if len(orgs) > 0 {
		memMessage.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	//jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	aExp := time.Duration(config.AppConf.JC.AccessExp) * 3600 * 24 * time.Second
	rExp := time.Duration(config.AppConf.JC.RefreshExp) * 3600 * 24 * time.Second

	token := jwts.CreateToken(memIdStr, config.AppConf.JC.AccessSecret, config.AppConf.JC.RefreshSecret, aExp, rExp)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}
	resp := &login.LoginResponse{
		Member:           memMessage,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}
	return resp, nil
}

func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.Replace(token, "bearer ", "", 1)
	}
	parseToken, err := jwts.ParseToken(token, config.AppConf.JC.AccessSecret)
	if err != nil {
		zap.L().Error("TokenVerify ParseToken err:", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	//数据库查询 优化点 登陆之后 把用户信息缓存起来
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	mem, err := ls.memberRepo.FindMemberById(context.Background(), id)
	if err != nil {
		zap.L().Error("TokenVerify DB FindMember err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	memMessage := &login.MemberMessage{}
	err = copier.Copy(&memMessage, mem)
	memMessage.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)

	orgs, err := ls.organizationRepo.FindOrganizationByMemId(context.Background(), mem.Id)
	if err != nil {
		zap.L().Error("Login DB FindOrganizationByMemId err:", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(orgs) > 0 {
		memMessage.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}

	return &login.LoginResponse{
		Member: memMessage,
	}, err
}

func (l *LoginService) MyOrgList(ctx context.Context, msg *login.UserMessage) (*login.OrgListResponse, error) {
	memId := msg.MemId
	orgs, err := l.organizationRepo.FindOrganizationByMemId(ctx, memId)
	if err != nil {
		zap.L().Error("MyOrgList FindOrganizationByMemId err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
	}
	return &login.OrgListResponse{OrganizationList: orgsMessage}, nil
}

func (ls *LoginService) FindMemInfoById(ctx context.Context, msg *login.UserMessage) (*login.MemberMessage, error) {
	memberById, err := ls.memberRepo.FindMemberById(context.Background(), msg.MemId)
	if err != nil {
		zap.L().Error("TokenVerify db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	orgs, err := ls.organizationRepo.FindOrganizationByMemId(context.Background(), memberById.Id)
	if err != nil {
		zap.L().Error("TokenVerify db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return memMsg, nil
}
