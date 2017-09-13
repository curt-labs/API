package main

import (
	"flag"

	"github.com/curt-labs/API/controllers/acesFile"
	"github.com/curt-labs/API/controllers/applicationGuide"
	"github.com/curt-labs/API/controllers/brand"
	"github.com/curt-labs/API/controllers/cache"
	"github.com/curt-labs/API/controllers/cartIntegration"
	"github.com/curt-labs/API/controllers/category"
	"github.com/curt-labs/API/controllers/contact"
	"github.com/curt-labs/API/controllers/customer"
	"github.com/curt-labs/API/controllers/dealers"
	"github.com/curt-labs/API/controllers/geography"
	"github.com/curt-labs/API/controllers/luverne"
	"github.com/curt-labs/API/controllers/middleware"
	"github.com/curt-labs/API/controllers/news"
	"github.com/curt-labs/API/controllers/part"
	"github.com/curt-labs/API/controllers/search"
	"github.com/curt-labs/API/controllers/vehicle"
	"github.com/curt-labs/API/controllers/videos"
	"github.com/curt-labs/API/controllers/vinLookup"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	// "github.com/martini-contrib/gzip"
	"log"
	"net/http"
	"time"

	"github.com/martini-contrib/sessions"
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

	m := martini.Classic()
	// gorelic.InitNewrelicAgent("5fbc49f51bd658d47b4d5517f7a9cb407099c08c", "API", false)
	// m.Use(gorelic.Handler)
	// m.Use(gzip.All())
	m.Use(middleware.Meddler())
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Access-Control-Allow-Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Type", "Accept", "Access-Control-Allow-Origin", "Authorization"},
		AllowCredentials: false,
	}))

	store := sessions.NewCookieStore([]byte("api_secret_session"))
	m.Use(sessions.Sessions("api_sessions", store))
	m.Use(encoding.MapEncoder)

	m.Group("/aces", func(r martini.Router) {
		r.Get("/:version", acesFile.GetAcesFile)
	})

	m.Group("/apiKeyTypes", func(r martini.Router) {
		r.Get("", Deprecated)
	})

	//Creating, updating, and deleting Appguides are all handled in GoAdmin directly
	m.Group("/applicationGuide", func(r martini.Router) {
		r.Get("/website/:id", applicationGuide.GetApplicationGuidesByWebsite)
		r.Get("/:id", applicationGuide.GetApplicationGuide)
		r.Delete("/:id", Deprecated)
		r.Post("", Deprecated)
	})

	//Creating, updating, and deleting all Blog related objects are handled in GoAdmin directly
	m.Group("/blogs", func(r martini.Router) {
		r.Get("", Deprecated)            //sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/categories", Deprecated) //all categories; sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/category/:id", Deprecated)
		r.Get("/search", Deprecated) //search field = value e.g. /blogs/search?key=8AEE0620-412E-47FC-900A-947820EA1C1D&slug=cyclo
		r.Post("/categories", Deprecated)
		r.Delete("/categories/:id", Deprecated)
		r.Get("/:id", Deprecated)    //get blog by {id}
		r.Put("/:id", Deprecated)    //create {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} returns new id
		r.Post("", Deprecated)       //update {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} required{id}
		r.Delete("/:id", Deprecated) //{?id=id}
		r.Delete("", Deprecated)     //{id}
	})

	//Creating, updating, and deleting Brands is not handled anywhere, but it does need to be
	//locked down for security.
	m.Group("/brands", func(r martini.Router) {
		r.Get("", middleware.InternalKeyAuthentication, brand_ctlr.GetAllBrands)
		r.Post("", middleware.InternalKeyAuthentication, brand_ctlr.CreateBrand)
		r.Get("/:id", brand_ctlr.GetBrand)
		r.Put("/:id", middleware.InternalKeyAuthentication, brand_ctlr.UpdateBrand)
		r.Delete("/:id", middleware.InternalKeyAuthentication, brand_ctlr.DeleteBrand)
	})

	m.Group("/category", func(r martini.Router) {
		r.Get("/:id/parts", category_ctlr.GetCategoryParts)
		r.Get("/:id", category_ctlr.GetCategory)
		r.Get("", category_ctlr.GetCategoryTree)
	})

	//Creating, updating, and deleting all Contact related entities is handled
	//in GoAdmin directly
	m.Group("/contact", func(r martini.Router) {
		m.Group("/types", func(r martini.Router) {
			r.Get("/receivers/:id", contact.GetReceiversByContactType)
			r.Get("", contact.GetAllContactTypes)
			r.Get("/:id", contact.GetContactType)
			r.Post("", Deprecated)
			r.Put("/:id", Deprecated)
			r.Delete("/:id", Deprecated)
		})
		m.Group("/receivers", func(r martini.Router) {
			r.Get("", Deprecated)
			r.Get("/:id", Deprecated)
			r.Post("", Deprecated)
			r.Put("/:id", Deprecated)
			r.Delete("/:id", Deprecated)
		})

		r.Get("", contact.GetAllContacts)
		r.Get("/:id", contact.GetContact)
		r.Post("/:contactTypeID", middleware.InternalKeyAuthentication, contact.AddDealerContact)
		r.Put("/:id", middleware.InternalKeyAuthentication, contact.UpdateContact)
		r.Delete("/:id", middleware.InternalKeyAuthentication, contact.DeleteContact)
	})

	//These shopify endpoints appear to not be used at all. Due to their customer related nature,
	//They are being locked down for security.
	m.Group("/shopify/customers", func(r martini.Router) {
		// Customers - shop endpoints
		r.Get("", Deprecated)
		r.Post("", Deprecated)
		r.Get("/search", Deprecated)
		r.Get("/:id", Deprecated)
		r.Put("/:id", Deprecated)
		r.Delete("/:id", Deprecated)
		r.Get("/:id/orders", Deprecated)

		// Addresses
		r.Get("/:id/addresses", Deprecated)
		r.Get("/:id/addresses/:address", Deprecated)
		r.Post("/:id/addresses", Deprecated)
		r.Put("/:id/addresses/:address/default", Deprecated)
		r.Put("/:id/addresses/:address", Deprecated)
		r.Delete("/:id/addresses/:address", Deprecated)

	})

	m.Group("/shopify/order", func(r martini.Router) {
		// Orders
		r.Post("/order", Deprecated)
	})

	m.Group("/shopify/account", func(r martini.Router) {
		// Account - user endpoints
		r.Get("", Deprecated)
		r.Post("", Deprecated)
		r.Put("", Deprecated)
		r.Post("/login", Deprecated)
	})

	//Used on the dealer site, no lockdown for now
	m.Group("/cartIntegration", func(r martini.Router) {
		r.Get("/part/:part", cartIntegration.GetPartPricesByPartID)
		r.Get("/part", cartIntegration.GetAllPartPrices)
		r.Get("/count", cartIntegration.GetPricingCount)
		r.Get("", cartIntegration.GetPricing)
		r.Get("/:page/:count", cartIntegration.GetPricingPaged)
		r.Post("/part", Deprecated)
		r.Put("/part", Deprecated)
		r.Get("/priceTypes", Deprecated)

		r.Post("/resetToMap", Deprecated)
		r.Post("/global/:type/:percentage", Deprecated)

		r.Post("/upload", cartIntegration.Upload)
		r.Post("/download", cartIntegration.Download)

	})

	//Cache should definitely be locked down
	m.Group("/cache", func(r martini.Router) { // different endpoint because partial matching matches this to another excused route
		r.Get("/key", cache.GetByKey)
		r.Get("/keys", cache.GetKeys)
		r.Delete("/keys", middleware.InternalKeyAuthentication, cache.DeleteKey)
	})

	//No lockdown of customer related endpoints for now
	m.Group("/cust", func(r martini.Router) { // different endpoint because partial matching matches this to another excused route
		r.Post("/user/changePassword", customer_ctlr.ChangePassword)
	})

	//No lockdown of customer related endpoints for now
	m.Group("/customer", func(r martini.Router) {
		r.Get("", customer_ctlr.GetCustomer)
		r.Post("", customer_ctlr.GetCustomer)

		r.Post("/auth", customer_ctlr.AuthenticateUser)
		r.Get("/auth", customer_ctlr.KeyedUserAuthentication)
		r.Post("/user/changePassword", customer_ctlr.ChangePassword)
		r.Post("/user", customer_ctlr.GetUser)
		r.Post("/user/register", customer_ctlr.RegisterUser)
		r.Post("/user/resetPassword", customer_ctlr.ResetPassword)
		r.Delete("/deleteKey", Deprecated)
		r.Post("/generateKey/user/:id/key/:type", Deprecated)
		r.Get("/user/:id", Deprecated)
		//r.Post("/user/:id", customer_ctlr.UpdateCustomerUser)
		//r.Delete("/user/:id", customer_ctlr.DeleteCustomerUser)
		// August 16th, 2017
		// If 6 months have passed with these being commented out, delete them and their functions
		r.Get("/users", Deprecated)

		r.Delete("/allUsersByCustomerID/:id", Deprecated) //Takes CustomerID (UUID)---danger!

		r.Put("/location/json", Deprecated)
		r.Put("/location/json/:id", Deprecated)
		r.Post("/location", Deprecated)
		r.Get("/location/:id", Deprecated)
		r.Put("/location/:id", Deprecated)
		r.Delete("/location/:id", Deprecated)

		r.Get("/locations", Deprecated)
		r.Post("/locations", Deprecated)

		r.Get("/price/:id", Deprecated)   //{part id}
		r.Get("/cartRef/:id", Deprecated) //{part id}

		// Customer CMS endpoints
		// All Customer Contents
		r.Get("/cms", Deprecated)
		// Content Types
		r.Get("/cms/content_types", Deprecated)

		// Customer Part Content
		r.Get("/cms/part", Deprecated)
		r.Get("/cms/part/:id", Deprecated)
		r.Put("/cms/part/:id", Deprecated) //partId
		r.Post("/cms/part/:id", Deprecated)
		r.Delete("/cms/part/:id", Deprecated)

		// Customer Category Content
		r.Get("/cms/category", Deprecated)
		r.Get("/cms/category/:id", Deprecated)
		r.Post("/cms/category/:id", Deprecated) //categoryId
		r.Delete("/cms/category/:id", Deprecated)

		// Customer Content By Content Id
		r.Get("/cms/:id", Deprecated)
		r.Get("/cms/:id/revisions", Deprecated)

		//Customer prices
		r.Get("/prices/part/:id", Deprecated)      //{id}; id refers to partId
		r.Post("/prices/sale", Deprecated)         //{start}{end}{id} -all required params; id refers to customerId
		r.Get("/prices/:id", Deprecated)           //{id}; id refers to {id} refers to customerPriceId
		r.Get("/prices", Deprecated)               //returns all {sort=field&direction=dir}
		r.Put("/prices/:id", Deprecated)           //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Post("/prices", Deprecated)              //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Delete("/prices/:id", Deprecated)        //{id} refers to customerPriceId
		r.Get("/pricesByCustomer/:id", Deprecated) //{id} refers to customerId; returns CustomerPrices

		r.Post("/:id", Deprecated)
		r.Delete("/:id", Deprecated)
		r.Put("", Deprecated)
	})

	m.Group("/dealers", func(r martini.Router) {
		r.Get("/business/classes", dealers_ctlr.GetAllBusinessClasses)
		r.Get("/etailer", dealers_ctlr.GetEtailers)
		r.Get("/local", dealers_ctlr.GetLocalDealers)
		r.Get("/local/regions", dealers_ctlr.GetLocalRegions)
		r.Get("/local/tiers", dealers_ctlr.GetLocalDealerTiers)
		r.Get("/local/types", dealers_ctlr.GetLocalDealerTypes)
		r.Get("/etailer/platinum", dealers_ctlr.PlatinumEtailers)
		r.Get("/location/:id", dealers_ctlr.GetLocationById)
		r.Get("/search/:search", dealers_ctlr.SearchLocations)
		r.Get("/search/type/:search", dealers_ctlr.SearchLocationsByType)
		r.Get("/search/geo/:latitude/:longitude", dealers_ctlr.SearchLocationsByLatLng)
	})

	//Creating, updating, and deleting FAQs are done in GoAdmin directly
	m.Group("/faqs", func(r martini.Router) {
		r.Get("", Deprecated)          //get all faqs; takes optional sort param {sort=true} to sort by question
		r.Get("/search", Deprecated)   //takes {question, answer, page, results} - all parameters are optional
		r.Get("/(:id)", Deprecated)    //get by id {id}
		r.Post("", Deprecated)         //takes {question, answer}; returns object with new ID
		r.Put("/(:id)", Deprecated)    //{id, question and/or answer}
		r.Delete("/(:id)", Deprecated) //{id}
		r.Delete("", Deprecated)       //{?id=id}
	})

	//All creating, updating, and deleting of things related to Forums
	//is done in GoAdmin directly
	m.Group("/forum", func(r martini.Router) {
		//groups
		r.Get("/groups", Deprecated)
		r.Get("/groups/:id", Deprecated)
		r.Post("/groups", Deprecated)
		r.Put("/groups/:id", Deprecated)
		r.Delete("/groups/:id", Deprecated)
		//topics
		r.Get("/topics", Deprecated)
		r.Get("/topics/:id", Deprecated)
		r.Post("/topics", Deprecated)
		r.Put("/topics/:id", Deprecated)
		r.Delete("/topics/:id", Deprecated)
		//threads
		r.Get("/threads", Deprecated)
		r.Get("/threads/:id", Deprecated)
		//posts
		r.Get("/posts", Deprecated)
		r.Get("/posts/:id", Deprecated)
		r.Post("/posts", Deprecated)
		r.Put("/posts/:id", Deprecated)
		r.Delete("/posts/:id", Deprecated)
	})

	m.Group("/geography", func(r martini.Router) {
		r.Get("/states", geography.GetAllStates)
		r.Get("/countries", geography.GetAllCountries)
		r.Get("/countrystates", geography.GetAllCountriesAndStates)
	})

	//Creating, updating, and deleting of News entites is done in GoAdmin directly
	m.Group("/news", func(r martini.Router) {
		r.Get("", news_controller.GetAll)           //get all news; takes optional sort param {sort=title||lead||content||startDate||endDate||active||slug} to sort by question
		r.Get("/titles", news_controller.GetTitles) //get titles!{page, results} - all parameters are optional
		r.Get("/leads", news_controller.GetLeads)   //get leads!{page, results} - all parameters are optional
		r.Get("/search", news_controller.Search)    //takes {title, lead, content, publishStart, publishEnd, active, slug, page, results, page, results} - all parameters are optional
		r.Get("/:id", news_controller.Get)          //get by id {id}
		r.Post("", Deprecated)                      //takes {question, answer}; returns object with new ID
		r.Post("/:id", Deprecated)                  //{id, question and/or answer}
		r.Delete("/:id", Deprecated)                //{id}
		r.Delete("", Deprecated)                    //{id}
	})

	m.Group("/part", func(r martini.Router) {
		r.Get("/featured", part_ctlr.Featured)
		r.Get("/latest", part_ctlr.Latest)
		r.Post("/multi", part_ctlr.GetMulti) //Actually a GET request, because of some "max length" myth
		r.Get("/:part/vehicles", part_ctlr.Vehicles)
		r.Get("/:part/attributes", part_ctlr.Attributes)
		r.Get("/:part/reviews", part_ctlr.ActiveApprovedReviews)
		r.Get("/:part/categories", part_ctlr.Categories)
		r.Get("/:part/content", part_ctlr.GetContent)
		r.Get("/:part/images", part_ctlr.Images)
		r.Get("/:part((.*?)\\.(PDF|pdf)$)", part_ctlr.InstallSheet)
		r.Get("/:part/packages", part_ctlr.Packaging)
		r.Get("/:part/pricing", part_ctlr.Prices)
		r.Get("/:part/related", part_ctlr.GetRelated)
		r.Get("/:part/videos", part_ctlr.Videos)
		r.Get("/:part/:year/:make/:model", Deprecated)
		r.Get("/:part/:year/:make/:model/:submodel", Deprecated)
		r.Get("/:part/:year/:make/:model/:submodel/:config(.+)", Deprecated)
		r.Get("/id/:part", part_ctlr.Get)
		r.Get("/identifiers", part_ctlr.Identifiers)
		r.Get("/:part", part_ctlr.PartNumber)
		r.Get("", part_ctlr.All)
	})

	//Creating, updating, and Deleting of salesRep entities is all done in GoAdmin directly
	m.Group("/salesrep", func(r martini.Router) {
		r.Get("", Deprecated)
		r.Post("", Deprecated)
		r.Get("/:id", Deprecated)
		r.Put("/:id", Deprecated)
		r.Delete("/:id", Deprecated)
	})

	m.Get("/search/:term", search_ctlr.Search)
	m.Get("/searchExactAndClose/:term", search_ctlr.SearchExactAndClose)

	//POST, PUT, and DELETE for these don't seem to be used, but even if they are,
	//they shouldn't, so they're getting locked down
	m.Group("/site", func(r martini.Router) {
		m.Group("/menu", func(r martini.Router) {
			r.Get("/all", Deprecated)
			r.Get("/:id", Deprecated)          //may pass id (int) or name(string)
			r.Get("/contents/:id", Deprecated) //may pass id (int) or name(string)
			r.Post("", Deprecated)
			r.Put("/:id", Deprecated)
			r.Delete("/:id", Deprecated)
		})
		m.Group("/content", func(r martini.Router) {
			r.Get("/all", Deprecated)
			r.Get("/:id", Deprecated) //may pass id (int) or slug(string)
			r.Get("/:id/revisions", Deprecated)
			r.Post("", Deprecated)
			r.Put("/:id", Deprecated)
			r.Delete("/:id", Deprecated)
		})
		r.Get("/details/:id", Deprecated)
		r.Post("", Deprecated)
		r.Put("/:id", Deprecated)
		r.Delete("/:id", Deprecated)
	})

	m.Group("/lp", func(r martini.Router) {
		r.Get("/:id", Deprecated)
	})

	//Creating of showcases is handled by GoAdmin directly
	m.Group("/showcase", func(r martini.Router) {
		r.Get("", Deprecated)
		r.Get("/:id", Deprecated)
		r.Post("", Deprecated)
	})

	//Completely unused
	m.Group("/techSupport", func(r martini.Router) {
		r.Get("/all", Deprecated)
		r.Get("/contact/:id", Deprecated)
		r.Get("/:id", Deprecated)
		r.Post("/:contactReceiverTypeID/:sendEmail", Deprecated) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", Deprecated)
	})

	//Creating, updating, and deleting of testimonials is done in GoAdmin directly
	m.Group("/testimonials", func(r martini.Router) {
		r.Get("", Deprecated)
		r.Get("/:id", Deprecated)
		r.Post("", Deprecated)
		r.Put("/:id", Deprecated)
		r.Delete("/:id", Deprecated)
	})

	//warranty related actions are handled in Survey
	m.Group("/warranty", func(r martini.Router) {
		r.Get("/all", Deprecated)
		r.Get("/contact/:id", Deprecated)
		r.Get("/:id", Deprecated)
		r.Post("/:contactReceiverTypeID/:sendEmail", Deprecated) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", Deprecated)
	})

	//This is unholy and should not exist
	m.Group("/webProperties", func(r martini.Router) {
		r.Post("/requirement/:id", Deprecated)
		r.Put("/requirement", Deprecated)
		r.Delete("/requirement/:id", Deprecated)
		r.Get("/requirement/:id", Deprecated)
		r.Get("/requirement", Deprecated)
		r.Post("/json/type", Deprecated)
		r.Post("/json/type/:id", Deprecated)
		r.Post("/json/requirement", Deprecated)
		r.Post("/json/requirement/:id", Deprecated)
		r.Post("/json/note", Deprecated)
		r.Post("/json/note/:id", Deprecated)
		r.Post("/json/:id", Deprecated)
		r.Put("/json", Deprecated)
		r.Post("/note/:id", Deprecated)   //updates when an id is present; otherwise, creates
		r.Put("/note", Deprecated)        //updates when an id is present; otherwise, creates
		r.Delete("/note/:id", Deprecated) //{id}
		r.Get("/note/:id", Deprecated)    //{id}
		r.Post("/type/:id", Deprecated)   //updates when an id is present; otherwise, creates
		r.Put("/type", Deprecated)        //updates when an id is present; otherwise, creates
		r.Delete("/type/:id", Deprecated) //{id}
		r.Get("/type/:id", Deprecated)    //{id}
		r.Get("/search", Deprecated)
		r.Get("/type", Deprecated)
		r.Get("/note", Deprecated)
		r.Get("/customer", Deprecated)
		r.Get("", Deprecated)
		r.Get("/:id", Deprecated)    //?id=id
		r.Delete("/:id", Deprecated) //{id}
		r.Post("/:id", Deprecated)   //
		r.Put("", Deprecated)        //can create notes(text) and requirements (requirement, by requirement=requirementID) while creating a property
	})

	// ARIES Year/Make/Model/Style
	m.Post("/vehicle", Deprecated)
	m.Post("/findVehicle", Deprecated)
	m.Post("/vehicle/inquire", Deprecated)

	// Used by ARIES ProductWidget
	m.Get("/vehicle/mongo/cols", Deprecated)

	// Used for ARIES Application Guides page
	m.Post("/vehicle/mongo/apps", vehicle.ByCategory)
	m.Post("/vehicle/mongo/allCollections", Deprecated)

	// Used by the ARIES website
	m.Get("/vehicle/category", vehicle.QueryCategoryStyle)
	m.Get("/vehicle/category/:year", vehicle.QueryCategoryStyle)
	m.Get("/vehicle/category/:year/:make", vehicle.QueryCategoryStyle)
	m.Get("/vehicle/category/:year/:make/:model", vehicle.QueryCategoryStyle)
	m.Get("/vehicle/category/:year/:make/:model/:category", vehicle.QueryCategoryStyle)

	// Used by the Luverne website
	m.Get("/luverne/vehicle", luverne.QueryCategoryStyle)
	m.Get("/luverne/vehicle/:year", luverne.QueryCategoryStyle)
	m.Get("/luverne/vehicle/:year/:make", luverne.QueryCategoryStyle)
	m.Get("/luverne/vehicle/:year/:make/:model", luverne.QueryCategoryStyle)
	m.Get("/luverne/vehicle/:year/:make/:model/:category", luverne.QueryCategoryStyle)

	// CURT Year/Make/Model/Style
	m.Post("/vehicle/curt", vehicle.CurtLookup)
	m.Get("/vehicle/curt", vehicle.CurtLookupGet)

	//videos are handled in GoAdmin
	m.Group("/videos", func(r martini.Router) {
		r.Get("/distinct", videos_ctlr.DistinctVideos) //old "videos" table - curtmfg?
		r.Get("/channel/type", videos_ctlr.GetAllChannelTypes)
		r.Get("/channel/type/:id", videos_ctlr.GetChannelType)
		r.Get("/channel", videos_ctlr.GetAllChannels)
		r.Get("/channels", videos_ctlr.GetAllChannels)
		r.Get("/channel/:id", videos_ctlr.GetChannel)
		r.Get("/cdn/type", videos_ctlr.GetAllCdnTypes)
		r.Get("/cdn/type/:id", videos_ctlr.GetCdnType)
		r.Get("/cdn", videos_ctlr.GetAllCdns)
		r.Get("/cdn/:id", videos_ctlr.GetCdn)
		r.Get("/type", videos_ctlr.GetAllVideoTypes)
		r.Get("/type/:id", videos_ctlr.GetVideoType)
		r.Get("", videos_ctlr.GetAllVideos)
		r.Get("/details/:id", videos_ctlr.GetVideoDetails)
		r.Get("/:id", videos_ctlr.Get)
	})

	m.Group("/vin", func(r martini.Router) {
		//option 1 - two calls - ultimately returns parts
		r.Get("/configs/:vin", vinLookup.GetConfigs)                    //returns vehicles - user must call vin/vehicle with vehicleID to get parts
		r.Get("/vehicleID/:vehicleID", vinLookup.GetPartsFromVehicleID) //returns an array of parts

		//option 2 - one call - returns vehicles with parts
		r.Get("/:vin", vinLookup.GetParts) //returns vehicles + configs with associates parts -or- an array of parts if only one vehicle config matches
	})

	m.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("running"))
	})

	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
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

func Deprecated(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGone)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This resource has been deprecated. If you use this resource, please contact blah blah blah"))
}
