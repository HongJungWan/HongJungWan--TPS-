package main

import (
	"go_chat/network"
	"log"
)

func init() {
	log.Println("init test")
}

func main() {
	network := network.NewServer()
	network.StartServer()
}
