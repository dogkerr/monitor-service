package grpc

import (
	"dogker/lintang/monitor-service/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer(
	monitorServer pb.MonitorServiceServer,
	listener net.Listener,
	ch chan *grpc.Server,
) error {
	// GRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterMonitorServiceServer(grpcServer, monitorServer)
	reflection.Register(grpcServer)

	ch <- grpcServer

	return grpcServer.Serve(listener)
}
