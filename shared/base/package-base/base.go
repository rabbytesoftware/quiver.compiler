package packagebase

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/rabbytesoftware/quiver.compiler/shared/package"
	"google.golang.org/grpc"
)

// PackageBase is the base implementation for all packages
type PackageBase struct {
	// Define the methods that must be overridden by implementers
	Impl	PackageImplementor
	
	// gRPC server
	server	*grpc.Server
	port	int
	pb.UnimplementedPackageServiceServer
}

// RunPackage starts the gRPC server for a package
func RunPackage(impl PackageImplementor) {
	// Check if we have enough command-line arguments
	if len(os.Args) < 2 {
		log.Fatalf("Expected port number as argument")
	}
	
	// Parse the port number
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	
	// Create the package base
	packageObj := &PackageBase{
		Impl: impl,
		port: port,
	}
	
	// Start the gRPC server
	if err := packageObj.Start(); err != nil {
		log.Fatalf("Failed to start package: %v", err)
	}
}

// Start runs the gRPC server
func (p *PackageBase) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", p.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	
	// Create a new gRPC server
	p.server = grpc.NewServer()
	pb.RegisterPackageServiceServer(p.server, p)
	
	// Start serving
	// log.Printf("Package %s v%s starting on port %d", p.Impl.GetName(), p.Impl.GetVersion(), p.port)
	return p.server.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (p *PackageBase) GracefulStop() {
	if p.server != nil {
		p.server.GracefulStop()
	}
}
