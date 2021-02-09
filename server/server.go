package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	pb "github.com/matey97/grpc_test/grpc_test"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
)

type grpcServer struct {
	pb.UnimplementedGRPCTestServer
	savedMessages []*pb.Message
}

type MessageToBQ struct {
	Message     string
	FromName    string
	FromSurname string
	ToName      string
	ToSurname   string
}

/* For custom mapping
func (m *MessageToBQ) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"message": m.Message,
		"from": m.From,
		"to": m.To,
	}, bigquery.NoDedupeID, nil
}*/

func GetClient(ctx context.Context) (*bigquery.Client, error) {
	return bigquery.NewClient(ctx, "activity-detection-55cc7")
}

func (s *grpcServer) SendMessage(ctx context.Context, message *pb.Message) (*pb.ACK, error) {
	log.Printf("Message Received: %v", message)
	s.savedMessages = append(s.savedMessages, message)

	ctxt := context.Background()
	client, err := GetClient(ctxt)
	if err != nil {
		log.Printf("Bigquery.NewClient error: %v", err)
		return &pb.ACK{Received: false}, err
	}

	defer client.Close()
	inserter := client.Dataset("test_dataset").Table("test_table").Inserter()

	messageToBQ := MessageToBQ{
		Message:     message.GetMessage(),
		FromName:    message.GetFrom().GetName(),
		FromSurname: message.GetFrom().GetSurname(),
		ToName:      message.GetTo().GetName(),
		ToSurname:   message.GetTo().GetSurname(),
	}

	if err := inserter.Put(ctxt, messageToBQ); err != nil {
		log.Printf("Inserter.Put error: %v", err)
		return &pb.ACK{Received: false}, err
	}

	return &pb.ACK{Received: true}, nil
}

func (s *grpcServer) GetMessagesTo(person *pb.Person, stream pb.GRPCTest_GetMessagesToServer) error {
	ctxt := context.Background()
	client, err := GetClient(ctxt)
	if err != nil {
		log.Printf("Bigquery.NewClient error: %v", err)
		return err
	}

	defer client.Close()

	q := client.Query(`SELECT * FROM ` +
		"`activity-detection-55cc7.test_dataset.test_table` " +
		`WHERE ` + "`toName` " + `= @name AND` + "`toSurname`" + `=@surname`)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "name", Value: person.GetName()},
		{Name: "surname", Value: person.GetSurname()},
	}

	it, err := q.Read(ctxt)
	if err != nil {
		log.Printf("client.Read error: %v", err)
		return err
	}

	messages, err := parseDataFromBQ(it)

	if err != nil {
		log.Printf("parseDataFromBQ error: %v", err)
		return err
	}

	for _, message := range messages {
		err = stream.Send(&message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *grpcServer) GetAllMessages(_ *emptypb.Empty, stream pb.GRPCTest_GetAllMessagesServer) error {
	ctxt := context.Background()
	client, err := GetClient(ctxt)
	if err != nil {
		log.Printf("Bigquery.NewClient error: %v", err)
		return err
	}

	defer client.Close()

	q := client.Query(`SELECT * FROM ` +
		"`activity-detection-55cc7.test_dataset.test_table`")

	it, err := q.Read(ctxt)
	if err != nil {
		log.Printf("client.Read error: %v", err)
		return err
	}

	messages, err := parseDataFromBQ(it)

	if err != nil {
		log.Printf("parseDataFromBQ error: %v", err)
		return err
	}

	for _, message := range messages {
		err = stream.Send(&message)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseDataFromBQ(it *bigquery.RowIterator) ([]pb.Message, error) {
	messages := make([]pb.Message, 0)
	for {
		var values MessageToBQ
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		messages = append(messages, buildMessage(values))
	}

	return messages, nil
}

func buildMessage(m MessageToBQ) pb.Message {
	return pb.Message{
		Message: m.Message,
		From: &pb.Person{
			Name:    m.FromName,
			Surname: m.FromSurname,
		},
		To: &pb.Person{
			Name:    m.ToName,
			Surname: m.ToSurname,
		},
	}
}

func newServer() *grpcServer {
	s := &grpcServer{savedMessages: []*pb.Message{}}
	return s
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, error := net.Listen("tcp", ":"+port)
	if error != nil {
		log.Fatalf("Failed to listen: %v", error)
	}

	server := grpc.NewServer()
	pb.RegisterGRPCTestServer(server, newServer())
	error = server.Serve(listener)
	if error != nil {
		log.Fatalf("Failed to serve: %v", error)
	}
}
