package main

import (
	"context"
	"crypto/tls"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/matey97/grpc_test/grpc_test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	conn, err := grpc.Dial("api-gateway-q5t3vl3tfa-ew.a.run.app:443", opts...)
	if err != nil {
		log.Printf("Could not create connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewGRPCTestClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.GetAllMessages(ctx, &empty.Empty{})
	if err != nil {
		log.Printf("Request failed: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Done")
		}
		if err != nil {
			log.Printf("Failed to receive a message: %v", err)
			break
		}
		log.Printf("Message: %v", message)
	}
}
