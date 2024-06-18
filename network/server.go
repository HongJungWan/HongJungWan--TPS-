// 서버의 API 엔드포인트를 정의하고 요청 핸들러를 등록하는 역할

package network

import (
	"github.com/gin-gonic/gin"
	"go_chat/types"
	"net/http"
)

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
	if res, err := a.server.service.RoomList(); err != nil {
		Response(context, http.StatusInternalServerError, err.Error())
	} else {
		Response(context, http.StatusOK, res)
	}
}

func (a *api) makeRoom(context *gin.Context) {
	var req types.BodyRoomReq

	if err := context.ShouldBindJSON(&req); err != nil {
		Response(context, http.StatusUnprocessableEntity, err.Error())
	} else if err = a.server.service.MakeRoom(req.Name); err != nil {
		Response(context, http.StatusInternalServerError, err.Error())
	} else {
		Response(context, http.StatusOK, "Success")
	}
}

func (a *api) room(context *gin.Context) {
	var req types.FormRoomReq

	if err := context.ShouldBindQuery(&req); err != nil {
		Response(context, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := a.server.service.Room(req.Name); err != nil {
		Response(context, http.StatusInternalServerError, err.Error())
	} else {
		Response(context, http.StatusOK, res)
	}
}

func (a *api) enterRoom(context *gin.Context) {
	var req types.FormRoomReq

	if err := context.ShouldBindQuery(&req); err != nil {
		Response(context, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := a.server.service.EnterRoom(req.Name); err != nil {
		Response(context, http.StatusInternalServerError, err.Error())
	} else {
		Response(context, http.StatusOK, res)
	}
}
