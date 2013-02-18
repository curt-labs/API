package main

import (
	"./controllers"
	"./helpers/auth"
	"./plate"
	"net/http"
)

func main() {
	server := plate.NewServer("doughboy")
	plate.DefaultAuthHandler = auth.AuthHandler

	server.Get("/", controllers.Index).Secure()

	//session_key := "your key here"

	http.Handle("/", server)
	http.ListenAndServe(":8080", nil)
}
