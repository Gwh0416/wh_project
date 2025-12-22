package repo

import (
	"context"

	"gwh.com/project-project/internal/data"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*data.ProjectMenu, error)
}
