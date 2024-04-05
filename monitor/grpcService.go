package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PrometheusApi interface {
	GetUserContainerResourceUsageRequest(ctx context.Context, userId string, fromTimeIn *timestamppb.Timestamp,
	) (domain.Prometheus, error)
}

type MonitorServerImpl struct {
	pb.UnimplementedMonitorServiceServer
	prome PrometheusApi
}

func NewMonitorServer(prome PrometheusApi) *MonitorServerImpl {
	return &MonitorServerImpl{prome: prome}
}

func (server *MonitorServerImpl) GetAllUserContainerResourceUsage(
	ctx context.Context,
	req *pb.GetUserContainerResourceUsageRequest,
) (*pb.GetAllUserContainerResourceUsageResponse, error) {

	userId := req.GetUserId()
	fromTime := req.GetFromTime()

	promeQueryRes, err := server.prome.GetUserContainerResourceUsageRequest(ctx, userId, fromTime)
	if err != nil {
		zap.L().Error("gagal query prometheus!", zap.Error(err))
		return &pb.GetAllUserContainerResourceUsageResponse{}, status.Errorf(codes.InvalidArgument, "fromTime is not in the right format (RFC3399): %v", err)
	}
	now := time.Now()
	start := now.AddDate(0, -1, 0)

	var usersContainer []*pb.Container
	var ctLifecycles []*pb.ContainerLifeCycles
	ctLifecycles = append(ctLifecycles, &pb.ContainerLifeCycles{
		Id:          "asd",
		ContainerId: "12",
		StartTime:   timestamppb.New(start),
		StopTime:    timestamppb.Now(),
		CpuCore:     4.0,
		MemCapacity: 1.0,
		Replica:     2,
	})
	usersContainer = append(usersContainer, &pb.Container{
		Id:                     "tes",
		ImageUrl:               "tes",
		Status:                 pb.ContainerStatus_RUN,
		Name:                   "tes",
		ContainerPort:          1000,
		PublicPort:             1000,
		CreatedTime:            timestamppb.Now(),
		CpuUsage:               0,
		MemoryUsage:            0,
		NetworkIngressUsage:    0,
		NetworkEgressUsage:     0,
		AllContainerLifecycles: ctLifecycles,
	})

	res := &pb.GetAllUserContainerResourceUsageResponse{
		CurrentTime:            promeQueryRes.CurrentTime,
		AllCpuUsage:            promeQueryRes.AllCpuUsage,
		AllMemoryUsage:         promeQueryRes.AllMemoryUsage,
		AllNetworkIngressUsage: promeQueryRes.AllNetworkIngressUsage,
		AllNetworkEgressUsage:  promeQueryRes.AllNetworkEgressUsage,
		UserContainer:          usersContainer,
		FromTime:               promeQueryRes.FromTime,
	}
	return res, nil

}

func (server *MonitorServerImpl) GetSpecificContainerResourceUsage(ctx context.Context, req *pb.GetSpecificContainerResourceUsageRequest) (*pb.GetSpecificContainerResourceUsageResponse, error) {
	return &pb.GetSpecificContainerResourceUsageResponse{}, nil
}
