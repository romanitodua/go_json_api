package main

import "go-api/server"

func main() {
	apiServer := server.NewAPIServer(":8080")
	apiServer.StartServer()
}
