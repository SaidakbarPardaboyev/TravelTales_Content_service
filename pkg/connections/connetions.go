package connections

import (
	"log"
	"travel/config"
	pb "travel/genproto/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient() pb.UsersClient {
	conn, err := grpc.NewClient(config.Load().USER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return pb.NewUsersClient(conn)
}
