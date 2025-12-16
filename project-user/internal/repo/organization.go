package repo

import (
	"context"

	"gwh.com/project-user/internal/data/organization"
	"gwh.com/project-user/internal/database"
)

type OrganizationRepo interface {
	SaveOrganization(conn database.DBConn, ctx context.Context, org *organization.Organization) error
	FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error)
}
