package main

import (
	"bufio"
	"fmt"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	eventMsg = make(chan string)
)

func broadcaster() {
	//broadcaster是唯一能访问clientMap的goroutine, 不同的goroutine通过chan来通信, 而chan是并发安全的
	//这保证了程序的并发安全。
	clientMap := make(map[client]bool)
	for {
		select {
		case cli := <-entering:
			clientMap[cli] = true
		case cli := <-leaving:
			delete(clientMap, cli)
			close(cli)
		case msg := <-eventMsg:
			for cli := range clientMap {
				//向客户端连接写入, clientWriter函数会写回客户端
				cli <- msg
			}
		}
	}
}

func handlerClient(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who

	//发送客户端进入事件
	eventMsg <- who + " has arrived"
	//发送客户端
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		//记录写入事件
		eventMsg <- who + ": " + input.Text()
	}

	leaving <- ch
	eventMsg <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

//聊天服务
func main() {
	listen, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		return
	}
	go broadcaster()
	for {
		accept, err := listen.Accept()
		if err != nil {
			continue
		}
		go handlerClient(accept)
	}
}
