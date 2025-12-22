package main

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "gwh.com/project-api/api"
	"gwh.com/project-api/config"
	"gwh.com/project-api/router"
	common "gwh.com/project-common"
)

func main() {
	r := gin.Default()
	//r.Use(midd.RequestLog())
	r.StaticFS("/upload", http.Dir("upload"))
	router.InitRouter(r)
	//开启pprof 默认访问路径是/debug/pprof dev:"/"
	pprof.Register(r)
	common.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, nil)
}
