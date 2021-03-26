package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"grpc-course/pb"
	"grpc-course/service"
	"log"
	"net"
)


func main(){
	flag.Parse()

	port := ":5000"
	fmt.Printf("server running at port %v\n", port)

	grpcServer := grpc.NewServer()
	var laptopServer = service.NewLaptopServer(service.NewInMemoryLaptopStore())
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
