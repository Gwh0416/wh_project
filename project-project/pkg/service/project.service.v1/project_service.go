package project_service_v1

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gwh.com/project-common/encrypts"
	"gwh.com/project-common/errs"
	"gwh.com/project-common/tms"
	"gwh.com/project-grpc/project"
	"gwh.com/project-grpc/user/login"
	"gwh.com/project-project/internal/dao"
	"gwh.com/project-project/internal/data/menu"
	"gwh.com/project-project/internal/data/pro"
	"gwh.com/project-project/internal/data/task"
	"gwh.com/project-project/internal/database"
	"gwh.com/project-project/internal/database/tran"
	"gwh.com/project-project/internal/repo"
	"gwh.com/project-project/internal/rpc"
	"gwh.com/project-project/pkg/model"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	menuRepo               repo.MenuRepo
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransactionImpl(),
		menuRepo:               dao.NewMenuDao(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
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

func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.MyProjectResponse, error) {
	membetId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	var pms []*pro.ProjectAndMember
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, membetId, "and a.deleted = 0 ", page, pageSize)
	}
	if msg.SelectBy == "archive" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, membetId, "and a.archive = 1 ", page, pageSize)
	}
	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, membetId, "and a.deleted = 1 ", page, pageSize)
	}
	if msg.SelectBy == "collect" {
		pms, total, err = p.projectRepo.FindCollectProjectByMemId(ctx, membetId, page, pageSize)
		for _, pm := range pms {
			pm.Collected = model.Collected
		}
	} else {
		collectPms, _, err := p.projectRepo.FindCollectProjectByMemId(ctx, membetId, page, pageSize)
		if err != nil {
			zap.L().Error("Project FindCollectProjectByMemId error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		var cMap = make(map[int64]*pro.ProjectAndMember)
		for _, pm := range collectPms {
			cMap[pm.Id] = pm
		}
		for _, pm := range pms {
			if cMap[pm.ProjectCode] != nil {
				pm.Collected = model.Collected
			}
		}
	}

	if err != nil {
		zap.L().Error("Project FindProjectByMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if pms == nil {
		return &project.MyProjectResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}

	var pmm []*project.ProjectMessage
	err = copier.Copy(&pmm, pms)
	for _, v := range pmm {
		v.Code, _ = encrypts.EncryptInt64(v.ProjectCode, model.AESKey)
		pam := pro.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &project.MyProjectResponse{Pm: pmm, Total: total}, err
}

func (p *ProjectService) FindProjectTemplate(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectTemplateResponse, error) {
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	var err error
	var pts []pro.ProjectTemplate
	var total int64
	if msg.ViewType == -1 {
		pts, total, err = p.projectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
	}
	if msg.ViewType == 0 {
		pts, total, err = p.projectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
	}
	if msg.ViewType == 1 {
		pts, total, err = p.projectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
	}
	if err != nil {
		zap.L().Error("FindProjectTemplate FindProjectTemplate err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	//拿到id列表， 去任务步骤模板表进行查询
	tsts, err := p.taskStagesTemplateRepo.FindInProTemIds(ctx, pro.ToProjectTemplateIds(pts))
	if err != nil {
		zap.L().Error("FindProjectTemplate FindInProTemIds err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var ptas []*pro.ProjectTemplateAll
	for _, pt := range pts {
		ptas = append(ptas, pt.Convert(task.CovertProjectMap(tsts)[pt.Id]))
	}

	//组装数据 返回
	var pmMsg []*project.ProjectTemplateMessage
	err = copier.Copy(&pmMsg, ptas)
	return &project.ProjectTemplateResponse{
		Ptm:   pmMsg,
		Total: total,
	}, err
}

func (ps *ProjectService) SaveProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.SaveProjectMessage, error) {
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKey)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	templateCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, model.AESKey)
	templateCode, _ := strconv.ParseInt(templateCodeStr, 10, 64)
	pr := &pro.Project{
		Name:              msg.Name,
		Description:       msg.Description,
		TemplateCode:      int(templateCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://gwh0416.oss-cn-hangzhou.aliyuncs.com/Typora/fix-dir%2F2025%2F12%2F%E5%A4%B4%E5%83%8F3-87a77db3afb7bcbe64f1fdb9177dc691.png",
		Deleted:           model.NoDeleted,
		Archive:           model.NoArchive,
		OrganizationCode:  organizationCode,
		AccessControlType: model.Open,
		TaskBoardTheme:    model.Simple,
	}
	var rsp *project.SaveProjectMessage
	err := ps.transaction.Action(func(conn database.DBConn) error {
		err := ps.projectRepo.SaveProject(ctx, conn, pr)
		if err != nil {
			zap.L().Error("SaveProject SaveProject error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		pm := &pro.MemberProject{
			ProjectCode: pr.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   "",
		}
		err = ps.projectRepo.SaveProjectMember(ctx, conn, pm)
		if err != nil {
			zap.L().Error("SaveProject SaveProjectMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		code, _ := encrypts.EncryptInt64(pr.Id, model.AESKey)
		rsp = &project.SaveProjectMessage{
			Id:               pr.Id,
			Code:             code,
			OrganizationCode: organizationCodeStr,
			Name:             pr.Name,
			Cover:            pr.Cover,
			CreateTime:       tms.FormatByMill(pr.CreateTime),
			TaskBoardTheme:   pr.TaskBoardTheme,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// 1. 查项目表
// 2. 项目和成员的关联表 查到项目的拥有者 去member表查名字
// 3. 查收藏表 判断收藏状态
func (ps *ProjectService) FindProjectDetail(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectDetailMessage, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	memberId := msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectAndMember, err := ps.projectRepo.FindProjectByPIdAndMemId(c, projectCode, memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindProjectByPIdAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	ownerId := projectAndMember.IsOwner
	member, err := rpc.LoginServiceClient.FindMemInfoById(c, &login.UserMessage{MemId: ownerId})
	if err != nil {
		zap.L().Error("project rpc FindProjectDetail FindMemInfoById error", zap.Error(err))
		return nil, err
	}
	//去user模块去找了
	//TODO 优化 收藏的时候 可以放入redis
	isCollect, err := ps.projectRepo.FindCollectByPidAndMemId(c, projectCode, memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindCollectByPidAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if isCollect {
		projectAndMember.Collected = model.Collected
	}
	var detailMsg = &project.ProjectDetailMessage{}
	copier.Copy(detailMsg, projectAndMember)
	detailMsg.OwnerAvatar = member.Avatar
	detailMsg.OwnerName = member.Name
	detailMsg.Code, _ = encrypts.EncryptInt64(projectAndMember.ProjectCode, model.AESKey)
	detailMsg.AccessControlType = projectAndMember.GetAccessControlType()
	detailMsg.OrganizationCode, _ = encrypts.EncryptInt64(projectAndMember.OrganizationCode, model.AESKey)
	detailMsg.Order = int32(projectAndMember.Sort)
	detailMsg.CreateTime = tms.FormatByMill(projectAndMember.CreateTime)
	return detailMsg, nil
}

func (ps *ProjectService) UpdateDeleteProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.DeletedProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := ps.projectRepo.UpdateDeleteProject(c, projectCode, msg.Deleted)
	if err != nil {
		zap.L().Error("project RecycleProject DeleteProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.DeletedProjectResponse{}, nil
}

func (ps *ProjectService) UpdateProject(ctx context.Context, msg *project.UpdateProjectMessage) (*project.UpdateProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	proj := &pro.Project{
		Id:                 projectCode,
		Name:               msg.Name,
		Description:        msg.Description,
		Cover:              msg.Cover,
		TaskBoardTheme:     msg.TaskBoardTheme,
		Prefix:             msg.Prefix,
		Private:            int(msg.Private),
		OpenPrefix:         int(msg.OpenPrefix),
		OpenBeginTime:      int(msg.OpenBeginTime),
		OpenTaskPrivate:    int(msg.OpenTaskPrivate),
		Schedule:           msg.Schedule,
		AutoUpdateSchedule: int(msg.AutoUpdateSchedule),
	}
	err := ps.projectRepo.UpdateProject(c, proj)
	if err != nil {
		zap.L().Error("project UpdateProject::UpdateProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.UpdateProjectResponse{}, nil
}
