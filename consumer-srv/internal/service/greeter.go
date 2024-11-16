package service

import (
	"consumer-srv/internal/biz"
	"context"
	v1 "kratos_first/consumer-srv-interface/api/helloworld/v1"
	"log"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
)

type GreeterService struct {
	v1.UnimplementedGreeterServer
	uc *biz.GreeterUsecase
}

func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

func (s *GreeterService) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500"
	consulClient, err := api.NewClient(consulConfig)
	//获取服务发现管理器
	dis := consul.New(consulClient)
	if err != nil {
		log.Fatal(err)
	}
	//连接目标grpc服务器
	endpoint := "discovery:///provider-srv"
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(dis),
	)

	if err != nil {
		return &v1.HelloReply{
			Message: "provider-svr is down",
		}, nil
	}

	client := v1.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &v1.HelloRequest{Name: "QurreChan"})
	if err != nil {
		return &v1.HelloReply{
			Message: "grpc call of provider-svr is failed",
		}, nil
	}

	return &v1.HelloReply{Message: "grpc call from  consumer to provider: " + resp.Message}, nil
}
