package main

import (
	"github.com/gin-gonic/gin"
	common "gwh.com/project-common"
	"gwh.com/project-user/config"
	"gwh.com/project-user/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	stop := func() {
		gc.Stop()
	}
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, stop)
}
