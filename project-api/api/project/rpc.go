package project

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"gwh.com/project-api/config"
	"gwh.com/project-common/discovery"
	"gwh.com/project-common/logs"
	"gwh.com/project-grpc/project"
)

var ProjectServiceClient project.ProjectServiceClient

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.NewClient("etcd:///project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
}
