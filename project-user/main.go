package main

import (
	"github.com/gin-gonic/gin"
	common "gwh.com/project-common"
	_ "gwh.com/project-user/api"
	"gwh.com/project-user/config"
	"gwh.com/project-user/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr)
}
