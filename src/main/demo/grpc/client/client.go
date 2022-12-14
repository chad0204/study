package main

import (
	"context"
	"fmt"
	"log"
	"study/src/main/demo/grpc/server/win"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := win.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &win.String{Value: "world"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}
