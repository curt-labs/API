#### Brands

---
*Get All Brands*
	
	GET - http://API.curtmfg.com/brands?key=[public api key]

*Get Brand*

	GET - http://API.curtmfg.com/brands/<brand id>?key=[public api key]

*Create Brand (internal)*

	POST - http://API.curtmfg.com/brands?key=[public api key]

	Form Payload:

		"name" : <the brand name (ex. CURT) (string)>,
		"code" : <the brand code (ex. CURT) (string)>

*Update Brand (internal)*

	PUT - http://API.curtmfg.com/brands/<brand id>?key=[public api key]

	Form Payload:

		"name" : <the brand name (ex. CURT) (string)>,
		"code" : <the brand code (ex. CURT) (string)>

*Delete Brand (internal)*

	DELETE - http://API.curtmfg.com/brands/<brand id>?key=[public api key]

