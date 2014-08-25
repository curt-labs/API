package main

import (
	"flag"
	"github.com/curt-labs/GoAPI/controllers/blog"
	"github.com/curt-labs/GoAPI/controllers/category"
	"github.com/curt-labs/GoAPI/controllers/customer"
	"github.com/curt-labs/GoAPI/controllers/dealers"
	"github.com/curt-labs/GoAPI/controllers/faq"
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/controllers/news"
	"github.com/curt-labs/GoAPI/controllers/part"
	"github.com/curt-labs/GoAPI/controllers/search"
	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/controllers/videos"
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
	m.Use(middleware.Meddler())
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
		r.Get("", category_ctlr.Parents)
		r.Get("/:id", category_ctlr.GetCategory)
		r.Get("/:id/subs", category_ctlr.SubCategories)
		r.Get("/:id/parts", category_ctlr.GetParts)
		r.Get("/:id/parts/:page/:count", category_ctlr.GetParts)
	})

	m.Group("/part", func(r martini.Router) {
		r.Get("/:part/vehicles", part_ctlr.Vehicles)
		r.Get("/:part/attributes", part_ctlr.Attributes)
		r.Get("/:part/reviews", part_ctlr.Reviews)
		r.Get("/:part/categories", part_ctlr.Categories)
		r.Get("/:part/content", part_ctlr.GetContent)
		r.Get("/:part/images", part_ctlr.Images)
		r.Get("/:part((.*?)\\.(PDF|pdf)$)", part_ctlr.InstallSheet) // Resolves: /part/11000.pdf
		r.Get("/:part/packages", part_ctlr.Packaging)
		r.Get("/:part/pricing", part_ctlr.Prices)
		r.Get("/:part/related", part_ctlr.GetRelated)
		r.Get("/:part/videos", part_ctlr.Videos)
		r.Get("/:part/:year/:make/:model", part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel", part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel/:config(.+)", part_ctlr.GetWithVehicle)
		r.Get("/:part", part_ctlr.Get)
	})

	m.Group("/customer", func(r martini.Router) {
		r.Post("", customer_ctlr.GetCustomer)
		r.Post("/user", customer_ctlr.GetUser)
		r.Post("/locations", customer_ctlr.GetLocations)
		r.Post("/users", customer_ctlr.GetUsers) // requires no user to be marked as sudo
		// Customer CMS endpoints

		// Content Types
		r.Get("/cms/content_types", customer_ctlr.GetAllContentTypes)

		// All Customer Content
		r.Get("/cms", customer_ctlr.GetAllContent)

		// Customer Part Content
		r.Get("/cms/part", customer_ctlr.AllPartContent)
		r.Get("/cms/part/:id", customer_ctlr.UniquePartContent)
		r.Post("/cms/part/:id", customer_ctlr.UpdatePartContent)
		r.Delete("/cms/part/:id", customer_ctlr.DeletePartContent)

		// Customer Category Content
		r.Get("/cms/category", customer_ctlr.AllCategoryContent)
		r.Get("/cms/category/:id", customer_ctlr.UniqueCategoryContent)
		r.Post("/cms/category/:id", customer_ctlr.UpdateCategoryContent)
		r.Delete("/cms/category/:id", customer_ctlr.DeleteCategoryContent)

		// Customer Content By Content Id
		r.Get("/cms/:id", customer_ctlr.GetContentById)
		r.Get("/cms/:id/revisions", customer_ctlr.GetContentRevisionsById)
		r.Post("/auth", customer_ctlr.UserAuthentication)
		r.Get("/auth", customer_ctlr.KeyedUserAuthentication)
	})

	m.Group("/faqs", func(r martini.Router) {
		r.Get("", faq_controller.GetAll)                        //get all faqs; takes optional sort param {sort=true} to sort by question
		r.Get("/questions", faq_controller.GetQuestions)        //get questions!{page, results} - all parameters are optional
		r.Get("/answers", faq_controller.GetAnswers)            //get answers!{page, results} - all parameters are optional
		r.Get("/search", faq_controller.Search)                 //takes {question, answer, page, results} - all parameters are optional
		r.Get("/(:id)", faq_controller.Get)                     //get by id {id}
		r.Put("", internalCors, faq_controller.Create)          //takes {question, answer}; returns object with new ID
		r.Post("/(:id)", internalCors, faq_controller.Update)   //{id, question and/or answer}
		r.Delete("/(:id)", internalCors, faq_controller.Delete) //{id}
		r.Delete("", internalCors, faq_controller.Delete)       //{?id=id}

	})
	m.Group("/blogs", func(r martini.Router) {
		r.Get("", blog_controller.GetAll)                      //sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/categories", blog_controller.GetAllCategories) //all categories; sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/category/:id", blog_controller.GetBlogCategory)
		r.Get("/search", blog_controller.Search) //search field = value e.g. /blogs/search?key=8AEE0620-412E-47FC-900A-947820EA1C1D&slug=cyclo
		r.Post("/categories", internalCors, blog_controller.CreateBlogCategory)
		r.Get("/:id", blog_controller.GetBlog)                     //get blog by {id}
		r.Post("/:id", internalCors, blog_controller.UpdateBlog)   //create {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} returns new id
		r.Put("", internalCors, blog_controller.CreateBlog)        //update {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} required{id}
		r.Delete("/:id", internalCors, blog_controller.DeleteBlog) //{?id=id}
		r.Delete("", internalCors, blog_controller.DeleteBlog)     //{id}
	})
	m.Group("/news", func(r martini.Router) {
		r.Get("", news_controller.GetAll)                      //get all news; takes optional sort param {sort=title||lead||content||startDate||endDate||active||slug} to sort by question
		r.Get("/titles", news_controller.GetTitles)            //get titles!{page, results} - all parameters are optional
		r.Get("/leads", news_controller.GetLeads)              //get leads!{page, results} - all parameters are optional
		r.Get("/search", news_controller.Search)               //takes {title, lead, content, publishStart, publishEnd, active, slug, page, results, page, results} - all parameters are optional
		r.Get("/:id", news_controller.Get)                     //get by id {id}
		r.Put("", internalCors, news_controller.Create)        //takes {question, answer}; returns object with new ID
		r.Post("/:id", internalCors, news_controller.Update)   //{id, question and/or answer}
		r.Delete("/:id", internalCors, news_controller.Delete) //{id}
		r.Delete("", internalCors, news_controller.Delete)     //{id}

	})

	m.Get("/search/part/:term", search_ctlr.SearchPart)
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
