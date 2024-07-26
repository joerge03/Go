package main

func main() {
	newServer := NewAPIServer(":8082")
	newServer.Run()
}
