# Categories
--

List of endpoints

 - [Get All Categories](#all-categories)
 - [Get Single Category](#single-category)
 - [Get Sub-Categories](#sub-categories)
 - [Get Category Parts](#category-parts)
 

## <a name="all-categories"></a>Get All Categories `GET  - http://goapi.curtmfg.com/category`
All Categories.

*Example:*

	http://goapi.curtmfg.com/category?key=[public api key]&brandID=[Brand ID]

#### Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brandID *(optional)* | The Brand ID of the categories your looking for |


#### Response
Returns an unnamed array of Caregory objets. A Category object is described in the Get Single Category response.

| Property Name  |  Value |  Description |
|---|---|---|
|  | []object  | Array of Category Objects  |


## <a name="single-category"></a>Get Single Category `GET  - http://goapi.curtmfg.com/category/:partId`
Information about the part.

*Example:*

	http://goapi.curtmfg.com/category/254?key=[API Key]&brandID=[Brand ID]


#### Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brand **(required)** | Brand querying that part number for (1=CURT, 3=ARIES, 4=Luverne) |

#### Response

| Property Name  |  Value |  Description |
|---|---|---|
| id   			| int  |  Unique Part Identifier |
| parent_identifier   	| string  |  |
| parent_id   		| int  |  The parent Category, if this is not a top-level category |
| children   		| []Category  | An array of child Category objects  |
| sort   		| int  |  This is the sort order when this is viewed as a child node |
| date_added   		| date  |  Date Category was created |
| title   		| string  |  Title |
| short_description   	| string  |  Short Description |
| long_description   	| string  |  Long Description |
| color_code   		| string  |  Gradient Background Color Code |
| font_code   		| string  |  Color of Font  |
| image   		| {}  |  Image |
| icon   		| {}  |  Icon |
| lifestyle   		| bool  |   |
| vehicle_specific   	| bool  |   |
| vehicle_required   	| bool  |   |
| meta_title   		| string  |  HTML Metadata |
| meta_description   	| string  |  HTML Metadata |
| meta_keywords   	| string  |  HTML Metadata |
| product_listing   	| {}  |   |
| content   		| []  |   |
| videos   		| []  |   |
| brand   		| Brand Object  |  Brand |
| part_identifiers   	| []  |   |
| pdf_path   		| string  |  Path of PDF |
| xls_path   		| string  |  Path of XLS |


*Get Sub Categories*

	GET - http://goapi.curtmfg.com/category/<parent category id>/subs?key=[public api key]

*Get Category Parts*

	GET - http://goapi.curtmfg.com/category/<parent category id>/parts?key=[public api key]

	GET (paged) - http://goapi.curtmfg.com/category/<category id>/parts?page=[page]&count=[count]&key=[public api key]

	POST - http://goapi.curtmfg.com/category/<parent category id>/parts?key=[public api key]
