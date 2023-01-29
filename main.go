package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/e-commerce-microservices/user-service/pb"
	"github.com/e-commerce-microservices/user-service/repository"
	"github.com/e-commerce-microservices/user-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// postgres driver
	_ "github.com/lib/pq"
)

func main() {
	// create grpc server
	grpcServer := grpc.NewServer()

	// init user db connection
	pgDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("DB_DBNAME"),
	)

	userDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer userDB.Close()
	if err := userDB.Ping(); err != nil {
		log.Fatal("can't ping to user db", err)
	}

	// init user queries
	userQueries := repository.New(userDB)

	// dial auth client
	authServiceConn, err := grpc.Dial("auth-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can't dial user service: ", err)
	}
	// create auth client
	authClient := pb.NewAuthServiceClient(authServiceConn)

	// create user service
	userService := service.NewUserService(userQueries, authClient)
	// register user service
	pb.RegisterUserServiceServer(grpcServer, userService)

	// listen and serve
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
