#### Categories

---
*Get Parent Categories*
	
	GET - http://goapi.curtmfg.com/category?key=[public api key]

*Get Category*

	GET - http://goapi.curtmfg.com/category/<category id>?key=[public api key]

	POST - http://goapi.curtmfg.com/category/<category id>?key=[public api key]

	Optional Query Parameters:

		page - results page
		count - count on each page

*Get Sub Categories*

	GET - http://goapi.curtmfg.com/category/<parent category id>/subs?key=[public api key]

*Get Category Parts*

	GET - http://goapi.curtmfg.com/category/<parent category id>/parts?key=[public api key]

	GET (paged) - http://goapi.curtmfg.com/category/<category id>/parts/<page>/<count>?key=[public api key]

	POST - http://goapi.curtmfg.com/category/<parent category id>/parts?key=[public api key]