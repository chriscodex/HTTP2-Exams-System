package main

import (
	"log"
	"net"

	"github.com/ChrisCodeX/gRPC/database"
	"github.com/ChrisCodeX/gRPC/exampb"
	"github.com/ChrisCodeX/gRPC/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":5070")
	if err != nil {
		log.Fatal(err)
	}

	var postgresURL string = "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable"
	// Instance of Postgres Repository
	repo, err := database.NewPostgresRepository(postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	// Instance of Student Server
	server := server.NewExamServer(repo)

	/*gRPC*/
	// gRPC Server
	s := grpc.NewServer()

	// Register server on the gRPC
	exampb.RegisterExamServiceServer(s, server)

	// Register server reflection service on the gRPC
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
