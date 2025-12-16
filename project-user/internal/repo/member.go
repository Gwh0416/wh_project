package repo

import (
	"context"

	"gwh.com/project-user/internal/data/member"
	"gwh.com/project-user/internal/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, name string) (bool, error)
	GetMemberByMobile(ctx context.Context, phone string) (bool, error)
	SaveMember(conn database.DBConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account, pwd string) (*member.Member, error)
	FindMemberById(background context.Context, id int64) (*member.Member, error)
}
