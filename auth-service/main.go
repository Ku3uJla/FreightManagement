package main

import (
	"auth-service/internal/controller"
	"auth-service/internal/repository"
	"auth-service/internal/repository/db"
	"auth-service/internal/routes"
	"auth-service/internal/service"
	"log"
	"net"
	"os"

	pb1 "auth-service/proto/authpb"
	pb "auth-service/proto/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := db.InitDB()
	router := gin.Default()
	gRPCaddr := os.Getenv("USER_APP_URL")
	if gRPCaddr == "" {
		log.Fatalf("failed to extract user_app_url")
	}
	conn, err := grpc.NewClient(gRPCaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if conn == nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)

	}
	userClient := pb.NewUserServiceClient(conn)
	if userClient == nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	authStore := repository.NewAuthRepository(db)
	authSrv := service.NewAuthService(authStore, userClient)
	authController := controller.NewAuthController(authSrv)
	routes.AuthRoutes(router, authController)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "auth-service"})
	})
	go startGrpc(authSrv)
	go func() {
		if err := router.Run(":8083"); err != nil {
			log.Fatal(err)
		}
	}()

	select {}

}

func startGrpc(authSrv service.AuthService) {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcController := controller.NewAuthGrpcController(authSrv)
	pb1.RegisterAuthServiceServer(s, grpcController)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
