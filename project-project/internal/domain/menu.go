package domain

import (
	"context"
	"time"

	"gwh.com/project-common/errs"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/data"
	"gwh.com/project-project/internal/repo"
	"gwh.com/project-project/pkg/model"
)

type MenuDomain struct {
	menuRepo repo.MenuRepo
}

func (d *MenuDomain) MenuTreeList() ([]*data.ProjectMenuChild, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	menus, err := d.menuRepo.FindMenus(c)
	if err != nil {
		return nil, model.DBError
	}
	menuChildren := data.CovertChild(menus)
	return menuChildren, nil
}

func NewMenuDomain() *MenuDomain {
	return &MenuDomain{
		menuRepo: dao.NewMenuDao(),
	}
}
