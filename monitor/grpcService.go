package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pb"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PrometheusApi interface {
	GetUserContainerResourceUsageRequest(ctx context.Context, userId string, fromTimeIn *timestamppb.Timestamp,
	) (*domain.Prometheus, error)
	GetMetricsByServiceId(ctx context.Context, serviceId string, fromTimeIn *timestamppb.Timestamp) (*domain.Prometheus, error)
}

type ContainerRepository interface {
	Get(ctx context.Context, cId string) (*domain.Container, error)
	GetAllUserContainer(ctx context.Context, userId string) (*[]domain.Container, error)
}

type MonitorServerImpl struct {
	pb.UnimplementedMonitorServiceServer
	prome         PrometheusApi
	containerRepo ContainerRepository
}

func NewMonitorServer(prome PrometheusApi, cRepo ContainerRepository) *MonitorServerImpl {
	return &MonitorServerImpl{prome: prome, containerRepo: cRepo}
}

// grpc service buat billing service
// intinya dapetin metrics usage untuk keseluruhan container user dan metrics setiap contianer yang dipunya user
func (server *MonitorServerImpl) GetAllUserContainerResourceUsage(
	ctx context.Context,
	req *pb.GetUserContainerResourceUsageRequest,
) (*pb.GetAllUserContainerResourceUsageResponse, error) {

	userId := req.GetUserId()
	fromTime := req.GetFromTime()

	promeQueryRes, err := server.prome.GetUserContainerResourceUsageRequest(ctx, userId, fromTime)
	if err != nil {
		zap.L().Error("server.prome.GetUserContainerResourceUsageRequest", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "server.prome.GetUSerContainerRequest: %v", err)
	}

	var usersContainer []*pb.Container
	allUserCtr, err := server.containerRepo.GetAllUserContainer(ctx, userId)
	if err != nil {
		zap.L().Error("server.containerRepo.GetAllUserContainer", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "User tidak pernah membuat container ke dogker %v", err)
	}

	for _, ctr := range *allUserCtr {
		ctrMetrics, err := server.prome.GetMetricsByServiceId(ctx, ctr.ServiceId, fromTime)
		if err != nil {
			zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "Gagal mendapatkan metrics dari container %v", err)
		}
		var ctLifecycles []*pb.ContainerLifeCycles
		for _, life := range ctr.ContainerLifecycles {
			ctLifecycles = append(ctLifecycles, &pb.ContainerLifeCycles{
				Id:          life.ID.String(),
				ContainerId: ctr.ID.String(),
				StartTime:   timestamppb.New(life.StartTime),
				StopTime:    timestamppb.New(life.StopTime),

				Replica: life.Replica,
				Status:  life.Status.String(),
			})
		}

		usersContainer = append(usersContainer, &pb.Container{
			Id:                     ctr.ID.String(),
			ImageUrl:               ctr.ImageUrl,
			Status:                 pb.ContainerStatus(ctr.Status),
			Name:                   ctr.Name,
			ContainerPort:          uint64(ctr.ContainerPort),
			PublicPort:             uint64(ctr.PublicPort),
			CreatedTime:            timestamppb.New(ctr.CreatedTime),
			CpuUsage:               ctrMetrics.AllCpuUsage,
			MemoryUsage:            ctrMetrics.AllMemoryUsage,
			NetworkIngressUsage:    ctrMetrics.AllNetworkIngressUsage,
			NetworkEgressUsage:     ctrMetrics.AllNetworkEgressUsage,
			ServiceId:              ctr.ServiceId,
			TerminatedTime:         timestamppb.New(ctr.TerminatedTime),
			AllContainerLifecycles: ctLifecycles,
		})

	}

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

func (server *MonitorServerImpl) GetSpecificContainerResourceUsage(
	ctx context.Context,
	req *pb.GetSpecificContainerResourceUsageRequest,
) (*pb.GetSpecificContainerResourceUsageResponse, error) {
	userId := req.UserId
	fromTime := req.FromTime
	containerId := req.ContainerId

	ctr, err := server.containerRepo.Get(ctx, containerId)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "container not found %v", err)
	}
	if ctr.UserId.String() != userId {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("anda bukan pemilik container dg id %v", containerId))
	}

	ctrMetrics, err := server.prome.GetMetricsByServiceId(ctx, ctr.ServiceId, fromTime)
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
			Status:  life.Status.String(),
		})
	}

	res := &pb.GetSpecificContainerResourceUsageResponse{

		CurrentTime: timestamppb.New(time.Now()),
		UserContainer: &pb.Container{
			Id:                     ctr.ID.String(),
			ImageUrl:               ctr.ImageUrl,
			Status:                 pb.ContainerStatus(ctr.Status),
			Name:                   ctr.Name,
			ContainerPort:          uint64(ctr.ContainerPort),
			PublicPort:             uint64(ctr.PublicPort),
			CreatedTime:            timestamppb.New(ctr.CreatedTime),
			CpuUsage:               ctrMetrics.AllCpuUsage,
			MemoryUsage:            ctrMetrics.AllMemoryUsage,
			NetworkIngressUsage:    ctrMetrics.AllNetworkIngressUsage,
			NetworkEgressUsage:     ctrMetrics.AllNetworkEgressUsage,
			ServiceId:              ctr.ServiceId,
			TerminatedTime:         timestamppb.New(ctr.TerminatedTime),
			AllContainerLifecycles: ctLifes,
		},
		FromTime: fromTime,
	}

	return res, nil

}
