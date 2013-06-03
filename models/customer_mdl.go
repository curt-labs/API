package models

import (
	"../helpers/api"
	"../helpers/database"
	"../helpers/mymysql/mysql"
	"../helpers/redis"
	"../helpers/sortutil"
	"encoding/json"
	"math"
	"net/url"
	"strconv"
	"strings"
)

var (
	customerPriceStmt_Grouped = `select distinct cp.price, cp.partID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CustomerPricing cp on c.customerID = cp.cust_id
					where api_key = '%s'
					and cp.partID IN (%s)`

	customerPartStmt_Grouped = `select distinct ci.custPartID, ci.partID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CartIntegration ci on c.customerID = ci.custID
					where ak.api_key = '%s'
					and ci.partID IN (%s)`

	etailersStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				c.latitude, c.longitude, c.searchURL, c.logo, c.website,
				c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
				dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
				dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
				mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
				mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from Customer as c
				join DealerTypes as dt on c.dealer_type = dt.dealer_type
				join DealerTiers dtr on c.tier = dtr.ID
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where dt.online = 1 && c.isDummy = 0`
)

type Customer struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State                                *State
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude                  float64
	Website                              *url.URL
	Parent                               *Customer
	SearchUrl, Logo                      *url.URL
	DealerType                           DealerType
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
	Locations                            *[]CustomerLocation
	Users                                []CustomerUser
}

type CustomerLocation struct {
	Id                                     int
	Name, Email, Address, City, PostalCode string
	State                                  *State
	Phone, Fax                             string
	Latitude, Longitude                    float64
	CustomerId                             int
	ContactPerson                          string
	IsPrimary, ShippingDefault             bool
}

type DealerType struct {
	Id           int
	Type, Label  string
	Online, Show bool
	MapIcon      MapIcon
}

type DealerTier struct {
	Id   int
	Tier string
	Sort int
}

type MapIcon struct {
	Id, TierId             int
	MapIcon, MapIconShadow *url.URL
}

type MapGraphics struct {
	DealerTier DealerTier
	DealerType DealerType
	MapIcon    MapIcon
}

type GeoLocation struct {
	Latitude, Longitude float64
}

type DealerLocation struct {
	Id, LocationId                       int
	Name, Email, Address, Address2, City string
	State                                *State
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude, Distance        float64
	Website                              *url.URL
	Parent                               *Customer
	SearchUrl, Logo                      *url.URL
	DealerType                           DealerType
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
}

func (c *Customer) GetCustomer() (err error) {

	locationChan := make(chan int)
	go func() {
		if locErr := c.GetLocations(); locErr != nil {
			err = locErr
		}
		locationChan <- 1
	}()

	err = c.Basics()

	<-locationChan

	return err
}

func (c *Customer) Basics() error {

	qry, err := database.GetStatement("CustomerStmt")
	if database.MysqlError(err) {
		return err
	}

	row, res, err := qry.ExecFirst(c.Id)
	if database.MysqlError(err) {
		return err
	}

	customerID := res.Map("customerID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	address2 := res.Map("address2")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postal_code")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	iconId := res.Map("iconID")
	icon := res.Map("mapicon")
	shadow := res.Map("mapiconshadow")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	sURL, _ := url.Parse(row.Str(search))
	websiteURL, _ := url.Parse(row.Str(site))
	logoURL, _ := url.Parse(row.Str(logo))
	iconUrl, _ := url.Parse(row.Str(icon))
	shadowUrl, _ := url.Parse(row.Str(shadow))

	c.Id = row.Int(customerID)
	c.Name = row.Str(name)
	c.Email = row.Str(email)
	c.Address = row.Str(address)
	c.Address2 = row.Str(address2)
	c.City = row.Str(city)
	c.PostalCode = row.Str(zip)
	c.Phone = row.Str(phone)
	c.Fax = row.Str(fax)
	c.ContactPerson = row.Str(contact)
	c.Latitude = row.ForceFloat(lat)
	c.Longitude = row.ForceFloat(lon)
	c.Website = websiteURL
	c.SearchUrl = sURL
	c.Logo = logoURL
	c.DealerType = DealerType{
		Id:     row.Int(dealerTypeId),
		Type:   row.Str(dealerType),
		Label:  row.Str(typeLabel),
		Online: row.ForceBool(typeOnline),
		Show:   row.ForceBool(typeShow),
		MapIcon: MapIcon{
			Id:            row.Int(iconId),
			TierId:        row.Int(tierID),
			MapIcon:       iconUrl,
			MapIconShadow: shadowUrl,
		},
	}
	c.DealerTier = DealerTier{
		Id:   row.Int(tierID),
		Tier: row.Str(tier),
		Sort: row.Int(tierSort),
	}
	c.SalesRepresentative = row.Str(rep_name)
	c.SalesRepresentativeCode = row.Str(rep_code)
	c.MapixCode = row.Str(mpx_code)
	c.MapixDescription = row.Str(mpx_desc)

	ctry := Country{
		Id:           row.Int(countryID),
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	c.State = &State{
		Id:           row.Int(stateID),
		State:        row.Str(state),
		Abbreviation: row.Str(state_abbr),
		Country:      &ctry,
	}

	if row.Int(parentID) != 0 {
		parent := Customer{
			Id: row.Int(parentID),
		}
		if err = parent.GetCustomer(); err == nil {
			c.Parent = &parent
		}
	}

	return nil
}

func (c *Customer) GetLocations() error {

	qry, err := database.GetStatement("CustomerLocationStmt")
	if database.MysqlError(err) {
		return err
	}

	rows, res, err := qry.Exec(c.Id)
	if database.MysqlError(err) {
		return err
	}

	locationID := res.Map("locationID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	zip := res.Map("postalCode")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("cty_name")
	country_abbr := res.Map("cty_abbr")
	customerID := res.Map("cust_id")
	isPrimary := res.Map("isprimary")
	shipDefault := res.Map("ShippingDefault")

	var locs []CustomerLocation
	for _, row := range rows {
		l := CustomerLocation{
			Id:              row.Int(locationID),
			Name:            row.Str(name),
			Email:           row.Str(email),
			Address:         row.Str(address),
			City:            row.Str(city),
			PostalCode:      row.Str(zip),
			Phone:           row.Str(phone),
			Fax:             row.Str(fax),
			ContactPerson:   row.Str(contact),
			CustomerId:      row.Int(customerID),
			Latitude:        row.ForceFloat(lat),
			Longitude:       row.ForceFloat(lon),
			IsPrimary:       row.ForceBool(isPrimary),
			ShippingDefault: row.ForceBool(shipDefault),
		}

		ctry := Country{
			Id:           row.Int(countryID),
			Country:      row.Str(country),
			Abbreviation: row.Str(country_abbr),
		}

		l.State = &State{
			Id:           row.Int(stateID),
			State:        row.Str(state),
			Abbreviation: row.Str(state_abbr),
			Country:      &ctry,
		}
		locs = append(locs, l)
	}
	c.Locations = &locs
	return nil
}

func (c *Customer) GetUsers() (users []CustomerUser, err error) {

	qry, err := database.GetStatement("CustomerUserStmt")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec(c.Id)
	if database.MysqlError(err) {
		return
	}
	user_id := res.Map("id")
	name := res.Map("name")
	mail := res.Map("email")
	date := res.Map("date_added")
	active := res.Map("active")
	sudo := res.Map("isSudo")

	for _, row := range rows {
		var u CustomerUser
		u.Name = row.Str(name)
		u.Email = row.Str(mail)
		u.Active = row.Int(active) == 1
		u.Sudo = row.Int(sudo) == 1
		u.Current = false
		u.Id = row.Str(user_id)
		u.DateAdded = row.ForceLocaltime(date)

		users = append(users, u)
	}
	return
}

func GetCustomerPrice(api_key string, part_id int) (price float64, err error) {
	qry, err := database.GetStatement("CustomerPriceStmt")
	if database.MysqlError(err) {
		return
	}

	row, _, err := qry.ExecFirst(api_key, part_id)
	if database.MysqlError(err) {
		return
	}
	if len(row) == 1 {
		price = row.Float(0)
	}

	return
}

func (lookup *Lookup) GetCustomerPrice(api_key string) (prices map[int]float64, err error) {
	prices = make(map[int]float64, len(lookup.Parts))

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}
	if len(ids) == 0 {
		return
	}

	rows, res, err := database.Db.Query(customerPriceStmt_Grouped, api_key, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	} else if len(rows) == 0 {
		return
	}

	price := res.Map("price")
	partID := res.Map("partID")

	for _, row := range rows {
		pId := row.Int(partID)
		pr := row.Float(price)
		prices[pId] = pr
	}

	return
}

func GetCustomerCartReference(api_key string, part_id int) (ref int, err error) {
	qry, err := database.GetStatement("CustomerPartStmt")
	if database.MysqlError(err) {
		return
	}

	row, _, err := qry.ExecFirst(api_key, part_id)
	if database.MysqlError(err) {
		return
	}

	if len(row) == 1 {
		ref = row.Int(0)
	}

	return
}

func (lookup *Lookup) GetCustomerCartReference(api_key string) (references map[int]int, err error) {

	references = make(map[int]int, len(lookup.Parts))

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}
	if len(ids) == 0 {
		return
	}

	rows, res, err := database.Db.Query(customerPartStmt_Grouped, api_key, strings.Join(ids, ","))
	if err != nil {
		return
	} else if len(rows) == 0 {
		return
	}

	partID := res.Map("partID")
	custPartID := res.Map("custPartID")

	for _, row := range rows {
		pId := row.Int(partID)
		ref := row.Int(custPartID)
		references[pId] = ref
	}

	return
}

/* Internal Use Only */

func GetEtailers() (dealers []Customer, err error) {

	redis_key := "goapi:dealers:etailers"

	// Attempt to get the etailers from Redis
	etailer_bytes, err := redis.RedisClient.Get(redis_key)
	if len(etailer_bytes) > 0 {
		err = json.Unmarshal(etailer_bytes, &dealers)
		if err == nil {
			return
		}
	}

	rows, res, err := database.Db.Query(etailersStmt)
	if database.MysqlError(err) {
		return
	}

	customerID := res.Map("customerID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	address2 := res.Map("address2")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postal_code")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	iconId := res.Map("iconID")
	icon := res.Map("mapicon")
	shadow := res.Map("mapiconshadow")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	c := make(chan int)

	for _, row := range rows {
		go func(r mysql.Row, ch chan int) {

			sURL, _ := url.Parse(r.Str(search))
			websiteURL, _ := url.Parse(r.Str(site))
			logoURL, _ := url.Parse(r.Str(logo))
			iconUrl, _ := url.Parse(row.Str(icon))
			shadowUrl, _ := url.Parse(row.Str(shadow))

			cust := Customer{
				Id:            r.Int(customerID),
				Name:          r.Str(name),
				Email:         r.Str(email),
				Address:       r.Str(address),
				Address2:      r.Str(address2),
				City:          r.Str(city),
				PostalCode:    r.Str(zip),
				Phone:         r.Str(phone),
				Fax:           r.Str(fax),
				ContactPerson: r.Str(contact),
				Latitude:      r.ForceFloat(lat),
				Longitude:     r.ForceFloat(lon),
				Website:       websiteURL,
				SearchUrl:     sURL,
				Logo:          logoURL,
				DealerType: DealerType{
					Id:     row.Int(dealerTypeId),
					Type:   row.Str(dealerType),
					Label:  row.Str(typeLabel),
					Online: row.ForceBool(typeOnline),
					Show:   row.ForceBool(typeShow),
					MapIcon: MapIcon{
						Id:            row.Int(iconId),
						TierId:        row.Int(tierID),
						MapIcon:       iconUrl,
						MapIconShadow: shadowUrl,
					},
				},
				DealerTier: DealerTier{
					Id:   row.Int(tierID),
					Tier: row.Str(tier),
					Sort: row.Int(tierSort),
				},
				SalesRepresentative:     r.Str(rep_name),
				SalesRepresentativeCode: r.Str(rep_code),
				MapixCode:               r.Str(mpx_code),
				MapixDescription:        r.Str(mpx_desc),
			}

			ctry := Country{
				Id:           r.Int(countryID),
				Country:      r.Str(country),
				Abbreviation: r.Str(country_abbr),
			}

			cust.State = &State{
				Id:           r.Int(stateID),
				State:        r.Str(state),
				Abbreviation: r.Str(state_abbr),
				Country:      &ctry,
			}

			_ = cust.GetLocations()

			if r.Int(parentID) != 0 {
				parent := Customer{
					Id: r.Int(parentID),
				}
				if err = parent.GetCustomer(); err == nil {
					cust.Parent = &parent
				}
			}
			dealers = append(dealers, cust)

			ch <- 1
		}(row, c)

	}

	for _, _ = range rows {
		<-c
	}

	if etailer_bytes, err = json.Marshal(dealers); err == nil {
		redis.RedisClient.Setex(redis_key, 86400, etailer_bytes)
	}

	return
}

func GetLocalDealers(center string, latlng string) (dealers []DealerLocation, err error) {

	qry, err := database.GetStatement("LocalDealersStmt")
	if database.MysqlError(err) {
		return
	}

	var latlngs []string
	var center_latlngs []string

	clat := api_helpers.CENTER_LATITUDE
	clong := api_helpers.CENTER_LONGITUDE
	swlat := api_helpers.SOUTWEST_LATITUDE
	swlong := api_helpers.SOUTHWEST_LONGITUDE
	swlong2 := api_helpers.SOUTHWEST_LONGITUDE
	nelat := api_helpers.NORTHEAST_LATITUDE
	nelong := api_helpers.NORTHEAST_LONGITUDE
	nelong2 := api_helpers.NORTHEAST_LONGITUDE

	// Get the center point
	if center != "" {
		center_latlngs = strings.Split(center, ",")
		if len(center_latlngs) == 2 {
			center_lat, err := strconv.ParseFloat(center_latlngs[0], 64)
			if err == nil {
				center_lon, err := strconv.ParseFloat(center_latlngs[1], 64)
				if err == nil {
					clat = center_lat
					clong = center_lon
				}
			}
		}
	}

	// Get the boundary points
	if latlng != "" {
		latlngs = strings.Split(latlng, ",")
		if len(latlngs) == 4 {
			sw_lat, err := strconv.ParseFloat(latlngs[0], 64)
			if err == nil {
				sw_lon, err := strconv.ParseFloat(latlngs[1], 64)
				if err == nil {
					ne_lat, err := strconv.ParseFloat(latlngs[2], 64)
					if err == nil {
						ne_lon, err := strconv.ParseFloat(latlngs[3], 64)
						if err == nil {
							swlat = sw_lat
							swlong = sw_lon
							swlong2 = sw_lon
							nelat = ne_lat
							nelong = ne_lon
							nelong2 = ne_lon
						}
					}
				}
			}
		}
	}

	if swlong > nelong {
		swlong = -180
		nelong2 = 180
	}

	distance_a := getViewPortWidth(swlat, swlong, clat, clong)
	distance_b := getViewPortWidth(nelat, nelong2, clat, clong)

	view_distance := distance_b
	if distance_a > distance_b {
		view_distance = distance_a
	}

	params := struct {
		Earth        float64
		Clat1        float64
		Clat2        float64
		Clong1       float64
		Clong2       float64
		Clat3        float64
		Clat4        float64
		Clat5        float64
		Clong3       float64
		Clong4       float64
		Clat6        float64
		ViewDistance float64
		SWLat        float64
		NELat        float64
		SWLong       float64
		NELong       float64
		SWLong2      float64
		NELong2      float64
	}{
		api_helpers.EARTH,
		clat,
		clat,
		clong,
		clong,
		clat,
		clat,
		clat,
		clong,
		clong,
		clat,
		view_distance,
		swlat,
		nelat,
		swlong,
		nelong,
		swlong2,
		nelong2,
	}

	rows, res, err := qry.Exec(params)
	if database.MysqlError(err) {
		return
	}

	customerID := res.Map("customerID")
	locationID := res.Map("locationID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postalCode")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	iconId := res.Map("iconID")
	icon := res.Map("mapicon")
	shadow := res.Map("mapiconshadow")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")

	for _, r := range rows {
		sURL, _ := url.Parse(r.Str(search))
		websiteURL, _ := url.Parse(r.Str(site))
		logoURL, _ := url.Parse(r.Str(logo))
		iconUrl, _ := url.Parse(r.Str(icon))
		shadowUrl, _ := url.Parse(r.Str(shadow))

		cust := DealerLocation{
			Id:            r.Int(customerID),
			LocationId:    r.Int(locationID),
			Name:          r.Str(name),
			Email:         r.Str(email),
			Address:       r.Str(address),
			City:          r.Str(city),
			PostalCode:    r.Str(zip),
			Phone:         r.Str(phone),
			Fax:           r.Str(fax),
			ContactPerson: r.Str(contact),
			Latitude:      r.ForceFloat(lat),
			Longitude:     r.ForceFloat(lon),
			Website:       websiteURL,
			SearchUrl:     sURL,
			Logo:          logoURL,
			DealerType: DealerType{
				Id:     r.Int(dealerTypeId),
				Type:   r.Str(dealerType),
				Label:  r.Str(typeLabel),
				Online: r.ForceBool(typeOnline),
				Show:   r.ForceBool(typeShow),
				MapIcon: MapIcon{
					Id:            r.Int(iconId),
					TierId:        r.Int(tierID),
					MapIcon:       iconUrl,
					MapIconShadow: shadowUrl,
				},
			},
			DealerTier: DealerTier{
				Id:   r.Int(tierID),
				Tier: r.Str(tier),
				Sort: r.Int(tierSort),
			},
			SalesRepresentative:     r.Str(rep_name),
			SalesRepresentativeCode: r.Str(rep_code),
			MapixCode:               r.Str(mpx_code),
			MapixDescription:        r.Str(mpx_desc),
		}

		cust.Distance = api_helpers.EARTH * (2 * math.Atan2(
			math.Sqrt((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))),
			math.Sqrt(1-((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))))))

		ctry := Country{
			Id:           r.Int(countryID),
			Country:      r.Str(country),
			Abbreviation: r.Str(country_abbr),
		}

		cust.State = &State{
			Id:           r.Int(stateID),
			State:        r.Str(state),
			Abbreviation: r.Str(state_abbr),
			Country:      &ctry,
		}

		dealers = append(dealers, cust)

	}
	sortutil.AscByField(dealers, "Distance")
	return
}

type StateRegion struct {
	Id                 int
	Name, Abbreviation string
	Count              int
	Polygons           []MapPolygon
}

type MapPolygon struct {
	Id          int
	Coordinates []GeoLocation
}

func GetLocalRegions() (regions []StateRegion, err error) {

	redis_key := "goapi:local:regions"

	// Attempt to get the local regions from Redis
	regions_bytes, err := redis.RedisClient.Get(redis_key)
	if len(regions_bytes) > 0 {
		err = json.Unmarshal(regions_bytes, &regions)
		if err == nil {
			return
		}
	}

	polyQuery, err := database.GetStatement("PolygonStmt")
	if database.MysqlError(err) {
		return
	}

	coordQry, err := database.GetStatement("MapPolygonCoordinatesForStateStmt")
	if database.MysqlError(err) {
		return
	}

	// Get the local regions from the database
	_, _, _ = database.Db.Query("SET SESSION group_concat_max_len = 100024")
	rows, res, err := polyQuery.Exec()
	_, _, _ = database.Db.Query("SET SESSION group_concat_max_len = 1024")

	if database.MysqlError(err) || rows == nil {
		return
	}

	ch := make(chan int)

	for _, row := range rows {
		go func(c chan int, regRow mysql.Row, regRes mysql.Result) {

			// Populate the StateRegion with state data
			stateID := regRes.Map("stateID")
			state := regRes.Map("state")
			abbr := regRes.Map("abbr")
			count := regRes.Map("count")

			reg := StateRegion{
				Id:           regRow.Int(stateID),
				Name:         regRow.Str(state),
				Abbreviation: regRow.Str(abbr),
				Count:        regRow.Int(count),
			}

			// Build out the polygons for this state
			// including latitude and longitude
			coordRows, coordRes, err := coordQry.Exec(reg.Id)
			if err == nil {
				polyId := coordRes.Map("ID")
				lat := coordRes.Map("latitude")
				lon := coordRes.Map("longitude")

				polygons := make(map[int]MapPolygon, 0)

				// Loops the coordinates (latitude, longitude)
				for _, coordRow := range coordRows {
					// Check if we have an index for this polygon created
					if _, ok := polygons[coordRow.Int(polyId)]; !ok {
						// First time hitting this polygon
						// so we'll create one
						polygons[coordRow.Int(polyId)] = MapPolygon{
							Id:          coordRow.Int(polyId),
							Coordinates: make([]GeoLocation, 0),
						}
					}

					// Add the GeoLocartion info to our polygon
					poly := polygons[coordRow.Int(polyId)]
					poly.Coordinates = append(poly.Coordinates, GeoLocation{coordRow.ForceFloat(lat), coordRow.ForceFloat(lon)})
					polygons[coordRow.Int(polyId)] = poly
				}

				// We need to drop the key/value pair
				// our end user doesn't need that
				var polys []MapPolygon
				for _, poly := range polygons {
					polys = append(polys, poly)
				}

				reg.Polygons = polys
			}

			regions = append(regions, reg)
			c <- 1
		}(ch, row, res)
	}

	for _, _ = range rows {
		<-ch
	}

	if regions_bytes, err = json.Marshal(regions); err == nil {
		// We're not going to set the expiration on this
		// it won't ever change...until the San Andreas fault
		// completely drops the western part of CA anyway :/
		redis.RedisClient.Set(redis_key, regions_bytes)
	}

	return
}

func getViewPortWidth(lat1 float64, lon1 float64, lat2 float64, long2 float64) float64 {
	dlat := (lat2 - lat1) * (math.Pi / 180)
	dlon := (long2 - lon1) * (math.Pi / 180)

	lat1 = lat1 * (math.Pi / 180)
	lat2 = lat2 * (math.Pi / 180)

	a := (math.Sin(dlat/2) * math.Sin(dlat/2)) + ((math.Sin(dlon/2))*(math.Sin(dlon/2)))*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return api_helpers.EARTH * c
}

func GetLocalDealerTiers() (tiers []DealerTier) {

	redis_key := "goapi:local:tiers"

	// Attempt to get the local regions from Redis
	redis_bytes, err := redis.RedisClient.Get(redis_key)
	if len(redis_bytes) > 0 {
		err = json.Unmarshal(redis_bytes, &tiers)
		if err == nil {
			return
		}
	}

	qry, err := database.GetStatement("GetLocalDealerTiers")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	id := res.Map("ID")
	tier := res.Map("tier")
	sort := res.Map("sort")

	for _, row := range rows {
		dTier := DealerTier{
			Id:   row.Int(id),
			Tier: row.Str(tier),
			Sort: row.Int(sort),
		}
		tiers = append(tiers, dTier)
	}

	if redis_bytes, err = json.Marshal(tiers); err == nil {
		redis.RedisClient.Setex(redis_key, 86400, redis_bytes)
	}

	return
}

func GetLocalDealerTypes() (graphics []MapGraphics) {

	redis_key := "goapi:local:types"

	// Attempt to get the local regions from Redis
	redis_bytes, err := redis.RedisClient.Get(redis_key)
	if len(redis_bytes) > 0 {
		err = json.Unmarshal(redis_bytes, &graphics)
		if err == nil {
			return
		}
	}

	qry, err := database.GetStatement("GetLocalDealerTypes")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	iconId := res.Map("iconId")
	mapicon := res.Map("mapicon")
	mapiconshadow := res.Map("mapiconshadow")
	dealerTypeId := res.Map("dealerTypeId")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("online")
	typeShow := res.Map("show")
	typeLabel := res.Map("label")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")

	for _, row := range rows {
		iconUrl, _ := url.Parse(row.Str(mapicon))
		shadowUrl, _ := url.Parse(row.Str(mapiconshadow))

		gx := MapGraphics{
			DealerTier: DealerTier{
				Id:   row.Int(tierID),
				Tier: row.Str(tier),
				Sort: row.Int(tierSort),
			},
			DealerType: DealerType{
				Id:     row.Int(dealerTypeId),
				Type:   row.Str(dealerType),
				Label:  row.Str(typeLabel),
				Online: row.ForceBool(typeOnline),
				Show:   row.ForceBool(typeShow),
			},
			MapIcon: MapIcon{
				Id:            row.Int(iconId),
				TierId:        row.Int(tierID),
				MapIcon:       iconUrl,
				MapIconShadow: shadowUrl,
			},
		}
		graphics = append(graphics, gx)
	}

	if redis_bytes, err = json.Marshal(graphics); err == nil {
		redis.RedisClient.Setex(redis_key, 86400, redis_bytes)
	}

	return
}

func GetWhereToBuyDealers() (customers []Customer) {

	redis_key := "goapi:dealers:wheretobuy"

	// Attempt to get the local regions from Redis
	redis_bytes, err := redis.RedisClient.Get(redis_key)
	if len(redis_bytes) > 0 {
		err = json.Unmarshal(redis_bytes, &customers)
		if err == nil {
			return
		}
	}

	qry, err := database.GetStatement("WhereToBuyDealersStmt")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	customerID := res.Map("customerID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	address2 := res.Map("address2")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postal_code")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	iconId := res.Map("iconID")
	icon := res.Map("mapicon")
	shadow := res.Map("mapiconshadow")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	c := make(chan int)

	for _, row := range rows {
		go func(r mysql.Row, ch chan int) {

			sURL, _ := url.Parse(r.Str(search))
			websiteURL, _ := url.Parse(r.Str(site))
			logoURL, _ := url.Parse(r.Str(logo))
			iconUrl, _ := url.Parse(row.Str(icon))
			shadowUrl, _ := url.Parse(row.Str(shadow))

			cust := Customer{
				Id:            r.Int(customerID),
				Name:          r.Str(name),
				Email:         r.Str(email),
				Address:       r.Str(address),
				Address2:      r.Str(address2),
				City:          r.Str(city),
				PostalCode:    r.Str(zip),
				Phone:         r.Str(phone),
				Fax:           r.Str(fax),
				ContactPerson: r.Str(contact),
				Latitude:      r.ForceFloat(lat),
				Longitude:     r.ForceFloat(lon),
				Website:       websiteURL,
				SearchUrl:     sURL,
				Logo:          logoURL,
				DealerType: DealerType{
					Id:     row.Int(dealerTypeId),
					Type:   row.Str(dealerType),
					Label:  row.Str(typeLabel),
					Online: row.ForceBool(typeOnline),
					Show:   row.ForceBool(typeShow),
					MapIcon: MapIcon{
						Id:            row.Int(iconId),
						TierId:        row.Int(tierID),
						MapIcon:       iconUrl,
						MapIconShadow: shadowUrl,
					},
				},
				DealerTier: DealerTier{
					Id:   row.Int(tierID),
					Tier: row.Str(tier),
					Sort: row.Int(tierSort),
				},
				SalesRepresentative:     r.Str(rep_name),
				SalesRepresentativeCode: r.Str(rep_code),
				MapixCode:               r.Str(mpx_code),
				MapixDescription:        r.Str(mpx_desc),
			}

			ctry := Country{
				Id:           r.Int(countryID),
				Country:      r.Str(country),
				Abbreviation: r.Str(country_abbr),
			}

			cust.State = &State{
				Id:           r.Int(stateID),
				State:        r.Str(state),
				Abbreviation: r.Str(state_abbr),
				Country:      &ctry,
			}

			_ = cust.GetLocations()

			if r.Int(parentID) != 0 {
				parent := Customer{
					Id: r.Int(parentID),
				}
				if err = parent.GetCustomer(); err == nil {
					cust.Parent = &parent
				}
			}
			customers = append(customers, cust)

			ch <- 1
		}(row, c)

	}

	for _, _ = range rows {
		<-c
	}

	if redis_bytes, err = json.Marshal(customers); err == nil {
		redis.RedisClient.Setex(redis_key, 86400, redis_bytes)
	}

	return
}

func GetLocationById(id int) (location DealerLocation, err error) {

	qry, err := database.GetStatement("CustomerLocationByIdStmt")
	if database.MysqlError(err) {
		return
	}

	row, res, err := qry.ExecFirst(id)
	if database.MysqlError(err) {
		return
	}

	stateID := res.Map("stateID")
	state := res.Map("state")
	abbr := res.Map("abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	locationID := res.Map("locationID")
	name := res.Map("name")
	address := res.Map("address")
	city := res.Map("city")
	postalCode := res.Map("postalCode")
	email := res.Map("email")
	phone := res.Map("phone")
	fax := res.Map("fax")
	latitude := res.Map("latitude")
	longitude := res.Map("longitude")
	cust_id := res.Map("cust_id")
	contactPerson := res.Map("contact_person")
	showWebsite := res.Map("showWebsite")
	website := res.Map("website")
	elocal := res.Map("eLocalURL")

	site := row.Str(website)
	var siteUrl *url.URL
	if row.ForceBool(showWebsite) {
		if site == "" {
			site = row.Str(elocal)
		}
		if site != "" {
			siteUrl, _ = url.Parse(site)
		}
	}

	dType := DealerType{
		Id:     row.Int(dealerTypeId),
		Type:   row.Str(dealerType),
		Label:  row.Str(typeLabel),
		Online: row.ForceBool(typeOnline),
		Show:   row.ForceBool(typeShow),
	}

	dealerTier := DealerTier{
		Id:   row.Int(tierID),
		Tier: row.Str(tier),
		Sort: row.Int(tierSort),
	}

	location = DealerLocation{
		Id:         row.Int(cust_id),
		LocationId: row.Int(locationID),
		Name:       row.Str(name),
		Address:    row.Str(address),
		City:       row.Str(city),
		PostalCode: row.Str(postalCode),
		State: &State{
			Id:           row.Int(stateID),
			State:        row.Str(state),
			Abbreviation: row.Str(abbr),
		},
		Email:         row.Str(email),
		Phone:         row.Str(phone),
		Fax:           row.Str(fax),
		Latitude:      row.ForceFloat(latitude),
		Longitude:     row.ForceFloat(longitude),
		ContactPerson: row.Str(contactPerson),
		Website:       siteUrl,
		DealerType:    dType,
		DealerTier:    dealerTier,
	}

	return
}

func SearchLocations(term string) (locations []DealerLocation, err error) {

	qry, err := database.GetStatement("SearchDealerLocations")
	if database.MysqlError(err) {
		return
	}

	term = "%" + term + "%"
	rows, res, err := qry.Exec(term, term)
	if database.MysqlError(err) {
		return
	}

	stateID := res.Map("stateID")
	state := res.Map("state")
	abbr := res.Map("abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	locationID := res.Map("locationID")
	name := res.Map("name")
	address := res.Map("address")
	city := res.Map("city")
	postalCode := res.Map("postalCode")
	email := res.Map("email")
	phone := res.Map("phone")
	fax := res.Map("fax")
	latitude := res.Map("latitude")
	longitude := res.Map("longitude")
	cust_id := res.Map("cust_id")
	contactPerson := res.Map("contact_person")
	showWebsite := res.Map("showWebsite")
	website := res.Map("website")
	elocal := res.Map("eLocalURL")

	for _, row := range rows {

		site := row.Str(website)
		var siteUrl *url.URL
		if row.ForceBool(showWebsite) {
			if site == "" {
				site = row.Str(elocal)
			}
			if site != "" {
				siteUrl, _ = url.Parse(site)
			}
		}

		dealerType := DealerType{
			Id:     row.Int(dealerTypeId),
			Type:   row.Str(dealerType),
			Label:  row.Str(typeLabel),
			Online: row.ForceBool(typeOnline),
			Show:   row.ForceBool(typeShow),
		}

		dealerTier := DealerTier{
			Id:   row.Int(tierID),
			Tier: row.Str(tier),
			Sort: row.Int(tierSort),
		}

		loc := DealerLocation{
			Id:         row.Int(cust_id),
			LocationId: row.Int(locationID),
			Name:       row.Str(name),
			Address:    row.Str(address),
			City:       row.Str(city),
			PostalCode: row.Str(postalCode),
			State: &State{
				Id:           row.Int(stateID),
				State:        row.Str(state),
				Abbreviation: row.Str(abbr),
			},
			Email:         row.Str(email),
			Phone:         row.Str(phone),
			Fax:           row.Str(fax),
			Latitude:      row.ForceFloat(latitude),
			Longitude:     row.ForceFloat(longitude),
			ContactPerson: row.Str(contactPerson),
			Website:       siteUrl,
			DealerType:    dealerType,
			DealerTier:    dealerTier,
		}
		locations = append(locations, loc)
	}

	return
}

func SearchLocationsByType(term string) (locations []DealerLocation, err error) {

	qry, err := database.GetStatement("SearchDealerLocationsByType")
	if database.MysqlError(err) {
		return
	}

	term = "%" + term + "%"
	rows, res, err := qry.Exec(term, term)
	if database.MysqlError(err) {
		return
	}

	stateID := res.Map("stateID")
	state := res.Map("state")
	abbr := res.Map("abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	locationID := res.Map("locationID")
	name := res.Map("name")
	address := res.Map("address")
	city := res.Map("city")
	postalCode := res.Map("postalCode")
	email := res.Map("email")
	phone := res.Map("phone")
	fax := res.Map("fax")
	latitude := res.Map("latitude")
	longitude := res.Map("longitude")
	cust_id := res.Map("cust_id")
	contactPerson := res.Map("contact_person")
	showWebsite := res.Map("showWebsite")
	website := res.Map("website")
	elocal := res.Map("eLocalURL")

	for _, row := range rows {

		site := row.Str(website)
		var siteUrl *url.URL
		if row.ForceBool(showWebsite) {
			if site == "" {
				site = row.Str(elocal)
			}
			if site != "" {
				siteUrl, _ = url.Parse(site)
			}
		}

		dealerType := DealerType{
			Id:     row.Int(dealerTypeId),
			Type:   row.Str(dealerType),
			Label:  row.Str(typeLabel),
			Online: row.ForceBool(typeOnline),
			Show:   row.ForceBool(typeShow),
		}

		dealerTier := DealerTier{
			Id:   row.Int(tierID),
			Tier: row.Str(tier),
			Sort: row.Int(tierSort),
		}

		loc := DealerLocation{
			Name:       row.Str(name),
			Website:    siteUrl,
			DealerType: dealerType,
			DealerTier: dealerTier,
			Id:         row.Int(cust_id),
			LocationId: row.Int(locationID),
			Address:    row.Str(address),
			City:       row.Str(city),
			PostalCode: row.Str(postalCode),
			State: &State{
				Id:           row.Int(stateID),
				State:        row.Str(state),
				Abbreviation: row.Str(abbr),
			},
			Email:         row.Str(email),
			Phone:         row.Str(phone),
			Fax:           row.Str(fax),
			Latitude:      row.ForceFloat(latitude),
			Longitude:     row.ForceFloat(longitude),
			ContactPerson: row.Str(contactPerson),
		}
		locations = append(locations, loc)
	}

	return
}

func SearchLocationsByLatLng(loc GeoLocation) (locations []DealerLocation, err error) {

	qry, err := database.GetStatement("SearchDealerLocationsByLatLng")
	if database.MysqlError(err) {
		return
	}

	params := struct {
		Earth       float64
		Lat1        float64
		Lat2        float64
		Lon1        float64
		Lon2        float64
		LatRadians1 float64
		Lat3        float64
		Lat4        float64
		Lon3        float64
		Lon4        float64
		LatRadians2 float64
	}{
		api_helpers.EARTH,
		loc.Latitude,
		loc.Latitude,
		loc.Longitude,
		loc.Longitude,
		loc.LatitudeRadians(),
		loc.Latitude,
		loc.Latitude,
		loc.Longitude,
		loc.Longitude,
		loc.LatitudeRadians(),
	}

	rows, res, err := qry.Exec(params)
	if database.MysqlError(err) {
		return
	}

	stateID := res.Map("stateID")
	state := res.Map("state")
	abbr := res.Map("abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	locationID := res.Map("locationID")
	name := res.Map("name")
	address := res.Map("address")
	city := res.Map("city")
	postalCode := res.Map("postalCode")
	email := res.Map("email")
	phone := res.Map("phone")
	fax := res.Map("fax")
	latitude := res.Map("latitude")
	longitude := res.Map("longitude")
	cust_id := res.Map("cust_id")
	contactPerson := res.Map("contact_person")
	showWebsite := res.Map("showWebsite")
	website := res.Map("website")
	elocal := res.Map("eLocalURL")

	for _, row := range rows {

		site := row.Str(website)
		var siteUrl *url.URL
		if row.ForceBool(showWebsite) {
			if site == "" {
				site = row.Str(elocal)
			}
			if site != "" {
				siteUrl, _ = url.Parse(site)
			}
		}

		dealerType := DealerType{
			Id:     row.Int(dealerTypeId),
			Type:   row.Str(dealerType),
			Label:  row.Str(typeLabel),
			Online: row.ForceBool(typeOnline),
			Show:   row.ForceBool(typeShow),
		}

		dealerTier := DealerTier{
			Id:   row.Int(tierID),
			Tier: row.Str(tier),
			Sort: row.Int(tierSort),
		}

		loc := DealerLocation{
			Name:       row.Str(name),
			Website:    siteUrl,
			DealerType: dealerType,
			DealerTier: dealerTier,
			Id:         row.Int(cust_id),
			LocationId: row.Int(locationID),
			Address:    row.Str(address),
			City:       row.Str(city),
			PostalCode: row.Str(postalCode),
			State: &State{
				Id:           row.Int(stateID),
				State:        row.Str(state),
				Abbreviation: row.Str(abbr),
			},
			Email:         row.Str(email),
			Phone:         row.Str(phone),
			Fax:           row.Str(fax),
			Latitude:      row.ForceFloat(latitude),
			Longitude:     row.ForceFloat(longitude),
			ContactPerson: row.Str(contactPerson),
		}
		locations = append(locations, loc)
	}

	return
}

func (g *GeoLocation) LatitudeRadians() float64 {
	return (g.Latitude * (math.Pi / 180))
}

func (g *GeoLocation) LongitudeRadians() float64 {
	return (g.Longitude * (math.Pi / 180))
}
