package grpc

import (
	"context"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/kitex_gen/container-service/pb"
	"dogker/lintang/monitor-service/kitex_gen/container-service/pb/containergrpcservice"

	"time"

	"github.com/cloudwego/kitex/client"
	"go.uber.org/zap"
)

type ContainerClient struct {
	// service pb.ContainerGRPCServiceClient
	service containergrpcservice.Client
}

func NewContainerClient(cfg *config.Config) *ContainerClient {
	c, err := containergrpcservice.NewClient("containerGRPCService", client.WithHostPorts(cfg.GRPC.ContainerURL))
	if err != nil {
		zap.L().Fatal("containergrpcservice.NewClient")
	}
	// service := pb.NewContainerGRPCServiceClient(cc)
	return &ContainerClient{service: c}
}

func (c *ContainerClient) SendDownContainerToCtrService(ctx context.Context, terminatedInstances []string) error {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	zap.L().Info("serviceIDs", zap.Strings("serviceID", terminatedInstances))

	req := &pb.ContainerTerminatedAccidentallyReq{
		ServiceIDs: terminatedInstances,
	}

	if len(terminatedInstances) != 0 {
		// pastiin yang dikirim ke ctr service kalau beneran ada terminated instance
		_, err := c.service.ContainerTerminatedAccidentally(grpcCtx, req)
		if err != nil {
			zap.L().Error("c.service.ContainerTerminatedAccidentally (SendTerminatedContainerToCtrService) (containerGRPCServiceClient)", zap.Error(err))
			return err
		}
	}

	return nil

}
