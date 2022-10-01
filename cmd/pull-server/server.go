package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Big-Kotik/ivt-pull-api/pkg/api"
	"github.com/Big-Kotik/ivt-pull/pkg/server"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8080, "The grpc server port")

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterPullerServer(grpcServer, &server.PullServer{Client: http.Client{}, Logger: log.Default()})
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
