package main

import (
	"auth-service/internal/controller"
	"auth-service/internal/repository"
	"auth-service/internal/repository/db"
	"auth-service/internal/routes"
	"auth-service/internal/service"
	"log"

	pb "auth-service/proto/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := db.InitDB()
	router := gin.Default()
	gRPCaddr := "localhost:50051"
	conn, err := grpc.NewClient(gRPCaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)

	}
	userClient := pb.NewUserServiceClient(conn)
	authStore := repository.NewAuthRepository(db)
	authSrv := service.NewAuthService(authStore, userClient)
	authController := controller.NewAuthController(authSrv)
	routes.AuthRoutes(router, authController)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "user-service"})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
