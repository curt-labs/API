#### Sites

---

*Get All Menus*

	GET - http://API.curtmfg.com/site/menu/all?key=[public api key]

*Get All Contents*

	GET - http://API.curtmfg.com/site/content/all?key=[public api key]

*Get Menu*

	GET - http://API.curtmfg.com/site/menu/<menu id>?key=[public api key]

*Get Content*

	GET - http://API.curtmfg.com/site/content/<content id>?key=[public api key]

*Get Content Revisions*

	GET - http://API.curtmfg.com/site/content/<content id>/revisions?key=[public api key]

*Get Menu With Contents*

	GET - http://API.curtmfg.com/site/menu/contents/<content id or name>?key=[public api key]

*Get Website Details*

	GET - http://API.curtmfg.com/site/details/<website id>?key=[public api key]

*Create Website*

	PUT - http://API.curtmfg.com/site?key=[public api key]

	JSON Payload:

	{
		"url"         : <website url (string)>,
		"description" : <website description (string)>, 
	}

*Create Menu*

	PUT - http://API.curtmfg.com/site/menu?key=[public api key]

	JSON Payload:

	{
		"name"                  : <menu name (string)>,
		"isPrimary"             : <menu is primary? (ex. "true" or "false") (string)>,
		"active"                : <memu is active? (ex. "true" or "false") (string)>,
		"displayName"           : <menu display name (string)>,
		"requireAuthentication" : <menu requires authentication? (ex. "true" or "false") (string)>,
		"showOnSitemap"         : <menu is shown on sitemap? (ex. "true" or "false") (string)>,
		"sort"                  : <menu sort order (ex. "1") (string)>,
		"websiteId"             : <menu website id it is tied to (ex. "1") (string)>
	}

*Create Content*

	PUT - http://API.curtmfg.com/site/content?key=[public api key]

	JSON Payload:

	{
		"type"                  : <content type (string)>,
		"title"                 : <content title (string)>,
		"createdDate"           : <content created on date (ex. "2013-10-02") (string)>,
		"metaTitle"             : <content meta title (string)>,
		"metaDescription"       : <content meta description (string)>,
		"keywords"              : <content keywords (string)>,
		"isPrimary"             : <content is primary? (ex. "true" or "false") (string)>,
		"published"             : <content is published? (ex. "true" or "false") (string)>,
		"active"                : <content is active? (ex. "true" or "false") (string)>,
		"slug"                  : <content slug (string)>,
		"requireAuthentication" : <content requires authentication? (ex. "true" or "false") (string),
		"canonical"             : <content canonical (string)>,
		"websiteId"             : <content website id (ex. "1") (string)>
	}

*Update Website*

	POST - http://API.curtmfg.com/site/<website id>?key=[public api key]

	JSON Payload:

	{
		"url"         : <website url (string)>,
		"description" : <website description (string)>, 
	}

*Update Menu*

	POST - http://API.curtmfg.com/site/menu/<menu id>?key=[public api key]

	JSON Payload:

	{
		"name"                  : <menu name (string)>,
		"isPrimary"             : <menu is primary? (ex. "true" or "false") (string)>,
		"active"                : <memu is active? (ex. "true" or "false") (string)>,
		"displayName"           : <menu display name (string)>,
		"requireAuthentication" : <menu requires authentication? (ex. "true" or "false") (string)>,
		"showOnSitemap"         : <menu is shown on sitemap? (ex. "true" or "false") (string)>,
		"sort"                  : <menu sort order (ex. "1") (string)>,
		"websiteId"             : <menu website id it is tied to (ex. "1") (string)>
	}

*Update Content*

	POST - http://API.curtmfg.com/site/content/<content id>?key=[public api key]

	JSON Payload:

	{
		"type"                  : <content type (string)>,
		"title"                 : <content title (string)>,
		"createdDate"           : <content created on date (ex. "2013-10-02") (string)>,
		"metaTitle"             : <content meta title (string)>,
		"metaDescription"       : <content meta description (string)>,
		"keywords"              : <content keywords (string)>,
		"isPrimary"             : <content is primary? (ex. "true" or "false") (string)>,
		"published"             : <content is published? (ex. "true" or "false") (string)>,
		"active"                : <content is active? (ex. "true" or "false") (string)>,
		"slug"                  : <content slug (string)>,
		"requireAuthentication" : <content requires authentication? (ex. "true" or "false") (string),
		"canonical"             : <content canonical (string)>,
		"websiteId"             : <content website id (ex. "1") (string)>
	}

*Delete Website*

	DELETE - http://API.curtmfg.com/site/<website id>?key=[public api key]

*Delete Menu*

	DELETE - http://API.curtmfg.com/site/menu/<menu id>?key=[public api key]

*Delete Content*

	DELETE - http://API.curtmfg.com/site/content/<content id>?key=[public api key]
