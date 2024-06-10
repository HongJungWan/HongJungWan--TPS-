package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	network := &Network{engin: gin.New()}

	network.engin.Use(gin.Logger())
	network.engin.Use(gin.Recovery())
	network.engin.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	return network
}

func (network *Network) StartServer() error {
	log.Println("Start Server")
	return network.engin.Run(":8080")
}
