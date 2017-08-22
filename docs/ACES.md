ACES
===
ACES (Aftermarket Catalog Exchange Standard) is the North American industry standard for the management and exchange of automotive catalog applications data. With ACES CURT Group publishes part information using standardized vehicle attributes, parts classifications, and qualifier statements. 

[ACES v3.2 Resources](https://www.autocare.org/ProductDetail.aspx?id=288&gmssopc=1)

List of endpoints

 - [Get ACES v3.2 XML](#aces-file)


## <a name="aces-file"></a>Get ACES Files `GET  - http://goapi.curtmfg.com/aces/3.2`
ACES v3.2 XML File

*Example:*

	http://goapi.curtmfg.com/aces/3.2?key=[API Key]&brandID=1


#### Parameters


| Paramter  |  Description |
|---|---|
| key **(required)** | Provide your API key  |
| brandID **(required)** | Brand querying that part number for (1=CURT, 3=ARIES, 4=Luverne) |



#### Response
Returns an ACES v3.2 XML file defined by Auto Care Association. 


