#### Blogs

---
*Get All Blog Posts*

	GET - http://goapi.curtmfg.com/blogs?key=[public api key]

*Get All Blog Categories*
	
	GET - http://goapi.curtmfg.com/blogs/categories?key=[public api key]

*Get Blog Post*

	GET - http://goapi.curtmfg.com/blogs/<blog post id>?key=[public api key]

*Get Blog Category*
	
	GET - http://goapi.curtmfg.com/blogs/category/<blog category id>?key=[public api key]

*Create Blog Post*

	POST - http://goapi.curtmfg.com/blogs?key=[public api key]

	Form Post Payload:

		"title" : <post title (string)>,
		"slug"  : <post slug  (string)>,
		"text"  : <post text  (string)>,
		"publishedDate" : <post publish date (ex. "2014-12-01") (string)>,
		"userID" : <user id for poster of this blog post (ex. "1") (string)>,
		"metaTitle": <post meta title (string)>,
		"metaDescription": <post meta description (string)>,
		"keywords": <post keywords (string)>,
		"active": <post active? "true" or "false" (string)>,
		"categoryID": <post category ID (ex. "42") (string)>


*Update Blog Post*

	PUT - http://goapi.curtmfg.com/blogs/<blog post id>?key=[public api key]

	Form Post Payload:

		"title" : <title here (string)>,
		"slug"  : <slug here (string)>,
		"text"  : <blog post text (string)>,
		"publishedDate" : <post publish date (ex. "2014-12-01") (string)>,
		"userID" : <user id for poster of this blog post (ex. "1") (string)>,
		"metaTitle": <post meta title (string)>,
		"metaDescription": <post meta description (string)>,
		"keywords": <post keywords (string)>,
		"active": <post active? "true" or "false" (string)>,
		"categoryID": <post category ID (ex. "42") (string)>

*Create Blog Category*

	POST - http://goapi.curtmfg.com/blogs/categories?key=[public api key]

	Form Post Payload:

		"name": <blog category name>,
		"slug": <blog category slug>,
		"active": <is active? "true" || "false">

*Delete Blog Post*

	DELETE - http://goapi.curtmfg.com/blogs/<blog post id>?key=[public api key]

*Delete Blog Category*

	DELETE - http://goapi.curtmfg.com/blogs/categories/<category id>?key=[public api key]

*Search Blog Posts*

	GET - http://goapi.curtmfg.com/blogs/search?key=[public api key]

	Search by any combination of fields:

		"title" : <blog title (string)>,
		"slug"  : <blog slug (string)>,
		"text"  : <blog post text (string)>,
		"createdDate" : <blog post created date (string)>,
		"publishedDate": <blog post publish date (string)>,
		"lastModified" : <blog post last modified (string)>,
		"userID" : <userID for someone that posted or replied to a blog post (string)>,
		"metaTitle": <meta title for a blog posting (string)>,
		"metaDescription": <meta description for a blog posting (string)>,
		"keywords": <keywords for a blog post (string)>,
		"active" : <blog posts that are active "true" or inactive "false">,
		"page" : <paged results - "1","2","4",etc (string)>,
		"results": <limit results - "1","2","10",etc. (string)>