package main

import (
	"fmt"
	"log"
	"net"
	"travel/config"
	pbInter "travel/genproto/interactions"
	pbItiner "travel/genproto/itineraries"

	pb "travel/genproto/stories"
	"travel/service"
	"travel/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", config.Load().CONTENT_SERVICE_PORT)
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	db, err := postgres.ConnectDB()
	if err != nil {
		log.Panic(err)
	}

	u := service.NewContentService(db)
	interactions := service.NewInterationsService(db)
	itiner := service.NewItinerariesService(db)
	server := grpc.NewServer()
	pb.RegisterStoriesServer(server, u)
	pbInter.RegisterInteractionsServer(server, interactions)
	pbItiner.RegisterItinerariesServer(server, itiner)

	fmt.Printf("Content service is listening on port %s...\n", config.Load().CONTENT_SERVICE_PORT)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error with listening content server: %s", err)
	}
}
