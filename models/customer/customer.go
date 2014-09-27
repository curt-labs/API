package customer

import (
	"database/sql"
	//"encoding/json"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	//"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getCustomerBasicsStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
                            c.latitude, c.longitude, c.searchURL, c.logo, c.website,
                            c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
                            dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
                            dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
                            mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
                            mpx.code as mapix_code, mpx.description as mapic_desc,
                            sr.name as rep_name, sr.code as rep_code, c.parentID
                            from Customer as c
                            left join States as s on c.stateID = s.stateID
                            left join Country as cty on s.countryID = cty.countryID
                            left join DealerTypes as dt on c.dealer_type = dt.dealer_type
                            left join MapIcons as mi on dt.dealer_type = mi.dealer_type
                            left join DealerTiers as dtr on c.tier = dtr.ID
                            left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
                            left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
                            where c.customerID = ?`
	getCustomerLocationsStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
                                cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
                                cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
                                s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
                                from CustomerLocations as cl
                                left join States as s on cl.stateID = s.stateID
                                left join Country as cty on s.countryID = cty.countryID
                                where cl.cust_id = ?`
	getCustomerLocationByIdStmt = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
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
	getCustomerUsersStmt = `select cu.* from CustomerUser as cu
                            join Customer as c on cu.cust_ID = c.cust_id
                            where c.customerID = ? && cu.active = 1`
	getCustomerPriceStmt = `select distinct cp.price from ApiKey as ak
                            join CustomerUser cu on ak.user_id = cu.id
                            join Customer c on cu.cust_ID = c.cust_id
                            join CustomerPricing cp on c.customerID = cp.cust_id
                            where api_key = ? && cp.partID = ?`
	getCustomerPartStmt = `select distinct ci.custPartID from ApiKey as ak
                           join CustomerUser cu on ak.user_id = cu.id
                           join Customer c on cu.cust_ID = c.cust_id
                           join CartIntegration ci on c.customerID = ci.custID
                           where ak.api_key = ? && ci.partID = ?`
	getEtailersStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
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
	getLocalDealersStmt = `select cl.locationID, c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
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
                          group by cl.locationID
                          order by dtr.sort desc`
	getLocalDealerTiersStmt = `select distinct dtr.* from DealerTiers as dtr
                              join Customer as c on dtr.ID = c.tier
                              join DealerTypes as dt on c.dealer_type = dt.dealer_type
                              where dt.online = false and dt.show = true
                              order by dtr.sort`
	getLocalDealerTypesStmt = `select m.ID as iconId, m.mapicon, m.mapiconshadow,
                              dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
                              dt.dealer_type as dealerTypeId, dt.type as dealerType, dt.online, dt.show, dt.label
                              from MapIcons as m
                              join DealerTypes as dt on m.dealer_type = dt.dealer_type
                              join DealerTiers as dtr on m.tier = dtr.ID
                              where dt.show = true
                              order by dtr.sort desc`
	getWhereToBuyDealersStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
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
                               where dt.dealer_type = 1 && dtr.ID = 4 && c.isDummy = false && length(c.searchURL) > 1`
	getMapPolygonStmt = `select s.stateID, s.state, s.abbr,
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
	getMapPolygonCoordinatesForStateStmt = `select mp.ID,mpc.latitude, mpc.longitude
                                            from MapPolygonCoordinates as mpc
                                            join MapPolygon as mp on mpc.MapPolygonID = mp.ID
                                            where mp.stateID = ?`
	searchDealerLocationsStmt = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
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
                                where (dt.dealer_type = 2 or dt.dealer_type = 3) and c.isDummy = false and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	searchDealerLocationsByTypeStmt = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
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
                                      where dt.online = false and c.isDummy = false and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	searchDealerLocationsByLatLngStmt = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
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

type Customer struct {
	Id                                   int
	Name, Email, Address, Address2, City string
	State                                *geography.State
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
	State                                  *geography.State
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

type GeoLocation struct {
	Latitude, Longitude float64
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

type MapPolygon struct {
	Id          int
	Coordinates []GeoLocation
}

type StateRegion struct {
	Id                 int
	Name, Abbreviation string
	Count              int
	Polygons           []MapPolygon
}

type DealerLocation struct {
	Id, LocationId                       int
	Name, Email, Address, Address2, City string
	State                                *geography.State
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
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerBasicsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var parentID *int
	var searchUrl, logoUrl, websiteUrl, mapIconUrl, mapShadowUrl *string

	c.State = &geography.State{}
	c.State.Country = &geography.Country{}

	err = stmt.QueryRow(c.Id).Scan(
		&c.Id,
		&c.Name,
		&c.Email,
		&c.Address,
		&c.Address2,
		&c.City,
		&c.Phone,
		&c.Fax,
		&c.ContactPerson,
		&c.Latitude,
		&c.Longitude,
		&searchUrl,
		&logoUrl,
		&websiteUrl,
		&c.PostalCode,
		&c.State.Id,
		&c.State.State,
		&c.State.Abbreviation,
		&c.State.Country.Id,
		&c.State.Country.Country,
		&c.State.Country.Abbreviation,
		&c.DealerType.Id,
		&c.DealerType.Type,
		&c.DealerType.Online,
		&c.DealerType.Show,
		&c.DealerType.Label,
		&c.DealerTier.Id,
		&c.DealerTier.Tier,
		&c.DealerTier.Sort,
		&c.DealerType.MapIcon.Id,
		&mapIconUrl,
		&mapShadowUrl,
		&c.MapixCode,
		&c.MapixDescription,
		&c.SalesRepresentative,
		&c.SalesRepresentativeCode,
		&parentID,
	)
	defer stmt.Close()

	if searchUrl != nil {
		c.SearchUrl, _ = url.Parse(*searchUrl)
	}

	if logoUrl != nil {
		c.Logo, _ = url.Parse(*logoUrl)
	}

	if websiteUrl != nil {
		c.Website, _ = url.Parse(*websiteUrl)
	}

	if mapIconUrl != nil {
		c.DealerType.MapIcon.MapIcon, _ = url.Parse(*mapIconUrl)
	}

	if mapShadowUrl != nil {
		c.DealerType.MapIcon.MapIconShadow, _ = url.Parse(*mapShadowUrl)
	}

	if parentID != nil && *parentID != 0 {
		c.Parent = &Customer{Id: *parentID}
		c.Parent.GetCustomer()
	}

	return nil
}

func (c *Customer) GetLocations() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerLocationsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(c.Id)
	if err != nil {
		return err
	}

	var locs []CustomerLocation
	for rows.Next() {
		var l CustomerLocation
		l.State = &geography.State{}
		l.State.Country = &geography.Country{}
		err = rows.Scan(
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
		locs = append(locs, l)
	}

	c.Locations = &locs

	return nil
}

func (c *Customer) GetUsers() (users []CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUsersStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(c.Id)
	if err != nil {
		return
	}

	var dbpass string
	var customerID, custID, notCustomer, passConverted int

	for rows.Next() {
		var u CustomerUser

		u.Location = &CustomerLocation{}

		err = rows.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&dbpass,
			&customerID,
			&u.DateAdded,
			&u.Active,
			&u.Location.Id,
			&u.Sudo,
			&custID,
			&notCustomer,
			&passConverted,
		)

		c.Users = append(c.Users, u)
	}

	return
}

func GetCustomerPrice(api_key string, part_id int) (price float64, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerPriceStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(api_key, part_id).Scan(
		&price,
	)

	return
}

func GetCustomerCartReference(api_key string, part_id int) (ref int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerPartStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(api_key, part_id).Scan(
		&ref,
	)

	return
}

/* Internal Use Only */

func GetEtailers() (dealers []Customer, err error) {
	//redis_key := "dealers:etailers"
	//data, err := redis.Get(redis_key)
	//if len(data) > 0 && err != nil {
	//  err = json.Unmarshal(data, &dealers)
	//  if err == nil {
	//      return
	//  }
	//}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getEtailersStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var c Customer
	var parentID *int
	var searchUrl, logoUrl, websiteUrl, mapIconUrl, mapShadowUrl *string

	for rows.Next() {
		c = Customer{}
		c.State = &geography.State{}
		c.State.Country = &geography.Country{}

		err = rows.Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Address,
			&c.Address2,
			&c.City,
			&c.Phone,
			&c.Fax,
			&c.ContactPerson,
			&c.Latitude,
			&c.Longitude,
			&searchUrl,
			&logoUrl,
			&websiteUrl,
			&c.PostalCode,
			&c.State.Id,
			&c.State.State,
			&c.State.Abbreviation,
			&c.State.Country.Id,
			&c.State.Country.Country,
			&c.State.Country.Abbreviation,
			&c.DealerType.Id,
			&c.DealerType.Type,
			&c.DealerType.Online,
			&c.DealerType.Show,
			&c.DealerType.Label,
			&c.DealerTier.Id,
			&c.DealerTier.Tier,
			&c.DealerTier.Sort,
			&c.DealerType.MapIcon.Id,
			&mapIconUrl,
			&mapShadowUrl,
			&c.MapixCode,
			&c.MapixDescription,
			&c.SalesRepresentative,
			&c.SalesRepresentativeCode,
			&parentID,
		)

		if searchUrl != nil {
			c.SearchUrl, _ = url.Parse(*searchUrl)
		}

		if logoUrl != nil {
			c.Logo, _ = url.Parse(*logoUrl)
		}

		if websiteUrl != nil {
			c.Website, _ = url.Parse(*websiteUrl)
		}

		if mapIconUrl != nil {
			c.DealerType.MapIcon.MapIcon, _ = url.Parse(*mapIconUrl)
		}

		if mapShadowUrl != nil {
			c.DealerType.MapIcon.MapIconShadow, _ = url.Parse(*mapShadowUrl)
		}

		if parentID != nil && *parentID != 0 {
			c.Parent = &Customer{Id: *parentID}
			c.Parent.GetCustomer()
		}

		dealers = append(dealers, c)
	}

	//go redis.Setex(redis_key, dealers, 86400)

	return
}

func GetLocalDealers(center string, latlng string) (dealers []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getLocalDealersStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

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

	rows, err := stmt.Query(
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
	)
	if err != nil {
		return
	}

	var parentID *int
	var searchUrl, logoUrl, websiteUrl, mapIconUrl, mapShadowUrl *string

	for rows.Next() {
		var l DealerLocation

		l.State = &geography.State{}
		l.State.Country = &geography.Country{}

		err = rows.Scan(
			&l.LocationId,
			&l.Id,
			&l.Name,
			&l.Email,
			&l.Address,
			&l.City,
			&l.Phone,
			&l.Fax,
			&l.ContactPerson,
			&l.Latitude,
			&l.Longitude,
			&searchUrl,
			&logoUrl,
			&websiteUrl,
			&l.PostalCode,
			&l.State.Id,
			&l.State.State,
			&l.State.Abbreviation,
			&l.State.Country.Id,
			&l.State.Country.Country,
			&l.State.Country.Abbreviation,
			&l.DealerType.Id,
			&l.DealerType.Type,
			&l.DealerType.Online,
			&l.DealerType.Show,
			&l.DealerType.Label,
			&l.DealerTier.Id,
			&l.DealerTier.Tier,
			&l.DealerTier.Sort,
			&l.DealerType.MapIcon.Id,
			&mapIconUrl,
			&mapShadowUrl,
			&l.MapixCode,
			&l.MapixDescription,
			&l.SalesRepresentative,
			&l.SalesRepresentativeCode,
			&parentID,
		)

		if searchUrl != nil {
			l.SearchUrl, _ = url.Parse(*searchUrl)
		}

		if logoUrl != nil {
			l.Logo, _ = url.Parse(*logoUrl)
		}

		if websiteUrl != nil {
			l.Website, _ = url.Parse(*websiteUrl)
		}

		if mapIconUrl != nil {
			l.DealerType.MapIcon.MapIcon, _ = url.Parse(*mapIconUrl)
		}

		if mapShadowUrl != nil {
			l.DealerType.MapIcon.MapIconShadow, _ = url.Parse(*mapShadowUrl)
		}

		l.Distance = api_helpers.EARTH * (2 * math.Atan2(
			math.Sqrt((math.Sin(((l.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((l.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((l.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((l.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(l.Latitude*(math.Pi/180))),
			math.Sqrt(1-((math.Sin(((l.Latitude-clat)*(math.Pi/180))/2)*math.Sin(((l.Latitude-clat)*(math.Pi/180))/2))+((math.Sin(((l.Longitude-clong)*(math.Pi/180))/2))*(math.Sin(((l.Longitude-clong)*(math.Pi/180))/2)))*math.Cos(clat*(math.Pi/180))*math.Cos(l.Latitude*(math.Pi/180))))))

		dealers = append(dealers, l)
	}

	sortutil.AscByField(dealers, "Distance")
	return
}

func GetLocalRegions() (regions []StateRegion, err error) {

	//redis_key := "local:regions"

	// Attempt to get the local regions from Redis
	//data, err := redis.Get(redis_key)
	//if len(data) > 0 && err != nil {
	//  err = json.Unmarshal(data, &regions)
	//  if err == nil {
	//      return
	//  }
	//}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	polyStmt, err := db.Prepare(getMapPolygonStmt)
	if err != nil {
		return
	}
	defer polyStmt.Close()

	coordStmt, err := db.Prepare(getMapPolygonCoordinatesForStateStmt)
	if err != nil {
		return
	}
	defer coordStmt.Close()

	// Get the local regions from the database
	_, err = db.Exec("SET SESSION group_concat_max_len = 100024")
	if err != nil {
		return
	}
	polygonRows, err := polyStmt.Query()
	if err != nil {
		return
	}
	_, err = db.Exec("SET SESSION group_concat_max_len = 1024")
	if err != nil {
		return
	}

	for polygonRows.Next() {
		var region StateRegion
		err = polygonRows.Scan(
			&region.Id,
			&region.Name,
			&region.Abbreviation,
			&region.Count,
		)
		if err != nil {
			return
		}

		//build out the polygons for this state, including latitude/longitude
		coordRows, err := coordStmt.Query(region.Id)
		if err != nil {
			return regions, err
		}

		polygons := make(map[int]MapPolygon, 0)
		for coordRows.Next() {
			var tmpPolygon MapPolygon
			var tmpGeo GeoLocation
			err = coordRows.Scan(
				&tmpPolygon.Id,
				&tmpGeo.Latitude,
				&tmpGeo.Longitude,
			)

			//check if we have an index for this polygon created
			if _, ok := polygons[tmpPolygon.Id]; !ok {
				polygons[tmpPolygon.Id] = MapPolygon{
					Id:          tmpPolygon.Id,
					Coordinates: make([]GeoLocation, 0),
				}
			}

			//add the geolocation info to our polygon
			poly := polygons[tmpPolygon.Id]
			poly.Coordinates = append(poly.Coordinates, tmpGeo)
			polygons[tmpPolygon.Id] = poly
		}

		//drop the key-value pair (our end user doesn't need that)
		var polys []MapPolygon
		for _, poly := range polygons {
			polys = append(polys, poly)
		}

		region.Polygons = polys

		regions = append(regions, region)
	}

	// We're not going to set the expiration on this
	// it won't ever change...until the San Andreas fault
	// completely drops the western part of CA anyway :/
	//go redis.Set(redis_key, regions)

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
	//redis_key := "local:tiers"

	// Attempt to get the local tiers from Redis
	//data, err := redis.Get(redis_key)
	//if len(data) > 0 && err != nil {
	//	err = json.Unmarshal(data, &tiers)
	//	if err == nil {
	//		return
	//	}
	//}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getLocalDealerTiersStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var tier DealerTier
		err = rows.Scan(
			&tier.Id,
			&tier.Tier,
			&tier.Sort,
		)
		if err != nil {
			return
		}
		tiers = append(tiers, tier)
	}

	//go redis.Setex(redis_key, tiers, 86400)

	return
}

func GetLocalDealerTypes() (graphics []MapGraphics) {
	//redis_key := "local:types"

	// Attempt to get the local types from Redis
	//data, err := redis.Get(redis_key)
	//if len(data) > 0 && err != nil {
	//	err = json.Unmarshal(data, &graphics)
	//	if err == nil {
	//		return
	//	}
	//}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getLocalDealerTypesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var mapIconUrl, mapShadowUrl *string

	for rows.Next() {
		var g MapGraphics
		err = rows.Scan(
			&g.MapIcon.Id,
			&mapIconUrl,
			&mapShadowUrl,
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
			return
		}

		if mapIconUrl != nil {
			g.DealerType.MapIcon.MapIcon, _ = url.Parse(*mapIconUrl)
		}

		if mapShadowUrl != nil {
			g.DealerType.MapIcon.MapIconShadow, _ = url.Parse(*mapShadowUrl)
		}

		graphics = append(graphics, g)
	}

	//go redis.Setex(redis_key, graphics, 86400)

	return
}

func GetWhereToBuyDealers() (customers []Customer) {
	//redis_key := "dealers:wheretobuy"

	//data, err := redis.Get(redis_key)
	//if len(data) > 0 && err != nil {
	//	err = json.Unmarshal(data, &customers)
	//	if err == nil {
	//		return
	//	}
	//}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getWhereToBuyDealersStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var parentID *int
	var logoUrl, searchUrl, websiteUrl, mapIconUrl, mapShadowUrl *string

	for rows.Next() {
		var c Customer
		c.State = &geography.State{}
		c.State.Country = &geography.Country{}
		err = rows.Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Address,
			&c.Address2,
			&c.City,
			&c.Phone,
			&c.Fax,
			&c.ContactPerson,
			&c.Latitude,
			&c.Longitude,
			&searchUrl,
			&logoUrl,
			&websiteUrl,
			&c.PostalCode,
			&c.State.Id,
			&c.State.State,
			&c.State.Abbreviation,
			&c.State.Country.Id,
			&c.State.Country.Country,
			&c.State.Country.Abbreviation,
			&c.DealerType.Id,
			&c.DealerType.Type,
			&c.DealerType.Online,
			&c.DealerType.Show,
			&c.DealerType.Label,
			&c.DealerTier.Id,
			&c.DealerTier.Tier,
			&c.DealerTier.Sort,
			&c.DealerType.MapIcon.Id,
			&mapIconUrl,
			&mapShadowUrl,
			&c.MapixCode,
			&c.MapixDescription,
			&c.SalesRepresentative,
			&c.SalesRepresentativeCode,
			&parentID,
		)

		if searchUrl != nil {
			c.SearchUrl, _ = url.Parse(*searchUrl)
		}

		if logoUrl != nil {
			c.Logo, _ = url.Parse(*logoUrl)
		}

		if websiteUrl != nil {
			c.Website, _ = url.Parse(*websiteUrl)
		}

		if mapIconUrl != nil {
			c.DealerType.MapIcon.MapIcon, _ = url.Parse(*mapIconUrl)
		}

		if mapShadowUrl != nil {
			c.DealerType.MapIcon.MapIconShadow, _ = url.Parse(*mapShadowUrl)
		}

		err = c.GetLocations()

		if parentID != nil && *parentID != 0 {
			c.Parent = &Customer{Id: *parentID}
			c.Parent.GetCustomer()
		}

		customers = append(customers, c)
	}

	//go redis.Setex(redis_key, customers, 86400)

	return
}

func GetLocationById(id int) (location DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerLocationByIdStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	location.State = &geography.State{}
	location.State.Country = &geography.Country{}

	var isPrimary, shippingDefault, showWeb bool
	var elocalUrl, websiteUrl *string

	err = stmt.QueryRow(id).Scan(
		&location.State.Id,
		&location.State.State,
		&location.State.Abbreviation,
		&location.State.Country.Id,
		&location.DealerType.Id,
		&location.DealerType.Type,
		&location.DealerType.Online,
		&location.DealerType.Show,
		&location.DealerType.Label,
		&location.DealerTier.Id,
		&location.DealerTier.Tier,
		&location.DealerTier.Sort,
		&location.LocationId,
		&location.Name,
		&location.Email,
		&location.Address,
		&location.City,
		&location.PostalCode,
		&location.Phone,
		&location.Fax,
		&location.Latitude,
		&location.Longitude,
		&location.Id,
		&isPrimary,
		&shippingDefault,
		&location.ContactPerson,
		&showWeb,
		&websiteUrl,
		&elocalUrl,
	)
	if err != nil {
		return
	}

	if showWeb {
		if websiteUrl != nil {
			location.Website, _ = url.Parse(*websiteUrl)
		} else if elocalUrl != nil {
			location.Website, _ = url.Parse(*elocalUrl)
		}
	}

	return
}

func SearchLocations(term string) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(searchDealerLocationsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	term = "%" + term + "%"

	rows, err := stmt.Query(term, term)
	if err != nil {
		return
	}

	var isPrimary, shippingDefault, showWeb bool
	var elocalUrl, websiteUrl *string

	for rows.Next() {
		var loc DealerLocation
		loc.State = &geography.State{}
		loc.State.Country = &geography.Country{}
		err = rows.Scan(
			&loc.State.Id,
			&loc.State.State,
			&loc.State.Abbreviation,
			&loc.State.Country.Id,
			&loc.DealerType.Id,
			&loc.DealerType.Type,
			&loc.DealerType.Online,
			&loc.DealerType.Show,
			&loc.DealerType.Label,
			&loc.DealerTier.Id,
			&loc.DealerTier.Tier,
			&loc.DealerTier.Sort,
			&loc.LocationId,
			&loc.Name,
			&loc.Address,
			&loc.City,
			&loc.PostalCode,
			&loc.Email,
			&loc.Phone,
			&loc.Fax,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Id,
			&isPrimary,
			&shippingDefault,
			&loc.ContactPerson,
			&showWeb,
			&websiteUrl,
			&elocalUrl,
		)
		if err != nil {
			return
		}

		if showWeb {
			if websiteUrl != nil {
				loc.Website, _ = url.Parse(*websiteUrl)
			} else if elocalUrl != nil {
				loc.Website, _ = url.Parse(*elocalUrl)
			}
		}

		locations = append(locations, loc)
	}

	return
}

func SearchLocationsByType(term string) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(searchDealerLocationsByTypeStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	term = "%" + term + "%"

	rows, err := stmt.Query(term, term)
	if err != nil {
		return
	}

	var isPrimary, shippingDefault, showWeb bool
	var elocalUrl, websiteUrl *string

	for rows.Next() {
		var loc DealerLocation
		loc.State = &geography.State{}
		loc.State.Country = &geography.Country{}
		err = rows.Scan(
			&loc.State.Id,
			&loc.State.State,
			&loc.State.Abbreviation,
			&loc.State.Country.Id,
			&loc.DealerType.Id,
			&loc.DealerType.Type,
			&loc.DealerType.Online,
			&loc.DealerType.Show,
			&loc.DealerType.Label,
			&loc.DealerTier.Id,
			&loc.DealerTier.Tier,
			&loc.DealerTier.Sort,
			&loc.LocationId,
			&loc.Name,
			&loc.Address,
			&loc.City,
			&loc.PostalCode,
			&loc.Email,
			&loc.Phone,
			&loc.Fax,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Id,
			&isPrimary,
			&shippingDefault,
			&loc.ContactPerson,
			&showWeb,
			&websiteUrl,
			&elocalUrl,
		)
		if err != nil {
			return
		}

		if showWeb {
			if websiteUrl != nil {
				loc.Website, _ = url.Parse(*websiteUrl)
			} else if elocalUrl != nil {
				loc.Website, _ = url.Parse(*elocalUrl)
			}
		}

		locations = append(locations, loc)
	}

	return
}

func SearchLocationsByLatLng(loc GeoLocation) (locations []DealerLocation, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(searchDealerLocationsByLatLngStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
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
	)
	if err != nil {
		return
	}

	var isPrimary, shippingDefault, showWeb bool
	var elocalUrl, websiteUrl *string

	for rows.Next() {
		var loc DealerLocation
		loc.State = &geography.State{}
		loc.State.Country = &geography.Country{}
		err = rows.Scan(
			&loc.State.Id,
			&loc.State.State,
			&loc.State.Abbreviation,
			&loc.State.Country.Id,
			&loc.DealerType.Id,
			&loc.DealerType.Type,
			&loc.DealerType.Online,
			&loc.DealerType.Show,
			&loc.DealerType.Label,
			&loc.DealerTier.Id,
			&loc.DealerTier.Tier,
			&loc.DealerTier.Sort,
			&loc.LocationId,
			&loc.Name,
			&loc.Address,
			&loc.City,
			&loc.PostalCode,
			&loc.Email,
			&loc.Phone,
			&loc.Fax,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Id,
			&isPrimary,
			&shippingDefault,
			&loc.ContactPerson,
			&showWeb,
			&websiteUrl,
			&elocalUrl,
		)
		if err != nil {
			return
		}

		if showWeb {
			if websiteUrl != nil {
				loc.Website, _ = url.Parse(*websiteUrl)
			} else if elocalUrl != nil {
				loc.Website, _ = url.Parse(*elocalUrl)
			}
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
