package rpc

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"gwh.com/project-api/config"
	"gwh.com/project-common/discovery"
	"gwh.com/project-common/logs"
	"gwh.com/project-grpc/account"
	"gwh.com/project-grpc/auth"
	"gwh.com/project-grpc/department"
	"gwh.com/project-grpc/menu"
	"gwh.com/project-grpc/project"
	"gwh.com/project-grpc/task"
)

var ProjectServiceClient project.ProjectServiceClient

var TaskServiceClient task.TaskServiceClient

var AccountServiceClient account.AccountServiceClient

var DepartmentServiceClient department.DepartmentServiceClient

var AuthServiceClient auth.AuthServiceClient

var MenuServiceClient menu.MenuServiceClient

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.NewClient("etcd:///project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
	TaskServiceClient = task.NewTaskServiceClient(conn)
	AccountServiceClient = account.NewAccountServiceClient(conn)
	DepartmentServiceClient = department.NewDepartmentServiceClient(conn)
	AuthServiceClient = auth.NewAuthServiceClient(conn)
	MenuServiceClient = menu.NewMenuServiceClient(conn)

}
