package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/muzzlol/gRPCwgolangez/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	start := time.Now()

	// Connection setup timing
	connStart := time.Now()
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	log.Printf("Connection setup took: %v", time.Since(connStart))

	client := pb.NewCoffeeShopClient(conn)

	// Menu streaming timing
	menuStart := time.Now()
	stream, err := client.GetMenu(context.Background(), &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("failed to get menu: %v", err)
	}

	done := make(chan bool)

	var items []*pb.Item
	go func() {
		for {
			itemStart := time.Now()
			menu, err := stream.Recv()

			if err == io.EOF {
				done <- true
				break
			}

			if err != nil {
				log.Fatalf("failed to receive menu: %v", err)
			}

			items = append(items, menu.Items...)
			log.Printf("Received menu item in: %v", time.Since(itemStart))
		}
	}()
	<-done
	log.Printf("Menu streaming took: %v", time.Since(menuStart))

	// Order placement timing
	orderStart := time.Now()
	reciept, err := client.PlaceOrder(context.Background(), &pb.Order{Items: items})
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}
	log.Printf("Order placement took: %v", time.Since(orderStart))
	log.Printf("Reciept: %v", reciept)

	// Total time
	log.Printf("Total operation took: %v", time.Since(start))
}
