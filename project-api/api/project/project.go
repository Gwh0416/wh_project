package project

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gwh.com/project-api/api/rpc"
	"gwh.com/project-api/pkg/model"
	"gwh.com/project-api/pkg/model/menu"
	"gwh.com/project-api/pkg/model/pro"
	common "gwh.com/project-common"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/project"
)

type HandlerProject struct {
}

func NewHandlerProject() *HandlerProject {
	return &HandlerProject{}
}

func (h *HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := rpc.ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var ms []*menu.Menu
	err = copier.Copy(&ms, indexResponse.Menus)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	c.JSON(http.StatusOK, result.Success(ms))
}

func (h *HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		Page:       page.Page,
		PageSize:   page.PageSize,
		MemberName: memberName,
		SelectBy:   selectBy,
	}
	myProjectResponse, err := rpc.ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectAndMember
	err = copier.Copy(&pms, myProjectResponse.Pm)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": myProjectResponse.Total,
	}))

}

func (h *HandlerProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	organizationCode := c.GetString("organizationCode")
	var page = &model.Page{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	projectTemplateRsp, err := rpc.ProjectServiceClient.FindProjectTemplate(ctx,
		&project.ProjectRpcMessage{
			MemberId:         memberId,
			MemberName:       memberName,
			OrganizationCode: organizationCode,
			Page:             page.Page,
			PageSize:         page.PageSize,
			ViewType:         int32(viewType)})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var pts []*pro.ProjectTemplate
	err = copier.Copy(&pts, projectTemplateRsp.Ptm)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	if pts == nil {
		pts = []*pro.ProjectTemplate{}
	}
	for _, pt := range pts {
		if pt.TaskStages == nil {
			pt.TaskStages = []*pro.TaskStagesOnlyName{}
		}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pts,
		"total": projectTemplateRsp.Total,
	}))
}

func (h *HandlerProject) projectSave(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	organizationCode := c.GetString("organizationCode")
	var req *pro.SaveProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		Name:             req.Name,
		OrganizationCode: organizationCode,
		Description:      req.Description,
		TemplateCode:     req.TemplateCode,
		Id:               int64(req.Id)}
	saveProjectMessage, err := rpc.ProjectServiceClient.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	sp := &pro.SaveProject{}
	err = copier.Copy(&sp, saveProjectMessage)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	c.JSON(http.StatusOK, result.Success(sp))
}

func (h *HandlerProject) readProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detail, err := rpc.ProjectServiceClient.FindProjectDetail(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	pd := &pro.ProjectDetail{}
	err = copier.Copy(pd, detail)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy参数格式有误"))
		return
	}
	c.JSON(http.StatusOK, result.Success(pd))
}

func (h *HandlerProject) recycleProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateDeleteProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: true})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (h *HandlerProject) recoveryProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateDeleteProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, Deleted: false})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) collectProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	collectType := c.PostForm("type")
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := rpc.ProjectServiceClient.UpdateCollectProject(ctx, &project.ProjectRpcMessage{ProjectCode: projectCode, CollectType: collectType, MemberId: memberId})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) editProject(c *gin.Context) {
	result := &common.Result{}
	var req *pro.ProjectUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	memberId := c.GetInt64("memberId")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.UpdateProjectMessage{}
	copier.Copy(msg, req)
	msg.MemberId = memberId
	_, err := rpc.ProjectServiceClient.UpdateProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{}))
}

func (p *HandlerProject) getLogBySelfProject(c *gin.Context) {
	result := &common.Result{}
	var page = &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.ProjectRpcMessage{
		MemberId: c.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	projectLogResponse, err := rpc.ProjectServiceClient.GetLogBySelfProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var list []*model.ProjectLog
	copier.Copy(&list, projectLogResponse.List)
	if list == nil {
		list = []*model.ProjectLog{}
	}
	c.JSON(http.StatusOK, result.Success(list))
}
