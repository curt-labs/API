#### Brands

---
*Get All Brands*
	
	GET - http://goapi.curtmfg.com/brands?key=[public api key]

*Get Brand*

	GET - http://goapi.curtmfg.com/brands/<brand id>?key=[public api key]

*Create Brand (internal)*

	POST - http://goapi.curtmfg.com/brands?key=[public api key]

	Form Payload:

		"name" : <the brand name (ex. CURT) (string)>,
		"code" : <the brand code (ex. CURT) (string)>

*Update Brand (internal)*

	PUT - http://goapi.curtmfg.com/brands/<brand id>?key=[public api key]

	Form Payload:

		"name" : <the brand name (ex. CURT) (string)>,
		"code" : <the brand code (ex. CURT) (string)>

*Delete Brand (internal)*

	DELETE - http://goapi.curtmfg.com/brands/<brand id>?key=[public api key]

