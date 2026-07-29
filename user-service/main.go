package main

import (
	"log"
	"net"

	"user-service/internal/controller"
	"user-service/internal/repository"
	"user-service/internal/repository/db"
	"user-service/internal/routes"
	"user-service/internal/service"

	pb "user-service/proto/userpb"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
)

func main() {

	database := db.InitDB()

	userStore := repository.NewUserRepository(database)

	userService := service.NewUserService(userStore)

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
	go startGRPC(userController)

	select {}
}

func startGRPC(
	userController *controller.UserController,
) {

	listener, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(
		grpcServer,
		userController,
	)

	log.Println(
		"gRPC server running :50051",
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
