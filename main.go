package main

import (
	"time"
	"file-service.com/server"
)

func main() {
	time.Sleep(1 * time.Second)
	server.New(":4040").Listen()
}