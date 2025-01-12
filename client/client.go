package main

import (
	"context"
	"log"
	"time"

	pb "github.com/muzzlol/gRPCwgolangez/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewCoffeeShopClient(conn)

	stream, err := client.GetMenu(context.Background(), &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("failed to get menu: %v", err)
	}

	for {
		start := time.Now()
		menu, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to receive menu: %v", err)
		}
		duration := time.Since(start)
		log.Printf("menu: %v", menu)
		log.Printf("Time taken to receive menu: %v", duration)
	}
}
