package main

import (
	"flag"
	"fmt"
	config "go_chat/config"
	"go_chat/repository"
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
		fmt.Println(rep)
	}
}
