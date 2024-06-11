package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_chat/types"
	"net/http"
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
}
