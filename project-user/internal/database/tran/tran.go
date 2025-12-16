package tran

import "gwh.com/project-user/internal/database"

type Transaction interface {
	Action(func(conn database.DBConn) error) error
}
