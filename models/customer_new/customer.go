package customer_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Customer struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State                                geography.State_New
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude                  float64
	Website                              url.URL
	Parent                               *Customer
	SearchUrl, Logo                      url.URL
	DealerType                           DealerType
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
	Locations                            []CustomerLocation
	Users                                []CustomerUser
}

type CustomerLocation struct {
	Id                                     int
	Name, Email, Address, City, PostalCode string
	State                                  geography.State_New
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
	MapIcon, MapIconShadow url.URL
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
	State                                geography.State_New
	PostalCode                           string
	Phone, Fax                           string
	ContactPerson                        string
	Latitude, Longitude, Distance        float64
	Website                              url.URL
	Parent                               Customer
	SearchUrl, Logo                      url.URL
	DealerType                           DealerType
	DealerTier                           DealerTier
	SalesRepresentative                  string
	SalesRepresentativeCode              string
	MapixCode, MapixDescription          string
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
	localDealerTiers = `select distinct dtr.* from DealerTiers as dtr
							join Customer as c on dtr.ID = c.tier
							join DealerTypes as dt on c.dealer_type = dt.dealer_type
							where dt.online = false and dt.show = true
							order by dtr.sort`
	localDealerTypes = `select m.ID as iconId, m.mapicon, m.mapiconshadow,
							dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
							dt.dealer_type as dealerTypeId, dt.type as dealerType, dt.online, dt.show, dt.label
							from MapIcons as m
							join DealerTypes as dt on m.dealer_type = dt.dealer_type
							join DealerTiers as dtr on m.tier = dtr.ID
							where dt.show = true
							order by dtr.sort desc`
	whereToBuyDealers = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
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
							where dt.dealer_type = 1 and dtr.ID = 4 and c.isDummy = false and length(c.searchURL) > 1`

	customerByLocation = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
							dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
							cl.locationID, cl.name, cl.address,cl.city,
							cl.postalCode, cl.email, cl.phone,cl.fax,
							cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
							c.showWebsite, c.website, c.eLocalURL
							from CustomerLocations as cl
							join States as cls on cl.stateID = cls.stateID
							join Customer as c on cl.cust_id = c.cust_id
							join DealerTypes as dt on c.dealer_type = dt.dealer_type
							join DealerTiers as dtr on c.tier = dtr.ID
							where cl.locationID = ? limit 1`
	searchDealerLocations = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
								dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
								cl.locationID, cl.name, cl.address,cl.city,
								cl.postalCode, cl.email, cl.phone,cl.fax,
								cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
								c.showWebsite, c.website, c.eLocalURL
								from CustomerLocations as cl
								join States as cls on cl.stateID = cls.stateID
								join Customer as c on cl.cust_id = c.cust_id
								join DealerTypes as dt on c.dealer_type = dt.dealer_type
								join DealerTiers as dtr on c.tier = dtr.ID
								where (dt.dealer_type = 2 or dt.dealer_type = 3) and c.isDummy = false
								and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	dealerLocationsByType = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
								dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
								cl.locationID, cl.name, cl.address,cl.city,
								cl.postalCode, cl.email, cl.phone,cl.fax,
								cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
								c.showWebsite, c.website, c.eLocalURL
								from CustomerLocations as cl
								join States as cls on cl.stateID = cls.stateID
								join Customer as c on cl.cust_id = c.cust_id
								join DealerTypes as dt on c.dealer_type = dt.dealer_type
								join DealerTiers as dtr on c.tier = dtr.ID
								where dt.online = false and c.isDummy = false
								and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	searchDealerLocationsByLatLng = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
										dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
										cl.locationID, cl.name, cl.address,cl.city,
										cl.postalCode, cl.email, cl.phone,cl.fax,
										cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
										c.showWebsite, c.website, c.eLocalURL
										from CustomerLocations as cl
										join States as cls on cl.stateID = cls.stateID
										join Customer as c on cl.cust_id = c.cust_id
										join DealerTypes as dt on c.dealer_type = dt.dealer_type
										join DealerTiers as dtr on c.tier = dtr.ID
										where dt.online = false and c.isDummy = false
										and dt.show = true and
										( ? * (
											2 * ATAN2(
												SQRT((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))),
												SQRT(1 - ((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))))
											)
										) < 100.0)`
)

func (c *Customer) GetCustomer() (err error) {

	locationChan := make(chan int)
	basicsChan := make(chan int)

	go func() {
		if locErr := c.GetLocations(); locErr != nil {
			err = locErr
		}
		locationChan <- 1
	}()
	go func() {
		if basErr := c.Basics(); basErr != nil {
			err = basErr
		}
		basicsChan <- 1
	}()

	<-locationChan
	<-basicsChan

	return err
}

//TODO - I hate the hacks in this scan/query!!!
func (c *Customer) Basics() error {
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

	var name, email, address, address2, city, phone, fax, contactPerson []byte
	var logo, web, lat, lon, url, icon, shadow, mapIconId []byte
	var stateId, state, stateAbbr, countryId, country, countryAbbr, parentId, postalCode, mapixCode, mapixDesc, rep, repCode []byte
	err = stmt.QueryRow(c.Id).Scan(
		&c.Id,          //c.customerID,
		&name,          //c.name
		&email,         //c.email
		&address,       //c.address
		&address2,      //c.address2
		&city,          //c.city,
		&phone,         //phone,
		&fax,           //c.fax
		&contactPerson, //c.contact_person,
		&lat,           //c.latitude
		&lon,           //c.longitude
		&url,
		&logo,
		&web,
		&postalCode,          //c.postal_code
		&stateId,             //s.stateID
		&state,               //s.state
		&stateAbbr,           //s.abbr as state_abbr
		&countryId,           //cty.countryID,
		&country,             //cty.name as country_name
		&countryAbbr,         //cty.abbr as country_abbr,
		&c.DealerType.Id,     //dt.dealer_type as typeID
		&c.DealerType.Type,   // dt.type as dealerType
		&c.DealerType.Online, // dt.online as typeOnline,
		&c.DealerType.Show,   //dt.show as typeShow
		&c.DealerType.Label,  //dt.label as typeLabel,
		&c.DealerTier.Id,     //dtr.ID as tierID,
		&c.DealerTier.Tier,   //dtr.tier as tier
		&c.DealerTier.Sort,   //dtr.sort as tierSort
		&mapIconId,
		&icon,
		&shadow,    //mi.ID as iconID
		&mapixCode, //mpx.code as mapix_code
		&mapixDesc, //mpx.description as mapic_desc,
		&rep,       //sr.name as rep_name
		&repCode,   // sr.code as rep_code,
		&parentId,  //c.parentID
	)
	if err != nil {
		return err
	}
	c.Name, err = conversions.ByteToString(name)
	c.Address, err = conversions.ByteToString(address)
	c.City, err = conversions.ByteToString(city)
	c.Email, err = conversions.ByteToString(email)
	c.Phone, err = conversions.ByteToString(phone)
	c.Fax, err = conversions.ByteToString(fax)
	c.ContactPerson, err = conversions.ByteToString(contactPerson)

	c.Latitude, err = conversions.ByteToFloat(lat)
	c.Longitude, err = conversions.ByteToFloat(lon)
	c.SearchUrl, err = conversions.ByteToUrl(url)
	c.Logo, err = conversions.ByteToUrl(logo)
	c.Website, err = conversions.ByteToUrl(web)
	c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
	c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)

	c.PostalCode, err = conversions.ByteToString(postalCode)
	c.State.Id, err = conversions.ByteToInt(stateId)
	c.State.State, err = conversions.ByteToString(state)
	c.State.Abbreviation, err = conversions.ByteToString(stateAbbr)
	c.State.Country.Id, err = conversions.ByteToInt(countryId)
	c.State.Country.Country, err = conversions.ByteToString(country)
	c.State.Country.Abbreviation, err = conversions.ByteToString(countryAbbr)
	c.DealerType.MapIcon.Id, err = conversions.ByteToInt(mapIconId)
	c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
	c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)
	c.MapixCode, err = conversions.ByteToString(mapixCode)
	c.MapixDescription, err = conversions.ByteToString(mapixDesc)
	c.SalesRepresentative, err = conversions.ByteToString(rep)
	c.SalesRepresentativeCode, err = conversions.ByteToString(repCode)

	parentInt, err := conversions.ByteToInt(parentId)
	if err != nil {
		return err
	}
	if parentInt != 0 {
		par := Customer{Id: parentInt}
		par.GetCustomer()
		c.Parent = &par
	}

	return nil
}

func (c *Customer) GetLocations() error {
	var err error
	c.Locations = make([]CustomerLocation, 0)
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
	var stateId, state, stateAbbr, countryId, country, countryAbbr, postalCode []byte
	var name, email, address, city, phone, fax, contactPerson []byte
	res, err := stmt.Query(c.Id)
	for res.Next() {
		var l CustomerLocation
		err = res.Scan(
			&l.Id,
			&name,       //c.name
			&email,      //c.email
			&address,    //c.address
			&city,       //c.city,
			&postalCode, //c.postal_code
			&phone,      //phone,
			&fax,        //c.fax
			&l.Latitude,
			&l.Longitude,
			&l.CustomerId,
			&contactPerson, //c.contact_person,
			&l.IsPrimary,
			&l.ShippingDefault,
			&stateId,     //s.stateID
			&state,       //s.state
			&stateAbbr,   //s.abbr as state_abbr
			&countryId,   //cty.countryID,
			&country,     //cty.name as country_name
			&countryAbbr, //cty.abbr as country_ab

		)
		if err != nil {
			return err
		}
		c.Locations = append(c.Locations, l)
	}

	return nil
}

func (c *Customer) GetUsers() (users []CustomerUser, err error) {
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
	var name []byte
	for res.Next() {
		var u CustomerUser
		err = res.Scan(
			&u.Id,
			&name,
			&u.Email,
			&u.DateAdded,
			&u.Active,
			&u.Sudo,
		)
		if err != nil {
			return users, err
		}
		u.Name, err = conversions.ByteToString(name)
		users = append(users, u)
	}
	if err != nil {
		return users, err
	}
	return users, err
}

func GetCustomerPrice(api_key string, part_id int) (price float64, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return price, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerPrice)
	if err != nil {
		return price, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(api_key, part_id).Scan(&price)
	return price, err
}

func GetCustomerCartReference(api_key string, part_id int) (ref int, err error) {
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
func GetEtailers() (dealers []Customer, err error) {

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
	var name, email, address, address2, city, phone, fax, contactPerson []byte
	var url, logo, web, icon, shadow, lat, lon []byte
	var stateId, state, stateAbbr, countryId, country, countryAbbr, parentId, postalCode, mapixCode, mapixDesc, rep, repCode []byte
	res, err := stmt.Query()
	for res.Next() {
		var c Customer
		err = res.Scan(
			&c.Id,          //c.customerID,
			&name,          //c.name
			&email,         //c.email
			&address,       //c.address
			&address2,      //c.address2
			&city,          //c.city,
			&phone,         //phone,
			&fax,           //c.fax
			&contactPerson, //c.contact_person,
			&lat,           //c.latitude
			&lon,           //c.longitude
			&url,
			&logo,
			&web,
			&postalCode,          //c.postal_code
			&stateId,             //s.stateID
			&state,               //s.state
			&stateAbbr,           //s.abbr as state_abbr
			&countryId,           //cty.countryID,
			&country,             //cty.name as country_name
			&countryAbbr,         //cty.abbr as country_abbr,
			&c.DealerType.Id,     //dt.dealer_type as typeID
			&c.DealerType.Type,   // dt.type as dealerType
			&c.DealerType.Online, // dt.online as typeOnline,
			&c.DealerType.Show,   //dt.show as typeShow
			&c.DealerType.Label,  //dt.label as typeLabel,
			&c.DealerTier.Id,     //dtr.ID as tierID,
			&c.DealerTier.Tier,   //dtr.tier as tier
			&c.DealerTier.Sort,   //dtr.sort as tierSort
			&c.DealerType.MapIcon.Id,
			&icon,
			&shadow,
			&mapixCode, //mpx.code as mapix_code
			&mapixDesc, //mpx.description as mapic_desc,
			&rep,       //sr.name as rep_name
			&repCode,   // sr.code as rep_code,
			&parentId,  //c.parentID
		)
		if err != nil {
			return dealers, err
		}

		c.Name, err = conversions.ByteToString(name)
		c.Address, err = conversions.ByteToString(address)
		c.Address2, err = conversions.ByteToString(address2)
		c.City, err = conversions.ByteToString(city)
		c.Email, err = conversions.ByteToString(email)
		c.Phone, err = conversions.ByteToString(phone)
		c.Fax, err = conversions.ByteToString(fax)
		c.ContactPerson, err = conversions.ByteToString(contactPerson)

		c.Latitude, err = conversions.ByteToFloat(lat)
		c.Longitude, err = conversions.ByteToFloat(lon)
		c.SearchUrl, err = conversions.ByteToUrl(url)
		c.Logo, err = conversions.ByteToUrl(logo)
		c.Website, err = conversions.ByteToUrl(web)
		c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
		c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)
		c.MapixCode, err = conversions.ByteToString(mapixCode)
		c.MapixDescription, err = conversions.ByteToString(mapixDesc)
		c.SalesRepresentative, err = conversions.ByteToString(rep)
		c.SalesRepresentativeCode, err = conversions.ByteToString(repCode)
		if err != nil {
			return dealers, err
		}
		err = c.GetLocations()
		parentInt, err := conversions.ByteToInt(parentId)
		if err != nil {
			return dealers, err
		}
		if parentInt != 0 {
			par := Customer{Id: parentInt}
			par.GetCustomer()
			c.Parent = &par
		}

		dealers = append(dealers, c)
	}
	// go redis.Setex(redis_key, dealers, 86400)
	return dealers, err
}

func GetLocalDealers(center string, latlng string) (dealers []DealerLocation, err error) {

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

	var ur, logo, web, lat, lon, icon, shadow, mapixCode, mapixDesc, rep, repCode []byte
	var name, email, address, city, phone, fax, contactPerson []byte
	for res.Next() {
		var cust DealerLocation
		res.Scan(

			&cust.LocationId,
			&cust.Id,
			&name,          //c.name
			&email,         //c.email
			&address,       //c.address
			&city,          //c.city,
			&phone,         //phone,
			&fax,           //c.fax
			&contactPerson, //c.contact_person,
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
			&mapixCode, //mpx.code as mapix_code
			&mapixDesc, //mpx.description as mapic_desc,
			&rep,       //sr.name as rep_name
			&repCode,   // sr.code as rep_code,
			&cust.Parent.Id,
		)

		cust.Name, err = conversions.ByteToString(name)
		cust.Address, err = conversions.ByteToString(address)
		cust.City, err = conversions.ByteToString(city)
		cust.Email, err = conversions.ByteToString(email)
		cust.Phone, err = conversions.ByteToString(phone)
		cust.Fax, err = conversions.ByteToString(fax)
		cust.ContactPerson, err = conversions.ByteToString(contactPerson)

		cust.Latitude, err = conversions.ByteToFloat(lat)
		cust.Longitude, err = conversions.ByteToFloat(lon)
		cust.SearchUrl, err = conversions.ByteToUrl(ur)
		cust.Logo, err = conversions.ByteToUrl(logo)
		cust.Website, err = conversions.ByteToUrl(web)
		cust.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
		cust.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)

		cust.MapixCode, err = conversions.ByteToString(mapixCode)
		cust.MapixDescription, err = conversions.ByteToString(mapixDesc)
		cust.SalesRepresentative, err = conversions.ByteToString(rep)
		cust.SalesRepresentativeCode, err = conversions.ByteToString(repCode)
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

func GetLocalRegions() (regions []StateRegion, err error) {

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

func GetLocalDealerTiers() (tiers []DealerTier, err error) {
	// redis_key := "goapi:local:tiers"
	// data, err := redis.Get(redis_key)
	// if len(data) > 0 && err != nil {
	// 	err = json.Unmarshal(data, &tiers)
	// 	if err == nil {
	// 		return
	// 	}
	// }

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return tiers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(localDealerTiers)
	if err != nil {
		return tiers, err
	}
	res, err := stmt.Query()
	for res.Next() {
		var t DealerTier
		err = res.Scan(&t.Id, &t.Tier, &t.Sort)
		if err != nil {
			return tiers, err
		}
		tiers = append(tiers, t)
	}
	// go redis.Setex(redis_key, tiers, 86400)
	return
}

func GetLocalDealerTypes() (graphics []MapGraphics, err error) {
	// redis_key := "goapi:local:types"
	// data, err := redis.Get(redis_key)
	// if len(data) > 0 && err != nil {
	// 	err = json.Unmarshal(data, &graphics)
	// 	if err == nil {
	// 		return
	// 	}
	// }
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return graphics, err
	}
	defer db.Close()

	stmt, err := db.Prepare(localDealerTypes)
	if err != nil {
		return graphics, err
	}
	res, err := stmt.Query()
	var icon, shadow []byte
	for res.Next() {
		var g MapGraphics
		err = res.Scan(
			&g.MapIcon.Id,
			&icon,
			&shadow,
			&g.DealerTier.Id,
			&g.DealerTier.Tier,
			&g.DealerTier.Sort,
			&g.DealerType.Id,
			&g.DealerType.Type,
			&g.DealerType.Online,
			&g.DealerType.Show,
			&g.DealerType.Label,
		)
		if err != nil {
			return graphics, err
		}
		g.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
		g.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)
		graphics = append(graphics, g)
	}
	// go redis.Setex(redis_key, graphics, 86400)
	return
}

func GetWhereToBuyDealers() (customers []Customer, err error) {
	// redis_key := "goapi:dealers:wheretobuy"
	// data, err := redis.Get(redis_key)
	// if len(data) > 0 && err != nil {
	// 	err = json.Unmarshal(data, &customers)
	// 	if err == nil {
	// 		return
	// 	}
	// }

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return customers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(whereToBuyDealers)
	if err != nil {
		return customers, err
	}
	var ur, logo, web, icon, shadow, lat, lon, postalCode, mapixCode, mapixDesc, rep, repCode, parentId []byte
	var stateId, state, stateAbbr, countryId, country, countryAbbr, mapIconId []byte
	var name, email, address, address2, city, phone, fax, contactPerson []byte

	res, err := stmt.Query()
	if err != nil {
		return customers, err
	}
	for res.Next() {
		var c Customer
		err = res.Scan(
			&c.Id,    //c.customerID,
			&name,    //c.name
			&email,   //c.email
			&address, //c.address
			&address2,
			&city,          //c.city,
			&phone,         //phone,
			&fax,           //c.fax
			&contactPerson, //c.contact_person
			&lat,
			&lon,
			&ur,
			&logo,
			&web,
			&postalCode,          //c.postal_code
			&stateId,             //s.stateID
			&state,               //s.state
			&stateAbbr,           //s.abbr as state_abbr
			&countryId,           //cty.countryID,
			&country,             //cty.name as country_name
			&countryAbbr,         //cty.abbr as country_abbr,
			&c.DealerType.Id,     //dt.dealer_type as typeID
			&c.DealerType.Type,   // dt.type as dealerType
			&c.DealerType.Online, // dt.online as typeOnline,
			&c.DealerType.Show,   //dt.show as typeShow
			&c.DealerType.Label,  //dt.label as typeLabel,
			&c.DealerTier.Id,     //dtr.ID as tierID,
			&c.DealerTier.Tier,   //dtr.tier as tier
			&c.DealerTier.Sort,   //dtr.sort as tierSort
			&mapIconId,
			&icon,
			&shadow,
			&mapixCode, //mpx.code as mapix_code
			&mapixDesc, //mpx.description as mapic_desc,
			&rep,       //sr.name as rep_name
			&repCode,   // sr.code as rep_code,
			&parentId,  //c.parentID
		)
		if err != nil {
			return customers, err
		}

		c.Name, err = conversions.ByteToString(name)
		c.Address, err = conversions.ByteToString(address)
		c.City, err = conversions.ByteToString(city)
		c.Email, err = conversions.ByteToString(email)
		c.Phone, err = conversions.ByteToString(phone)
		c.Fax, err = conversions.ByteToString(fax)
		c.ContactPerson, err = conversions.ByteToString(contactPerson)

		c.Latitude, err = conversions.ByteToFloat(lat)
		c.Longitude, err = conversions.ByteToFloat(lon)
		c.SearchUrl, err = conversions.ByteToUrl(ur)
		c.Logo, err = conversions.ByteToUrl(logo)
		c.Website, err = conversions.ByteToUrl(web)
		c.PostalCode, err = conversions.ByteToString(postalCode)
		c.State.Id, err = conversions.ByteToInt(stateId)
		c.State.State, err = conversions.ByteToString(state)
		c.State.Abbreviation, err = conversions.ByteToString(stateAbbr)
		c.DealerType.MapIcon.Id, err = conversions.ByteToInt(mapIconId)
		c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(icon)
		c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(shadow)
		c.MapixCode, err = conversions.ByteToString(mapixCode)
		c.MapixDescription, err = conversions.ByteToString(mapixDesc)
		c.SalesRepresentative, err = conversions.ByteToString(rep)
		c.SalesRepresentativeCode, err = conversions.ByteToString(repCode)
		_ = c.GetLocations()
		parentInt, err := conversions.ByteToInt(parentId)
		if err != nil {
			return customers, err
		}
		if parentInt != 0 {
			par := Customer{Id: parentInt}
			par.GetCustomer()
			c.Parent = &par
		}

		customers = append(customers, c)
	}

	// go redis.Setex(redis_key, customers, 86400)
	return
}

func GetLocationById(id int) (location DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return location, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerByLocation)
	if err != nil {
		return location, err
	}
	var website, eLocal, isPrimary, shippingDefault []byte //ununsed, but in the original query
	var name, email, address, city, phone, fax, contactPerson, postalCode []byte
	var showWeb bool
	err = stmt.QueryRow(id).Scan(
		&location.State.Id,           //s.stateID
		&location.State.State,        //s.state
		&location.State.Abbreviation, //s.abbr as state_abbr
		&location.State.Country.Id,   //cty.countryID,
		&location.DealerType.Id,      //dt.dealer_type as typeID
		&location.DealerType.Type,    // dt.type as dealerType
		&location.DealerType.Online,  // dt.online as typeOnline,
		&location.DealerType.Show,    //dt.show as typeShow
		&location.DealerType.Label,   //dt.label as typeLabel,
		&location.DealerTier.Id,      //dtr.ID as tierID,
		&location.DealerTier.Tier,    //dtr.tier as tier
		&location.DealerTier.Sort,    //dtr.sort as tierSort
		&location.LocationId,
		&name,    //c.name
		&email,   //c.email
		&address, //c.address
		&city,    //c.city,
		&postalCode,
		&phone, //phone,
		&fax,   //c.fax
		&location.Latitude,
		&location.Longitude,
		&location.Id,
		&isPrimary,       //Unused
		&shippingDefault, //Unused
		&location.ContactPerson,
		&showWeb,
		&website,
		&eLocal,
	)
	location.Name, err = conversions.ByteToString(name)
	location.Address, err = conversions.ByteToString(address)
	location.City, err = conversions.ByteToString(city)
	location.PostalCode, err = conversions.ByteToString(postalCode)
	location.Email, err = conversions.ByteToString(email)
	location.Phone, err = conversions.ByteToString(phone)
	location.Fax, err = conversions.ByteToString(fax)
	location.ContactPerson, err = conversions.ByteToString(contactPerson)
	if showWeb {
		if website == nil {
			website = eLocal
		}
		if website != nil {
			location.Website, err = conversions.ByteToUrl(website)
			if err != nil {
				return location, err
			}
		}
	}
	return
}

func SearchLocations(term string) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return locations, err
	}
	defer db.Close()

	stmt, err := db.Prepare(searchDealerLocations)
	if err != nil {
		return locations, err
	}
	term = "%" + term + "%"

	res, err := stmt.Query(term, term)
	if err != nil {
		return locations, err
	}
	var name, email, address, city, phone, fax, contactPerson, postalCode []byte
	var website, eLocal, isPrimary, shippingDefault []byte //ununsed, but in the original query
	var showWeb bool
	for res.Next() {
		var location DealerLocation
		err = res.Scan(
			&location.State.Id,           //s.stateID
			&location.State.State,        //s.state
			&location.State.Abbreviation, //s.abbr as state_abbr
			&location.State.Country.Id,   //cty.countryID,
			&location.DealerType.Id,      //dt.dealer_type as typeID
			&location.DealerType.Type,    // dt.type as dealerType
			&location.DealerType.Online,  // dt.online as typeOnline,
			&location.DealerType.Show,    //dt.show as typeShow
			&location.DealerType.Label,   //dt.label as typeLabel,
			&location.DealerTier.Id,      //dtr.ID as tierID,
			&location.DealerTier.Tier,    //dtr.tier as tier
			&location.DealerTier.Sort,    //dtr.sort as tierSort
			&location.LocationId,
			&name,    //c.name
			&address, //c.address
			&city,    //c.city,
			&postalCode,
			&email, //c.email
			&phone, //phone,
			&fax,   //c.fax
			&location.Latitude,
			&location.Longitude,
			&location.Id,
			&isPrimary,       //Unused
			&shippingDefault, //Unused
			&contactPerson,   //c.contact_person
			&showWeb,
			&website,
			&eLocal,
		)
		if err != nil {
			return locations, err
		}
		location.Name, err = conversions.ByteToString(name)
		location.Address, err = conversions.ByteToString(address)
		location.City, err = conversions.ByteToString(city)
		location.PostalCode, err = conversions.ByteToString(postalCode)
		location.Email, err = conversions.ByteToString(email)
		location.Phone, err = conversions.ByteToString(phone)
		location.Fax, err = conversions.ByteToString(fax)
		location.ContactPerson, err = conversions.ByteToString(contactPerson)
		if showWeb {
			if website == nil {
				website = eLocal
			}
			if website != nil {
				location.Website, err = conversions.ByteToUrl(website)
				if err != nil {
					return locations, err
				}
			}
		}
		locations = append(locations, location)
	}
	return
}

func SearchLocationsByType(term string) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return locations, err
	}
	defer db.Close()

	stmt, err := db.Prepare(dealerLocationsByType)
	if err != nil {
		return locations, err
	}
	term = "%" + term + "%"

	res, err := stmt.Query(term, term)
	if err != nil {
		return locations, err
	}
	var name, email, address, city, phone, fax, contactPerson, postalCode []byte
	var website, eLocal, isPrimary, shippingDefault []byte //ununsed, but in the original query
	var showWeb bool
	for res.Next() {

		var location DealerLocation
		err = res.Scan(
			&location.State.Id,           //s.stateID
			&location.State.State,        //s.state
			&location.State.Abbreviation, //s.abbr as state_abbr
			&location.State.Country.Id,   //cty.countryID,
			&location.DealerType.Id,      //dt.dealer_type as typeID
			&location.DealerType.Type,    // dt.type as dealerType
			&location.DealerType.Online,  // dt.online as typeOnline,
			&location.DealerType.Show,    //dt.show as typeShow
			&location.DealerType.Label,   //dt.label as typeLabel,
			&location.DealerTier.Id,      //dtr.ID as tierID,
			&location.DealerTier.Tier,    //dtr.tier as tier
			&location.DealerTier.Sort,    //dtr.sort as tierSort
			&location.LocationId,
			&name,    //c.name
			&address, //c.address
			&city,    //c.city,
			&postalCode,
			&email, //c.email
			&phone, //phone,
			&fax,   //c.fax
			&location.Latitude,
			&location.Longitude,
			&location.Id,
			&isPrimary,       //Unused
			&shippingDefault, //Unused
			&contactPerson,   //c.contact_person
			&showWeb,
			&website,
			&eLocal,
		)
		if err != nil {
			return locations, err
		}
		location.Name, err = conversions.ByteToString(name)
		location.Address, err = conversions.ByteToString(address)
		location.City, err = conversions.ByteToString(city)
		location.PostalCode, err = conversions.ByteToString(postalCode)
		location.Email, err = conversions.ByteToString(email)
		location.Phone, err = conversions.ByteToString(phone)
		location.Fax, err = conversions.ByteToString(fax)
		location.ContactPerson, err = conversions.ByteToString(contactPerson)
		if showWeb {
			if website == nil {
				website = eLocal
			}
			if website != nil {
				location.Website, err = conversions.ByteToUrl(website)
				if err != nil {
					return locations, err
				}
			}
		}
		locations = append(locations, location)
	}
	return
}

func SearchLocationsByLatLng(loc GeoLocation) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return locations, err
	}
	defer db.Close()

	stmt, err := db.Prepare(searchDealerLocationsByLatLng)
	if err != nil {
		return locations, err
	}
	params := []interface{}{ //all are float64
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
	res, err := stmt.Query(params...)
	if err != nil {
		return locations, err
	}
	var name, email, address, city, phone, fax, contactPerson, postalCode []byte
	var website, eLocal, isPrimary, shippingDefault []byte //ununsed, but in the original query
	var showWeb bool
	for res.Next() {
		var location DealerLocation
		err = res.Scan(
			&location.State.Id,
			&location.State.State,        //s.state
			&location.State.Abbreviation, //s.abbr as state_abbr
			&location.State.Country.Id,   //cty.countryID,
			&location.DealerType.Id,      //dt.dealer_type as typeID
			&location.DealerType.Type,    // dt.type as dealerType
			&location.DealerType.Online,  // dt.online as typeOnline,
			&location.DealerType.Show,    //dt.show as typeShow
			&location.DealerType.Label,   //dt.label as typeLabel,
			&location.DealerTier.Id,      //dtr.ID as tierID,
			&location.DealerTier.Tier,    //dtr.tier as tier
			&location.DealerTier.Sort,    //dtr.sort as tierSort
			&location.LocationId,
			&name,    //c.name
			&address, //c.address
			&city,    //c.city,
			&postalCode,
			&email, //c.email
			&phone, //phone,
			&fax,   //c.fax
			&location.Latitude,
			&location.Longitude,
			&location.Id,
			&isPrimary,       //Unused
			&shippingDefault, //Unused
			&contactPerson,   //c.contact_person
			&showWeb,
			&website,
			&eLocal,
		)
		if err != nil {
			return locations, err
		}
		location.Name, err = conversions.ByteToString(name)
		location.Address, err = conversions.ByteToString(address)
		location.City, err = conversions.ByteToString(city)
		location.PostalCode, err = conversions.ByteToString(postalCode)
		location.Email, err = conversions.ByteToString(email)
		location.Phone, err = conversions.ByteToString(phone)
		location.Fax, err = conversions.ByteToString(fax)
		location.ContactPerson, err = conversions.ByteToString(contactPerson)

		if showWeb {
			if website == nil {
				website = eLocal
			}
			if website != nil {
				location.Website, err = conversions.ByteToUrl(website)
				if err != nil {
					return locations, err
				}
			}
		}
		locations = append(locations, location)
	}
	return
}

func (g *GeoLocation) LatitudeRadians() float64 {
	return (g.Latitude * (math.Pi / 180))
}

func (g *GeoLocation) LongitudeRadians() float64 {
	return (g.Longitude * (math.Pi / 180))
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
