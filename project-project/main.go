package main

import (
	"github.com/gin-gonic/gin"
	common "gwh.com/project-common"
	"gwh.com/project-project/config"
	"gwh.com/project-project/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	stop := func() {
		gc.Stop()
	}
	//初始化rpc
	router.InitUserRpc()
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, stop)
}
