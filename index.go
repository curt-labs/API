package main

import (
	"flag"
	"github.com/curt-labs/GoAPI/controllers/category"
	"github.com/curt-labs/GoAPI/controllers/customer"
	"github.com/curt-labs/GoAPI/controllers/dealers"
	// "github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/controllers/part"
	"github.com/curt-labs/GoAPI/controllers/search"
	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/controllers/videos"
	"github.com/curt-labs/GoAPI/helpers/auth"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/gorelic"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	listenAddr  = flag.String("http", ":8080", "http listen address")
	CorsHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		return
	}
)

/**
 * All GET routes require either public or private api keys to be passed in.
 *
 * All POST routes require private api keys to be passed in.
 */
func main() {
	flag.Parse()

	err := database.PrepareAll()
	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()
	gorelic.InitNewrelicAgent("5fbc49f51bd658d47b4d5517f7a9cb407099c08c", "GoAPI", false)
	m.Use(gorelic.Handler)
	m.Use(gzip.All())
	// m.Use(middleware.Meddler())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	store := sessions.NewCookieStore([]byte("api_secret_session"))
	m.Use(sessions.Sessions("api_sessions", store))
	m.Use(MapEncoder)

	internalCors := cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*.curtmfg.com", "http://*.curtmfg.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	m.Post("/vehicle", vehicle.Query)

	m.Group("/category", func(r martini.Router) {
		r.Get("", auth.AuthHandler, category_ctlr.Parents)
		r.Get("/:id", auth.AuthHandler, category_ctlr.GetCategory)
		r.Get("/:id/subs", auth.AuthHandler, category_ctlr.SubCategories)
		r.Get("/:id/parts", auth.AuthHandler, category_ctlr.GetParts)
		r.Get("/:id/parts/:page/:count", auth.AuthHandler, category_ctlr.GetParts)
	})

	m.Group("/part", func(r martini.Router) {
		r.Get("/:part/vehicles", auth.AuthHandler, part_ctlr.Vehicles)
		r.Get("/:part/attributes", auth.AuthHandler, part_ctlr.Attributes)
		r.Get("/:part/reviews", auth.AuthHandler, part_ctlr.Reviews)
		r.Get("/:part/categories", auth.AuthHandler, part_ctlr.Categories)
		r.Get("/:part/content", auth.AuthHandler, part_ctlr.GetContent)
		r.Get("/:part/images", auth.AuthHandler, part_ctlr.Images)
		r.Get("/:part((.*?)\\.(PDF|pdf)$)", auth.AuthHandler, part_ctlr.InstallSheet) // Resolves: /part/11000.pdf
		r.Get("/:part/packages", auth.AuthHandler, part_ctlr.Packaging)
		r.Get("/:part/pricing", auth.AuthHandler, part_ctlr.Prices)
		r.Get("/:part/related", auth.AuthHandler, part_ctlr.GetRelated)
		r.Get("/:part/videos", auth.AuthHandler, part_ctlr.Videos)
		r.Get("/:part/:year/:make/:model", auth.AuthHandler, part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel", auth.AuthHandler, part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel/:config(.+)", auth.AuthHandler, part_ctlr.GetWithVehicle)
		r.Get("/:part", auth.AuthHandler, part_ctlr.Get)
	})

	m.Group("/customer", func(r martini.Router) {
		r.Post("", auth.AuthHandler, customer_ctlr.GetCustomer)
		r.Post("/user", auth.AuthHandler, customer_ctlr.GetUser)
		r.Post("/locations", auth.AuthHandler, customer_ctlr.GetLocations)
		r.Post("/users", auth.AuthHandler, customer_ctlr.GetUsers) // requires no user to be marked as sudo
		// Customer CMS endpoints

		// Content Types
		r.Get("/cms/content_types", auth.AuthHandler, customer_ctlr.GetAllContentTypes)

		// All Customer Content
		r.Get("/cms", auth.AuthHandler, customer_ctlr.GetAllContent)

		// Customer Part Content
		r.Get("/cms/part", auth.AuthHandler, customer_ctlr.AllPartContent)
		r.Get("/cms/part/:id", auth.AuthHandler, customer_ctlr.UniquePartContent)
		r.Post("/cms/part/:id", auth.AuthHandler, customer_ctlr.UpdatePartContent)
		r.Delete("/cms/part/:id", auth.AuthHandler, customer_ctlr.DeletePartContent)

		// Customer Category Content
		r.Get("/cms/category", auth.AuthHandler, customer_ctlr.AllCategoryContent)
		r.Get("/cms/category/:id", auth.AuthHandler, customer_ctlr.UniqueCategoryContent)
		r.Post("/cms/category/:id", auth.AuthHandler, customer_ctlr.UpdateCategoryContent)
		r.Delete("/cms/category/:id", auth.AuthHandler, customer_ctlr.DeleteCategoryContent)

		// Customer Content By Content Id
		r.Get("/cms/:id", auth.AuthHandler, customer_ctlr.GetContentById)
		r.Get("/cms/:id/revisions", auth.AuthHandler, customer_ctlr.GetContentRevisionsById)
		r.Post("/auth", customer_ctlr.UserAuthentication)
		r.Get("/auth", customer_ctlr.KeyedUserAuthentication)
	})

	m.Get("/search/part/:term", auth.AuthHandler, search_ctlr.SearchPart)
	m.Get("/videos", videos_ctlr.DistinctVideos)

	/**** INTERNAL USE ONLY ****/
	// These endpoints will not work to the public eye when deployed on CURT's
	// servers. We will have restrictions in place to prevent access...sorry :/
	m.Group("/dealers", func(r martini.Router) {
		r.Get("/etailer", internalCors, dealers_ctlr.Etailers)
		r.Get("/etailer/platinum", internalCors, dealers_ctlr.PlatinumEtailers)
		r.Get("/local", internalCors, dealers_ctlr.LocalDealers)
		r.Get("/local/regions", internalCors, dealers_ctlr.LocalRegions)
		r.Get("/local/tiers", internalCors, dealers_ctlr.LocalDealerTiers)
		r.Get("/local/types", internalCors, dealers_ctlr.LocalDealerTypes)
		r.Get("/search", internalCors, dealers_ctlr.SearchLocations)
		r.Get("/search/:search", internalCors, dealers_ctlr.SearchLocations)
		r.Get("/search/type", internalCors, dealers_ctlr.SearchLocationsByType)
		r.Get("/search/type/:search", internalCors, dealers_ctlr.SearchLocationsByType)
		r.Get("/search/geo", internalCors, dealers_ctlr.SearchLocationsByLatLng)
		r.Get("/search/geo/:latitude/:longitude", internalCors, dealers_ctlr.SearchLocationsByLatLng)
	})
	m.Get("/dealer/location/:id", internalCors, dealers_ctlr.GetLocation)

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
	})

	m.Any("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hit any")
	})

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      m,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on 127.0.0.1%s\n", *listenAddr)
	log.Fatal(srv.ListenAndServe())
}

var rxAccept = regexp.MustCompile(`(?:xml|html|plain|json)\/?$`)

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept == "*/*" {
		accept = r.Header.Get("Content-Type")
	}
	matches := rxAccept.FindStringSubmatch(accept)

	dt := "json"
	if len(matches) == 1 {
		dt = matches[0]
	}
	switch dt {
	case "xml":

		c.MapTo(encoding.XmlEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml")
	case "plain":
		c.MapTo(encoding.TextEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain")
	case "html":
		c.MapTo(encoding.TextEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "text/html")
	default:
		c.MapTo(encoding.JsonEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}
