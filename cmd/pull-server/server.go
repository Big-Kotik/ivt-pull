package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Big-Kotik/ivt-pull-api/pkg/api"
	"github.com/Big-Kotik/ivt-pull/pkg/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 7272))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterPullerServer(grpcServer, &server.PullServer{Client: http.Client{}, Logger: log.Default()})
	grpcServer.Serve(lis)
}
