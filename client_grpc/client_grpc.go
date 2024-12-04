package clientGrpc

// import (
// 	"context"
// 	"google.golang.org/grpc"
// 	"log"
// 	pb "pad/proto/message.proto"
// )
//
// // ClientGRPC defines the gRPC client structure
// type ClientGRPC struct {
// 	client pb.BrokerServiceClient
// }
//
// // NewClientGRPC initializes a new ClientGRPC instance
// func NewClientGRPC(conn *grpc.ClientConn) *ClientGRPC {
// 	return &ClientGRPC{
// 		client: pb.NewBrokerServiceClient(conn),
// 	}
// }
//
// // SendMessage sends a message to the server
// func (c *ClientGRPC) SendMessage(sender, message string) {
// 	req := &pb.MessageRequest{
// 		Sender:  sender,
// 		Message: message,
// 	}
//
// 	res, err := c.client.SendMessage(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Error sending message: %v", err)
// 	}
//
// 	log.Printf("Server response: %s", res.Confirmation)
// }
//
// // Run starts the gRPC client
// //func Run() {
// //	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// //	if err != nil {
// //		log.Fatalf("Failed to connect to server: %v", err)
// //	}
// //	defer conn.Close()
// //
// //	client := NewClientGRPC(conn)
// //
// //	// Send a message every second
// //	ticker := time.NewTicker(1 * time.Second)
// //	defer ticker.Stop()
// //
// //	for range ticker.C {
// //		client.SendMessage("Client", "Hello from the gRPC client!")
// //	}
// //}
