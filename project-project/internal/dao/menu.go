package dao

import (
	"context"

	"gwh.com/project-project/internal/data/menu"
	"gwh.com/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func (m *MenuDao) FindMenus(ctx context.Context) (pms []*menu.ProjectMenu, err error) {
	err = m.conn.Session(ctx).Find(&pms).Error
	return pms, err
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		conn: gorms.New(),
	}
}
