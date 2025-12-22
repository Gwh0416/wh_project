package project

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gwh.com/project-api/api/rpc"
	"gwh.com/project-api/pkg/model"
	common "gwh.com/project-common"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/auth"
)

type HandlerAuth struct {
}

func (a *HandlerAuth) authList(c *gin.Context) {
	result := &common.Result{}
	organizationCode := c.GetString("organizationCode")
	var page = &model.Page{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		OrganizationCode: organizationCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
	}
	response, err := rpc.AuthServiceClient.AuthList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	}))
}

func (a *HandlerAuth) apply(c *gin.Context) {
	result := &common.Result{}
	var req *model.ProjectAuthReq
	c.ShouldBind(&req)
	var nodes []string
	if req.Nodes != "" {
		json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthReqMessage{
		Action: req.Action,
		AuthId: req.Id,
		Nodes:  nodes,
	}
	applyResponse, err := rpc.AuthServiceClient.Apply(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var list []*model.ProjectNodeAuthTree
	copier.Copy(&list, applyResponse.List)
	var checkedList []string
	copier.Copy(&checkedList, applyResponse.CheckedList)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":        list,
		"checkedList": checkedList,
	}))
}

func (a *HandlerAuth) GetAuthNodes(c *gin.Context) ([]string, error) {
	memberId := c.GetInt64("memberId")
	msg := &auth.AuthReqMessage{
		MemberId: memberId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.AuthServiceClient.AuthNodesByMemberId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		return nil, errs.NewError(errs.ErrorCode(code), msg)
	}
	return response.List, err
}

func NewAuth() *HandlerAuth {
	return &HandlerAuth{}
}
