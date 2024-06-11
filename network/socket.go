package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_chat/types"
	"log"
	"net/http"
	"time"
)

var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	Forward chan *Message // 수신되는 메시지를 보관하는 값, 들어오는 메시지를 다른 클라이언트에게 전송을 한다.

	Join  chan *Client // Socket이 연결되는 경우에 동작
	Leave chan *Client // Socket이 끊어지는 경우에 동작

	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

type Message struct {
	Name            string
	ChattingMessage string
	Time            int64
}

type Client struct {
	Send   chan *Message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *Message),

		Join:  make(chan *Client),
		Leave: make(chan *Client),

		Clients: map[*Client]bool{},
	}
}

func (client *Client) Read() {
	// 클라이언트가 들어오는 메시지를 읽는 함수
	defer client.Socket.Close()
	for {
		var message *Message
		err := client.Socket.ReadJSON(&message)
		if err != nil {
			// 클라이언트와의 WebSocket 연결이 예상치 못하게 종료되는 경우를 처리
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			}
			panic(err)
		} else {
			log.Println("READ : ", message, "Client", client.Name)
			log.Println()

			message.Time = time.Now().Unix()
			message.Name = client.Name

			client.Room.Forward <- message
		}
	}
}

func (client *Client) Write() {
	// 클라이언트가 메시지를 전송하는 함수
	defer client.Socket.Close()

	for message := range client.Send {
		log.Println("WRITE : ", message, "Client", client.Name)
		log.Println()

		err := client.Socket.WriteJSON(message)
		if err != nil {
			panic(err)
		}
	}
}

func (room *Room) RunInit() {
	// Room에 있는 모든 채널 값들을 받는 역할
	for {
		select {
		case client := <-room.Join:
			room.Clients[client] = true
		case client := <-room.Leave:
			room.Clients[client] = false
			close(client.Send)
			delete(room.Clients, client)
		case message := <-room.Forward:
			for client := range room.Clients {
				client.Send <- message
			}
		}
	}
}

func (room *Room) SocketServe(context *gin.Context) {
	socket, err := Upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := context.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Send:   make(chan *Message, types.MessageBufferSize),
		Room:   room,
		Name:   userCookie.Value,
		Socket: socket,
	}

	room.Join <- client

	defer func() { room.Leave <- client }()

	go client.Write()

	client.Read()
}
