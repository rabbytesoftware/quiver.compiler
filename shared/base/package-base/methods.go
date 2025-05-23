package packagebase

import (
	"context"

	pb "github.com/rabbytesoftware/quiver.compiler/shared/package"
)

// gRPC Service Implementation

func (p *PackageBase) SetPorts(ctx context.Context, req *pb.SetPortsRequest) (*pb.BoolResponse, error) {
	return &pb.BoolResponse{Success: p.Impl.SetPorts(req.Ports)}, nil
}

func (p *PackageBase) Install(ctx context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {
	return &pb.BoolResponse{Success: p.Impl.Install()}, nil
}

func (p *PackageBase) Run(ctx context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {
	return &pb.BoolResponse{Success: p.Impl.Run()}, nil
}

func (p *PackageBase) Exit(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	p.Impl.Exit()
	return &pb.Empty{}, nil
}
