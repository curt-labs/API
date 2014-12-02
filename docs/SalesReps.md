#### Sales Reps

---

*Get All Sales Reps*

	GET - http://goapi.curtmfg.com/salesrep?key=[public api key]

*Get Sales Rep*

	GET - http://goapi.curtmfg.com/salesrep/<sales rep id>?key=[public api key]

*Add Sales Rep (internal)*

	POST - http://goapi.curtmfg.com/salesrep?key=[public api key]

	Form Payload:

		"name" : <sales rep name (ex. "Wilfred Smith") (string)>,
		"code" : <sales rep code (ex. "1234") (string)>

*Update Sales Rep (internal)*

	PUT - http://goapi.curtmfg.com/salesrep/<sales rep id>?key=[public api key]

	Form Payload:

		"name" : <sales rep name (ex. "Wilfred Smith") (string)>,
		"code" : <sales rep code (ex. "1234") (string)>

*Delete Sales Rep (internal)*

	DELETE - http://goapi.curtmfg.com/salesrep/<sales rep id>?key=[public api key]

