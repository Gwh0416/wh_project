package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"gwh.com/project-user/router"
)

func init(){
	log.Println("init user")
	router.RegisterRouters(&RouterUser{})
}

type RouterUser struct {

}

func (*RouterUser) Register(r *gin.Engine) {
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}