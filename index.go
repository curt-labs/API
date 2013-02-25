package main

import (
	"./controllers/part"
	"./controllers/vehicle"
	"./helpers/auth"
	"./plate"
	"flag"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("http", ":8080", "http listen address")
)

const (
	port = "80"
)

func main() {
	flag.Parse()

	server := plate.NewServer("doughboy")

	server.AddFilter(auth.AuthHandler)
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
	}).NoFilter()

	server.Get("/vehicle", vehicle_ctlr.Year)
	server.Get("/vehicle/:year", vehicle_ctlr.Make)
	server.Get("/vehicle/:year/:make", vehicle_ctlr.Model)
	server.Get("/vehicle/:year/:make/:model", vehicle_ctlr.Submodel)
	server.Get("/vehicle/:year/:make/:model/:submodel", vehicle_ctlr.Config)
	server.Get("/vehicle/:year/:make/:model/:submodel/:config(.+)", vehicle_ctlr.Config)

	server.Get("/part/:part/vehicles", part_ctlr.Vehicles)
	server.Get("/part/:part", part_ctlr.Get)


	http.Handle("/", server)
	http.ListenAndServe(*listenAddr, nil)

	log.Println("Server running on port " + *listenAddr)
}
