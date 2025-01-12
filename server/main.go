package main

import (
	"log"
	"net"

	pb "github.com/muzzlol/gRPCwgolangez/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
	menuItems   []*pb.Item
	menuUpdates chan *pb.Item
}

func (s *server) GetMenu(req *pb.MenuRequest, stream pb.CoffeeShop_GetMenuServer) error {
	menu := &pb.Menu{
		Items: s.menuItems,
	}

	if err := stream.Send(menu); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case newItem := <-s.menuUpdates:
			s.menuItems = append(s.menuItems, newItem)
			// sending de whole menu again
			if err := stream.Send(&pb.Menu{Items: s.menuItems}); err != nil {
				return err
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := &server{
		// menuItems is a slice of pointers to Item structs
		menuItems: []*pb.Item{
			{Name: "Espresso", Price: 2.50},
			{Name: "Latte", Price: 3.00},
			{Name: "Cappuccino", Price: 3.00},
			{Name: "Mocha", Price: 3.50},
			{Name: "Americano", Price: 2.00},
		},
		// menuUpdates is a channel that receives pointers to Item structs in order to update the menu
		menuUpdates: make(chan *pb.Item),
	}

	s := grpc.NewServer()
	pb.RegisterCoffeeShopServer(s, srv)

	log.Printf("Server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
