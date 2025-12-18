package midd

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gwh.com/project-api/api/rpc"
	common "gwh.com/project-common"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/user/login"
)

func TokenVerify() func(c *gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		token := c.GetHeader("Authorization")
		//验证用户是否已经登录
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		member, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token})
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}
		c.Set("memberId", member.Member.Id)
		c.Set("memberName", member.Member.Name)
		c.Set("organizationCode", member.Member.OrganizationCode)
		c.Next()
	}
}
