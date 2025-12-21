package dao

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gwh.com/project-project/internal/data"
	"gwh.com/project-project/internal/database"
	"gwh.com/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p *ProjectDao) FindProjectByIds(ctx context.Context, pids []int64) (list []*data.Project, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&data.Project{}).Where("id in (?)", pids).Find(&list).Error
	return
}

func (p *ProjectDao) FindProjectById(ctx context.Context, projectCode int64) (pj *data.Project, err error) {
	err = p.conn.Session(ctx).Where("id=?", projectCode).Find(&pj).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (p *ProjectDao) FindMemberInfoByProjectCode(ctx context.Context, projectCode int64, page int64, size int64) ([]*data.ProjectMemberInfo, int64, error) {
	sql := "select a.project_code,a.member_code,a.is_owner,b.`name`,b.avatar,b.email from ms_project_member a, ms_member b where a.member_code=b.id and project_code = ? limit ?,?"
	session := p.conn.Session(ctx)
	db := session.Raw(sql, projectCode, (page-1)*size, size)
	var mis []*data.ProjectMemberInfo
	err := db.Scan(&mis).Error
	var total int64
	sqlCount := "select count(*) from ms_project_member a, ms_member b where a.member_code=b.id and project_code = ?"
	dbCount := session.Raw(sqlCount, projectCode)
	err = dbCount.Scan(&total).Error
	return mis, total, err
}

func (p *ProjectDao) UpdateProject(ctx context.Context, proj *data.Project) error {
	return p.conn.Session(ctx).Updates(proj).Error
}

func (p *ProjectDao) DeleteProjectCollect(ctx context.Context, memId int64, projectCode int64) error {
	return p.conn.Session(ctx).Where("member_code=? and project_code=?", memId, projectCode).Delete(&data.CollectionProject{}).Error
}

func (p *ProjectDao) SaveProjectCollect(ctx context.Context, pc *data.CollectionProject) error {

	return p.conn.Session(ctx).Save(&pc).Error
}

func (p *ProjectDao) UpdateDeleteProject(ctx context.Context, code int64, deleted bool) error {
	var err error
	if deleted {
		err = p.conn.Session(ctx).Model(&data.Project{}).Where("id=?", code).Update("deleted", 1).Error
	} else {
		err = p.conn.Session(ctx).Model(&data.Project{}).Where("id=?", code).Update("deleted", 0).Error
	}
	return err
}

func (p *ProjectDao) FindProjectByPIdAndMemId(ctx context.Context, projectCode int64, memberId int64) (*data.ProjectAndMember, error) {
	var pms *data.ProjectAndMember
	session := p.conn.Session(ctx)
	sql := fmt.Sprintf("select a.*,b.project_code,b.member_code,b.join_time,b.is_owner,b.authorize from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code=? and b.project_code=? limit 1")
	raw := session.Raw(sql, memberId, projectCode)
	err := raw.Scan(&pms).Error
	return pms, err
}

func (p *ProjectDao) FindCollectByPidAndMemId(ctx context.Context, projectCode int64, memberId int64) (bool, error) {
	var count int64
	session := p.conn.Session(ctx)
	sql := fmt.Sprintf("select count(*) from ms_project_collection where member_code=? and project_code=?")
	raw := session.Raw(sql, memberId, projectCode)
	err := raw.Scan(&count).Error
	return count > 0, err
}

func (p *ProjectDao) SaveProject(ctx context.Context, conn database.DBConn, pr *data.Project) error {
	p.conn = conn.(*gorms.GormConn)
	return p.conn.Tx(ctx).Save(pr).Error
}

func (p *ProjectDao) SaveProjectMember(ctx context.Context, conn database.DBConn, pm *data.MemberProject) error {
	p.conn = conn.(*gorms.GormConn)
	return p.conn.Tx(ctx).Save(pm).Error
}

func (p *ProjectDao) FindCollectProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*data.ProjectAndMember, int64, error) {
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id in (select project_code from ms_project_collection where member_code=? ) and a.id = b.project_code order by sort limit ?,?")

	db := session.Raw(sql, memId, index, size)
	var mp []*data.ProjectAndMember
	err := db.Scan(&mp).Error
	var total int64
	query := fmt.Sprintf("member_code=?")
	session.Model(&data.CollectionProject{}).Where(query, memId).Count(&total)
	return mp, total, err
}

func (p ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, condition string, page int64, size int64) ([]*data.ProjectAndMember, int64, error) {
	var pms []*data.ProjectAndMember
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code=? %s order by sort limit ?,?", condition)
	raw := session.Raw(sql, memId, index, size)
	raw.Scan(&pms)
	var total int64
	query := fmt.Sprintf("select count(*) from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code=? %s", condition)
	tx := session.Raw(query, memId)
	err := tx.Scan(&total).Error
	return pms, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.New(),
	}
}
