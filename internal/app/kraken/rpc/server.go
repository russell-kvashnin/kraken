package rpc

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/model/config"
	"github.com/russell-kvashnin/kraken/pkg/api/v1"
	"google.golang.org/grpc"
	"net"
)

// gRPC server struct
type Server struct {
	grpc *grpc.Server
	cfg  config.RpcConfig

	mirroring v1.FileMirroringServer
}

// gRPC server constructor
func NewServer(cfg config.RpcConfig, mirroring v1.FileMirroringServer) *Server {
	srv := new(Server)
	srv.cfg = cfg
	srv.mirroring = mirroring

	return srv
}

// Configure gRPC server
func (server *Server) Configure() error {
	server.grpc = grpc.NewServer()
	v1.RegisterFileMirroringServer(server.grpc, server.mirroring)

	return nil
}

// Run gRPC server
func (server *Server) Run() error {
	listen, err := net.Listen("tcp", server.cfg.GetRpcListenAddress())
	if err != nil {
		return err
	}

	err = server.grpc.Serve(listen)
	if err != nil {
		return err
	}

	return nil
}

// Stop gRPC server
func (server *Server) Stop() error {
	server.grpc.Stop()

	return nil
}
