package main

import (
	"flag"
	"fmt"
	config "go_chat/config"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()
	configSetUp := config.NewConfig(*pathFlag)
	fmt.Println(configSetUp)
}
