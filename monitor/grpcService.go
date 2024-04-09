package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pb"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PrometheusAPI interface {
	GetUserContainerResourceUsageRequest(ctx context.Context, userID string, fromTimeIn *timestamppb.Timestamp) (*domain.Prometheus, error)
	GetMetricsByServiceID(ctx context.Context, serviceID string, fromTimeIn *timestamppb.Timestamp) (*domain.Metric, error)
}

type ContainerRepository interface {
	Get(ctx context.Context, serviceID string) (*domain.Container, error)
	GetAllUserContainer(ctx context.Context, userID string) (*[]domain.Container, error)
}

type MonitorServerImpl struct {
	pb.UnimplementedMonitorServiceServer
	prome         PrometheusAPI
	containerRepo ContainerRepository
}

func NewMonitorServer(prome PrometheusAPI, cRepo ContainerRepository) *MonitorServerImpl {
	return &MonitorServerImpl{prome: prome, containerRepo: cRepo}
}

// grpc service buat billing service
// intinya dapetin metrics usage untuk keseluruhan container user dan metrics setiap container yang dipunya user
func (server *MonitorServerImpl) GetAllUserContainerResourceUsage(
	ctx context.Context,
	req *pb.GetUserContainerResourceUsageRequest,
) (*pb.GetAllUserContainerResourceUsageResponse, error) {
	userID := req.GetUserId()
	fromTime := req.GetFromTime()

	promeQueryRes, err := server.prome.GetUserContainerResourceUsageRequest(ctx, userID, fromTime)
	if err != nil {
		zap.L().Error("server.prome.GetUserContainerResourceUsageRequest", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "server.prome.GetUSerContainerRequest: %v", err)
	}

	var usersContainer []*pb.Container
	allUserCtr, err := server.containerRepo.GetAllUserContainer(ctx, userID)
	if err != nil {
		zap.L().Error("server.containerRepo.GetAllUserContainer", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "User tidak pernah membuat container ke dogker %v", err)
	}

	for i := 0; i < len(*allUserCtr); i++ {
		ctr := *allUserCtr

		ctrMetrics, err := server.prome.GetMetricsByServiceID(ctx, ctr[i].ServiceID, fromTime)
		if err != nil {
			zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "Gagal mendapatkan metrics dari container %v", err)
		}
		var ctLifecycles []*pb.ContainerLifeCycles
		for _, life := range ctr[i].ContainerLifecycles {
			ctLifecycles = append(ctLifecycles, &pb.ContainerLifeCycles{
				Id:          life.ID.String(),
				ContainerId: ctr[i].ID.String(),
				StartTime:   timestamppb.New(life.StartTime),
				StopTime:    timestamppb.New(life.StopTime),

				Replica: life.Replica,
				Status:  pb.ContainerStatus(life.Status),
			})
		}

		usersContainer = append(usersContainer, &pb.Container{
			Id:                     ctr[i].ID.String(),
			ImageUrl:               ctr[i].ImageURL,
			Status:                 pb.ContainerStatus(ctr[i].Status),
			Name:                   ctr[i].Name,
			ContainerPort:          uint64(ctr[i].ContainerPort),
			PublicPort:             uint64(ctr[i].PublicPort),
			CreatedTime:            timestamppb.New(ctr[i].CreatedTime),
			CpuUsage:               ctrMetrics.CpuUsage,
			MemoryUsage:            ctrMetrics.MemoryUsage,
			NetworkIngressUsage:    ctrMetrics.NetworkIngressUsage,
			NetworkEgressUsage:     ctrMetrics.NetworkEgressUsage,
			ServiceId:              ctr[i].ServiceID,
			TerminatedTime:         timestamppb.New(ctr[i].TerminatedTime),
			AllContainerLifecycles: ctLifecycles,
		})
	}

	res := &pb.GetAllUserContainerResourceUsageResponse{
		CurrentTime:            promeQueryRes.CurrentTime,
		AllCpuUsage:            promeQueryRes.AllCPUUsage,
		AllMemoryUsage:         promeQueryRes.AllMemoryUsage,
		AllNetworkIngressUsage: promeQueryRes.AllNetworkIngressUsage,
		AllNetworkEgressUsage:  promeQueryRes.AllNetworkEgressUsage,
		UserContainer:          usersContainer,
		FromTime:               promeQueryRes.FromTime,
	}
	return res, nil
}

func (server *MonitorServerImpl) GetSpecificContainerResourceUsage(
	ctx context.Context,
	req *pb.GetSpecificContainerResourceUsageRequest,
) (*pb.GetSpecificContainerResourceUsageResponse, error) {
	userID := req.UserId
	fromTime := req.FromTime
	containerID := req.ContainerId

	ctr, err := server.containerRepo.Get(ctx, containerID)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "container not found %v", err)
	}
	if ctr.UserID.String() != userID {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("anda bukan pemilik container dg id %v", containerID))
	}

	ctrMetrics, err := server.prome.GetMetricsByServiceID(ctx, ctr.ServiceID, fromTime)
	if err != nil {
		zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Gagal mendapatkan metrics dari container %v", err)
	}

	var ctLifes []*pb.ContainerLifeCycles
	for _, life := range ctr.ContainerLifecycles {
		ctLifes = append(ctLifes, &pb.ContainerLifeCycles{
			Id:          life.ID.String(),
			ContainerId: ctr.ID.String(),
			StartTime:   timestamppb.New(life.StartTime),
			StopTime:    timestamppb.New(life.StopTime),

			Replica: life.Replica,
			Status:  pb.ContainerStatus(life.Status),
		})
	}

	res := &pb.GetSpecificContainerResourceUsageResponse{

		CurrentTime: timestamppb.New(time.Now()),
		UserContainer: &pb.Container{
			Id:                     ctr.ID.String(),
			ImageUrl:               ctr.ImageURL,
			Status:                 pb.ContainerStatus(ctr.Status),
			Name:                   ctr.Name,
			ContainerPort:          uint64(ctr.ContainerPort),
			PublicPort:             uint64(ctr.PublicPort),
			CreatedTime:            timestamppb.New(ctr.CreatedTime),
			CpuUsage:               ctrMetrics.CpuUsage,
			MemoryUsage:            ctrMetrics.MemoryUsage,
			NetworkIngressUsage:    ctrMetrics.NetworkIngressUsage,
			NetworkEgressUsage:     ctrMetrics.NetworkEgressUsage,
			ServiceId:              ctr.ServiceID,
			TerminatedTime:         timestamppb.New(ctr.TerminatedTime),
			AllContainerLifecycles: ctLifes,
		},
		FromTime: fromTime,
	}
	return res, nil
}
