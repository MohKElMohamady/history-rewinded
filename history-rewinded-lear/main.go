package main

import (
	_ "history-rewinded-lear/cassandra"
	"history-rewinded-lear/lear"
	"history-rewinded-regan/pb"
	"log"
	"net"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	address = ":1606"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Lear...")

	godotenv.Load()

	king_lear := lear.New()

	lis, err := net.Listen("tcp", address)	
	if err != nil {
		log.Fatalf("failed to listen to port %v, reason: %v", address, err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterLearServer(s, &king_lear)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start gRPC server, reason: %v\n", err.Error())
	}

}
