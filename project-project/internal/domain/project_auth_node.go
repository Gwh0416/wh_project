package domain

import (
	"context"
	"fmt"

	"gwh.com/project-common/errs"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/database"
	"gwh.com/project-project/internal/repo"
	"gwh.com/project-project/pkg/model"
)

type ProjectAuthNodeDomain struct {
	projectAuthNodeRepo repo.ProjectAuthNodeRepo
}

func NewProjectAuthNodeDomain() *ProjectAuthNodeDomain {
	return &ProjectAuthNodeDomain{
		projectAuthNodeRepo: dao.NewProjectAuthNodeDao(),
	}
}

func (d *ProjectAuthNodeDomain) AuthNodeList(authId int64) ([]string, *errs.BError) {
	list, err := d.projectAuthNodeRepo.FindNodeStringList(context.Background(), authId)
	if err != nil {
		return nil, model.DBError
	}
	return list, nil
}

func (d *ProjectAuthNodeDomain) Save(conn database.DBConn, authId int64, nodes []string) *errs.BError {
	fmt.Println(nodes)
	err := d.projectAuthNodeRepo.DeleteByAuthId(context.Background(), conn, authId)
	if err != nil {
		return model.DBError
	}
	err = d.projectAuthNodeRepo.Save(context.Background(), conn, authId, nodes)
	if err != nil {
		return model.DBError
	}
	return nil
}
