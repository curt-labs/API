package main

import (
	"flag"
	"github.com/curt-labs/GoAPI/controllers/blog"
	"github.com/curt-labs/GoAPI/controllers/category"
	"github.com/curt-labs/GoAPI/controllers/customer"
	"github.com/curt-labs/GoAPI/controllers/customer_new"
	"github.com/curt-labs/GoAPI/controllers/dealers"
	"github.com/curt-labs/GoAPI/controllers/dealers_new"
	"github.com/curt-labs/GoAPI/controllers/faq"
	"github.com/curt-labs/GoAPI/controllers/forum"
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/controllers/news"
	"github.com/curt-labs/GoAPI/controllers/part"
	"github.com/curt-labs/GoAPI/controllers/salesrep"
	"github.com/curt-labs/GoAPI/controllers/search"
	"github.com/curt-labs/GoAPI/controllers/site"
	"github.com/curt-labs/GoAPI/controllers/site_new"
	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/controllers/videos"
	"github.com/curt-labs/GoAPI/controllers/webProperty"
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
	listenAddr = flag.String("http", ":8080", "http listen address")
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
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
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
		r.Post("/:id", category_ctlr.GetCategory)
		r.Get("/:id/subs", category_ctlr.SubCategories)
		r.Get("/:id/parts", category_ctlr.GetParts)
		r.Get("/:id/parts/:page/:count", category_ctlr.GetParts)
	})

	m.Group("/part", func(r martini.Router) {
		r.Get("/:part/vehicles", part_ctlr.Vehicles)
		r.Get("/:part/attributes", part_ctlr.Attributes)
		r.Get("/:part/reviews", part_ctlr.ActiveApprovedReviews)
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
		r.Get("", part_ctlr.All)
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
		//Customer prices
		r.Get("/prices/part/:id", internalCors, customer_ctlr_new.GetPricesByPart)         //{id}; id refers to partId
		r.Get("/prices/sale", internalCors, customer_ctlr_new.GetSales)                    //{start}{end}{id} -all required params; id refers to customerId
		r.Get("/prices/:id", internalCors, customer_ctlr_new.GetPrice)                     //{id}; id refers to {id} refers to customerPriceId
		r.Get("/prices", internalCors, customer_ctlr_new.GetAllPrices)                     //returns all {sort=field&direction=dir}
		r.Post("/prices/:id", internalCors, customer_ctlr_new.CreateUpdatePrice)           //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Put("/prices", internalCors, customer_ctlr_new.CreateUpdatePrice)                //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Delete("/prices/:id", internalCors, customer_ctlr_new.DeletePrice)               //{id} refers to customerPriceId
		r.Get("/pricesByCustomer/:id", internalCors, customer_ctlr_new.GetPriceByCustomer) //{id} refers to customerId; returns CustomerPrices

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
		r.Delete("/categories/:id", internalCors, blog_controller.DeleteBlogCategory)
		r.Get("/:id", blog_controller.GetBlog)                     //get blog by {id}
		r.Post("/:id", internalCors, blog_controller.UpdateBlog)   //create {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} returns new id
		r.Put("", internalCors, blog_controller.CreateBlog)        //update {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} required{id}
		r.Delete("/:id", internalCors, blog_controller.DeleteBlog) //{?id=id}
		r.Delete("", internalCors, blog_controller.DeleteBlog)     //{id}
	})

	m.Group("/forum", func(r martini.Router) {
		//groups
		r.Get("/groups", forum_ctlr.GetAllGroups)
		r.Get("/groups/:id", forum_ctlr.GetGroup)
		r.Post("/groups", forum_ctlr.AddGroup)
		r.Put("/groups/:id", forum_ctlr.UpdateGroup)
		r.Delete("/groups/:id", forum_ctlr.DeleteGroup)
		//topics
		r.Get("/topics", forum_ctlr.GetAllTopics)
		r.Get("/topics/:id", forum_ctlr.GetTopic)
		r.Post("/topics", forum_ctlr.AddTopic)
		r.Put("/topics/:id", forum_ctlr.UpdateTopic)
		r.Delete("/topics/:id", forum_ctlr.DeleteTopic)
		//threads
		r.Get("/threads", forum_ctlr.GetAllThreads)
		r.Get("/threads/:id", forum_ctlr.GetThread)
		r.Delete("/threads/:id", forum_ctlr.DeleteThread)
		//posts
		r.Get("/posts", forum_ctlr.GetAllPosts)
		r.Get("/posts/:id", forum_ctlr.GetPost)
		r.Post("/posts", forum_ctlr.AddPost)
		r.Put("/posts/:id", forum_ctlr.UpdatePost)
		r.Delete("/posts/:id", forum_ctlr.DeletePost)
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

	m.Group("/reviews", func(r martini.Router) {
		r.Get("", part_ctlr.GetAllReviews)
		r.Get("/:id", part_ctlr.GetReview)
		r.Put("", part_ctlr.SaveReview)
		r.Post("/:id", part_ctlr.SaveReview)
		r.Delete("/:id", part_ctlr.DeleteReview)
	})

	m.Group("/site", func(r martini.Router) {
		r.Get("/page/get", site.GetContentPage)
		r.Get("/page/primary", site.GetPrimaryContentPage)
		r.Get("/page/sitemap", site.GetSitemapCP)
		r.Get("/menu/primary", site.GetPrimaryMenu)
		r.Get("/menu/footer/sitemap", site.GetFooterSitemap)
		r.Get("/menu/get", site.GetMenuWithContent)
		r.Get("/menu/get/id", site.GetMenuByContentId)
		r.Get("/menu/sitemap", site.GetMenuSitemap)

	})

	m.Group("/landingPage", func(r martini.Router) {
		r.Get("/get", site.GetLandingPage)
	})
	m.Group("/webProperties", func(r martini.Router) {
		r.Post("/note/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyNote)                                       //updates when an id is present; otherwise, creates
		r.Put("/note", internalCors, webProperty_controller.CreateUpdateWebPropertyNote)                                            //updates when an id is present; otherwise, creates
		r.Delete("/note/:id", internalCors, webProperty_controller.DeleteWebPropertyNote)                                           //{id}
		r.Get("/note/:id", webProperty_controller.GetWebPropertyNote)                                                               //{id}
		r.Post("/requirementCheck/:id", internalCors, internalCors, webProperty_controller.CreateUpdateWebPropertyRequirementCheck) //updates when an id is present; otherwise, creates
		r.Put("/requirementCheck", internalCors, webProperty_controller.CreateUpdateWebPropertyRequirementCheck)                    //updates when an id is present; otherwise, creates
		r.Delete("/requirementCheck/:id", internalCors, webProperty_controller.DeleteWebPropertyRequirementCheck)                   //{id}
		r.Get("/requirementCheck/:id", webProperty_controller.GetWebPropertyRequirementCheck)                                       //{id}
		r.Post("/type/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyType)                                       //updates when an id is present; otherwise, creates
		r.Put("/type", internalCors, webProperty_controller.CreateUpdateWebPropertyType)                                            //updates when an id is present; otherwise, creates
		r.Delete("/type/:id", internalCors, webProperty_controller.DeleteWebPropertyType)                                           //{id}
		r.Get("/type/:id", webProperty_controller.GetWebPropertyType)                                                               //{id}
		r.Post("/requirement/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyRequirement)                         //updates when an id is present; otherwise, creates
		r.Put("/requirement", internalCors, webProperty_controller.CreateUpdateWebPropertyRequirement)                              //updates when an id is present; otherwise, creates
		r.Delete("/requirement/:id", internalCors, webProperty_controller.DeleteWebPropertyRequirement)                             //{id}
		r.Get("/requirement/:id", webProperty_controller.GetWebPropertyRequirement)                                                 //{id}
		r.Get("/search", internalCors, webProperty_controller.Search)
		r.Get("/type", webProperty_controller.GetAllTypes)               //all tyeps
		r.Get("/note", webProperty_controller.GetAllNotes)               //all notes
		r.Get("/requirement", webProperty_controller.GetAllRequirements) //requirements
		r.Get("", internalCors, webProperty_controller.GetAll)
		r.Get("/:id", internalCors, webProperty_controller.Get)                      //?id=id
		r.Delete("/:id", internalCors, webProperty_controller.DeleteWebProperty)     //{id}
		r.Post("/:id", internalCors, webProperty_controller.CreateUpdateWebProperty) //
		r.Put("", internalCors, webProperty_controller.CreateUpdateWebProperty)      //can create notes(text) and requirements (requirement, by requirement=requirementID) while creating a property

	})

	m.Group("/salesrep", func(r martini.Router) {
		r.Get("", salesrep.GetAllSalesReps)
		r.Post("", salesrep.AddSalesRep)
		r.Get("/:id", salesrep.GetSalesRep)
		r.Put("/:id", salesrep.UpdateSalesRep)
		r.Delete("/:id", salesrep.DeleteSalesRep)
	})

	m.Get("/search/:term", search_ctlr.Search)
	m.Get("/videos", videos_ctlr.DistinctVideos)

	//NEW Customer & Dealer endpoints - Seems to work. Feeling brave?
	m.Group("/new", func(r martini.Router) {
		m.Group("/customer", func(r martini.Router) {
			// r.Get("", internalCors, customer_ctlr_new.GetCustomer)
			r.Post("", internalCors, customer_ctlr_new.GetCustomer)
			r.Post("/user", customer_ctlr_new.GetUser)
			r.Get("/user/:id", customer_ctlr_new.GetUserById)
			r.Get("/users", customer_ctlr_new.GetUsers)
			r.Post("/users", customer_ctlr_new.GetUsers)
			r.Post("/user/register", customer_ctlr_new.RegisterUser)
			r.Post("/user/resetPassword", customer_ctlr_new.ResetPassword)
			r.Post("/user/changePassword", customer_ctlr_new.ChangePassword)
			r.Delete("/user/:id", customer_ctlr_new.DeleteCustomerUser)
			r.Delete("/allUsersByCustomerID/:id", customer_ctlr_new.DeleteCustomerUsersByCustomerID) //Takes CustomerID (UUID)---danger!

			// r.Get("/locations", internalCors, customer_ctlr_new.GetLocations)
			r.Post("/locations", internalCors, customer_ctlr_new.GetLocations)
			r.Get("/price/:id", internalCors, customer_ctlr_new.GetCustomerPrice)           //{part id}
			r.Get("/cartRef/:id", internalCors, customer_ctlr_new.GetCustomerCartReference) //{part id}
			r.Post("/auth", customer_ctlr_new.UserAuthentication)
			r.Get("/auth", customer_ctlr_new.KeyedUserAuthentication)

			// Customer CMS endpoints
			// All Customer Content
			r.Get("/cms", customer_ctlr_new.GetAllContent)
			// Content Types
			r.Get("/cms/content_types", customer_ctlr.GetAllContentTypes)

			// Customer Part Content
			r.Get("/cms/part", customer_ctlr_new.AllPartContent)
			r.Get("/cms/part/:id", customer_ctlr_new.UniquePartContent)
			r.Post("/cms/part/:id", customer_ctlr_new.UpdatePartContent)
			r.Delete("/cms/part/:id", customer_ctlr_new.DeletePartContent)

			// Customer Category Content
			r.Get("/cms/category", customer_ctlr_new.AllCategoryContent)
			r.Get("/cms/category/:id", customer_ctlr_new.UniqueCategoryContent)
			r.Post("/cms/category/:id", customer_ctlr_new.UpdateCategoryContent)
			r.Delete("/cms/category/:id", customer_ctlr_new.DeleteCategoryContent)

			// Customer Content By Content Id
			r.Get("/cms/:id", customer_ctlr_new.GetContentById)
			r.Get("/cms/:id/revisions", customer_ctlr_new.GetContentRevisionsById)

			//Customer prices
			r.Get("/prices/part/:id", internalCors, customer_ctlr_new.GetPricesByPart)         //{id}; id refers to partId
			r.Get("/prices/sale", internalCors, customer_ctlr_new.GetSales)                    //{start}{end}{id} -all required params; id refers to customerId
			r.Get("/prices/:id", internalCors, customer_ctlr_new.GetPrice)                     //{id}; id refers to {id} refers to customerPriceId
			r.Get("/prices", internalCors, customer_ctlr_new.GetAllPrices)                     //returns all {sort=field&direction=dir}
			r.Post("/prices/:id", internalCors, customer_ctlr_new.CreateUpdatePrice)           //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
			r.Put("/prices", internalCors, customer_ctlr_new.CreateUpdatePrice)                //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
			r.Delete("/prices/:id", internalCors, customer_ctlr_new.DeletePrice)               //{id} refers to customerPriceId
			r.Get("/pricesByCustomer/:id", internalCors, customer_ctlr_new.GetPriceByCustomer) //{id} refers to customerId; returns CustomerPrices

		})
		m.Group("/dealers", func(r martini.Router) {
			r.Get("/etailer", internalCors, dealers_ctlr_new.GetEtailers)
			r.Get("/local", internalCors, dealers_ctlr_new.GetLocalDealers)
			r.Get("/local/regions", internalCors, dealers_ctlr_new.GetLocalRegions)     //move to dealers
			r.Get("/local/tiers", internalCors, dealers_ctlr_new.GetLocalDealerTiers)   //move to dealers
			r.Get("/local/types", internalCors, dealers_ctlr_new.GetLocalDealerTypes)   //move to dealers
			r.Get("/etailer/platinum", internalCors, dealers_ctlr_new.PlatinumEtailers) //move to dealers
			r.Get("/location/:id", internalCors, dealers_ctlr_new.GetLocationById)      //move to dealers
			r.Get("/search/:search", internalCors, dealers_ctlr_new.SearchLocations)
			r.Get("/search/type", internalCors, dealers_ctlr_new.SearchLocationsByType)
			r.Get("/search/type/:search", internalCors, dealers_ctlr_new.SearchLocationsByType)
			r.Get("/search/geo", internalCors, dealers_ctlr_new.SearchLocationsByLatLng)
			r.Get("/search/geo/:latitude/:longitude", internalCors, dealers_ctlr_new.SearchLocationsByLatLng)
		})
		m.Get("/dealer/location/:id", internalCors, dealers_ctlr_new.GetLocationById)

		//New, more better site endpoints - will not work with current .mfg travesty
		m.Group("/site", func(r martini.Router) {
			m.Group("/menu", func(r martini.Router) {
				r.Get("/all", site_new.GetAllMenus)
				r.Get("/:id", site_new.GetMenu)                      //may pass id (int) or name(string)
				r.Get("/contents/:id", site_new.GetMenuWithContents) //may pass id (int) or name(string)
				r.Put("", site_new.SaveMenu)
				r.Post("/:id", site_new.SaveMenu)
				r.Delete("/:id", site_new.DeleteMenu)
			})
			m.Group("/content", func(r martini.Router) {
				r.Get("/all", site_new.GetAllContents)
				r.Get("/:id", site_new.GetContent) //may pass id (int) or slug(string)
				r.Get("/:id/revisions", site_new.GetContentRevisions)
				r.Put("", site_new.SaveContent)
				r.Post("/:id", site_new.SaveContent)
				r.Delete("/:id", site_new.DeleteContent)
			})
		})

	})

	// m.Get("/customer/auth/keyedUser", internalCors, customer_ctlr_new.KeyedUserAuthentication)
	// m.Post("/customer/auth/new", internalCors, customer_ctlr_new.UserAuthentication)
	// // m.Get("/customer/auth/oldKeyedUser", internalCors, customer_ctlr.KeyedUserAuthentication)
	// m.Get("/customer/auth/resetAuth", internalCors, customer_ctlr_new.ResetAuthentication)

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
		// log.Println("hit any")
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
