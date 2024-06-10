package network

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	network := &Network{engin: gin.New()}
	return network
}

func (network *Network) StartServer() error {
	log.Println("Start Server")
	return network.engin.Run(":8080")
}
