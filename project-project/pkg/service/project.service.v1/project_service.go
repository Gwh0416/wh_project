package project_service_v1

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/project"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/data/menu"
	"gwh.com/project-project/internal/database/tran"
	"gwh.com/project-project/internal/repo"
	"gwh.com/project-project/pkg/model"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuRepo    repo.MenuRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		transaction: dao.NewTransactionImpl(),
		menuRepo:    dao.NewMenuDao(),
	}
}

func (p *ProjectService) Index(context.Context, *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index db FindMenus err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	childs := menu.CovertChild(pms)
	var mms []*project.MenuMessage
	err = copier.Copy(&mms, childs)
	if err != nil {
		return nil, errors.New("copy参数格式有误")
	}
	return &project.IndexResponse{
		Menus: mms,
	}, nil
}
