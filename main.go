package main

func main() {
	apiServer := newAPIServer(":8080")
	apiServer.startServer()
}
