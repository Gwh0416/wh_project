package repo

import (
	"context"

	"gwh.com/project-project/internal/data"
	"gwh.com/project-project/internal/database"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, ids []int) ([]data.MsTaskStagesTemplate, error)
	FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) ([]*data.MsTaskStagesTemplate, error)
}

type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn database.DBConn, stages *data.TaskStages) error
	FindStagesByProjectId(ctx context.Context, projectCode int64, page int64, size int64) ([]*data.TaskStages, int64, error)
	FindById(ctx context.Context, id int) (ts *data.TaskStages, err error)
}

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int) (list []*data.Task, err error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *data.TaskMember, err error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error)
	SaveTask(ctx context.Context, conn database.DBConn, ts *data.Task) error
	SaveTaskMember(ctx context.Context, conn database.DBConn, tm *data.TaskMember) error
	FindTaskById(ctx context.Context, taskCode int64) (ts *data.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DBConn, ts *data.Task) error
	FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *data.Task, err error)
	FindTaskByAssignTo(ctx context.Context, memberId int64, done int, page int64, size int64) ([]*data.Task, int64, error)
	FindTaskByMemberCode(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*data.Task, total int64, err error)
	FindTaskByCreateBy(ctx context.Context, memberId int64, done int, page int64, size int64) (tList []*data.Task, total int64, err error)
	FindTaskMemberPage(ctx context.Context, taskCode int64, page int64, size int64) (list []*data.TaskMember, total int64, err error)
	FindTaskByIds(ctx context.Context, taskIdList []int64) (list []*data.Task, err error)
	TaskDone(ctx context.Context, taskCode int64, done int32) error
}
