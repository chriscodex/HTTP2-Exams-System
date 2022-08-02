package main

import (
	"log"
	"net"

	"github.com/ChrisCodeX/gRPC/database"
	"github.com/ChrisCodeX/gRPC/server"
	"github.com/ChrisCodeX/gRPC/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var postgresURL string = "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable"

func main() {
	// Listen announces
	listener, err := net.Listen("tcp", ":5060")
	if err != nil {
		log.Fatal(err)
	}

	// Instance of Postgres Repository
	repo, err := database.NewPostgresRepository(postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	// Instance of Student Server
	server := server.NewStudentServer(repo)

	/*gRPC*/
	// gRPC Server
	s := grpc.NewServer()

	// Register server on the gRPC
	studentpb.RegisterStudentServiceServer(s, server)

	// Register server reflection service on the gRPC
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
