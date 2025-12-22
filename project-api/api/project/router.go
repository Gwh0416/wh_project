package project

import (
	"log"

	"github.com/gin-gonic/gin"
	"gwh.com/project-api/api/midd"
	"gwh.com/project-api/api/rpc"
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
	rpc.InitRpcProjectClient()
	h := NewHandlerProject()
	group := r.Group("/project")
	group.Use(midd.TokenVerify())
	group.POST("/index", h.index)
	group.POST("/project/selfList", h.myProjectList)
	group.POST("/project", h.myProjectList)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)
	group.POST("/project/recovery", h.recoveryProject)
	group.POST("/project_collect/collect", h.collectProject)
	group.POST("/project/edit", h.editProject)
	group.POST("/project/getLogBySelfProject", h.getLogBySelfProject)

	t := NewTask()
	group.POST("/task_stages", t.taskStages)
	group.POST("/project_member/index", t.taskMemberList)
	group.POST("/task_stages/tasks", t.taskList)
	group.POST("/task/save", t.saveTask)
	group.POST("/task/sort", t.taskSort)
	group.POST("/task/selfList", t.myTaskList)
	group.POST("/task/read", t.readTask)
	group.POST("/task_member", t.listTaskMember)
	group.POST("/task/taskLog", t.taskLog)
	group.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	group.POST("/task/saveTaskWorkTime", t.saveTaskWorkTime)
	group.POST("/file/uploadFiles", t.uploadFiles)
	group.POST("/task/taskSources", t.taskSources)
	group.POST("/task/createComment", t.createComment)
	group.POST("/task/taskDone", t.taskDone)

	a := NewAccount()
	group.POST("/account", a.account)
	d := NewDepartment()
	group.POST("/department", d.department)
	group.POST("/department/save", d.save)
	group.POST("/department/read", d.read)

	auth := NewAuth()
	group.POST("/auth", auth.authList)
}
