package main

import (
	"flag"
	"github.com/curt-labs/GoAPI/controllers/applicationGuide"
	"github.com/curt-labs/GoAPI/controllers/blog"
	"github.com/curt-labs/GoAPI/controllers/brand"
	"github.com/curt-labs/GoAPI/controllers/cart"
	"github.com/curt-labs/GoAPI/controllers/cartIntegration"
	"github.com/curt-labs/GoAPI/controllers/category"
	"github.com/curt-labs/GoAPI/controllers/contact"
	"github.com/curt-labs/GoAPI/controllers/customer"
	"github.com/curt-labs/GoAPI/controllers/dealers"
	"github.com/curt-labs/GoAPI/controllers/faq"
	"github.com/curt-labs/GoAPI/controllers/forum"
	"github.com/curt-labs/GoAPI/controllers/geography"
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/controllers/news"
	"github.com/curt-labs/GoAPI/controllers/part"
	"github.com/curt-labs/GoAPI/controllers/salesrep"
	"github.com/curt-labs/GoAPI/controllers/search"
	"github.com/curt-labs/GoAPI/controllers/site"
	"github.com/curt-labs/GoAPI/controllers/techSupport"
	"github.com/curt-labs/GoAPI/controllers/testimonials"
	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/controllers/videos"
	"github.com/curt-labs/GoAPI/controllers/vinLookup"
	"github.com/curt-labs/GoAPI/controllers/warranty"
	"github.com/curt-labs/GoAPI/controllers/webProperty"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/gorelic"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
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

	m := martini.Classic()
	gorelic.InitNewrelicAgent("5fbc49f51bd658d47b4d5517f7a9cb407099c08c", "GoAPI", false)
	m.Use(gorelic.Handler)
	m.Use(gzip.All())
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

	m.Group("/applicationGuide", func(r martini.Router) {
		r.Get("/website/:id", applicationGuide.GetApplicationGuidesByWebsite)
		r.Get("/:id", applicationGuide.GetApplicationGuide)
		r.Delete("/:id", applicationGuide.DeleteApplicationGuide)
		r.Post("", applicationGuide.CreateApplicationGuide)
	})

	m.Group("/blogs", func(r martini.Router) {
		r.Get("", blog_controller.GetAll)                      //sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/categories", blog_controller.GetAllCategories) //all categories; sort on any field e.g. ?sort=Name&direction=descending
		r.Get("/category/:id", blog_controller.GetBlogCategory)
		r.Get("/search", blog_controller.Search) //search field = value e.g. /blogs/search?key=8AEE0620-412E-47FC-900A-947820EA1C1D&slug=cyclo
		r.Post("/categories", blog_controller.CreateBlogCategory)
		r.Delete("/categories/:id", blog_controller.DeleteBlogCategory)
		r.Get("/:id", blog_controller.GetBlog)       //get blog by {id}
		r.Put("/:id", blog_controller.UpdateBlog)    //create {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} returns new id
		r.Post("", blog_controller.CreateBlog)       //update {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} required{id}
		r.Delete("/:id", blog_controller.DeleteBlog) //{?id=id}
		r.Delete("", blog_controller.DeleteBlog)     //{id}
	})

	m.Group("/brands", func(r martini.Router) {
		r.Get("", brand_ctlr.GetAllBrands)
		r.Post("", brand_ctlr.CreateBrand)
		r.Get("/:id", brand_ctlr.GetBrand)
		r.Put("/:id", brand_ctlr.UpdateBrand)
		r.Delete("/:id", brand_ctlr.DeleteBrand)
	})

	m.Group("/category", func(r martini.Router) {
		r.Get("", category_ctlr.Parents)
		r.Get("/:id", category_ctlr.GetCategory)
		r.Post("/:id", category_ctlr.GetCategory)
		r.Get("/:id/subs", category_ctlr.SubCategories)
		r.Get("/:id/parts", category_ctlr.GetParts)
		r.Post("/:id/parts", category_ctlr.GetParts)
		r.Get("/:id/parts/:page/:count", category_ctlr.GetParts)
	})

	m.Group("/contact", func(r martini.Router) {
		m.Group("/types", func(r martini.Router) {
			r.Get("/receivers/:id", contact.GetReceiversByContactType)
			r.Get("", contact.GetAllContactTypes)
			r.Get("/:id", contact.GetContactType)
			r.Post("", contact.AddContactType)
			r.Put("/:id", contact.UpdateContactType)
			r.Delete("/:id", contact.DeleteContactType)
		})
		m.Group("/receivers", func(r martini.Router) {
			r.Get("", contact.GetAllContactReceivers)
			r.Get("/:id", contact.GetContactReceiver)
			r.Post("", contact.AddContactReceiver)
			r.Put("/:id", contact.UpdateContactReceiver)
			r.Delete("/:id", contact.DeleteContactReceiver)
		})
		// r.Post("/sendmail/:id", contact.SendEmail)
		r.Get("", contact.GetAllContacts)
		r.Get("/:id", contact.GetContact)
		r.Post("/:contactTypeID", contact.AddDealerContact)
		r.Put("/:id", contact.UpdateContact)
		r.Delete("/:id", contact.DeleteContact)
	})

	m.Group("/shopify/customers", func(r martini.Router) {
		// Customers - shop endpoints
		r.Get("", cart_ctlr.GetCustomers)
		r.Post("", cart_ctlr.AddCustomer)
		r.Get("/search", cart_ctlr.SearchCustomer)
		r.Get("/:id", cart_ctlr.GetCustomer)
		r.Put("/:id", cart_ctlr.EditCustomer)
		r.Delete("/:id", cart_ctlr.DeleteCustomer)
		r.Get("/:id/orders", cart_ctlr.GetCustomerOrders)

		// Addresses
		r.Get("/:id/addresses", cart_ctlr.GetAddresses)
		r.Get("/:id/addresses/:address", cart_ctlr.GetAddress)
		r.Post("/:id/addresses", cart_ctlr.AddAddress)
		r.Put("/:id/addresses/:address/default", cart_ctlr.SetDefaultAddress)
		r.Put("/:id/addresses/:address", cart_ctlr.EditAddress)
		r.Delete("/:id/addresses/:address", cart_ctlr.DeleteAddress)

	})

	m.Group("/shopify/order", func(r martini.Router) {
		// Orders
		r.Post("/order", cart_ctlr.CreateOrder)
	})

	m.Group("/shopify/account", func(r martini.Router) {
		// Account - user endpoints
		r.Get("", cart_ctlr.GetAccount)
		r.Post("", cart_ctlr.AddAccount)
		r.Put("", cart_ctlr.EditAccount)
		r.Post("/login", cart_ctlr.AccountLogin)

		// m.Group("/shopify/account/address", func(r martini.Router) {
		// 	r.Get("", cart_ctlr.GetAccountAddresses)
		// 	r.Post("", cart_ctlr.AddAccountAddress)
		// 	r.Put("", cart_ctlr.EditAccountAddress)
		// 	r.Delete("", cart_ctlr.DeleteAccountAddress)
		// })

	})

	m.Group("/cart", func(r martini.Router) {
		r.Get("/customer/pricing/:custID/:page/:count", cartIntegration.GetCustomerPricingPaged)
		r.Get("/customer/pricing/:custID", cartIntegration.GetCustomerPricing)
		r.Get("/customer/count/:custID", cartIntegration.GetCustomerPricingCount)
		r.Get("/part/:id", cartIntegration.GetCIbyPart)
		r.Get("/customer/:id", cartIntegration.GetCIbyCustomer) //shallower object than GetCustomerPricing
		r.Get("/:id", cartIntegration.GetCI)
		r.Put("/:id", cartIntegration.SaveCI)
		r.Post("", cartIntegration.SaveCI)
		r.Delete("/:id", cartIntegration.DeleteCI)
	})

	m.Group("/customer", func(r martini.Router) {
		r.Get("", customer_ctlr.GetCustomer)
		r.Post("", customer_ctlr.GetCustomer)

		r.Post("/auth", customer_ctlr.AuthenticateUser)
		r.Get("/auth", customer_ctlr.KeyedUserAuthentication)

		r.Post("/user", customer_ctlr.GetUser)
		r.Post("/user/register", customer_ctlr.RegisterUser)
		r.Post("/user/resetPassword", customer_ctlr.ResetPassword)
		r.Post("/user/changePassword", customer_ctlr.ChangePassword)
		r.Post("/user/:id/key/:type", customer_ctlr.GenerateApiKey)
		r.Get("/user/:id", customer_ctlr.GetUserById)
		r.Post("/user/:id", customer_ctlr.UpdateCustomerUser)
		r.Delete("/user/:id", customer_ctlr.DeleteCustomerUser)
		r.Any("/users", customer_ctlr.GetUsers)

		r.Delete("/allUsersByCustomerID/:id", customer_ctlr.DeleteCustomerUsersByCustomerID) //Takes CustomerID (UUID)---danger!

		r.Put("/location/json", customer_ctlr.SaveLocationJson)
		r.Put("/location/json/:id", customer_ctlr.SaveLocationJson)
		r.Post("/location", customer_ctlr.SaveLocation)
		r.Get("/location/:id", customer_ctlr.GetLocation)
		r.Put("/location/:id", customer_ctlr.SaveLocation)
		r.Delete("/location/:id", customer_ctlr.DeleteLocation)

		r.Get("/locations", customer_ctlr.GetLocations)
		r.Post("/locations", customer_ctlr.GetLocations)

		r.Get("/price/:id", customer_ctlr.GetCustomerPrice)           //{part id}
		r.Get("/cartRef/:id", customer_ctlr.GetCustomerCartReference) //{part id}

		// Customer CMS endpoints
		// All Customer Contents
		r.Get("/cms", customer_ctlr.GetAllContent)
		// Content Types
		r.Get("/cms/content_types", customer_ctlr.GetAllContentTypes)

		// Customer Part Content
		r.Get("/cms/part", customer_ctlr.AllPartContent)
		r.Get("/cms/part/:id", customer_ctlr.UniquePartContent)
		r.Put("/cms/part/:id", customer_ctlr.UpdatePartContent) //partId
		r.Post("/cms/part/:id", customer_ctlr.CreatePartContent)
		r.Delete("/cms/part/:id", customer_ctlr.DeletePartContent)

		// Customer Category Content
		r.Get("/cms/category", customer_ctlr.AllCategoryContent)
		r.Get("/cms/category/:id", customer_ctlr.UniqueCategoryContent)
		r.Post("/cms/category/:id", customer_ctlr.UpdateCategoryContent) //categoryId
		r.Delete("/cms/category/:id", customer_ctlr.DeleteCategoryContent)

		// Customer Content By Content Id
		r.Get("/cms/:id", customer_ctlr.GetContentById)
		r.Get("/cms/:id/revisions", customer_ctlr.GetContentRevisionsById)

		//Customer prices
		r.Get("/prices/part/:id", customer_ctlr.GetPricesByPart)         //{id}; id refers to partId
		r.Post("/prices/sale", customer_ctlr.GetSales)                   //{start}{end}{id} -all required params; id refers to customerId
		r.Get("/prices/:id", customer_ctlr.GetPrice)                     //{id}; id refers to {id} refers to customerPriceId
		r.Get("/prices", customer_ctlr.GetAllPrices)                     //returns all {sort=field&direction=dir}
		r.Put("/prices/:id", customer_ctlr.CreateUpdatePrice)            //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Post("/prices", customer_ctlr.CreateUpdatePrice)               //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
		r.Delete("/prices/:id", customer_ctlr.DeletePrice)               //{id} refers to customerPriceId
		r.Get("/pricesByCustomer/:id", customer_ctlr.GetPriceByCustomer) //{id} refers to customerId; returns CustomerPrices

		r.Post("/:id", customer_ctlr.SaveCustomer)
		r.Delete("/:id", customer_ctlr.DeleteCustomer)
		r.Put("", customer_ctlr.SaveCustomer)
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

	m.Group("/faqs", func(r martini.Router) {
		r.Get("", faq_controller.GetAll)          //get all faqs; takes optional sort param {sort=true} to sort by question
		r.Get("/search", faq_controller.Search)   //takes {question, answer, page, results} - all parameters are optional
		r.Get("/(:id)", faq_controller.Get)       //get by id {id}
		r.Post("", faq_controller.Create)         //takes {question, answer}; returns object with new ID
		r.Put("/(:id)", faq_controller.Update)    //{id, question and/or answer}
		r.Delete("/(:id)", faq_controller.Delete) //{id}
		r.Delete("", faq_controller.Delete)       //{?id=id}
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

	m.Group("/geography", func(r martini.Router) {
		r.Get("/states", geography.GetAllStates)
		r.Get("/countries", geography.GetAllCountries)
		r.Get("/countrystates", geography.GetAllCountriesAndStates)
	})

	m.Group("/news", func(r martini.Router) {
		r.Get("", news_controller.GetAll)           //get all news; takes optional sort param {sort=title||lead||content||startDate||endDate||active||slug} to sort by question
		r.Get("/titles", news_controller.GetTitles) //get titles!{page, results} - all parameters are optional
		r.Get("/leads", news_controller.GetLeads)   //get leads!{page, results} - all parameters are optional
		r.Get("/search", news_controller.Search)    //takes {title, lead, content, publishStart, publishEnd, active, slug, page, results, page, results} - all parameters are optional
		r.Get("/:id", news_controller.Get)          //get by id {id}
		r.Post("", news_controller.Create)          //takes {question, answer}; returns object with new ID
		r.Post("/:id", news_controller.Update)      //{id, question and/or answer}
		r.Delete("/:id", news_controller.Delete)    //{id}
		r.Delete("", news_controller.Delete)        //{id}
	})

	m.Group("/part", func(r martini.Router) {
		r.Get("/featured", part_ctlr.Featured)
		r.Get("/latest", part_ctlr.Latest)
		r.Get("/old/:part", part_ctlr.OldPartNumber)
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
		r.Get("/:part/:year/:make/:model", part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel", part_ctlr.GetWithVehicle)
		r.Get("/:part/:year/:make/:model/:submodel/:config(.+)", part_ctlr.GetWithVehicle)
		r.Get("/:part", part_ctlr.Get)
		r.Get("", part_ctlr.All)
		r.Put("/:id", part_ctlr.UpdatePart)
		r.Post("", part_ctlr.CreatePart)
		r.Delete("/:id", part_ctlr.DeletePart)
	})

	m.Group("/price", func(r martini.Router) {
		r.Get("/:id", part_ctlr.GetPrice)
		r.Post("", part_ctlr.SavePrice)
		r.Put("/:id", part_ctlr.SavePrice)
		r.Delete("/:id", part_ctlr.DeletePrice)
	})

	m.Group("/reviews", func(r martini.Router) {
		r.Get("", part_ctlr.GetAllReviews)
		r.Get("/:id", part_ctlr.GetReview)
		r.Put("", part_ctlr.SaveReview)
		r.Post("/:id", part_ctlr.SaveReview)
		r.Delete("/:id", part_ctlr.DeleteReview)
	})

	m.Group("/salesrep", func(r martini.Router) {
		r.Get("", salesrep.GetAllSalesReps)
		r.Post("", salesrep.AddSalesRep)
		r.Get("/:id", salesrep.GetSalesRep)
		r.Put("/:id", salesrep.UpdateSalesRep)
		r.Delete("/:id", salesrep.DeleteSalesRep)
	})

	m.Get("/search/:term", search_ctlr.Search)

	m.Group("/site", func(r martini.Router) {
		m.Group("/menu", func(r martini.Router) {
			r.Get("/all", site.GetAllMenus)
			r.Get("/:id", site.GetMenu)                      //may pass id (int) or name(string)
			r.Get("/contents/:id", site.GetMenuWithContents) //may pass id (int) or name(string)
			r.Post("", site.SaveMenu)
			r.Put("/:id", site.SaveMenu)
			r.Delete("/:id", site.DeleteMenu)
		})
		m.Group("/content", func(r martini.Router) {
			r.Get("/all", site.GetAllContents)
			r.Get("/:id", site.GetContent) //may pass id (int) or slug(string)
			r.Get("/:id/revisions", site.GetContentRevisions)
			r.Post("", site.SaveContent)
			r.Put("/:id", site.SaveContent)
			r.Delete("/:id", site.DeleteContent)
		})
		r.Get("/details/:id", site.GetSiteDetails)
		r.Post("", site.SaveSite)
		r.Put("/:id", site.SaveSite)
		r.Delete("/:id", site.DeleteSite)
	})

	m.Group("/techSupport", func(r martini.Router) {
		r.Get("/all", techSupport.GetAllTechSupport)
		r.Get("/contact/:id", techSupport.GetTechSupportByContact)
		r.Get("/:id", techSupport.GetTechSupport)
		r.Post("/:contactReceiverTypeID/:sendEmail", techSupport.CreateTechSupport) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", techSupport.DeleteTechSupport)
	})

	m.Group("/testimonials", func(r martini.Router) {
		r.Get("", testimonials.GetAllTestimonials)
		r.Get("/:id", testimonials.GetTestimonial)
		r.Post("", testimonials.Save)
		r.Put("/:id", testimonials.Save)
		r.Delete("/:id", testimonials.Delete)
	})

	m.Group("/warranty", func(r martini.Router) {
		r.Get("/all", warranty.GetAllWarranties)
		r.Get("/contact/:id", warranty.GetWarrantyByContact)
		r.Get("/:id", warranty.GetWarranty)
		r.Post("/:contactReceiverTypeID/:sendEmail", warranty.CreateWarranty) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", warranty.DeleteWarranty)
	})

	m.Group("/webProperties", func(r martini.Router) {
		r.Post("/json/type", webProperty_controller.CreateUpdateWebPropertyType)
		r.Post("/json/type/:id", webProperty_controller.CreateUpdateWebPropertyType)
		r.Post("/json/requirement", webProperty_controller.CreateUpdateWebPropertyRequirement)
		r.Post("/json/requirement/:id", webProperty_controller.CreateUpdateWebPropertyRequirement)
		r.Post("/json/note", webProperty_controller.CreateUpdateWebPropertyNote)
		r.Post("/json/note/:id", webProperty_controller.CreateUpdateWebPropertyNote)
		r.Post("/json/:id", webProperty_controller.CreateUpdateWebProperty)
		r.Put("/json", webProperty_controller.CreateUpdateWebProperty)
		r.Post("/note/:id", webProperty_controller.CreateUpdateWebPropertyNote)               //updates when an id is present; otherwise, creates
		r.Put("/note", webProperty_controller.CreateUpdateWebPropertyNote)                    //updates when an id is present; otherwise, creates
		r.Delete("/note/:id", webProperty_controller.DeleteWebPropertyNote)                   //{id}
		r.Get("/note/:id", webProperty_controller.GetWebPropertyNote)                         //{id}
		r.Post("/type/:id", webProperty_controller.CreateUpdateWebPropertyType)               //updates when an id is present; otherwise, creates
		r.Put("/type", webProperty_controller.CreateUpdateWebPropertyType)                    //updates when an id is present; otherwise, creates
		r.Delete("/type/:id", webProperty_controller.DeleteWebPropertyType)                   //{id}
		r.Get("/type/:id", webProperty_controller.GetWebPropertyType)                         //{id}
		r.Post("/requirement/:id", webProperty_controller.CreateUpdateWebPropertyRequirement) //updates when an id is present; otherwise, creates
		r.Put("/requirement", webProperty_controller.CreateUpdateWebPropertyRequirement)      //updates when an id is present; otherwise, creates
		r.Delete("/requirement/:id", webProperty_controller.DeleteWebPropertyRequirement)     //{id}
		r.Get("/requirement/:id", webProperty_controller.GetWebPropertyRequirement)           //{id}
		r.Get("/search", webProperty_controller.Search)
		r.Get("/type", webProperty_controller.GetAllTypes)
		r.Get("/note", webProperty_controller.GetAllNotes)
		r.Get("/requirement", webProperty_controller.GetAllRequirements)
		r.Get("/customer", webProperty_controller.GetByPrivateKey)
		r.Get("", webProperty_controller.GetAll)
		r.Get("/:id", webProperty_controller.Get)                      //?id=id
		r.Delete("/:id", webProperty_controller.DeleteWebProperty)     //{id}
		r.Post("/:id", webProperty_controller.CreateUpdateWebProperty) //
		r.Put("", webProperty_controller.CreateUpdateWebProperty)      //can create notes(text) and requirements (requirement, by requirement=requirementID) while creating a property
	})

	m.Post("/vehicle", vehicle.Query)
	m.Post("/vehicle/inquire", vehicle.Inquire)

	m.Group("/videos", func(r martini.Router) {
		r.Get("/distinct", videos_ctlr.DistinctVideos) //old "videos" table - curtmfg?
		r.Get("/channel/type", videos_ctlr.GetAllChannelTypes)
		r.Get("/channel/type/:id", videos_ctlr.GetChannelType)
		r.Post("/channel/type/:id", videos_ctlr.SaveChannelType)
		r.Post("/channel/type", videos_ctlr.SaveChannelType)
		r.Delete("/channel/type/:id", videos_ctlr.DeleteChannelType)
		r.Get("/channel", videos_ctlr.GetAllChannels)
		r.Get("/channel/:id", videos_ctlr.GetChannel)
		r.Post("/channel/:id", videos_ctlr.SaveChannel)
		r.Post("/channel", videos_ctlr.SaveChannel)
		r.Delete("/channel/:id", videos_ctlr.DeleteChannel)
		r.Get("/cdn/type", videos_ctlr.GetAllCdnTypes)
		r.Get("/cdn/type/:id", videos_ctlr.GetCdnType)
		r.Post("/cdn/type/:id", videos_ctlr.SaveCdnType)
		r.Post("/cdn/type", videos_ctlr.SaveCdnType)
		r.Delete("/cdn/type/:id", videos_ctlr.DeleteCdnType)

		r.Get("/cdn", videos_ctlr.GetAllCdns)
		r.Get("/cdn/:id", videos_ctlr.GetCdn)
		r.Post("/cdn/:id", videos_ctlr.SaveCdn)
		r.Post("/cdn", videos_ctlr.SaveCdn)
		r.Delete("/cdn/:id", videos_ctlr.DeleteCdn)
		r.Get("/type", videos_ctlr.GetAllVideoTypes)
		r.Get("/type/:id", videos_ctlr.GetVideoType)
		r.Post("/type/:id", videos_ctlr.SaveVideoType)
		r.Post("/type", videos_ctlr.SaveVideoType)
		r.Delete("/type/:id", videos_ctlr.DeleteVideoType)

		r.Get("/part/:id", videos_ctlr.GetPartVideos)
		r.Get("", videos_ctlr.GetAllVideos)
		r.Get("/details/:id", videos_ctlr.GetVideoDetails)
		r.Get("/:id", videos_ctlr.Get)
		r.Post("/:id", videos_ctlr.SaveVideo)
		r.Post("", videos_ctlr.SaveVideo)
		r.Delete("/:id", videos_ctlr.DeleteVideo)
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
