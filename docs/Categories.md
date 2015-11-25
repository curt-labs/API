#### Categories

---
*Get Parent Categories*
	
	GET - http://API.curtmfg.com/category?key=[public api key]

*Get Category*

	GET - http://API.curtmfg.com/category/<category id>?key=[public api key]

	POST - http://API.curtmfg.com/category/<category id>?key=[public api key]

	Optional Query Parameters:

		page - results page
		count - count on each page

*Get Sub Categories*

	GET - http://API.curtmfg.com/category/<parent category id>/subs?key=[public api key]

*Get Category Parts*

	GET - http://API.curtmfg.com/category/<parent category id>/parts?key=[public api key]

	GET (paged) - http://API.curtmfg.com/category/<category id>/parts?page=[page]&count=[count]&key=[public api key]

	POST - http://API.curtmfg.com/category/<parent category id>/parts?key=[public api key]