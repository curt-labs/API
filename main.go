package main

import (
	//"./controllers/home"
	//"./controllers/vehicle"
	//"./filters/access_control"
	//"github.com/astaxie/beego"
	"fmt"
	"log"
	"net/http"
)

func Middleware(w http.ResponseWriter, r *http.Request) {
	authChan := make(chan int)
	queryChan := make(chan int)

	log.Println(r.URL)
	go func() {

		authChan <- 1
	}()

	go func() {

		queryChan <- 1
	}()

	<-authChan
	<-queryChan

	log.Println("passed it")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{}")
}

func main() {

	http.HandleFunc("/", Middleware)
	http.ListenAndServe(":8080", nil)

}
