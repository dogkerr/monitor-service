package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pb"
	"time"

	"go.uber.org/zap"
)

type PrometheusApi interface {
	GetUserContainerResourceUsageRequest(ctx context.Context, userId string, fromTimeIn string) (domain.Prometheus, error)
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
	}

	var usersContainer []*pb.Container
	var ctLifecycles []*pb.ContainerLifeCycles
	ctLifecycles = append(ctLifecycles, &pb.ContainerLifeCycles{
		Id:          "asd",
		ContainerId: "12",
		StartTimme:  "120",
		StopTime:    "12",
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
		CreatedTime:            time.Now().Format(time.RFC3339),
		CpuUsage:               0,
		MemoryUsage:            0,
		NetworkIngressUsage:    0,
		NetworkEgressUsage:     0,
		AllContainerLifecycles: ctLifecycles,
	})

	res := &pb.GetAllUserContainerResourceUsageResponse{
		CurrentTime:            promeQueryRes.CurrentTime.Format(time.RFC3339),
		AllCpuUsage:            promeQueryRes.AllCpuUsage,
		AllMemoryUsage:         promeQueryRes.AllMemoryUsage,
		AllNetworkIngressUsage: promeQueryRes.AllNetworkIngressUsage,
		AllNetworkEgressUsage:  promeQueryRes.AllNetworkEgressUsage,
		UserContainer:          usersContainer,
		FromTime:               promeQueryRes.FromTime.Format(time.RFC3339),
	}
	return res, nil

}

func (server *MonitorServerImpl) GetSpecificContainerResourceUsage(ctx context.Context, req *pb.GetSpecificContainerResourceUsageRequest) (*pb.GetSpecificContainerResourceUsageResponse, error) {
	return &pb.GetSpecificContainerResourceUsageResponse{}, nil
}


/*

	var usersContainer []*pb.Container
	var ctLifecycles []*pb.ContainerLifeCycles
	ctLifecycles = append(ctLifecycles, &pb.ContainerLifeCycles{
		Id:          "asd",
		ContainerId: "12",
		StartTimme:  "120",
		StopTime:    "12",
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
		CreatedTime:            time.Now().String(),
		CpuUsage:               0,
		MemoryUsage:            0,
		NetworkIngressUsage:    0,
		NetworkEgressUsage:     0,
		AllContainerLifecycles: ctLifecycles,
	})
*/
