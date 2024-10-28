package broker_grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
	"time"

	//pb "pad/message_tcp/"
	"google.golang.org/grpc"
)

type Subscriber struct {
	id          string
	topics      map[string]bool
	messageChan chan *pb.Message
}

type BrokerGRPC struct {
	pb.UnimplementedBrokerServiceServer
	mutex       *sync.RWMutex
	subscribers map[string]*Subscriber
	topics      map[string]bool
}

func NewBrokerGRPC() *BrokerGRPC {
	return &BrokerGRPC{
		mutex:       &sync.RWMutex{},
		subscribers: make(map[string]*Subscriber),
		topics:      make(map[string]bool),
	}
}

func (b *BrokerGRPC) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	b.mutex.RLock()
	if !b.topics[req.Topic] {
		b.mutex.RUnlock()
		return nil, status.Errorf(codes.NotFound, "Topic does not exist")
	}
	b.mutex.RUnlock()

	message := &pb.Message{
		Sender:    req.Sender,
		Topic:     req.Topic,
		Content:   req.Message,
		Timestamp: time.Now().Unix(),
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for _, subscriber := range b.subscribers {
		if subscriber.topics[req.Topic] {
			select {
			case subscriber.messageChan <- message:
				log.Printf("Message sent to subscriber %s on topic %s", subscriber.id, req.Topic)
			default:
				log.Printf("Failed to send message to subscriber %s: channel full", subscriber.id)
			}
		}
	}

	return &pb.MessageResponse{Confirmation: "Message sent"}, nil
}

func (b *BrokerGRPC) Subscribe(req *pb.SubscribeRequest, stream pb.BrokerService_SubscribeServer) error {
	b.mutex.Lock()
	if !b.topics[req.Topic] {
		b.mutex.Unlock()
		return status.Errorf(codes.NotFound, "Topic does not exist")
	}

	subscriber, exists := b.subscribers[req.ClientId]
	if !exists {
		subscriber = &Subscriber{
			id:          req.ClientId,
			topics:      make(map[string]bool),
			messageChan: make(chan *pb.Message, 100),
		}
		b.subscribers[req.ClientId] = subscriber
	}
	subscriber.topics[req.Topic] = true
	b.mutex.Unlock()

	defer func() {
		b.mutex.Lock()
		delete(subscriber.topics, req.Topic)
		if len(subscriber.topics) == 0 {
			delete(b.subscribers, req.ClientId)
			close(subscriber.messageChan)
		}
		b.mutex.Unlock()
	}()

	for {
		select {
		case msg := <-subscriber.messageChan:
			if err := stream.Send(msg); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (b *BrokerGRPC) GetTopics(ctx context.Context, req *pb.TopicsRequest) (*pb.TopicsResponse, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	topics := make([]string, 0, len(b.topics))
	for topic := range b.topics {
		topics = append(topics, topic)
	}

	return &pb.TopicsResponse{Topics: topics}, nil
}

func (b *BrokerGRPC) CreateTopic(ctx context.Context, req *pb.TopicRequest) (*pb.TopicResponse, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.topics[req.Topic] {
		return nil, status.Errorf(codes.AlreadyExists, "Topic already exists")
	}

	b.topics[req.Topic] = true
	return &pb.TopicResponse{Status: "Topic created"}, nil
}

func (b *BrokerGRPC) DeleteTopic(ctx context.Context, req *pb.TopicRequest) (*pb.TopicResponse, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.topics[req.Topic] {
		return nil, status.Errorf(codes.NotFound, "Topic does not exist")
	}

	delete(b.topics, req.Topic)
	return &pb.TopicResponse{Status: "Topic deleted"}, nil
}

func Run() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterBrokerServiceServer(server, NewBrokerGRPC())

	log.Println("Broker gRPC server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
