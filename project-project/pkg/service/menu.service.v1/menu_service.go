package menu_service_v1

import (
	"context"

	"github.com/jinzhu/copier"
	"gwh.com/project-common/errs"
	"gwh.com/project-grpc/menu"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/database/tran"
	"gwh.com/project-project/internal/domain"
	"gwh.com/project-project/internal/repo"
)

type MenuService struct {
	menu.UnimplementedMenuServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuDomain  *domain.MenuDomain
}

func New() *MenuService {
	return &MenuService{
		cache:       dao.Rc,
		transaction: dao.NewTransactionImpl(),
		menuDomain:  domain.NewMenuDomain(),
	}
}

func (d *MenuService) MenuList(ctx context.Context, msg *menu.MenuReqMessage) (*menu.MenuResponseMessage, error) {
	list, err := d.menuDomain.MenuTreeList()
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var mList []*menu.MenuMessage
	copier.Copy(&mList, list)
	return &menu.MenuResponseMessage{List: mList}, nil
}
