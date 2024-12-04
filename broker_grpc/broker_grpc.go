package brokerGrpc

// import (
// 	"context"
// 	"log"
// 	"net"
// 	"sync"
//
// 	pb "pad/proto/message.proto"
//
// 	"google.golang.org/grpc"
// )
//
// // BrokerGRPC defines the gRPC server structure
// type BrokerGRPC struct {
// 	pb.UnimplementedBrokerServiceServer
// 	mutex *sync.Mutex
// }
//
// // NewBrokerGRPC initializes a new BrokerGRPC instance
// func NewBrokerGRPC() *BrokerGRPC {
// 	return &BrokerGRPC{
// 		mutex: &sync.Mutex{},
// 	}
// }
//
// // SendMessage handles sending messages from client
// func (b *BrokerGRPC) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
// 	log.Printf("Received message from %s: %s", req.Sender, req.Message)
//
// 	response := &pb.MessageResponse{
// 		Confirmation: "Message received by the server!",
// 	}
//
// 	return response, nil
// }
//
// // Run starts the gRPC broker server
// func Run() {
// 	listener, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen on port 50051: %v", err)
// 	}
//
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterBrokerServiceServer(grpcServer, NewBrokerGRPC())
//
// 	log.Println("Broker gRPC server is running on port 50051...")
//
// 	if err := grpcServer.Serve(listener); err != nil {
// 		log.Fatalf("Failed to serve gRPC server: %v", err)
// 	}
// }
