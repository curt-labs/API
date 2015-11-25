#### News

---

*Get All News*

	GET - http://API.curtmfg.com/news?key=[public api key]

	Optional Query Parameters:

		"page"         : <pages the results. (ex. 1) (string)>
		"results"      : <results per page. (ex. 10) (string)>

*Get Titles*

	GET - http://API.curtmfg.com/news/titles?key=[public api key]

	Optional Query Parameters:

		"page"         : <pages the results. (ex. 1) (string)>
		"results"      : <results per page. (ex. 10) (string)>

*Get Leads*

	GET - http://API.curtmfg.com/news/leads?key=[public api key]

	Optional Query Parameters:

		"page"         : <pages the results. (ex. 1) (string)>
		"results"      : <results per page. (ex. 10) (string)>

*Get News*

	GET - http://API.curtmfg.com/news/<news id>?key=[public api key]

*Search News*

	GET - http://API.curtmfg.com/news/search?key=[public api key]

	Optional/Searchable Query Parameters:

		"title"        : <news title (string)>
		"lead"         : <news lead (string)>
		"content"      : <news content (string)>
		"publishStart" : <news publish start date (ex. "2013-10-02") (string)>
		"publishEnd"   : <news publish end date (ex. "2013-10-10") (string)>
		"active"       : <news is active? (ex. "true" or "false") (string)>
		"slug"         : <news slug (string)>
		"page"         : <pages the results. (ex. 1) (string)>
		"results"      : <results per page. (ex. 10) (string)>

*Create News (internal)*

	POST - http://API.curtmfg.com/news?key=[public api key]

	Form Payload:

		"title"   : <news title (string)>
		"lead"    : <news lead (string)>
		"content" : <news content (string)>
		"start"   : <news publish start date (ex. "2013-10-02") (string)>
		"end"     : <news publish end date (ex. "2013-10-02") (string)>
		"active"  : <news is active? (ex. "true" or "false") (string)>
		"slug"    : <news slug (string)>


*Update News (internal)*

	POST - http://API.curtmfg.com/news/<news id>?key=[public api key]

	Form Payload:

		"title"   : <news title (string)>
		"lead"    : <news lead (string)>
		"content" : <news content (string)>
		"start"   : <news publish start date (ex. "2013-10-02") (string)>
		"end"     : <news publish end date (ex. "2013-10-02") (string)>
		"active"  : <news is active? (ex. "true" or "false") (string)>
		"slug"    : <news slug (string)>


*Delete News (internal)*

	DELETE - http://API.curtmfg.com/news/<news id>?key=[public api key]

	- or -

	DELETE - http://API.curtmfg.com/news?id=<news id>&key=[public api key]

