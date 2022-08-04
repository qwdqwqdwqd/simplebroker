package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "simplebroker/broker/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%s", port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBrokerClient(conn)
	ctx := context.Background()
	ch, _ := c.Subscribe(ctx, &pb.SubscribeRequest{Queue: "queue"})
	for {
		response, err := ch.Recv()
		fmt.Println(response)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response)
		}
	}
}
