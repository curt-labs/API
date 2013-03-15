package main

import (
	"./controllers/category"
	"./controllers/customer"
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

	CorsHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		return
	}
)

const (
	port = "80"
)

func main() {
	flag.Parse()

	server := plate.NewServer("doughboy")
	server.Logging = true

	server.AddFilter(auth.AuthHandler)
	server.AddFilter(CorsHandler)
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
	}).NoFilter()

	server.Get("/.status", func(w http.ResponseWriter, r *http.Request) {
		server.StatusService.GetStatus(w, r)
	}).NoFilter()

	server.Get("/vehicle", vehicle_ctlr.Year)
	server.Get("/vehicle/:year", vehicle_ctlr.Make)
	server.Get("/vehicle/:year/:make", vehicle_ctlr.Model)
	server.Get("/vehicle/:year/:make/:model", vehicle_ctlr.Submodel)
	server.Get("/vehicle/:year/:make/:model/:submodel", vehicle_ctlr.Config)
	server.Get("/vehicle/:year/:make/:model/:submodel/:config(.+)", vehicle_ctlr.Config)

	server.Get("/category", category_ctlr.Parents)
	server.Get("/category/:id", category_ctlr.GetCategory)
	server.Get("/category/:id/subs", category_ctlr.SubCategories)
	server.Get("/category/:id/parts", category_ctlr.GetParts)
	server.Get("/category/:id/parts/:page/:count", category_ctlr.GetParts)

	server.Get("/part/:part/vehicles", part_ctlr.Vehicles)
	server.Get("/part/:part", part_ctlr.Get)

	server.Post("/customer/auth", customer_ctlr.UserAuthentication).NoFilter()
	server.Get("/customer/auth", customer_ctlr.KeyedUserAuthentication).NoFilter()

	http.Handle("/", server)
	http.ListenAndServe(*listenAddr, nil)

	log.Println("Server running on port " + *listenAddr)
}
