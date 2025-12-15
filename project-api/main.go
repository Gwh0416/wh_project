package main

import (
	"github.com/gin-gonic/gin"
	_ "gwh.com/project-api/api"
	"gwh.com/project-api/config"
	"gwh.com/project-api/router"
	common "gwh.com/project-common"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, nil)
}
