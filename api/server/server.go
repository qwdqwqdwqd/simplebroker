package server

import (
	"context"

	pb "simplebroker/broker/proto"
	brokerImp "simplebroker/internal/broker"
	"simplebroker/pkg/broker"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	broker broker.Broker
	pb.UnimplementedBrokerServer
}

func (s Server) Publish(ctx context.Context, publishReq *pb.PublishRequest) (*pb.PublishResponse, error) {
	message := broker.Message{
		Body: string(publishReq.GetBody()),
	}
	messageId, err := s.broker.Publish(ctx, publishReq.GetQueue(), message)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Broker has been closed.")
	}
	response := &pb.PublishResponse{MessageID: int64(messageId)}
	return response, nil

}
func (s Server) Subscribe(req *pb.SubscribeRequest, stream pb.Broker_SubscribeServer) error {
	ch, err := s.broker.Subscribe(context.Background(), req.GetQueue())
	if err != nil {
		return status.Errorf(codes.Unavailable, "Broker has been closed.")
	}
	for message := range ch {
		messageResponse := &pb.MessageResponse{Body: []byte(message.Body)}
		err := stream.Send(messageResponse)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetServer() pb.BrokerServer {
	return &Server{broker: brokerImp.NewBroker()}
}
