package client_grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"pad/proto"
	pb "pad/proto"
	"sync"
)

type ClientGRPC struct {
	id     string
	client pb.BrokerServiceClient
	conn   *grpc.ClientConn
	mutex  sync.Mutex
}

func NewClientGRPC(id string) *ClientGRPC {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	return &ClientGRPC{
		id:     id,
		client: pb.NewBrokerServiceClient(conn),
		conn:   conn,
	}
}

func (c *ClientGRPC) Publish(message, topic string) error {
	req := &proto.MessageRequest{
		Sender:  c.id,
		Topic:   topic,
		Message: message,
	}

	_, err := c.client.SendMessage(context.Background(), req)
	return err
}

func (c *ClientGRPC) Subscribe(topic string) error {
	stream, err := c.client.Subscribe(context.Background(), &proto.SubscribeRequest{
		ClientId: c.id,
		Topic:    topic,
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving message: %v", err)
		}

		log.Printf("[%s] %s: %s", message.Topic, message.Sender, message.Content)
	}
}

func (c *ClientGRPC) Topics() error {
	resp, err := c.client.GetTopics(context.Background(), &proto.TopicsRequest{})
	if err != nil {
		return err
	}

	fmt.Println("Available topics:")
	for _, topic := range resp.Topics {
		fmt.Printf("- %s\n", topic)
	}
	return nil
}

func (c *ClientGRPC) NewTopic(topic string) error {
	_, err := c.client.CreateTopic(context.Background(), &proto.TopicRequest{Topic: topic})
	return err
}

func (c *ClientGRPC) DeleteTopic(topic string) error {
	_, err := c.client.DeleteTopic(context.Background(), &proto.TopicRequest{Topic: topic})
	return err
}

func (c *ClientGRPC) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ClientGRPC) SendMessage(req *proto.MessageRequest) error {
	_, err := c.client.SendMessage(context.Background(), req)
	return err
}
