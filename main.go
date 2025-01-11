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

}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCoffeeShopServer(s, &server{})
}
