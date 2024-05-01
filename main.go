package main

import (
	"go-api/server"
)

func main() {
	//wg := sync.WaitGroup{}
	apiServer := server.NewAPIServer(":8080")
	apiServer.StartServer()
	//	wg.Add(1)
	//	go simulation.SimulateAccountCreation()
	//	wg.Wait()
}
