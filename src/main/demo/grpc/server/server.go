package main

import (
	"log"
	"net"
	"study/src/main/demo/grpc/server/win"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	win.RegisterHelloServiceServer(grpcServer, new(win.UnimplementedHelloServiceServer))
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
