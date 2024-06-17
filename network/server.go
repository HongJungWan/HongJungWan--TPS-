// 서버의 API 엔드포인트를 정의하고 요청 핸들러를 등록하는 역할

package network

import "github.com/gin-gonic/gin"

type api struct {
	server *Server
}

func registerServer(server *Server) {
	a := &api{server: server}

	server.engine.GET("/room-list", a.roomList)
	server.engine.GET("/room", a.room)
	server.engine.POST("/make-room", a.makeRoom)
	server.engine.GET("/enter-room", a.enterRoom)
}

func (a *api) roomList(context *gin.Context) {

}

func (a *api) makeRoom(context *gin.Context) {

}

func (a *api) room(context *gin.Context) {

}

func (a *api) enterRoom(context *gin.Context) {

}
