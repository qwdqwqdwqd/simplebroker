package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"log"
	"net"
	"simplebroker/api/server"
	pb "simplebroker/broker/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "net/http/pprof"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("couldn't load env file")
	}

	port, exists := os.LookupEnv("GPRC_PORT")
	if !exists {
		log.Fatalf("cannot read gRPC port")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBrokerServer(grpcServer, server.GetServer())
	fmt.Println("starting grpc server...")
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("starting gPRC server failed: %v\n", err)
	}

}
