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
	"time"
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
	c.ReadMessage(stream)
	return nil
}
func (c *ClientGRPC) ReadMessage(stream pb.BrokerService_SubscribeClient) {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server.")
			return // End the loop when the server closes the stream.
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		// Only log non-empty messages
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

func InitNewClientFunc(name, s string, ch chan struct{}, f func(c *ClientGRPC)) {
	defer close(ch)
	f(NewClientGRPC(name))
}

func RunGRPCClients() {
	fmt.Println("Started GRPC Client")

	ch := make(chan struct{})

	// Create first client
	go InitNewClientFunc("Sanea", "", ch, func(c *ClientGRPC) {
		//c.NewTopic("dsa")
		//time.Sleep(300 * time.Millisecond)
		//c.Subscribe("dsa")
		//time.Sleep(500 * time.Millisecond)
		//c.Topics()
		//time.Sleep(600 * time.Millisecond)
		c.Publish("Hello from Alex", "dsa")
		time.Sleep(600 * time.Millisecond)
	})

	// Create additional clients
	go func() {
		c := NewClientGRPC("Test")
		InitClientFunc(c, "", ch, func(c *ClientGRPC) {
			//time.Sleep(500 * time.Millisecond)
			//c.Subscribe("dsa")
			//time.Sleep(500 * time.Millisecond)
			//c.Topics()
			//time.Sleep(600 * time.Millisecond)
			c.Publish("Message 1 from Test", "dsa")
			time.Sleep(600 * time.Millisecond)
		})
	}()

	go func() {
		c := NewClientGRPC("Test2")
		InitClientFunc(c, "", ch, func(c *ClientGRPC) {
			//time.Sleep(500 * time.Millisecond)
			//c.Subscribe("dsa")
			//time.Sleep(500 * time.Millisecond)
			//c.Topics()
			//time.Sleep(600 * time.Millisecond)
			c.Publish("Message 2 from Test2", "dsa")
			time.Sleep(600 * time.Millisecond)
		})
	}()

	go func() {
		c := NewClientGRPC("Test3")
		InitClientFunc(c, "", ch, func(c *ClientGRPC) {

			// Publish a message
			if err := c.Publish("Message 3 from Test3", "dsa"); err != nil {
				log.Printf("Error publishing message: %v", err)
				return
			}
		})
	}()

	<-ch
}

func InitClientFunc(c *ClientGRPC, s string, ch chan struct{}, f func(c *ClientGRPC)) {
	defer close(ch)
	f(c)
}
