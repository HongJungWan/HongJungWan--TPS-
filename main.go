package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	config "go_chat/config"
	"go_chat/network"
	"go_chat/repository"
	"go_chat/service"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*pathFlag)
	fmt.Println(cfg)

	if rep, err := repository.NewRepository(cfg); err != nil {
		panic(err)
	} else {
		server := network.NewServer(service.NewService(rep), *port)
		server.StartServer()
	}
}
