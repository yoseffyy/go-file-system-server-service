package main

import (
	"github.com/yosef32/go-file-system-server-service/server"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	server.NewServer().Serve()
}
