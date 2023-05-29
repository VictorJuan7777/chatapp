package main

import (
	"chatapp/client"
	pb "chatapp/pb"
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	systemRoots, err := x509.SystemCertPool()
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	conn, err := grpc.Dial("chatapp-m4im6z4fda-uc.a.run.app:443", grpc.WithAuthority("chatapp-m4im6z4fda-uc.a.run.app:443"), grpc.WithTransportCredentials(cred))
	// conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to server : ", err)
	}
	defer conn.Close()
	clientServer := pb.NewChatappServerClient(conn)

	stream, err := clientServer.Chatroom(context.Background())
	if err != nil {
		log.Fatal("Failed to call function : ", err)
	}
	waitc := make(chan bool)
	ch := client.Clienthandle{Stream: stream, Waitc: waitc}
	ch.ClientConfig()
	go ch.ClientSendMsg()
	go ch.ClientRecMsg()

	<-waitc

}
