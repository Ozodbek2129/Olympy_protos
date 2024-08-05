package main

import (
	"log"
	"net"
	"user/api"
	"user/api/handler"
	"user/config"
	pb "user/genproto/userservice"
	"user/logger"
	"user/service"
	"user/storage/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	listener, err := net.Listen("tcp", config.Load().USER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	db, err := postgres.ConnectionDb()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLogger()
	service := service.NewUserService(db, logger)

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, service)

	log.Printf("Server listening at %v", listener.Addr())
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	hand := NewHandler()
	router := api.NewRouter(hand)
	err = router.Run(config.Load().USER_ROUTER)
	if err != nil {
		log.Fatal(err)
	}
}

func NewHandler() *handler.Handler {
	cfg := config.Load()
	conn, err := grpc.Dial(cfg.USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error while connecting to authentication service: %v", err)
	}

	return &handler.Handler{
		AuthUser: pb.NewUserServiceClient(conn),
		Log:      logger.NewLogger(),
	}
}