package main

import (
	"romano.com/server"
)

var serverManager server.Server

func main() {
	serverManager := server.New(":4040");
	serverManager.CreateListener();
}