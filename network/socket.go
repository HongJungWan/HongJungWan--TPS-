// WebSocket 통신을 처리하며 클라이언트와의 실시간 메시지 교환 및 방 관리하는 역할

package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_chat/service"
	"log"
	"net/http"
	"time"
)

type Room struct {
	Forward chan *Message // 수신되는 메시지를 보관하는 값, 들어오는 메시지를 다른 클라이언트에게 전송을 한다.

	Join  chan *Client // Socket이 연결되는 경우에 동작
	Leave chan *Client // Socket이 끊어지는 경우에 동작

	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장

	service *service.Service
}

type Client struct {
	Socket *websocket.Conn // Client의 웹 소켓
	Send   chan *Message   // 전송되는 채널
	Room   *Room           // 클라이언트가 속해 있는 방
	Name   string          `json:"name"`
}

type Message struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Room    string    `json:"room"`
	When    time.Time `json:"when"`
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: map[*Client]bool{},
		service: service,
	}
}

func (client *Client) Read() {
	defer client.Socket.Close()
	for {
		var message *Message
		err := client.Socket.ReadJSON(&message)
		if err != nil {
			return
		}

		message.When = time.Now()
		message.Name = client.Name
		client.Room.Forward <- message
	}
}

func (client *Client) Write() {
	defer client.Socket.Close()
	for msg := range client.Send {
		err := client.Socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

func (room *Room) Run() {
	for {
		select {
		case client := <-room.Join:
			room.Clients[client] = true // client가 새로 들어올 떄,
		case client := <-room.Leave:
			delete(room.Clients, client) // 나갈 때에는 map 값에서 client를 제거.
			close(client.Send)           // 이후 client의 socker을 닫는다.
		case msg := <-room.Forward: // 만약 특정 메시지가 방에 들어오면,

			go room.service.InsertChatting(msg.Name, msg.Message, msg.Room)

			for client := range room.Clients {
				client.Send <- msg // 모든 client에게 전달.
			}
		}
	}
}

const (
	SocketBufferSize  = 1024
	messageBufferSize = 256
)

// 기본적으로 HTTP에 웹 소켓을 사용하려면, 이와 같이 업그레이드해 주어야 한다.
// -> 재사용 가능하기 때문에 하나만 만들어도 된다.
var upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: messageBufferSize}

func (room *Room) ServeHTTP(c *gin.Context) {

	// 이후 요청이 이렇게 들어오게 된다면 Upgrade를 통해서 소켓을 가져온다.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	Socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("---- serveHTTP:", err)
		return
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("auth cookie is failed", err)
		return
	}

	// 문제가 없다면 client를 생성하여 방에 입장했다고 채널에 전송.
	client := &Client{
		Socket: Socket,
		Send:   make(chan *Message, messageBufferSize),
		Room:   room,
		Name:   authCookie.Value,
	}

	room.Join <- client

	// 또한 defer를 통해서 client가 끝날 때를 대비하여 퇴장하는 작업 연기.
	defer func() { room.Leave <- client }()

	// 이후 고루틴을 통해서 write를 실행.
	go client.Write()

	// 이후 메인 루틴에서 read를 실행함으로써 해당 요청을 닫는 것을 차단.
	// -> 연결을 활성화. 채널을 활용
	client.Read()
}
