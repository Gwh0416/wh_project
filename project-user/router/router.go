package router

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"gwh.com/project-common/discovery"
	"gwh.com/project-common/logs"
	"gwh.com/project-grpc/user/login"
	"gwh.com/project-user/config"
	"gwh.com/project-user/internal/interceptor"
	login_service_v1 "gwh.com/project-user/pkg/service/login.service.v1"
)

type Router interface {
	Register(r *gin.Engine)
}

type RegisterRouter struct {
}

func NewRegisterRouter() RegisterRouter {
	return RegisterRouter{}
}

func (RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Register(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := NewRegisterRouter()
	//rg.Route(&user.RouterUser{}, r)
	for _, router := range routers {
		router.Register(r)
	}
}

func RegisterRouters(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.AppConf.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			login.RegisterLoginServiceServer(g, login_service_v1.NewLoginService())
		}}
	cacheInterceptor := interceptor.New()
	s := grpc.NewServer(cacheInterceptor.Cache())
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		log.Printf("grpc server started as: %s \n", c.Addr)
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.AppConf.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    config.AppConf.GC.Name,
		Addr:    config.AppConf.GC.Addr,
		Version: config.AppConf.GC.Version,
		Weight:  config.AppConf.GC.Weight,
	}
	r := discovery.NewRegister(config.AppConf.EC.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
