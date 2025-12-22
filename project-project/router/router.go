package router

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"gwh.com/project-common/discovery"
	"gwh.com/project-common/logs"
	"gwh.com/project-grpc/account"
	"gwh.com/project-grpc/auth"
	"gwh.com/project-grpc/department"
	"gwh.com/project-grpc/project"
	"gwh.com/project-grpc/task"
	"gwh.com/project-project/config"
	"gwh.com/project-project/internal/interceptor"
	"gwh.com/project-project/internal/rpc"
	account_service_v1 "gwh.com/project-project/pkg/service/account.service.v1"
	auth_service_v1 "gwh.com/project-project/pkg/service/auth.service.v1"
	department_service_v1 "gwh.com/project-project/pkg/service/department.service.v1"
	project_service_v1 "gwh.com/project-project/pkg/service/project.service.v1"
	task_service_v1 "gwh.com/project-project/pkg/service/task.service.v1"
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
			project.RegisterProjectServiceServer(g, project_service_v1.NewProjectService())
			task.RegisterTaskServiceServer(g, task_service_v1.NewTaskService())
			account.RegisterAccountServiceServer(g, account_service_v1.New())
			department.RegisterDepartmentServiceServer(g, department_service_v1.New())
			auth.RegisterAuthServiceServer(g, auth_service_v1.New())
		}}
	s := grpc.NewServer(interceptor.New().Cache())
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

func InitUserRpc() {
	rpc.InitRpcUserClient()
}
