package models

import (
	"../helpers/database"
	"../helpers/redis"
	"../helpers/sortutil"
	"encoding/json"
	"github.com/ziutek/mymysql/mysql"
	"math"
	"net/url"
	"strconv"
	"strings"
)

const (
	EARTH               = 3963.1676 // radius of Earth in miles
	SOUTWEST_LATITUDE   = -90.00
	SOUTHWEST_LONGITUDE = -180.00
	NORTHEAST_LATITUDE  = 90.00
	NORTHEAST_LONGITUDE = 180.00
	CENTER_LATITUDE     = 44.79300
	CENTER_LONGITUDE    = -91.41048
)

var (
	customerPriceStmt = `select distinct cp.price from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CustomerPricing cp on c.customerID = cp.cust_id
					where api_key = ?
					and cp.partID = ?`

	customerPriceStmt_Grouped = `select distinct cp.price, cp.partID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CustomerPricing cp on c.customerID = cp.cust_id
					where api_key = '%s'
					and cp.partID IN (%s)`

	customerPartStmt = `select distinct ci.custPartID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CartIntegration ci on c.customerID = ci.custID
					where ak.api_key = ?
					and ci.partID = ?`

	customerPartStmt_Grouped = `select distinct ci.custPartID, ci.partID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CartIntegration ci on c.customerID = ci.custID
					where ak.api_key = '%s'
					and ci.partID IN (%s)`

	customerStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				c.latitude, c.longitude, c.searchURL, c.logo, c.website,
				c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
				d_types.type as dealer_type, d_tier.tier as dealer_tier, mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as d_types on c.dealer_type = d_types.dealer_type
				left join DealerTiers d_tier on c.tier = d_tier.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where c.customerID = ?`

	customerLocationsStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
					cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
					cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
					s.state, s.abbr as state_abbr, cty.name as cty_name, cty.abbr as cty_abbr
					from CustomerLocations as cl
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					where cl.cust_id = ?`

	customerUsersStmt = `select cu.* from CustomerUser as cu
					join Customer as c on cu.cust_ID = c.cust_id
					where c.customerID = '?'
					&& cu.active = 1`

	etailersStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				c.latitude, c.longitude, c.searchURL, c.logo, c.website,
				c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
				d_types.type as dealer_type, d_tier.tier as dealer_tier, mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from Customer as c
				join DealerTypes as d_types on c.dealer_type = d_types.dealer_type
				join DealerTiers d_tier on c.tier = d_tier.ID
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where d_types.online = 1 && c.isDummy = 0`

	// localDealersStmt(earth, center_latitude, center_latitude, center_longitude, center_longitude, center_latitude, center_latitude, center_latitude, center_longitude, center_longitude, center_latitude, view_distance, swlat, nelat, swlong, nelong, swlong2, nelong2)
	localDealersStmt = `select c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
				cl.latitude, cl.longitude, c.searchURL, c.logo, c.website,
				cl.postalCode, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
				dt.type as dealer_type, dtr.tier as dealer_tier, mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from CustomerLocations as cl
				join Customer as c on cl.cust_id = c.cust_id
				join DealerTypes as dt on c.dealer_type = dt.dealer_type
				join DealerTiers as dtr on c.tier = dtr.ID
				left join States as s on cl.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where dt.online = 0 && c.isDummy = 0 && dt.show = 1 &&
				( ? * (
                                                            2 * ATAN2(
                                                                SQRT((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))),
                                                                SQRT(1 - ((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))))
                                                            )
                                                        ) < ?)
				&& (
					(cl.latitude >= ? && cl.latitude <= ?) 
					&&
                                              		(cl.longitude >= ? && cl.longitude <= ?) 
                                              		||
                                              		(cl.longitude >= ? && cl.longitude <= ?)
                                              	)
				order by dtr.sort desc`

	mapCoordinatesForStateStmt = `select mpc.latitude, mpc.longitude
						from MapPolygonCoordinates as mpc
						join MapPolygon as mp on mpc.MapPolygonID = mp.ID
						where mp.stateID = ?`

	polygonStmt = `select s.stateID, s.state, s.abbr,(
					select COUNT(cl.locationID) from CustomerLocations as cl
					join Customer as c on cl.cust_id = c.cust_id
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					where dt.online = 0 && cl.stateID = s.stateID
				) as count, 
				(select group_concat(mpc.latitude)
				from MapPolygonCoordinates as mpc
				join MapPolygon as mp on mpc.MapPolygonID = mp.ID
				where mp.stateID = s.stateID
				order by mpc.ID) as latitudes,
				(select group_concat(mpc.longitude)
				from MapPolygonCoordinates as mpc
				join MapPolygon as mp on mpc.MapPolygonID = mp.ID
				where mp.stateID = s.stateID
				order by mpc.ID) as longitudes
				from States as s
				where (
					select COUNT(cl.locationID) from CustomerLocations as cl
					join Customer as c on cl.cust_id = c.cust_id
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					where dt.online = 0 && cl.stateID = s.stateID
				) > 0
				order by s.state`
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
	DealerType, DealerTier               string
	SalesRepresentative                  string
	SalesRepresentativeCode              int
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

	qry, err := database.Db.Prepare(customerStmt)
	if err != nil {
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
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealer_type := res.Map("dealer_type")
	dealer_tier := res.Map("dealer_tier")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	sURL, _ := url.Parse(row.Str(search))
	websiteURL, _ := url.Parse(row.Str(site))
	logoURL, _ := url.Parse(row.Str(logo))

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
	c.DealerType = row.Str(dealer_type)
	c.DealerTier = row.Str(dealer_tier)
	c.SalesRepresentative = row.Str(rep_name)
	c.SalesRepresentativeCode = row.Int(rep_code)
	c.MapixCode = row.Str(mpx_code)
	c.MapixDescription = row.Str(mpx_desc)

	ctry := Country{
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	c.State = &State{
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

	qry, err := database.Db.Prepare(customerLocationsStmt)
	if err != nil {
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
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
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
			Country:      row.Str(country),
			Abbreviation: row.Str(country_abbr),
		}

		l.State = &State{
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

	qry, err := database.Db.Prepare(customerUsersStmt)
	if err != nil {
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
	qry, err := database.Db.Prepare(customerPriceStmt)
	if err != nil {
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
	qry, err := database.Db.Prepare(customerPartStmt)
	if err != nil {
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
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealer_type := res.Map("dealer_type")
	dealer_tier := res.Map("dealer_tier")
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

			cust := Customer{
				Id:                      r.Int(customerID),
				Name:                    r.Str(name),
				Email:                   r.Str(email),
				Address:                 r.Str(address),
				Address2:                r.Str(address2),
				City:                    r.Str(city),
				PostalCode:              r.Str(zip),
				Phone:                   r.Str(phone),
				Fax:                     r.Str(fax),
				ContactPerson:           r.Str(contact),
				Latitude:                r.ForceFloat(lat),
				Longitude:               r.ForceFloat(lon),
				Website:                 websiteURL,
				SearchUrl:               sURL,
				Logo:                    logoURL,
				DealerType:              r.Str(dealer_type),
				DealerTier:              r.Str(dealer_tier),
				SalesRepresentative:     r.Str(rep_name),
				SalesRepresentativeCode: r.Int(rep_code),
				MapixCode:               r.Str(mpx_code),
				MapixDescription:        r.Str(mpx_desc),
			}

			ctry := Country{
				Country:      r.Str(country),
				Abbreviation: r.Str(country_abbr),
			}

			cust.State = &State{
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

	return
}

type DealerLocation struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State                                *State
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude, Distance        float64
	Website                              *url.URL
	Parent                               *Customer
	SearchUrl, Logo                      *url.URL
	DealerType, DealerTier               string
	SalesRepresentative                  string
	SalesRepresentativeCode              int
	MapixCode, MapixDescription          string
	Locations                            *[]CustomerLocation
	Users                                []CustomerUser
}

type DealerTier struct {
	Id, Sort int
	Tier     string
}

type DealerType struct {
	MapIcons     *[]MapIcon
	Type, Label  string
	Online, Show bool
}

type MapIcon struct {
	Type, Tier             int
	MapIcon, MapIconShadow *url.URL
}

type MapCoordinates struct {
	Latitude, Longitude float64
}

type StateRegion struct {
	Name, Abbreviation string
	Count              int
	Polygons           *[]MapCoordinates
}

func GetLocalDealers(center string, latlng string) (dealers []DealerLocation, err error) {

	qry, err := database.Db.Prepare(localDealersStmt)
	if err != nil {
		return
	}

	var latlngs []string
	var center_latlngs []string

	clat := CENTER_LATITUDE
	clong := CENTER_LONGITUDE
	swlat := SOUTWEST_LATITUDE
	swlong := SOUTHWEST_LONGITUDE
	swlong2 := SOUTHWEST_LONGITUDE
	nelat := NORTHEAST_LATITUDE
	nelong := NORTHEAST_LONGITUDE
	nelong2 := NORTHEAST_LONGITUDE

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
		EARTH,
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
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealer_type := res.Map("dealer_type")
	dealer_tier := res.Map("dealer_tier")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")

	for _, r := range rows {
		sURL, _ := url.Parse(r.Str(search))
		websiteURL, _ := url.Parse(r.Str(site))
		logoURL, _ := url.Parse(r.Str(logo))

		cust := DealerLocation{
			Id:                      r.Int(customerID),
			Name:                    r.Str(name),
			Email:                   r.Str(email),
			Address:                 r.Str(address),
			City:                    r.Str(city),
			PostalCode:              r.Str(zip),
			Phone:                   r.Str(phone),
			Fax:                     r.Str(fax),
			ContactPerson:           r.Str(contact),
			Latitude:                r.ForceFloat(lat),
			Longitude:               r.ForceFloat(lon),
			Website:                 websiteURL,
			SearchUrl:               sURL,
			Logo:                    logoURL,
			DealerType:              r.Str(dealer_type),
			DealerTier:              r.Str(dealer_tier),
			SalesRepresentative:     r.Str(rep_name),
			SalesRepresentativeCode: r.Int(rep_code),
			MapixCode:               r.Str(mpx_code),
			MapixDescription:        r.Str(mpx_desc),
		}

		cust.Distance = EARTH * (2 * math.Atan2(
			math.Sqrt((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))),
			math.Sqrt(1-((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))))))

		ctry := Country{
			Country:      r.Str(country),
			Abbreviation: r.Str(country_abbr),
		}

		cust.State = &State{
			State:        r.Str(state),
			Abbreviation: r.Str(state_abbr),
			Country:      &ctry,
		}

		dealers = append(dealers, cust)

	}
	sortutil.AscByField(dealers, "Distance")
	return
}

func GetLocalRegions() (regions []StateRegion, err error) {

	polyQuery, err := database.Db.Prepare(polygonStmt)
	if err != nil {
		return
	}

	coordQry, err := database.Db.Prepare(mapCoordinatesForStateStmt)
	if err != nil {
		return
	}

	regions_bytes, _ := redis.RedisClient.Get("local_regions")
	if len(regions_bytes) == 0 {
		_, _, _ = database.Db.Query("SET SESSION group_concat_max_len = 100024")
		rows, res, err := polyQuery.Exec()
		_, _, _ = database.Db.Query("SET SESSION group_concat_max_len = 1024")
		if !database.MysqlError(err) && rows != nil {
			ch := make(chan int)

			for _, row := range rows {
				go func(c chan int, regRow mysql.Row, regRes mysql.Result) {
					state := regRes.Map("state")
					abbr := regRes.Map("abbr")
					count := regRes.Map("count")
					id := regRes.Map("stateID")

					reg := StateRegion{
						Name:         regRow.Str(state),
						Abbreviation: regRow.Str(abbr),
						Count:        regRow.Int(count),
					}
					coordRows, coordRes, err := coordQry.Exec(regRow.Int(id))
					if err == nil {
						lat := coordRes.Map("latitude")
						lon := coordRes.Map("longitude")

						var coords []MapCoordinates
						for _, coordRow := range coordRows {
							coords = append(coords, MapCoordinates{coordRow.ForceFloat(lat), coordRow.ForceFloat(lon)})
						}
						reg.Polygons = &coords
					}

					regions = append(regions, reg)
					c <- 1
				}(ch, row, res)
			}

			for _, _ = range rows {
				<-ch
			}

			if regions_bytes, err = json.Marshal(regions); err == nil {
				redis.RedisClient.Set("local_regions", regions_bytes)
				redis.RedisClient.Expire("local_regions", 86400)
			}
		}

	} else {
		_ = json.Unmarshal(regions_bytes, &regions)
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
	return EARTH * c
}
