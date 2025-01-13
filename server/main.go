package main

import (
	"context"
	"log"
	"net"

	pb "github.com/muzzlol/gRPCwgolangez/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(req *pb.MenuRequest, stream pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{
		{Id: 1, Name: "Espresso", Price: 2.5},
		{Id: 2, Name: "Latte", Price: 3.0},
		{Id: 3, Name: "Cappuccino", Price: 3.0},
		{Id: 4, Name: "Mocha", Price: 3.5},
		{Id: 5, Name: "Americano", Price: 2.0},
		{Id: 6, Name: "Macchiato", Price: 3.0},
		{Id: 7, Name: "Flat White", Price: 3.5},
		{Id: 8, Name: "Long Black", Price: 2.5},
		{Id: 9, Name: "Affogato", Price: 4.0},
		{Id: 10, Name: "Irish Coffee", Price: 5.0},
		{Id: 11, Name: "Cheesecake", Price: 3.0},
		{Id: 12, Name: "Brownie", Price: 2.5},
		{Id: 13, Name: "Croissant", Price: 2.0},
		{Id: 14, Name: "Muffin", Price: 2.0},
		{Id: 15, Name: "Cookie", Price: 1.5},
		{Id: 16, Name: "Sandwich", Price: 4.0},
		{Id: 17, Name: "Salad", Price: 5.0},
		{Id: 18, Name: "Pasta", Price: 6.0},
		{Id: 19, Name: "Pizza", Price: 7.0},
		{Id: 20, Name: "Burger", Price: 5.0},
	}
	// simulaitng streaming the menu by sending each item one by one
	for i := range items {
		stream.Send(&pb.Menu{Items: []*pb.Item{items[i]}})
	}

	return nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Reciept, error) {
	var total float32
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, item := range order.Items {
			total += item.Price
		}
		return &pb.Reciept{Total: total}, nil
	}
}

// func (s *server) GetOrderStatus(ctx context.Context, reciept *pb.Reciept) (pb.OrderStatus, error) {

// }

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCoffeeShopServer(s, &server{})

	log.Printf("Server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
