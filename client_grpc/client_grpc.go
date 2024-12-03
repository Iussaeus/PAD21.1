package client_grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"pad/proto"
	pb "pad/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	id     string
	client pb.BrokerServiceClient
	conn   *grpc.ClientConn
}

func NewClientGRPC(id string) *ClientGRPC {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	go c.ReadMessage(stream)
	return nil
}

func (c *ClientGRPC) ReadMessage(stream pb.BrokerService_SubscribeClient) {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server.")
			return
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		if message != nil && message.Topic != "" && message.Content != "" {
			log.Printf("[%s] %s: %s", message.Topic, message.Sender, message.Content)
		}
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

func InitNewClientFunc(name string, wch chan struct{}, f func(c *ClientGRPC)) {
	c := NewClientGRPC(name)

	f(c)
	wch <- struct{}{}
}

func RunGRPCClients() {
	fmt.Println("Started GRPC Client")

	ch := make(chan struct{})
	defer close(ch)

	go InitNewClientFunc("Sanea", ch, func(c *ClientGRPC) {
		c.NewTopic("dsa")
		time.Sleep(300 * time.Millisecond)
		c.Subscribe("dsa")
		time.Sleep(500 * time.Millisecond)
		time.Sleep(500 * time.Millisecond)
		c.Topics()
		time.Sleep(600 * time.Millisecond)
		c.Publish("Hello from Alex", "dsa")
		time.Sleep(600 * time.Millisecond)
	})

	go InitNewClientFunc("Test2", ch, func(c *ClientGRPC) {
		time.Sleep(500 * time.Millisecond)
		time.Sleep(500 * time.Millisecond)
		c.Subscribe("dsa")
		time.Sleep(500 * time.Millisecond)
		c.Topics()
		time.Sleep(600 * time.Millisecond)
		c.Publish("Message 1 from Test", "dsa")
		time.Sleep(600 * time.Millisecond)
	})

	go InitNewClientFunc("Test3", ch, func(c *ClientGRPC) {
		time.Sleep(500 * time.Millisecond)
		time.Sleep(500 * time.Millisecond)
		c.Subscribe("dsa")
		time.Sleep(500 * time.Millisecond)
		c.Topics()
		time.Sleep(600 * time.Millisecond)
		c.Publish("Message 2 from Test2", "dsa")
		time.Sleep(600 * time.Millisecond)
	})

	<-ch
}
