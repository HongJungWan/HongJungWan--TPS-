// 서버 설정 및 초기화를 담당하고 CORS 설정을 적용하여 서버를 실행하는 역할

package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_chat/repository"
	"go_chat/service"
)

type Server struct {
	engine *gin.Engine

	service    *service.Service
	repository *repository.Repository

	port string
	ip   string
}

func NewServer(service *service.Service, port string) *Server {
	server := &Server{engine: gin.New(), service: service, port: port}

	server.engine.Use(gin.Logger())
	server.engine.Use(gin.Recovery())
	server.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"ORIGIN", "Content-Length", "Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		ExposeHeaders:    []string{"ORIGIN", "Content-Length", "Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	registerServer(server)

	return server
}

func (server *Server) StartServer() error {
	return server.engine.Run(server.port)
}
