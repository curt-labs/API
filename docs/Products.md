Parts
===
List of endpoints

 - [Get All Parts](#all-parts)
 - [Get Single Part](#single-part)
 - [Get Multiple Parts](#multi-parts)
 - [Get Last Added Parts](#last-added-parts)

##<a name="all-parts"></a>Get All Parts `GET  - http://goapi.curtmfg.com/part`
Information about the part.

*Example:*

	http://goapi.curtmfg.com/part?key=[public api key]&count=20&page=2&modified-from=2017-01-02


####Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| count *(optional)* | The number of parts you want returned |
| page *(optional)* | An offset multiplier based off the count |
| format *(optional)* | The format you wish the data to be in (only supports `json-obj`) |
| modified-from *(optional)* | Including this will only show products modified on or *after* this date |
| modified-to *(optional)* | Including this will only show products modified on or *before* this date |

####Response
Returns an unnamed array of part objets. A part object is described in the Get Single Part response.

| Property Name  |  Value |  Description |
|---|---|---|
|  | []object  | Array of part Objects  |

If the **format** is selected to be `json-obj`, the response will be formatted differently, like so:

| Property Name | Value | Description |
|---|---|---|
| items | []object | Array of part Objects |
| count | integer | The total number of results |


##<a name="single-part"></a>Get Single Part `GET  - http://goapi.curtmfg.com/part/:partId`
Information about the part.

*Example:*

	http://goapi.curtmfg.com/part/110003?key=[public api key]


####Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brand **(required)** | Brand querying that part number for (1=CURT, 3=ARIES, 4=Luverne) |



####Response

| Property Name  |  Value |  Description |
|---|---|---|
| id   				| int  |  Unique Part Identifier |
| part_number   		| string  |  Part ID (SKU) used for look-up |
| **[brand](https://github.com/curt-labs/API/tree/master/controllers/brand)**   	| object  |  Object describing brand e.g. CURT, ARIES, Luverne |
| **[status]()** 	| int  	|  Implementation Status Code describes "state" of product i.e. numerical representations for "Discontinued", "While Supplies Last" |
| price_code 			| int  |  ??? |
| related_count 		| int  |  Number of parts related to this part, size of array `related`  |
| average_review 	| float64  |  Average value of all reviews |
| date_modified 		| object  |  Date modified |
| date_added 			| object  |  Date created |
| short_description | string  |  Part title |
| **[install_sheet](https://golang.org/pkg/net/url/#URL)**   	| object  |  URL object with path to Install Sheet |
| **[attributes]()**   		| []object  | Unsorted list of key-value technical specifications |
| **[aces_vehicles](#aces-vehicle)**   		| []object  |  Array of vehicles that fit the part in ACES fitment format |
| vehicle_atttributes 	| []string  |  ??? |
| **[vehicle_applications]()** *(optional)* | []object  |  Array of vehicles that fit the part in ARIES fitment format |
| **[luverne_applications]()** *(optional)* | []object  |  Array of vehicles that fit the part in Luverne fitment format |
| **[content](#content)**  		| []object 	|  Array of product descriptions objects |
| **[pricing](#price)**   	| []object 	|  Array of pricing levels |
| **[reviews](#reviews)**   	| []object 	|  Array of product reviews |
| **[images](#image)**   		| []object 	|  An arry of object with the Image URL and meta-data |
| related   	| []int 		|  Array of Part Numbers of related projects |
| **[categories]()**   	| []object 	|  ??? |
| **[videos](#video)**   		| []object 	|  An array of Video objects. They have the path to the video and a lot of meta-data |
| **[packages]()**   	| []object 	|  ??? |
| **[customer]()** *(optional)*	| object |  ??? |
| **[class]()** *(optional)*		| object |  ??? |
| featured *(optional)*		| bool  		|  ??? |
| acesPartTypeId *(optional)* | int  	|  ??? |
| **[inventory]()** *(optional)*	| object  |  ??? |
| upc *(optional)* 			| string  	|  Universal Product Code |
| iconLayer   				| string  	|  ??? |
| mappedToVehicle  			| bool  		|  ??? |



##<a name="multi-parts"></a>Get Multiple Parts `POST  - http://goapi.curtmfg.com/part/multi`
Use a POST to query for multiple parts

*Example:*

```
curl --request POST \
  --url 'https://goapi.curtmfg.com/part/multi?brandID=3&key=[public api key]' \
  --header 'content-type: application/json' \
  --data '["1042","35-3014"]'
```

####Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brandID **(required)** | Brand querying part numbers for (1=CURT, 3=ARIES, 4=Luverne) |

####Request Body
| Paramter  | value | Description |
|---|---|---|
| [] | []string |An array of strings with product numbers (SKUs)  |

####Response
Returns an unnamed array of part objets. A part object is described in the Get Single Part response.

| Property Name  |  Value |  Description |
|---|---|---|
| [] | []object  | Array of part Objects  |


##<a name="last-added-parts"></a>Get Last Added Parts `GET  - http://goapi.curtmfg.com/part/latest`
Get the last added parts

*Example:*

	http://goapi.curtmfg.com/part/latest?key=[public api key]&count=20&brand=1


####Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brand *(optional)* | Brand querying part numbers for (1=CURT, 3=ARIES, 4=Luverne) |
| count *(optional)* | The number of parts you want returned |

####Response
Returns an unnamed array of part objets. A part object is described in the Get Single Part response.

| Property Name  |  Value |  Description |
|---|---|---|
| [] | []object  | Array of part Objects  |


##Product Objects
A list of Product Object definitions

#### <a name="aces-vehicle"></a> aces_vehicle ####

| Property Name  | Value | Description |
|---|---|---|
| [base]()  | object  | Base model of a Vehicle |
| submodel  | string  | Vehicle submodel |
| [configurations]()  | []object  | ACES Configurations |

#### <a name="content"></a> content ####

| Property Name  | Value | Description |
|---|---|---|
| text  | string  | Content string |
| [contentType](#contentType)  | object  | ContentType object that describes the content |
| sort  | int  | The order to sort this content |

#### <a name="contentType"></a> contentType ####

| Property Name  | Value | Description |
|---|---|---|
| id  | string  | Internal Id |
| type  | string  | Describes what the conten is i.e. "Description", "Bullet, "Content Brief" |
| allows_html  | bool  | Does the content include html |

#### <a name="image"></a> image ####

| Property Name  | Value | Description |
|---|---|---|
| id *(optional)* 	| int  | Internal Id |
| size *(optional)*	| string | Describes what the conten is i.e. "Description", "Bullet, "Content Brief" |
| sort *(optional)*		| string  | Sort Order |
| height *(optional)* 	| int  | Height in pixels |
| width *(optional)* 	| int  | Width in pixels |
| path *(optional)*		| object  | URL object with path to image |
| partId *(optional)* 	| int  | Part Id reference |

#### <a name="price"></a> price ####

| Property Name  | Value | Description |
|---|---|---|
| id *(optional)* | int  | Internal Id |
| partId *(optional)* | int  | Part Id reference |
| type *(optional)* | string  | Pricing type i.e. "Jobber", "List" |
| price | float64  | Price of product |
| enforced *(optional)* | bool  | ??? |
| DateModified *(optional)* | object  | Date Modified |

#### <a name="reviews"></a> reviews ####

| Property Name  | Value | Description |
|---|---|---|
| rating  | int  | Star Rating Value (1-5) |
| subject  | string  | Subject of review |
| review_text  | string  | Body of review |
| name  | string  | Reviewer's Name |
| email  | string  | Email address of reviewer |
| created_date  | object  | Time/Date created |
| active *(optional)*  | bool  | Review visible on product page |
| approved *(optional)*  | bool  | Review approved by moderator |

#### <a name="video"></a> video ####
A reresentation of a video. It contains information about the video itself, as well as any associated files.

| Property Name  | Value | Description |
|---|---|---|
| id *(optional)* 		| int  | Internal Id |
| title *(optional)* 	| string  | Video Title |
| subject_type  			| string  |  |
| videoType *(optional)* | object  | Reviewer's Name |
| description  		| string  | Email address of reviewer |
| date_added  		| object  | Date created |
| date_modified   	| object  | Date modified |
| thumbnail   		| string  | Thumbnail image URL |
| channel   			| []object  | Information related to where the file is uploaded to |
| cdn_file  			| []object  | Array of objects with CDN locations of video |
| isPrimary *(optional)* 	| bool  | Main video |
| categoryIds *(optional)* | []int  | Related Categories |
| partIds *(optional)* 		| []int  | Related part numbers |
| websiteId *(optional)* 	| int  | ??? |
| brands *(optional)* 		| object  | Brand this video was created for |
