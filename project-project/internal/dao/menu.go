package dao

import (
	"context"

	"gwh.com/project-project/internal/data"
	"gwh.com/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func (m *MenuDao) FindMenus(ctx context.Context) (pms []*data.ProjectMenu, err error) {
	err = m.conn.Session(ctx).Order("pid,sort asc, id asc").Find(&pms).Error
	return pms, err
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		conn: gorms.New(),
	}
}
