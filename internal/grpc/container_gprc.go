package grpc

import (
	"context"
	"dogker/lintang/monitor-service/pb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ContainerClient struct {
	service pb.ContainerGRPCServiceClient
}

func NewContainerClient(cc *grpc.ClientConn) *ContainerClient {
	service := pb.NewContainerGRPCServiceClient(cc)
	return &ContainerClient{service: service}
}

func (c *ContainerClient) SendTerminatedContainerToCtrService(ctx context.Context, terminatedInstances []string) error {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	zap.L().Info("serviceIDs", zap.Strings("serviceID", terminatedInstances))

	req := &pb.ContainerTerminatedAccidentallyReq{
		ServiceIDs: terminatedInstances,
	}

	_, err := c.service.ContainerTerminatedAccidentally(grpcCtx, req)
	if err != nil {
		zap.L().Error("c.service.ContainerTerminatedAccidentally (SendTerminatedContainerToCtrService) (containerGRPCServiceClient)", zap.Error(err))
		return err
	}

	return nil

}
