#### Testimonials

---

*Get All Testimonials*

	GET - http://goapi.curtmfg.com/testimonials?key=[public api key]

	Optional Query Parameters:

		"page"  : <pages results (ex. 1) (int)
		"count" : <limits results to count (ex. 10) (int)
		"randomize": <randomizes the results (ex. "true") (string)

>>Note: If randomize is used, then count must also be specified.

*Get Testimonial*

	GET - http://goapi.curtmfg.com/testimonials/<testimonial id>?key=[public api key]

*Create Testimonial*

	POST - http://goapi.curtmfg.com/testimonials?key=[public api key]

	JSON Payload:
	{
		"rating"    : <testimonial rating (Ex. 5.0) (string)>,
		"title"     : <testimonial title (ex. "good product") (string)>,
		"content"   : <testimonial content (ex. "I bought a really good product.") (string)>,
		"dateAdded" : <testimonial dateAdded (ex. "2013-10-02") (string)>,
		"approved"  : <testimonial is approved? (ex. "true" or "false") (string)>,
		"active"    : <testimonial is active? (ex. "true" or "false") (string)>,
		"firstName" : <testimonial poster/business first name (ex. "Fred" or "Tim's Automotive") (string)>,
		"lastName"  : <testimonial poster's last name (ex. "Smith") (string)>,
		"location"  : <testimonial poster's location (ex. "Madison, WI") (string)>,
		"brandId"   : <testimonial brand id (ex. "1") (string)>
	}

*Update Testimonial*

	PUT - http://goapi.curtmfg.com/testimonials/<testimonial id>?key=[public api key]

	JSON Payload:
	{
		"rating"    : <testimonial rating (Ex. 5.0) (string)>,
		"title"     : <testimonial title (ex. "good product") (string)>,
		"content"   : <testimonial content (ex. "I bought a really good product.") (string)>,
		"dateAdded" : <testimonial dateAdded (ex. "2013-10-02") (string)>,
		"approved"  : <testimonial is approved? (ex. "true" or "false") (string)>,
		"active"    : <testimonial is active? (ex. "true" or "false") (string)>,
		"firstName" : <testimonial poster/business first name (ex. "Fred" or "Tim's Automotive") (string)>,
		"lastName"  : <testimonial poster's last name (ex. "Smith") (string)>,
		"location"  : <testimonial poster's location (ex. "Madison, WI") (string)>,
		"brandId"   : <testimonial brand id (ex. "1") (string)>
	}	

*Delete Testimonial*

	DELETE - http://goapi.curtmfg.com/testimonials/<testimonial id>?key=[public api key]

