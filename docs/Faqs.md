#### Faqs

---
*Get All Faqs*

	GET - http://API.curtmfg.com/faqs?sort=1&direction=asc&key=[public api key]

	Optional Query Parameters:

		"sort" : Value can be "true", "1", or any non-empty value indicating the results need to be sorted.
		"direction": Sort by 'ascending' or 'descending'

*Get Faq*

	GET - http://API.curtmfg.com/faqs/<faq id>?key=[public api key]

*Search Faq*

	GET - http://API.curtmfg.com/faqs/search?key=[public api key]

	Searchable Query Parameters:

		"question" : <faq question (string)>,
		"answer"   : <faq answer (string)>,
		"page"     : <results page (int)>,
		"results"  : <reuslts per page (int)>

*Create Faq (internal)*

	POST - http://API.curtmfg.com/faqs?key=[public api key]

	Form Payload:

		"question" : <faq question (string)>,
		"answer"   : <faq answer (string)>

*Update Faq (internal)*

	PUT - http://API.curtmfg.com/faqs/<faq id>?key=[public api key]

	Form Payload:

		"question" : <faq question (string)>,
		"answer"   : <faq answer (string)>

*Delete Faq (internal)*

	DELETE - http://API.curtmfg.com/faqs/<faq id>?key=[public api key]

	- or -

	DELETE - http://API.curtmfg.com/faqs?id=<faq id>&key=[public api key]
