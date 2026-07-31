package main

import (
	"log"
	"net"

	"os"
	"user-service/internal/controller"
	"user-service/internal/repository"
	"user-service/internal/repository/db"
	"user-service/internal/routes"
	"user-service/internal/service"
	"user-service/publisher"

	pb "user-service/proto/userpb"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	pub, err := publisher.NewUserPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("Критическая ошибка: %v", err)
	}
	defer pub.Close()

	database := db.InitDB()

	userStore := repository.NewUserRepository(database)

	userService := service.NewUserService(userStore, pub)

	userController := controller.NewUserController(userService)

	// HTTP
	router := gin.Default()

	routes.UserRoutes(
		router,
		userController,
	)

	go func() {

		if err := router.Run(":8080"); err != nil {
			log.Fatal(err)
		}

	}()

	// gRPC
	go startGRPC(userService)

	select {}
}

func startGRPC(
	userSrv service.UserService,
) {

	listener, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	userGrpcController := controller.NewUserGrpcController(userSrv)
	pb.RegisterUserServiceServer(
		grpcServer,
		userGrpcController,
	)

	log.Println(
		"gRPC server running :50051",
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
