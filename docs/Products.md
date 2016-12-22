Parts
===
List of endpoints
 
 - [Get All Parts](#all-parts)
 
 - [Get Single Part](#single-part)


##<a name="all-parts"></a>Get All Parts `GET  - http://goapi.curtmfg.com/part`
Information about the part.

*Example:*

	http://goapi.curtmfg.com/part?key=[public api key]&count=20&page=2
	
	
####Parameters


| Paramter  |  Description | 
|---|---|
| key **(required)** | Provide your API key  |
| count *(optional)* | The number of parts you want returned |
| page *(optional)* | An offset multiplier based off the count |

####Response 
Returns an unnamed array of part objets. A part object is described in the Get Single Part response.

| Property Name  |  Value |  Description |
|---|---|---|
|  | []object  | Array of part Objects  |


##<a name="single-part"></a>Get Single Part `GET  - http://goapi.curtmfg.com/part/:partId`
Information about the part.

*Example:*

	http://goapi.curtmfg.com/part/110003?key=[public api key]
	
	
####Parameters


| Paramter  |  Description | 
|---|---|
| key **(required)** | Provide your API key  |
| brand **(required)** | Brand your querying that part number for (1=CURT, 3=ARIES, 4=Luverne)  |



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


####Response Schema
```
{  
   "id":2060021,
   "part_number":"1042",
   "brand":{  
      "id":3,
      "name":"ARIES",
      "code":"ARIES",
      "logo":{  
         "Scheme":"https",
         "Opaque":"",
         "User":null,
         "Host":"storage.googleapis.com",
         "Path":"/aries-logo/SVG_Logo (2c_white with black outline on transparent).svg",
         "RawPath":"/aries-logo/SVG_Logo%20(2c_white%20with%20black%20outline%20on%20transparent).svg",
         "RawQuery":"",
         "Fragment":""
      },
      "logo_alternate":{  
         "Scheme":"https",
         "Opaque":"",
         "User":null,
         "Host":"storage.googleapis.com",
         "Path":"/aries-logo/ARIES Logo (1c_red on transparent).png",
         "RawPath":"/aries-logo/ARIES%20Logo%20(1c_red%20on%20transparent).png",
         "RawQuery":"",
         "Fragment":""
      },
      "formal_name":"Aries Automotive",
      "long_name":"Aries Automotive",
      "primary_color":"#57111A",
      "autocareId":"BBRD",
      "websites":null
   },
   "status":800,
   "price_code":0,
   "related_count":0,
   "average_review":0,
   "date_modified":"2016-10-04T20:00:22Z",
   "date_added":"2015-10-30T15:27:24.77Z",
   "short_description":"Grille Guard",
   "install_sheet":{  
      "Scheme":"https",
      "Opaque":"",
      "User":null,
      "Host":"www.curtmfg.com",
      "Path":"/masterlibrary/01ARIES/1042/installsheet/1042_INS.pdf",
      "RawPath":"",
      "RawQuery":"",
      "Fragment":""
   },
   "attributes":[  
      {  
         "name":"Material",
         "value":"Carbon steel"
      },
      {  
         "name":"Finish",
         "value":"Semi-Gloss Black"
      },
      {  
         "name":"Shipping Weight",
         "value":"49.000"
      },
      {  
         "name":"Warranty",
         "value":"Three Years Limited"
      },
      {  
         "name":"Application",
         "value":"Vehicle-specific"
      }
   ],
   "aces_vehicles":[  
      {  
         "base":{  
            "year":1997,
            "make":"Jeep",
            "model":"Grand Cherokee"
         },
         "submodel":"",
         "configurations":[  

         ]
      },
   ],
   "vehicle_atttributes":null,
   "vehicle_applications":[  
      {  
         "year":"1998",
         "make":"jeep",
         "model":"grand cherokee",
         "style":"tow hook only (4wd only)",
         "exposed":"",
         "drilling":"",
         "install_time":""
      },
   ],
   "content":[  
      {  
         "text":"One-piece, 1 1/2\" diameter, heavy-wall tube design",
         "contentType":{  
            "id":0,
            "type":"Bullet",
            "allows_html":false
         },
         "sort":34
      },
      {  
         "text":"High-strength carbon steel construction",
         "contentType":{  
            "id":0,
            "type":"Bullet",
            "allows_html":false
         },
         "sort":32
      },
      {  
         "text":"NOTE: Grille guard may interfere with forward-facing cameras or sensors",
         "contentType":{  
            "id":0,
            "type":"Bullet",
            "allows_html":false
         },
         "sort":18
      },
      {  
         "text":"\u003cp\u003e \u003cstrong\u003eCustomizable design\u003c/strong\u003e\u003c/p\u003e \u003cp\u003eThe crossbar has two pre-drilled holes to accept auxiliary lights, and the headlight cages can be removed for a custom look\u003c/p\u003e \u003cp\u003e \u003cstrong\u003eVehicle-specific fit\u003c/strong\u003e\u003c/p\u003e \u003cp\u003eThe 1/4\" thick risers and 1 1/2\" mandrel-bent steel tubing contour to the profile of each specific vehicle for seamless integration\u003c/p\u003e \u003cp\u003e \u003cstrong\u003eSimple, secure mounting\u003c/strong\u003e\u003c/p\u003e \u003cp\u003eThe risers bolt onto pre-existing factory holes in the vehicle''s frame to eliminate the need for drilling and ensure a solid mount\u003c/p\u003e",
         "contentType":{  
            "id":0,
            "type":"CategoryBrief",
            "allows_html":true
         },
         "sort":0
      },
   ],
   "pricing":[  
      {  
         "type":"Jobber",
         "price":385,
         "dateModified":"0001-01-01T00:00:00Z"
      },
      {  
         "type":"List",
         "price":481.25,
         "dateModified":"0001-01-01T00:00:00Z"
      }
   ],
   "reviews":[  

   ],
   "images":[  
      {  
         "size":"Tall",
         "sort":"a",
         "height":75,
         "width":100,
         "path":{  
            "Scheme":"https",
            "Opaque":"",
            "User":null,
            "Host":"www.curtmfg.com",
            "Path":"/masterlibrary/01ARIES/1042/images/1042_100x75_a.jpg",
            "RawPath":"",
            "RawQuery":"",
            "Fragment":""
         }
      },
      {  
         "size":"Medio",
         "sort":"a",
         "height":238,
         "width":200,
         "path":{  
            "Scheme":"https",
            "Opaque":"",
            "User":null,
            "Host":"www.curtmfg.com",
            "Path":"/masterlibrary/01ARIES/1042/images/1042_200x238_a.jpg",
            "RawPath":"",
            "RawQuery":"",
            "Fragment":""
         }
      },
   ],
   "related":[  

   ],
   "categories":[  
      {  
         "id":330,
         "parent_identifier":"",
         "parent_id":351,
         "children":[  

         ],
         "sort":1,
         "date_added":"2014-12-23T15:48:39Z",
         "title":"Grille Guards",
         "short_description":"Grille Guards",
         "long_description":"Grille Guards",
         "color_code":"",
         "font_code":"",
         "image":{  
            "Scheme":"https",
            "Opaque":"",
            "User":null,
            "Host":"storage.googleapis.com",
            "Path":"/aries-category-images-png/grille-guards-Category.png",
            "RawPath":"",
            "RawQuery":"",
            "Fragment":""
         },
         "icon":{  
            "Scheme":"",
            "Opaque":"",
            "User":null,
            "Host":"",
            "Path":"",
            "RawPath":"",
            "RawQuery":"",
            "Fragment":""
         },
         "lifestyle":false,
         "vehicle_specific":true,
         "vehicle_required":false,
         "meta_title":"Grille Guards | Brush Guards | ARIES\r\n",
         "meta_description":"Grille Guards | Brush Guards | ARIES\r\n",
         "meta_keywords":"aries, grille guard, brush guard, truck, suv, automotive\r\n",
         "content":[  
            {  
               "text":"\u003cdiv class=\"col-xs-12 col-sm-6 col-md-6 col-lg-6\" style=\"text-align:right; float:right;margin-bottom:10px; padding-left:25px;\"\u003e\r\n\t\u003cdiv class=\"ytVideoWrapper\"\u003e\r\n\t\t\u003ciframe align=\"right\" allowfullscreen=\"\" frameborder=\"0\" height=\"315\" src=\"//www.youtube.com/embed/7tTMyHOjR68?rel=0\" width=\"420\"\u003e\u003c/iframe\u003e\u003c/div\u003e\r\n\u003c/div\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u003cspan style=\"font-size:20px;\"\u003e\u003cstrong\u003eWith a vehicle-specific fit for hundreds of applications\u003c/strong\u003e\u003c/span\u003e\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\tFor both trucks and SUVs, the ARIES grille guard offers added protection, a vehicle-specific fit and\u0026nbsp;easy customization. It is built with a one-piece, heavy-duty 1 1/2\u0026quot; diameter steel tube construction to provide a safer, tougher-looking ride. Each grille guard is also made for a specific vehicle make and model, allowing it to contour to the front end and provide a seamless aftermarket look. ARIES grille guards feature pre-drilled auxiliary light holes to accept fog lights, off-road lights or other aftermarket lighting options, and they come with removable headlight cages for further customization. The ARIES grille guard is available in polished stainless steel or carbon steel with a semi-gloss black powder coat finish and comes with a warranty to ensure quality in materials and craftsmanship.\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u0026nbsp;\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u003cspan style=\"font-size:20px;\"\u003e\u003cstrong\u003eVehicle-specific fit\u003c/strong\u003e\u003c/span\u003e\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\tBecause they are made vehicle-specific, ARIES grille guards fit better, install easier and attach more securely. The solid, one-piece frame mounts onto the front of your vehicle using a simple, four-point mounting system and no-drill application. This means pre-existing factory holes are used for installation, ensuring a secure fit with less vibration. The 1 1/2\u0026rdquo; diameter tube brush guards are custom-bent to contour to the front end and accent the profile of the vehicle and its unique features.\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u0026nbsp;\u003c/p\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u003cstrong style=\"font-size: 20px;\"\u003e\u003cimg alt=\"\" src=\"https://storage.googleapis.com/aries-dealer_portal/Category_Page-Grille_Guards_Image1.jpg\" style=\"width: 500px; height: 331px; float: left; margin-left: 10px; margin-right: 10px;\" /\u003e\u003c/strong\u003e\u003c/p\u003e\r\n\u003cdiv\u003e\r\n\t\u0026nbsp;\u003c/div\u003e\r\n\u003cp style=\"padding-right:1em;\"\u003e\r\n\t\u003cspan style=\"font-size:20px;\"\u003e\u003cstrong\u003eTwo finish options\u003c/strong\u003e\u003c/span\u003e\u003c/p\u003e\r\n\u003cp\u003e\r\n\tThe ARIES grille guard is offered in two finish options: stainless steel and black powder coat. The stainless grille guard is made with 304 stainless steel, making it high in nickel content and truly resistant to rust and corrosion. It also comes with a polished finish for a mirror-like shine. Our black powder coat grille guard is made with high-strength carbon steel to add extra protection for the front end of your truck, and it features a semi-gloss black powder coat finish that easily resists rust and scratches. For both options, each grille guard comes with rubber stripping along the risers to help protect the finish.\u003c/p\u003e\r\n\u003cp\u003e\r\n\t\u0026nbsp;\u003c/p\u003e\r\n\u003cp\u003e\r\n\t\u0026nbsp;\u003c/p\u003e\r\n\u003cp\u003e\r\n\t\u0026nbsp;\u003c/p\u003e\r\n\u003cp\u003e\r\n\t\u003cspan style=\"font-size:20px;\"\u003e\u003cstrong\u003eBacked by warranty\u003c/strong\u003e\u003c/span\u003e\u003c/p\u003e\r\n\u003cp\u003e\r\n\tThe ARIES grille guard is available for most makes and models of pickup trucks and SUVs. Whatever your truck of choice -- RAM 1500, Chevy Silverado, GMC Sierra, Ford F150 or none of the above -- we want to back you up. We offer a three-year warranty on our black powder coat finish and a limited lifetime warranty on our polished stainless steel.\u003c/p\u003e\r\n",
               "contentType":{  
                  "id":0,
                  "type":"CategoryContent",
                  "allows_html":false
               },
               "sort":0
            }
         ],
         "videos":[  
            {  
               "id":339,
               "title":"ARIES Grille Guards",
               "subject_type":"",
               "videoType":{  
                  "id":4,
                  "name":"Product Video",
                  "icon":"http://www.curtmfg.com/assets/db7d1511-0203-4602-adaa-0b382f3d97ac.png"
               },
               "description":"ARIES grille guards are built to work as hard as you do every time they climb into that truck for another 4-wheel-drive kind of day. Not only are our grille guards designed to give you a tough, trail-ready look for their truck or SUV, but they also provide front end protection and a great avenue for customization.\r\n\r\nARIES grille guards are constructed with a one-piece, heavy-duty, 1 1/2\" steel tube design, and each one is made vehicle-specific for easier installation and a custom fit. Our Pro Series grille guards offer maximum customization with a patent-pending light bar housing and interchangeable cover plate. The Pro Series is available in a textured black carbon steel, while our standard ARIES Bar grille guards come in polished 304 stainless steel or semi-gloss black carbon steel.\r\n\r\nRegardless of the terrain they choose, dare your customers\r\nto conquer it with an ARIES grille guard.",
               "date_added":"0001-01-01T00:00:00Z",
               "date_modified":"0001-01-01T00:00:00Z",
               "thumbnail":"",
               "channel":[  
                  {  
                     "id":319,
                     "type":{  
                        "id":1,
                        "name":"YouTube",
                        "description":"YouTube"
                     },
                     "link":"https://www.youtube.com/watch?v=7tTMyHOjR68",
                     "embed_code":"",
                     "foreign_id":"",
                     "date_added":"0001-01-01T00:00:00Z",
                     "date_modified":"0001-01-01T00:00:00Z",
                     "title":"ARIES Grille Guards",
                     "description":"ARIES grille guards are built to work as hard as you do every time they climb into that truck for another 4-wheel-drive kind of day. Not only are our grille guards designed to give you a tough, trail-ready look for their truck or SUV, but they also provide front end protection and a great avenue for customization.\n\nARIES grille guards are constructed with a one-piece, heavy-duty, 1 1/2\" steel tube design, and each one is made vehicle-specific for easier installation and a custom fit. Our Pro Series grille guards offer maximum customization with a patent-pending light bar housing and interchangeable cover plate. The Pro Series is available in a textured black carbon steel, while our standard ARIES Bar grille guards come in polished 304 stainless steel or semi-gloss black carbon steel.\n\nRegardless of the terrain they choose, dare your customers\nto conquer it with an ARIES grille guard.",
                     "duration":""
                  }
               ],
               "cdn_file":[  
                  {  
                     "id":287,
                     "type":{  
                        "id":1,
                        "mime_type":"",
                        "title":"MP4",
                        "description":"MPEG 4 files with H264 video codec and AAC audio codec"
                     },
                     "path":"http://curtmfg.com/MasterLibrary/01RESOURCES/ARIES_Video/Product Video/ARIES Grille Guards.mp4",
                     "bucket":"",
                     "object_name":"",
                     "file_size":"",
                     "date_added":"0001-01-01T00:00:00Z",
                     "date_modified":"0001-01-01T00:00:00Z",
                     "date_uploaded":""
                  }
               ]
            }
         ],
         "part_ids":[  
            2060021,
            2060031,
            1231121963
         ],
         "brand":{  
            "id":3,
            "name":"ARIES",
            "code":"ARIES",
            "logo":{  
               "Scheme":"https",
               "Opaque":"",
               "User":null,
               "Host":"storage.googleapis.com",
               "Path":"/aries-logo/SVG_Logo (2c_white with black outline on transparent).svg",
               "RawPath":"/aries-logo/SVG_Logo%20(2c_white%20with%20black%20outline%20on%20transparent).svg",
               "RawQuery":"",
               "Fragment":""
            },
            "logo_alternate":{  
               "Scheme":"https",
               "Opaque":"",
               "User":null,
               "Host":"storage.googleapis.com",
               "Path":"/aries-logo/ARIES Logo (1c_red on transparent).png",
               "RawPath":"/aries-logo/ARIES%20Logo%20(1c_red%20on%20transparent).png",
               "RawQuery":"",
               "Fragment":""
            },
            "formal_name":"Aries Automotive",
            "long_name":"Aries Automotive",
            "primary_color":"#57111A",
            "autocareId":"BBRD",
            "websites":null
         },
         "product_listing":null
      }
   ],
   "videos":[  
      {  
         "title":"ARIES Grille Guards",
         "subject_type":"Product Video",
         "videoType":{  
            "name":"",
            "icon":""
         },
         "description":"ARIES grille guards are built to work as hard as you do every time they climb into that truck for another 4-wheel-drive kind of day. Not only are our grille guards designed to give you a tough, trail-ready look for their truck or SUV, but they also provide front end protection and a great avenue for customization.\r\n\r\nARIES grille guards are constructed with a one-piece, heavy-duty, 1 1/2\" steel tube design, and each one is made vehicle-specific for easier installation and a custom fit. Our Pro Series grille guards offer maximum customization with a patent-pending light bar housing and interchangeable cover plate. The Pro Series is available in a textured black carbon steel, while our standard ARIES Bar grille guards come in polished 304 stainless steel or semi-gloss black carbon steel.\r\n\r\nRegardless of the terrain they choose, dare your customers\r\nto conquer it with an ARIES grille guard.",
         "date_added":"0001-01-01T00:00:00Z",
         "date_modified":"0001-01-01T00:00:00Z",
         "thumbnail":"https://i.ytimg.com/vi/7tTMyHOjR68/default.jpg",
         "channel":[  
            {  
               "type":{  
                  "name":"YouTube"
               },
               "link":"https://www.youtube.com/watch?v=7tTMyHOjR68",
               "embed_code":"\u003ciframe width=\"560\" height=\"315\" src=\"https://www.youtube.com/embed/7tTMyHOjR68\" frameborder=\"0\" allowfullscreen\u003e\u003c/iframe\u003e",
               "foreign_id":"7tTMyHOjR68",
               "date_added":"0001-01-01T00:00:00Z",
               "date_modified":"0001-01-01T00:00:00Z",
               "title":"ARIES Grille Guards",
               "description":"ARIES grille guards are built to work as hard as you do every time they climb into that truck for another 4-wheel-drive kind of day. Not only are our grille guards designed to give you a tough, trail-ready look for their truck or SUV, but they also provide front end protection and a great avenue for customization.\n\nARIES grille guards are constructed with a one-piece, heavy-duty, 1 1/2\" steel tube design, and each one is made vehicle-specific for easier installation and a custom fit. Our Pro Series grille guards offer maximum customization with a patent-pending light bar housing and interchangeable cover plate. The Pro Series is available in a textured black carbon steel, while our standard ARIES Bar grille guards come in polished 304 stainless steel or semi-gloss black carbon steel.\n\nRegardless of the terrain they choose, dare your customers\nto conquer it with an ARIES grille guard.",
               "duration":"54s"
            }
         ],
         "cdn_file":[  
            {  
               "type":{  
                  "mime_type":"",
                  "title":"MP4",
                  "description":"MPEG 4 files with H264 video codec and AAC audio codec"
               },
               "path":"http://curtmfg.com/MasterLibrary/01RESOURCES/ARIES_Video/Product Video/ARIES Grille Guards.mp4",
               "bucket":"",
               "object_name":"",
               "file_size":"",
               "date_added":"0001-01-01T00:00:00Z",
               "date_modified":"0001-01-01T00:00:00Z",
               "date_uploaded":""
            }
         ]
      }
   ],
   "packages":[  
      {  
         "height":74,
         "width":17,
         "length":28,
         "weight":69,
         "dimensionUnit":"IN",
         "dimensionUnitLabel":"Inch",
         "weightUnit":"LB",
         "weightUnitLabel":"Pound",
         "packageUnit":"EA",
         "packageUnitLabel":"Each",
         "quantity":1,
         "name":"ShippingContainer"
      }
   ],
   "customer":{  
      "price":0,
      "cart_reference":0
   },
   "class":{  

   },
   "acesPartTypeId":14085,
   "inventory":{  
      "total_availability":0
   },
   "upc":"812410010049",
   "iconLayer":"",
   "mappedToVehicle":false
}
```


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

