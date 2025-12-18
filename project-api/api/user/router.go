package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"gwh.com/project-api/api/midd"
	"gwh.com/project-api/api/rpc"
	"gwh.com/project-api/router"
)

func init() {
	log.Println("init user router")
	router.RegisterRouters(&RouterUser{})
}

type RouterUser struct {
}

func (*RouterUser) Register(r *gin.Engine) {
	//初始化grpc的客户端连接
	rpc.InitRpcUserClient()
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.register)
	r.POST("/project/login", h.login)
	org := r.Group("/project/organization")
	org.Use(midd.TokenVerify())
	org.POST("/_getOrgList", h.myOrgList)
}
