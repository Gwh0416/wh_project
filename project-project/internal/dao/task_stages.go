package dao

import (
	"context"

	"gwh.com/project-project/internal/data"
	"gwh.com/project-project/internal/database"
	"gwh.com/project-project/internal/database/gorms"
)

type TaskStagesDao struct {
	conn *gorms.GormConn
}

func (t *TaskStagesDao) FindById(ctx context.Context, id int) (ts *data.TaskStages, err error) {
	err = t.conn.Session(ctx).Where("id=?", id).Find(&ts).Error
	return
}

func (t *TaskStagesDao) FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, pageSize int64) (list []*data.TaskStages, total int64, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.TaskStages{}).
		Where("project_code=?", projectCode).
		Order("sort asc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	err = session.Model(&data.TaskStages{}).
		Where("project_code=?", projectCode).
		Count(&total).Error
	return
}

func (t *TaskStagesDao) SaveTaskStages(ctx context.Context, conn database.DBConn, stages *data.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	session := t.conn.Tx(ctx)
	return session.Save(&stages).Error
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{
		conn: gorms.New(),
	}
}
