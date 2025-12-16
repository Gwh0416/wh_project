package project

import (
	"log"

	"github.com/gin-gonic/gin"
	"gwh.com/project-api/api/midd"
	"gwh.com/project-api/router"
)

func init() {
	log.Println("init user router")
	router.RegisterRouters(&RouterProject{})
}

type RouterProject struct {
}

func (*RouterProject) Register(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()
	h := NewHandlerProject()
	group := r.Group("/project/index")
	group.Use(midd.TokenVerify())
	group.POST("", h.index)
}
