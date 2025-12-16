package tran

import "gwh.com/project-project/internal/database"

type Transaction interface {
	Action(func(conn database.DBConn) error) error
}
