#### Categories
===
List of endpoints

 - [Get All Categories](#all-categories)
 - [Get Single Category](#single-category)
 - [Get Sub-Categories](#sub-categories)
 - [Get Category Parts](#category-parts)
 
---

##<a name="all-categories"></a>Get All Categories `GET  - http://goapi.curtmfg.com/category`
All Categories.

*Example:*

	http://goapi.curtmfg.com/category?key=[public api key]&brandID=[Brand ID]

####Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brandID *(optional)* | The Brand ID of the categories your looking for |


####Response
Returns an unnamed array of Caregory objets. A Category object is described in the Get Single Category response.

| Property Name  |  Value |  Description |
|---|---|---|
|  | []object  | Array of Category Objects  |


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

	GET (paged) - http://goapi.curtmfg.com/category/<category id>/parts?page=[page]&count=[count]&key=[public api key]

	POST - http://API.curtmfg.com/category/<parent category id>/parts?key=[public api key]
