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
	router.InitUserRpc()
	gc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	c := config.InitKafkaWriter()
	//初始化kafka消费者
	reader := config.NewCacheReader()
	go reader.DeleteCache()
	stop := func() {
		gc.Stop()
		c()
		reader.R.Close()
	}
	//初始化rpc
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, stop)
}
