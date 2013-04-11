package main

import (
	"./controllers/category"
	"./controllers/customer"
	"./controllers/dealers"
	"./controllers/part"
	"./controllers/vehicle"
	"./controllers/videos"
	"./helpers/auth"
	"./helpers/plate"
	"flag"
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

/**
 * All GET routes require either public or private api keys to be passed in.
 *
 * All POST routes require private api keys to be passed in.
 */
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
	server.Get("/part/:part/attributes", part_ctlr.Attributes)
	server.Get("/part/:part/reviews", part_ctlr.Reviews)
	server.Get("/part/:part/categories", part_ctlr.Categories)
	server.Get("/part/:part/content", part_ctlr.GetContent)
	server.Get("/part/:part/images", part_ctlr.Images)
	server.Get("/part/:part((.*?)\\.(PDF|pdf)$)", part_ctlr.InstallSheet).NoFilter() // Resolves: /part/11000.pdf
	server.Get("/part/:part/packages", part_ctlr.Packaging)
	server.Get("/part/:part/pricing", part_ctlr.Prices)
	// server.Get("/part/:part/related", part_ctlr.Get)
	server.Get("/part/:part/videos", part_ctlr.Videos)
	server.Get("/part/:part", part_ctlr.Get)

	server.Post("/customer/auth", customer_ctlr.UserAuthentication).NoFilter()
	server.Get("/customer/auth", customer_ctlr.KeyedUserAuthentication).NoFilter()

	server.Post("/customer/locations", customer_ctlr.GetLocations)
	server.Post("/customer/users", customer_ctlr.GetUsers) // Requires a user to be marked as sudo

	/**
	 * Video
	 */
	server.Get("/videos", videos_ctlr.DistinctVideos).NoFilter()

	/**** INTERNAL USE ONLY ****/
	server.Get("/dealers/etailer", dealers_ctlr.Etailers).NoFilter()
	server.Get("/dealers/local", dealers_ctlr.LocalDealers).NoFilter()
	server.Get("/dealers/local/regions", dealers_ctlr.LocalRegions).NoFilter()

	http.Handle("/", server)
	http.ListenAndServe(*listenAddr, nil)
}
