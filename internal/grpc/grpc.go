package grpc

import (
	"dogker/lintang/monitor-service/pb"
	"net"

	"google.golang.org/grpc"
)


func RunGRPCServer (
	monitorServer pb.MonitorServiceServer,
	listener net.Listener,
) error {
// GRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterMonitorServiceServer(grpcServer, monitorServer)
	return grpcServer.Serve(listener)
}

