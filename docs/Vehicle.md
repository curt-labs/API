#### Vehicle

---

*Get Years*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:
		No additional form data required.


*Get Makes*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:

		"year" : <year (string)>


*Get Models*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:

		"year" : <Vehicle Year (string)>,
		"make"  : <Vehicle Make  (string)>

*Get SubModels*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:

		"year" : <Vehicle Year (string)>,
		"make"  : <Vehicle Make  (string)>,
		"model"  : <Vehicle Model (string)>


*Get Dynamic Configuration Option*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:

		"year" : <Vehicle Year (string)>,
		"make"  : <Vehicle Make  (string)>,
		"model"  : <Vehicle Model (string)>,
		"submodel" : <Vehicle Sub Model (string)>

*Get Next Dynamic Configuration Option*

	POST - http://API.curtmfg.com/vehicle?key=[public api key]

	Form Post Payload:

		"year" : <Vehicle Year (string)>,
		"make"  : <Vehicle Make  (string)>,
		"model"  : <Vehicle Model (string)>,
		"submodel" : <Vehicle Sub Model (string)>,
		"[config type]" : <Vehicle Config Option (string)>

	*Note: Each new selected config type is its own key/value.*