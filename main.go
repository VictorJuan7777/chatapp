package main

import (
	pb "chatapp/pb"
	"chatapp/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port string = "8080"
)

func main() {
	// log.Println("Enter room number")
	// reader := bufio.NewReader(os.Stdin)
	// roomN, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Println("Failed to read from console: ", err)
	// }
	// roomN = strings.Trim(roomN, "\r\n")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen on :", port)
	}

	log.Println("Listening on : ", port)

	s := grpc.NewServer()
	pb.RegisterChatappServerServer(s, &server.ChatServer{})
	// reflection.Register(s)
	log.Println("Start server")
	if err = s.Serve(lis); err != nil {
		log.Fatal("Failed to server: ", err)
	}

}
