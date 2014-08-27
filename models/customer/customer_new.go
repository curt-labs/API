package customer

import (
	"database/sql"
	// "encoding/json"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Customer_New struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State                                geography.State_New
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude                  float64
	Website                              url.URL
	Parent                               Customer
	SearchUrl, Logo                      url.URL
	DealerType                           DealerType_New
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
	Locations                            []CustomerLocation_New
	Users                                []CustomerUser
}

type CustomerLocation_New struct {
	Id                                     int
	Name, Email, Address, City, PostalCode string
	State                                  geography.State_New
	Phone, Fax                             string
	Latitude, Longitude                    float64
	CustomerId                             int
	ContactPerson                          string
	IsPrimary, ShippingDefault             bool
}

type DealerLocation_New struct {
	Id, LocationId                       int
	Name, Email, Address, Address2, City string
	State                                geography.State_New
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude, Distance        float64
	Website                              url.URL
	Parent                               Customer
	SearchUrl, Logo                      url.URL
	DealerType                           DealerType_New
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
}
type DealerType_New struct {
	Id           int
	Type, Label  string
	Online, Show bool
	MapIcon      MapIcon_New
}

type MapIcon_New struct {
	Id, TierId             int
	MapIcon, MapIconShadow url.URL
}

var (
	basics = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				COALESCE(c.latitude,0), COALESCE(c.longitude,0), c.searchURL, c.logo, c.website,
				c.postal_code, COALESCE(s.stateID,0), COALESCE(s.state,""), COALESCE(s.abbr,"") as state_abbr, COALESCE(cty.countryID,0), COALESCE(cty.name,"") as country_name, COALESCE(cty.abbr,"") as country_abbr,
				dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
				dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
				COALESCE(mi.ID,0) as iconID, COALESCE(mi.mapicon,""), COALESCE(mi.mapiconshadow,""),
				COALESCE(mpx.code,"") as mapix_code, COALESCE(mpx.description,"") as mapic_desc,
				COALESCE(sr.name,"") as rep_name, COALESCE(sr.code,"") as rep_code, c.parentID
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where c.customerID = ? ` //TODO - clumsy, shoud use cust_ID, not customerID

	customerLocation = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
							cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
							cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
							s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
							from CustomerLocations as cl
							left join States as s on cl.stateID = s.stateID
							left join Country as cty on s.countryID = cty.countryID
							where cl.cust_id = ?`

	customerUser = `select cu.id, cu.name, cu.email, cu.date_added, cu.active, cu.isSudo from CustomerUser as cu
						join Customer as c on cu.cust_ID = c.cust_id
						where c.customerID = ?
						&& cu.active = 1`
	customerPrice = `select distinct cp.price from ApiKey as ak
						join CustomerUser cu on ak.user_id = cu.id
						join Customer c on cu.cust_ID = c.cust_id
						join CustomerPricing cp on c.customerID = cp.cust_id
						where api_key = ?
						and cp.partID = ?`

	customerPart = `select distinct ci.custPartID from ApiKey as ak
						join CustomerUser cu on ak.user_id = cu.id
						join Customer c on cu.cust_ID = c.cust_id
						join CartIntegration ci on c.customerID = ci.custID
						where ak.api_key = ?
						and ci.partID = ?`
	etailers = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
					COALESCE(c.latitude,0), COALESCE(c.longitude,0), c.searchURL, c.logo, c.website,
					c.postal_code, COALESCE(s.stateID,0), COALESCE(s.state,""), COALESCE(s.abbr,"") as state_abbr, COALESCE(cty.countryID,0), COALESCE(cty.name,"") as country_name, COALESCE(cty.abbr,"") as country_abbr,
					dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
					dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
					COALESCE(mi.ID,0) as iconID, COALESCE(mi.mapicon,""), COALESCE(mi.mapiconshadow,""),
					COALESCE(mpx.code,"") as mapix_code, COALESCE(mpx.description,"") as mapic_desc,
					COALESCE(sr.name,"") as rep_name, COALESCE(sr.code,"") as rep_code, c.parentID
					from Customer as c
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					join DealerTiers dtr on c.tier = dtr.ID
					left join MapIcons as mi on dt.dealer_type = mi.dealer_type
					left join States as s on c.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
					left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
					where dt.online = 1 && c.isDummy = 0`
	localDealers = `select cl.locationID, c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
						cl.latitude, cl.longitude, c.searchURL, c.logo, c.website,
						cl.postalCode, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
						dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
						dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
						mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
						mpx.code as mapix_code, mpx.description as mapic_desc,
						sr.name as rep_name, sr.code as rep_code, c.parentID
						from CustomerLocations as cl
						join Customer as c on cl.cust_id = c.cust_id
						join DealerTypes as dt on c.dealer_type = dt.dealer_type
						left join MapIcons as mi on dt.dealer_type = mi.dealer_type
						join DealerTiers as dtr on c.tier = dtr.ID
						left join States as s on cl.stateID = s.stateID
						left join Country as cty on s.countryID = cty.countryID
						left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
						left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
						where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier &&
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
						group by cl.locationID`
	polygon = `select s.stateID, s.state, s.abbr,
					(
						select COUNT(cl.locationID) from CustomerLocations as cl
						join Customer as c on cl.cust_id = c.cust_id
						join DealerTypes as dt on c.dealer_type = dt.dealer_type
						where dt.online = 0 && cl.stateID = s.stateID
					) as count
					from States as s
					where (
						select COUNT(cl.locationID) from CustomerLocations as cl
						join Customer as c on cl.cust_id = c.cust_id
						join DealerTypes as dt on c.dealer_type = dt.dealer_type
						where dt.online = 0 && cl.stateID = s.stateID
					) > 0
					order by s.state`
	MapPolygonCoordinatesForState = `select mp.ID, mpc.latitude,mpc.longitude
										from MapPolygonCoordinates as mpc
										join MapPolygon as mp on mpc.MapPolygonID = mp.ID
										where mp.stateID = ?
										`
	//TODO delete ld before publishing - for reference
	ld = `select cl.locationID, c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
						COALESCE(cl.latitude,0), COALESCE(cl.longitude,0), c.searchURL, c.logo, c.website,
						COALESCE(c.postal_code,0), COALESCE(s.stateID,0), COALESCE(s.state,""), COALESCE(s.abbr,"") as state_abbr, COALESCE(cty.countryID,0),COALESCE(cty.name,"") as country_name, COALESCE(cty.abbr,"") as country_abbr,
						dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
						dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
						COALESCE(mi.ID,0) as iconID, COALESCE(mi.mapicon,""), COALESCE(mi.mapiconshadow,""),
						COALESCE(mpx.code,"") as mapix_code, COALESCE(mpx.description,"") as mapic_desc,
						COALESCE(sr.name,"") as rep_name, COALESCE(sr.code,"") as rep_code, c.parentID
						from CustomerLocations as cl
						join Customer as c on cl.cust_id = c.cust_id
						join DealerTypes as dt on c.dealer_type = dt.dealer_type
						left join MapIcons as mi on dt.dealer_type = mi.dealer_type
						join DealerTiers as dtr on c.tier = dtr.ID
						left join States as s on cl.stateID = s.stateID
						left join Country as cty on s.countryID = cty.countryID
						left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
						left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
						where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier &&
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
						group by cl.locationID
						order by dtr.sort desc`
)

func (c *Customer_New) GetCustomer_New() (err error) {

	locationChan := make(chan int)
	go func() {
		if locErr := c.GetLocations_New(); locErr != nil {
			err = locErr
		}
		locationChan <- 1
	}()

	err = c.Basics_New()

	<-locationChan

	return err
}

//TODO - I hate the hacks in this scan/query!!!
func (c *Customer_New) Basics_New() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(basics)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var ur, logo, web, icon, shadow, lat, lon string
	err = stmt.QueryRow(c.Id).Scan(
		&c.Id,            //c.customerID,
		&c.Name,          //c.name
		&c.Email,         //c.email
		&c.Address,       //c.address
		&c.Address2,      //c.address2
		&c.City,          //c.city,
		&c.Phone,         //phone,
		&c.Fax,           //c.fax
		&c.ContactPerson, //c.contact_person,
		&lat,             //c.latitude
		&lon,             //c.longitude
		&ur,
		&logo,
		&web,
		&c.PostalCode,                 //c.postal_code
		&c.State.Id,                   //s.stateID
		&c.State.State,                //s.state
		&c.State.Abbreviation,         //s.abbr as state_abbr
		&c.State.Country.Id,           //cty.countryID,
		&c.State.Country.Country,      //cty.name as country_name
		&c.State.Country.Abbreviation, //cty.abbr as country_abbr,
		&c.DealerType.Id,              //dt.dealer_type as typeID
		&c.DealerType.Type,            // dt.type as dealerType
		&c.DealerType.Online,          // dt.online as typeOnline,
		&c.DealerType.Show,            //dt.show as typeShow
		&c.DealerType.Label,           //dt.label as typeLabel,
		&c.DealerTier.Id,              //dtr.ID as tierID,
		&c.DealerTier.Tier,            //dtr.tier as tier
		&c.DealerTier.Sort,            //dtr.sort as tierSort
		&c.DealerType.MapIcon.Id,
		&icon,
		&shadow,                    //mi.ID as iconID
		&c.MapixCode,               //mpx.code as mapix_code
		&c.MapixDescription,        //mpx.description as mapic_desc,
		&c.SalesRepresentative,     //sr.name as rep_name
		&c.SalesRepresentativeCode, // sr.code as rep_code,
		&c.Parent.Id,               //c.parentID
	)
	if err != nil {
		return err
	}
	c.Latitude, err = floatParse(lat)
	c.Longitude, err = floatParse(lon)
	c.SearchUrl, err = urlParse(ur)
	c.Logo, err = urlParse(logo)
	c.Website, err = urlParse(web)
	c.DealerType.MapIcon.MapIcon, err = urlParse(icon)
	c.DealerType.MapIcon.MapIconShadow, err = urlParse(shadow)
	if err != nil {
		return err
	}

	if c.Parent.Id != 0 {
		parent := Customer{
			Id: c.Parent.Id,
		}
		if err = parent.GetCustomer(); err == nil {
			c.Parent = parent
		}
	}
	return nil
}

//helper funcs
func floatParse(input string) (float64, error) {
	var err error
	if input != "" {
		output, err := strconv.ParseFloat(input, 64)
		return output, err
	}
	return 0.0, err
}

func urlParse(input string) (url.URL, error) {
	var err error
	if input != "" {
		output, err := url.Parse(input)
		output2 := *output
		return output2, err
	}
	output, err := url.Parse("")
	output2 := *output
	return output2, err
}

func (c *Customer_New) GetLocations_New() error {
	var err error
	var ls []CustomerLocation_New
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerLocation)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(c.Id)
	for res.Next() {
		var l CustomerLocation_New
		err = res.Scan(
			&l.Id,
			&l.Name,
			&l.Email,
			&l.Address,
			&l.City,
			&l.PostalCode,
			&l.Phone,
			&l.Fax,
			&l.Latitude,
			&l.Longitude,
			&l.CustomerId,
			&l.ContactPerson,
			&l.IsPrimary,
			&l.ShippingDefault,
			&l.State.Id,
			&l.State.State,
			&l.State.Abbreviation,
			&l.State.Country.Id,
			&l.State.Country.Country,
			&l.State.Country.Abbreviation,
		)
		if err != nil {
			return err
		}
		ls = append(ls, l)
	}
	c.Locations = ls
	return nil
}

func (c *Customer_New) GetUsers_New() (users []CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return users, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUser)
	if err != nil {
		return users, err
	}
	defer stmt.Close()

	res, err := stmt.Query(c.Id)
	for res.Next() {
		var u CustomerUser
		res.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.DateAdded,
			&u.Active,
			&u.Sudo,
		)
		users = append(users, u)
	}
	if err != nil {
		return users, err
	}
	return users, err
}

func GetCustomerPrice_New(api_key string, part_id int) (price float64, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return price, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPrice)
	if err != nil {
		return price, err
	}

	err = stmt.QueryRow(api_key, part_id).Scan(&price)
	return price, err
}

func GetCustomerCartReference_New(api_key string, part_id int) (ref int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ref, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPart)
	if err != nil {
		return ref, err
	}

	err = stmt.QueryRow(api_key, part_id).Scan(&ref)
	return ref, err
}

/* Internal Use Only */
func GetEtailers_New() (dealers []Customer_New, err error) {

	// redis_key := "goapi:dealers:etailers"
	// data, err := redis.Get(redis_key)
	// if len(data) > 0 && err != nil {
	// 	err = json.Unmarshal(data, &dealers)
	// 	if err == nil {
	// 		return
	// 	}
	// }
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dealers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(etailers)
	if err != nil {
		return dealers, err
	}
	var ur, logo, web, icon, shadow, lat, lon string
	res, err := stmt.Query()
	for res.Next() {
		var c Customer_New
		err = res.Scan(
			&c.Id,            //c.customerID,
			&c.Name,          //c.name
			&c.Email,         //c.email
			&c.Address,       //c.address
			&c.Address2,      //c.address2
			&c.City,          //c.city,
			&c.Phone,         //phone,
			&c.Fax,           //c.fax
			&c.ContactPerson, //c.contact_person,
			&lat,
			&lon,
			&ur,
			&logo,
			&web,
			&c.PostalCode,                 //c.postal_code
			&c.State.Id,                   //s.stateID
			&c.State.State,                //s.state
			&c.State.Abbreviation,         //s.abbr as state_abbr
			&c.State.Country.Id,           //cty.countryID,
			&c.State.Country.Country,      //cty.name as country_name
			&c.State.Country.Abbreviation, //cty.abbr as country_abbr,
			&c.DealerType.Id,              //dt.dealer_type as typeID
			&c.DealerType.Type,            // dt.type as dealerType
			&c.DealerType.Online,          // dt.online as typeOnline,
			&c.DealerType.Show,            //dt.show as typeShow
			&c.DealerType.Label,           //dt.label as typeLabel,
			&c.DealerTier.Id,              //dtr.ID as tierID,
			&c.DealerTier.Tier,            //dtr.tier as tier
			&c.DealerTier.Sort,            //dtr.sort as tierSort
			&c.DealerType.MapIcon.Id,
			&icon,
			&shadow,
			&c.MapixCode,               //mpx.code as mapix_code
			&c.MapixDescription,        //mpx.description as mapic_desc,
			&c.SalesRepresentative,     //sr.name as rep_name
			&c.SalesRepresentativeCode, // sr.code as rep_code,
			&c.Parent.Id,               //c.parentID
		)
		if err != nil {
			return dealers, err
		}

		c.Latitude, err = floatParse(lat)
		c.Longitude, err = floatParse(lon)
		c.SearchUrl, err = urlParse(ur)
		c.Logo, err = urlParse(logo)
		c.Website, err = urlParse(web)
		c.DealerType.MapIcon.MapIcon, err = urlParse(icon)
		c.DealerType.MapIcon.MapIconShadow, err = urlParse(shadow)
		if err != nil {
			return dealers, err
		}

		dealers = append(dealers, c)
	}
	// go redis.Setex(redis_key, dealers, 86400)
	return dealers, err
}

func GetLocalDealers_New(center string, latlng string) (dealers []DealerLocation_New, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dealers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(localDealers)
	if err != nil {
		return dealers, err
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

	params := []interface{}{ //all are float64 type
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
		nelong2}
	res, err := stmt.Query(params...)
	if err != nil {
		return dealers, err
	}

	var ur, logo, web, lat, lon, icon, shadow string

	for res.Next() {
		var cust DealerLocation_New
		res.Scan(
			&cust.LocationId,
			&cust.Id,
			&cust.Name,
			&cust.Email,
			&cust.Address,
			&cust.City,
			&cust.Phone,
			&cust.Fax,
			&cust.ContactPerson,
			&lat,
			&lon,
			&ur,
			&logo,
			&web,
			&cust.PostalCode,
			&cust.State.Id,
			&cust.State.State,
			&cust.State.Abbreviation,
			&cust.State.Country.Id,
			&cust.State.Country.Country,
			&cust.State.Country.Abbreviation,
			&cust.DealerType.Id,
			&cust.DealerType.Type,
			&cust.DealerType.Online,
			&cust.DealerType.Show,
			&cust.DealerType.Label,
			&cust.DealerTier.Id,
			&cust.DealerTier.Tier,
			&cust.DealerTier.Sort,
			&cust.DealerType.MapIcon.Id,
			&icon,
			&shadow,
			&cust.MapixCode,
			&cust.MapixDescription,
			&cust.SalesRepresentative,
			&cust.SalesRepresentativeCode,
			&cust.Parent.Id,
		)
		cust.Latitude, err = floatParse(lat)
		cust.Longitude, err = floatParse(lon)
		cust.SearchUrl, err = urlParse(ur)
		cust.Logo, err = urlParse(logo)
		cust.Website, err = urlParse(web)
		cust.DealerType.MapIcon.MapIcon, err = urlParse(icon)
		cust.DealerType.MapIcon.MapIconShadow, err = urlParse(shadow)
		if err != nil {
			return dealers, err
		}

		cust.Distance = api_helpers.EARTH * (2 * math.Atan2(
			math.Sqrt((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))),
			math.Sqrt(1-((math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((cust.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((cust.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(cust.Latitude*(math.Pi/180))))))

		dealers = append(dealers, cust)
	}
	sortutil.AscByField(dealers, "Distance")
	return
}

//HERE

func GetLocalRegions_New() (regions []StateRegion, err error) {

	// redis_key := "goapi:local:regions"
	// data, err := redis.Get(redis_key)
	// if len(data) > 0 && err != nil {
	// 	err = json.Unmarshal(data, &regions)
	// 	if err == nil {
	// 		return
	// 	}
	// }
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return regions, err
	}
	defer db.Close()

	stmtPolygon, err := db.Prepare(polygon)
	if err != nil {
		return regions, err
	}
	stmtCoordinates, err := db.Prepare(MapPolygonCoordinatesForState)
	if err != nil {
		return regions, err
	}

	_, err = db.Exec("SET SESSION group_concat_max_len = 100024")
	res, err := stmtPolygon.Query()
	_, err = db.Exec("SET SESSION group_concat_max_len = 1024")

	for res.Next() {
		var reg StateRegion
		res.Scan(
			&reg.Id,
			&reg.Name,
			&reg.Abbreviation,
			&reg.Count,
		)
		// log.Print("REG", reg)
		//coor query
		coorRes, err := stmtCoordinates.Query(reg.Id)
		if err != nil {
			return regions, err
		}
		polygons := make(map[int]MapPolygon, 0)
		coordRows := make(map[int]GeoLocation)
		for coorRes.Next() {
			var tempMap MapPolygon
			var tempGeo GeoLocation
			err = coorRes.Scan(
				&tempMap.Id,
				&tempGeo.Latitude,
				&tempGeo.Longitude,
			)
			coordRows[tempMap.Id] = tempGeo
			for id, _ := range coordRows {
				// Check if we have an index for this polygon created
				if _, ok := polygons[id]; !ok {
					// First time hitting this polygon
					// so we'll create one
					polygons[id] = MapPolygon{
						Id:          tempMap.Id,
						Coordinates: make([]GeoLocation, 0),
					}
				}

				// Add the GeoLocartion info to our polygon
				poly := polygons[tempMap.Id]
				poly.Coordinates = append(poly.Coordinates, GeoLocation{tempGeo.Latitude, tempGeo.Longitude})
				polygons[tempMap.Id] = poly
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
	}
	// go redis.Set(redis_key, regions)
	return
}

//no db - same
func getViewPortWidth_New(lat1 float64, lon1 float64, lat2 float64, long2 float64) float64 {
	dlat := (lat2 - lat1) * (math.Pi / 180)
	dlon := (long2 - lon1) * (math.Pi / 180)

	lat1 = lat1 * (math.Pi / 180)
	lat2 = lat2 * (math.Pi / 180)

	a := (math.Sin(dlat/2) * math.Sin(dlat/2)) + ((math.Sin(dlon/2))*(math.Sin(dlon/2)))*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return api_helpers.EARTH * c
}
