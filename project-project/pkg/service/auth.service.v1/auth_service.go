package auth_service_v1

import (
	"context"

	"github.com/jinzhu/copier"
	"gwh.com/project-common/encrypts"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/auth"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/database/tran"
	"gwh.com/project-project/internal/domain"
	"gwh.com/project-project/internal/repo"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer
	cache             repo.Cache
	transaction       tran.Transaction
	projectAuthDomain *domain.ProjectAuthDomain
}

func New() *AuthService {
	return &AuthService{
		cache:             dao.Rc,
		transaction:       dao.NewTransactionImpl(),
		projectAuthDomain: domain.NewProjectAuthDomain(),
	}
}

func (a *AuthService) AuthList(ctx context.Context, msg *auth.AuthReqMessage) (*auth.ListAuthMessage, error) {
	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
	listPage, total, err := a.projectAuthDomain.AuthListPage(organizationCode, msg.Page, msg.PageSize)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var prList []*auth.ProjectAuth
	copier.Copy(&prList, listPage)
	return &auth.ListAuthMessage{List: prList, Total: total}, nil
}
