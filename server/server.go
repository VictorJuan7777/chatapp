package server

import (
	pb "chatapp/pb"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

type Room struct {
	ID     int
	Name   string
	Client map[int]*Client
	Msg    chan *Message
}

type Client struct {
	Coon     pb.ChatappServer_ChatroomServer
	ID       int
	Room     string
	Username string
}

type Message struct {
	Content  string
	RoomID   string
	Username string
}

type Hub struct {
	Rooms map[string]*Room
}

type ChatServer struct {
	pb.ChatappServerServer
}

var Allroom = make(map[string]*Room)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (server *ChatServer) Chatroom(stream pb.ChatappServer_ChatroomServer) error {
	clientUniqueCode := rand.Intn(1e6)
	//ctx, cancel := context.WithCancel(context.Background())
	errch := make(chan error)
	// receive message

	go receiveFromStream(stream, clientUniqueCode, errch)
	// sending message
	//go sendToStream(ctx, stream, clientUniqueCode, errch)
	// defer cancel()
	return <-errch
}

func receiveFromStream(stream pb.ChatappServer_ChatroomServer, clientUC int, errch chan error) {
	room := JoinRoom(stream, clientUC, errch)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			LeaveRoom(room, clientUC)
			errch <- err
			break
		}
		if err != nil {
			log.Println("Error in receiving message from client : ", err)
			errch <- err
		} else {
			m := &Message{
				Content:  msg.Msg,
				RoomID:   room,
				Username: msg.Name,
			}
			Allroom[room].Msg <- m
		}
	}
}

func JoinRoom(stream pb.ChatappServer_ChatroomServer, clientUC int, errch chan error) string {
	msg, err := stream.Recv()
	if err == io.EOF {
		errch <- err
	}
	if err != nil {
		log.Println("Error in receiving message from client : ", err)
		errch <- err
	}
	if _, ok := Allroom[msg.Msg]; !ok {
		Allroom[msg.Msg] = &Room{
			ID:     rand.Intn(1e6),
			Name:   msg.Msg,
			Client: make(map[int]*Client),
			Msg:    make(chan *Message),
		}
		go BroadcastRoom(msg.Msg)
	}
	Allroom[msg.Msg].Client[clientUC] = &Client{
		Coon:     stream,
		Username: msg.Name,
		ID:       clientUC,
		Room:     msg.Msg,
	}
	content := fmt.Sprintf("%s is joining this room.", msg.Name)
	Allroom[msg.Msg].Msg <- &Message{
		Content:  content,
		Username: "Announcement",
		RoomID:   msg.Msg,
	}
	return msg.Msg
}

func LeaveRoom(RoomID string, clientUC int) {
	name := Allroom[RoomID].Client[clientUC].Username
	content := fmt.Sprintf("%s is leaving this room.", name)

	Allroom[RoomID].Msg <- &Message{
		Content:  content,
		Username: "Announcement",
		RoomID:   RoomID,
	}
	delete(Allroom[RoomID].Client, clientUC)
}

func BroadcastRoom(RoomID string) {
	for {
		select {
		case m := <-Allroom[RoomID].Msg:
			for k, _ := range Allroom[m.RoomID].Client {
				stream := Allroom[RoomID].Client[k].Coon
				err := stream.Send(&pb.FromServer{Name: m.Username, Msg: m.Content})
				if err != nil {
					fmt.Println("sending error")
				}
			}
		}
	}
}
