package client

import (
	"bufio"
	pb "chatapp/pb"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Clienthandle struct {
	Stream     pb.ChatappServer_ChatroomClient
	ClientName string
	Room       string
	Waitc      chan bool
}

func (ch *Clienthandle) ClientConfig() {
	fmt.Println("Enter your name")
	reader := bufio.NewReader(os.Stdin)
	clientName, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read from console : ", err)
	}
	ch.ClientName = strings.Trim(clientName, "\r\n")
	fmt.Println("Enter room number")
	reader = bufio.NewReader(os.Stdin)
	roomname, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read from console : ", err)
	}
	ch.Room = strings.Trim(roomname, "\r\n")
	arg := pb.FromClient{
		Name: ch.ClientName,
		Msg:  ch.Room,
	}
	err = ch.Stream.Send(&arg)
	if err != nil {
		log.Println("Error while sending message to server : ", err)
	}
	fmt.Printf("Hello : %v ! You are entering room : %v \n", ch.ClientName, ch.Room)

}

func (ch *Clienthandle) ClientSendMsg() {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Failed to read from console : ", err)
		}
		msg = strings.Trim(msg, "\r\n")
		if msg == "/exit" {
			log.Println("Ending chatting")
			ch.Stream.CloseSend()
			ch.Waitc <- true
			break
		}
		arg := pb.FromClient{
			Name: ch.ClientName,
			Msg:  msg,
		}
		err = ch.Stream.Send(&arg)
		if err != nil {
			log.Println("Error while sending message to server : ", err)
		}
	}
}

func (ch *Clienthandle) ClientRecMsg() {
	for {
		msg, err := ch.Stream.Recv()
		if err == io.EOF {
			ch.Waitc <- true
			break
		}
		if err != nil {
			fmt.Println("Error in receiving message from server : ", err)
		}
		fmt.Printf("%s : %s \n", msg.Name, msg.Msg)
	}
}
