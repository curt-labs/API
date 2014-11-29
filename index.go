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
	"github.com/curt-labs/GoAPI/controllers/customer_new"
	"github.com/curt-labs/GoAPI/controllers/dealers_new"
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
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))
	store := sessions.NewCookieStore([]byte("api_secret_session"))
	m.Use(sessions.Sessions("api_sessions", store))
	m.Use(encoding.MapEncoder)

	internalCors := cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*.curtmfg.com", "http://*.curtmfg.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	m.Post("/vehicle", vehicle.Query)

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
		r.Put("/:id", internalCors, part_ctlr.UpdatePart)
		r.Post("", internalCors, part_ctlr.CreatePart)
		r.Delete("/:id", internalCors, part_ctlr.DeletePart)
	})

	m.Group("/price", func(r martini.Router) {
		r.Get("/:id", internalCors, part_ctlr.GetPrice)
		r.Post("", internalCors, part_ctlr.SavePrice)
		r.Put("/:id", internalCors, part_ctlr.SavePrice)
		r.Delete("/:id", internalCors, part_ctlr.DeletePrice)
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

	m.Group("/applicationGuide", func(r martini.Router) {
		r.Get("/website/:id", applicationGuide.GetApplicationGuidesByWebsite)
		r.Get("/:id", applicationGuide.GetApplicationGuide)
		r.Delete("/:id", internalCors, applicationGuide.DeleteApplicationGuide)
		r.Post("", internalCors, applicationGuide.CreateApplicationGuide)
	})

	m.Group("/faqs", func(r martini.Router) {
		r.Get("", faq_controller.GetAll)                        //get all faqs; takes optional sort param {sort=true} to sort by question
		r.Get("/search", faq_controller.Search)                 //takes {question, answer, page, results} - all parameters are optional
		r.Get("/(:id)", faq_controller.Get)                     //get by id {id}
		r.Post("", internalCors, faq_controller.Create)         //takes {question, answer}; returns object with new ID
		r.Put("/(:id)", internalCors, faq_controller.Update)    //{id, question and/or answer}
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
		r.Put("/:id", internalCors, blog_controller.UpdateBlog)    //create {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} returns new id
		r.Post("", internalCors, blog_controller.CreateBlog)       //update {post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active} required{id}
		r.Delete("/:id", internalCors, blog_controller.DeleteBlog) //{?id=id}
		r.Delete("", internalCors, blog_controller.DeleteBlog)     //{id}
	})

	m.Group("/shopify", func(r martini.Router) {
		// Customers
		r.Get("/customers", cart_ctlr.GetCustomers)
		r.Get("/customers/search", cart_ctlr.SearchCustomer)
		r.Get("/customers/:id", cart_ctlr.GetCustomer)
		r.Get("/customers/:id/orders", cart_ctlr.GetCustomerOrders)
		r.Post("/customers", cart_ctlr.AddCustomer)
		r.Put("/customers/:id", cart_ctlr.EditCustomer)
		r.Delete("/customers/:id", cart_ctlr.DeleteCustomer)

		// Addresses
		r.Get("/customers/:id/addresses", cart_ctlr.GetAddresses)
		r.Get("/customers/:id/addresses/:address", cart_ctlr.GetAddress)
		r.Post("/customers/:id/addresses", cart_ctlr.AddAddress)
		r.Put("/customers/:id/addresses/:address/default", cart_ctlr.SetDefaultAddress)
		r.Put("/customers/:id/addresses/:address", cart_ctlr.EditAddress)
		r.Delete("/customers/:id/addresses/:address", cart_ctlr.DeleteAddress)
	})

	m.Group("/cart", func(r martini.Router) {
		r.Get("/customer/pricing/:custID/:page/:count", internalCors, cartIntegration.GetCustomerPricingPaged)
		r.Get("/customer/pricing/:custID", internalCors, cartIntegration.GetCustomerPricing)
		r.Get("/customer/count/:custID", internalCors, cartIntegration.GetCustomerPricingCount)
		r.Get("/part/:id", internalCors, cartIntegration.GetCIbyPart)
		r.Get("/customer/:id", internalCors, cartIntegration.GetCIbyCustomer) //shallower object than GetCustomerPricing
		r.Get("/:id", internalCors, cartIntegration.GetCI)
		r.Put("/:id", internalCors, cartIntegration.SaveCI)
		r.Post("", internalCors, cartIntegration.SaveCI)
		r.Delete("/:id", internalCors, cartIntegration.DeleteCI)
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
		r.Get("", news_controller.GetAll)                      //get all news; takes optional sort param {sort=title||lead||content||startDate||endDate||active||slug} to sort by question
		r.Get("/titles", news_controller.GetTitles)            //get titles!{page, results} - all parameters are optional
		r.Get("/leads", news_controller.GetLeads)              //get leads!{page, results} - all parameters are optional
		r.Get("/search", news_controller.Search)               //takes {title, lead, content, publishStart, publishEnd, active, slug, page, results, page, results} - all parameters are optional
		r.Get("/:id", news_controller.Get)                     //get by id {id}
		r.Post("", internalCors, news_controller.Create)       //takes {question, answer}; returns object with new ID
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

	m.Group("/webProperties", func(r martini.Router) {
		//Passing JSON in the request body?

		r.Post("/json/type", internalCors, webProperty_controller.SaveType_Json)
		r.Post("/json/type/:id", internalCors, webProperty_controller.SaveType_Json)
		r.Post("/json/requirement", internalCors, webProperty_controller.SaveRequirement_Json)
		r.Post("/json/requirement/:id", internalCors, webProperty_controller.SaveRequirement_Json)
		r.Post("/json/note", internalCors, webProperty_controller.SaveNote_Json)
		r.Post("/json/note/:id", internalCors, webProperty_controller.SaveNote_Json)
		r.Post("/json/:id", internalCors, webProperty_controller.Save_Json)
		r.Put("/json", internalCors, webProperty_controller.Save_Json)
		r.Post("/note/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyNote)               //updates when an id is present; otherwise, creates
		r.Put("/note", internalCors, webProperty_controller.CreateUpdateWebPropertyNote)                    //updates when an id is present; otherwise, creates
		r.Delete("/note/:id", internalCors, webProperty_controller.DeleteWebPropertyNote)                   //{id}
		r.Get("/note/:id", webProperty_controller.GetWebPropertyNote)                                       //{id}
		r.Post("/type/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyType)               //updates when an id is present; otherwise, creates
		r.Put("/type", internalCors, webProperty_controller.CreateUpdateWebPropertyType)                    //updates when an id is present; otherwise, creates
		r.Delete("/type/:id", internalCors, webProperty_controller.DeleteWebPropertyType)                   //{id}
		r.Get("/type/:id", webProperty_controller.GetWebPropertyType)                                       //{id}
		r.Post("/requirement/:id", internalCors, webProperty_controller.CreateUpdateWebPropertyRequirement) //updates when an id is present; otherwise, creates
		r.Put("/requirement", internalCors, webProperty_controller.CreateUpdateWebPropertyRequirement)      //updates when an id is present; otherwise, creates
		r.Delete("/requirement/:id", internalCors, webProperty_controller.DeleteWebPropertyRequirement)     //{id}
		r.Get("/requirement/:id", webProperty_controller.GetWebPropertyRequirement)                         //{id}
		r.Get("/search", internalCors, webProperty_controller.Search)
		r.Get("/type", webProperty_controller.GetAllTypes)               //all tyeps
		r.Get("/note", webProperty_controller.GetAllNotes)               //all notes
		r.Get("/requirement", webProperty_controller.GetAllRequirements) //requirements
		r.Get("/customer", webProperty_controller.GetByPrivateKey)       //private key
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

	m.Group("/site", func(r martini.Router) {
		m.Group("/menu", func(r martini.Router) {
			r.Get("/all", site.GetAllMenus)
			r.Get("/:id", site.GetMenu)                      //may pass id (int) or name(string)
			r.Get("/contents/:id", site.GetMenuWithContents) //may pass id (int) or name(string)
			r.Put("", site.SaveMenu)
			r.Post("/:id", site.SaveMenu)
			r.Delete("/:id", site.DeleteMenu)
		})
		m.Group("/content", func(r martini.Router) {
			r.Get("/all", site.GetAllContents)
			r.Get("/:id", site.GetContent) //may pass id (int) or slug(string)
			r.Get("/:id/revisions", site.GetContentRevisions)
			r.Put("", site.SaveContent)
			r.Post("/:id", site.SaveContent)
			r.Delete("/:id", site.DeleteContent)
		})
		r.Get("/details/:id", site.GetSiteDetails)
		r.Put("", site.SaveSite)
		r.Post("/:id", site.SaveSite)
		r.Delete("/:id", site.DeleteSite)
	})

	m.Group("/techSupport", func(r martini.Router) {
		r.Get("/all", techSupport.GetAllTechSupport)
		r.Get("/contact/:id", techSupport.GetTechSupportByContact)
		r.Get("/:id", techSupport.GetTechSupport)
		r.Post("/:contactReceiverTypeID/:sendEmail", techSupport.CreateTechSupport) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", techSupport.DeleteTechSupport)
	})

	m.Group("/warranty", func(r martini.Router) {
		r.Get("/all", warranty.GetAllWarranties)
		r.Get("/contact/:id", warranty.GetWarrantyByContact)
		r.Get("/:id", warranty.GetWarranty)
		r.Post("/:contactReceiverTypeID/:sendEmail", warranty.CreateWarranty) //contactType determines who receives the email/sendEmail is a bool indicating if email should be sent
		r.Delete("/:id", warranty.DeleteWarranty)
	})

	m.Group("/testimonials", func(r martini.Router) {
		r.Get("", testimonials.GetAllTestimonials)
		r.Get("/:id", testimonials.GetTestimonial)
		r.Post("/:id", testimonials.Save)
		r.Put("", testimonials.Save)
		r.Delete("/:id", testimonials.Delete)
	})

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

	//NEW Customer & Dealer endpoints - Seems to work. Feeling brave?
	m.Group("/new", func(r martini.Router) {
		m.Group("/customer", func(r martini.Router) {
			r.Get("", customer_ctlr_new.GetCustomer)
			r.Post("", customer_ctlr_new.GetCustomer)

			r.Post("/auth", customer_ctlr_new.AuthenticateUser)
			r.Get("/auth", customer_ctlr_new.KeyedUserAuthentication)

			r.Post("/user", customer_ctlr_new.GetUser)
			r.Post("/user/register", customer_ctlr_new.RegisterUser)
			r.Post("/user/resetPassword", customer_ctlr_new.ResetPassword)
			r.Post("/user/changePassword", customer_ctlr_new.ChangePassword)
			r.Post("/user/:id/key/:type", customer_ctlr_new.GenerateApiKey)
			r.Get("/user/:id", customer_ctlr_new.GetUserById)
			r.Post("/user/:id", customer_ctlr_new.UpdateCustomerUser)
			r.Delete("/user/:id", internalCors, customer_ctlr_new.DeleteCustomerUser)
			r.Any("/users", customer_ctlr_new.GetUsers)

			r.Delete("/allUsersByCustomerID/:id", internalCors, customer_ctlr_new.DeleteCustomerUsersByCustomerID) //Takes CustomerID (UUID)---danger!

			r.Put("/location/json", customer_ctlr_new.SaveLocationJson)
			r.Put("/location/json/:id", customer_ctlr_new.SaveLocationJson)
			r.Post("/location", customer_ctlr_new.SaveLocation)
			r.Get("/location/:id", customer_ctlr_new.GetLocation)
			r.Put("/location/:id", customer_ctlr_new.SaveLocation)
			r.Delete("/location/:id", customer_ctlr_new.DeleteLocation)

			r.Get("/locations", customer_ctlr_new.GetLocations)
			r.Post("/locations", customer_ctlr_new.GetLocations)

			r.Get("/price/:id", customer_ctlr_new.GetCustomerPrice)           //{part id}
			r.Get("/cartRef/:id", customer_ctlr_new.GetCustomerCartReference) //{part id}

			// Customer CMS endpoints
			// All Customer Contents
			r.Get("/cms", customer_ctlr_new.GetAllContent)
			// Content Types
			r.Get("/cms/content_types", customer_ctlr_new.GetAllContentTypes)

			// Customer Part Content
			r.Get("/cms/part", customer_ctlr_new.AllPartContent)
			r.Get("/cms/part/:id", customer_ctlr_new.UniquePartContent)
			r.Put("/cms/part/:id", customer_ctlr_new.UpdatePartContent) //partId
			r.Post("/cms/part/:id", customer_ctlr_new.CreatePartContent)
			r.Delete("/cms/part/:id", customer_ctlr_new.DeletePartContent)

			// Customer Category Content
			r.Get("/cms/category", customer_ctlr_new.AllCategoryContent)
			r.Get("/cms/category/:id", customer_ctlr_new.UniqueCategoryContent)
			r.Post("/cms/category/:id", customer_ctlr_new.UpdateCategoryContent) //categoryId
			r.Delete("/cms/category/:id", customer_ctlr_new.DeleteCategoryContent)

			// Customer Content By Content Id
			r.Get("/cms/:id", customer_ctlr_new.GetContentById)
			r.Get("/cms/:id/revisions", customer_ctlr_new.GetContentRevisionsById)

			//Customer prices
			r.Get("/prices/part/:id", customer_ctlr_new.GetPricesByPart)         //{id}; id refers to partId
			r.Post("/prices/sale", customer_ctlr_new.GetSales)                   //{start}{end}{id} -all required params; id refers to customerId
			r.Get("/prices/:id", customer_ctlr_new.GetPrice)                     //{id}; id refers to {id} refers to customerPriceId
			r.Get("/prices", customer_ctlr_new.GetAllPrices)                     //returns all {sort=field&direction=dir}
			r.Put("/prices/:id", customer_ctlr_new.CreateUpdatePrice)            //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
			r.Post("/prices", customer_ctlr_new.CreateUpdatePrice)               //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
			r.Delete("/prices/:id", customer_ctlr_new.DeletePrice)               //{id} refers to customerPriceId
			r.Get("/pricesByCustomer/:id", customer_ctlr_new.GetPriceByCustomer) //{id} refers to customerId; returns CustomerPrices

			r.Post("/:id", customer_ctlr_new.SaveCustomer)
			r.Delete("/:id", customer_ctlr_new.DeleteCustomer)
			r.Put("", customer_ctlr_new.SaveCustomer)

		})
		m.Group("/dealers", func(r martini.Router) {
			r.Get("/business/classes", dealers_ctlr_new.GetAllBusinessClasses)
			r.Get("/etailer", internalCors, dealers_ctlr_new.GetEtailers)
			r.Get("/local", internalCors, dealers_ctlr_new.GetLocalDealers)
			r.Get("/local/regions", internalCors, dealers_ctlr_new.GetLocalRegions)     //move to dealers
			r.Get("/local/tiers", internalCors, dealers_ctlr_new.GetLocalDealerTiers)   //move to dealers
			r.Get("/local/types", internalCors, dealers_ctlr_new.GetLocalDealerTypes)   //move to dealers
			r.Get("/etailer/platinum", internalCors, dealers_ctlr_new.PlatinumEtailers) //move to dealers
			r.Get("/location/:id", internalCors, dealers_ctlr_new.GetLocationById)      //move to dealers
			r.Get("/search/:search", internalCors, dealers_ctlr_new.SearchLocations)
			r.Get("/search/type/:search", internalCors, dealers_ctlr_new.SearchLocationsByType)
			r.Get("/search/geo/:latitude/:longitude", internalCors, dealers_ctlr_new.SearchLocationsByLatLng)
		})
		m.Get("/dealer/location/:id", internalCors, dealers_ctlr_new.GetLocationById)
	})

	//option 1 - two calls - ultimately returns parts
	m.Get("/vin/configs/:vin", vinLookup.GetConfigs)                    //returns vehicles - user must call vin/vehicle with vehicleID to get parts
	m.Get("/vin/vehicleID/:vehicleID", vinLookup.GetPartsFromVehicleID) //returns array of parts
	//option 2 - one call - returns vehicles with parts
	m.Get("/vin/:vin", vinLookup.GetParts) //returns vehicles+configs with associated parts -or- an array of parts if only one vehicle config matches

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
