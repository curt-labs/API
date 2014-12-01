#### Contacts

---
*Get All Contacts*

	GET - http://goapi.curtmfg.com/contact?key=[public api key]

	Optional Query Parameters:

		page - paged results (defaults to 1)
		count - count per page (defaults to 50)

*Get All Contact Types*

	GET - http://goapi.curtmfg.com/contact/types?key=[public api key]

*Get All Contact Receivers*

	GET - http://goapi.curtmfg.com/contact/receivers?key=[public api key]

*Get Contact*

	GET - http://goapi.curtmfg.com/contact/<contact id>?key=[public api key]

*Get Contact Type*

	GET - http://goapi.curtmfg.com/contact/types/<contact type id>?key=[public api key]

*Get Contact Receiver*

	GET - http://goapi.curtmfg.com/contact/receivers/<contact receiver id>?key=[public api key]

*Get Contact Receivers By Contact Type*

	GET - http://goapi.curtmfg.com/contact/types/receivers/<contact type id>?key=[public api key]

*Add Contact*

	POST - http://goapi.curtmfg.com/contact/<contact type id>?key=[public api key]

	JSON Payload:

		{
			"firstName" : <contact first name (string)>,
			"lastName"  : <contact last name (string)>,
			"email"     : <contact email address (string)>,
			"phone"     : <contact phone number (string)>,
			"subject"   : <contact subject line (string)>,
			"message"   : <contact message (string)>,
			"created"   : <contact created (string)>,
			"type"      : <contact type (ex. Become a Dealer, Tech Services, etc.) (string)>,
			"address1"  : <contact address 1 (string)>,
			"address2"  : <contact address 2 (string)>,
			"city"      : <contact city (string)>,
			"state"     : <contact state (ex. "Wisconsin") (string)>,
			"postalCode": <contact postal code (ex. "54701") (string)>,
			"country"   : <contact country (ex. "United States") (string)>
		}

*Add Contact Type*

	POST - http://goapi.curtmfg.com/contact/types?key=[public api key]

	Form Payload:

		"name" : <contact type name (Ex. Customer Service, Tech Services, etc) (string)>

*Add Contact Receiver*
	
	POST - http://goapi.curtmfg.com/contact/receivers?key=[public api key]

	Form Payload:

		"first_name" : <contact receiver first name (string)>
		"last_name"  : <contact receiver last name (string)>
		"email"      : <contact receiver email (string)>
		"contact_types" : <contact types for this receiver (Ex. "Customer Service,Tech Services, Webmaster")>

>> Note the comma separation for contact types in the example.

*Update Contact*

	PUT - http://goapi.curtmfg.com/contact/<contact id>?key=[public api key]

	JSON Payload:

		{
			"firstName" : <contact first name (string)>,
			"lastName"  : <contact last name (string)>,
			"email"     : <contact email address (string)>,
			"phone"     : <contact phone number (string)>,
			"subject"   : <contact subject line (string)>,
			"message"   : <contact message (string)>,
			"created"   : <contact created (string)>,
			"type"      : <contact type (ex. Become a Dealer, Tech Services, etc.) (string)>,
			"address1"  : <contact address 1 (string)>,
			"address2"  : <contact address 2 (string)>,
			"city"      : <contact city (string)>,
			"state"     : <contact state (ex. "Wisconsin") (string)>,
			"postalCode": <contact postal code (ex. "54701") (string)>,
			"country"   : <contact country (ex. "United States") (string)>
		}

	Form Payload:

		"first_name" : <contact first name (string)>,
		"last_name"  : <contact last name (string)>,
		"email"      : <contact email address (string)>,
		"phone"      : <contact phone number (string)>,
		"subject"    : <contact subject line (string),
		"message"    : <contact message (string),
		"type"       : <contact type (ex. Become a Dealer, Tech Services, etc.) (string)>,
		"address1"   : <contact address 1 (string)>,
	    "address2"   : <contact address 2 (string)>,
		"city"       : <contact city (string)>,
		"state"      : <contact state (ex. "Wisconsin") (string)>,
		"postal_code": <contact postal code (ex. "54701") (string)>,
		"country"    : <contact country (ex. "United States") (string)>

*Update Contact Type*

	PUT - http://goapi.curtmfg.com/contact/types/<contact type id>?key=[public api key]

	Form Payload:

		"name" : <contact type name (Ex. Customer Service, Tech Services, etc) (string)>
		"show" : <contact shown on website? (Ex. "true", "false") (string)>

*Update Contact Receiver*

	PUT - http://goapi.curtmfg.com/contact/receivers/<contact receiver id>?key=[public api key]

	Form Payload:

		"first_name" : <contact receiver first name (Ex. "Fred") (string)>
		"last_name"  : <contact receiver last name (Ex. "Smith") (string)>
		"email"      : <contact receiver email  (Ex. fsmith@web.com) (string)>
		"contact_types" : <contact types for this receiver (Ex. "Customer Service,Tech Services, Webmaster")>

*Delete Contact*

	DELETE - http://goapi.curtmfg.com/contact/<contact id>?key=[public api key]

*Delete Contact Type*

	DELETE - http://goapi.curtmfg.com/contact/types/<contact type id>?key=[public api key]

*Delete Contact Receiver*

	DELETE - http://goapi.curtmfg.com/contact/receivers/<contact receiver id>?key=[public api key]