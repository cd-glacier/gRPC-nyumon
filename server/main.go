package main

import (
	"log"
	"net"

	"github.com/g-hyoga/gRPC-nyumon/downloader"
	"github.com/g-hyoga/gRPC-nyumon/echo"
	"github.com/g-hyoga/gRPC-nyumon/uploader"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {
	log.SetFlags(0)
	log.SetPrefix("[file] ")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	srv := grpc.NewServer()
	echo.RegisterEchoServiceServer(srv, &echoService{})
	downloader.RegisterFileServiceServer(srv, &downloadService{})
	uploader.RegisterFileServiceServer(srv, &uploadService{})

	log.Printf("start server on port%s\n", port)

	if err := srv.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
