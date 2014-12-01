#### Application Guides

---
*Get Application Guide by ID*

		GET - http://goapi.curtmfg.com/applicationGuide/<application guide id>?key=[public api key]

*Get Application Guides By Website*

		GET - http://goapi.curtmfg.com/applicationGuide/website/<website id>?key=[public api key]

*Create Application Guide*

		POST - http://goapi.curtmfg.com/applicationGuide?key=[public api key]

			JSON Payload:
			{
				"url" : <web url to app guide (string)>,
				"category": {
					"id": <category id (int)> 
				}
				"website" : {
					"id": <website id (int)>
				},
				"fileType": <file type (string)>
			}

			- or - 

			Form Post Payload:

				"url"         : <web url to app guide (string)>
				"category_id" : <category id (int)>
				"website_id"  : <website id (int)>
				"file_type"   :	<file type (string)>

*Delete Application Guide (internal)*

		DELETE - http://goapi.curtmfg.com/applicationGuide/<app guide id (int)>?key=[public api key]