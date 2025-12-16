package project

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	indexResponse, err := ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success(indexResponse.Menus))
}
