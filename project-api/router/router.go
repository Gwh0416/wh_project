package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
