package customer_new

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Customer struct {
	Id                  int                 `json:"id,omitempty" xml:"id,omitempty"`
	Name                string              `json:"name,omitempty" xml:"name,omitempty"`
	Email               string              `json:"email,omitempty" xml:"email,omitempty"`
	Address             string              `json:"address,omitempty" xml:"address,omitempty"`
	Address2            string              `json:"address2,omitempty" xml:"address2,omitempty"`
	City                string              `json:"city,omitempty" xml:"city,omitempty"`
	State               geography.State     `json:"state,omitempty" xml:"state,omitempty"`
	PostalCode          string              `json:"postalCode,omitempty" xml:"postalCode,omitempty"`
	Phone               string              `json:"phone,omitempty" xml:"phone,omitempty"`
	Fax                 string              `json:"fax,omitempty" xml:"fax,omitempty"`
	ContactPerson       string              `json:"contactPerson,omitempty" xml:"contactPerson,omitempty"`
	Latitude            float64             `json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude           float64             `json:"longitude,omitempty" xml:"longitude,omitempty"`
	Website             url.URL             `json:"website,omitempty" xml:"website,omitempty"`
	Parent              *Customer           `json:"parent,omitempty" xml:"parent,omitempty"`
	SearchUrl           url.URL             `json:"searchUrl,omitempty" xml:"searchUrl,omitempty"`
	Logo                url.URL             `json:"logo,omitempty" xml:"logo,omitempty"`
	DealerType          DealerType          `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	DealerTier          DealerTier          `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	Locations           []CustomerLocation  `json:"locations,omitempty" xml:"locations,omitempty"`
	Users               []CustomerUser      `json:"users,omitempty" xml:"users,omitempty"`
	CustomerId          int                 `json:"customerId,omitempty" xml:"customerId,omitempty"`
	IsDummy             bool                `json:"isDummy,omitempty" xml:"isDummy,omitempty"`
	ELocalUrl           url.URL             `json:"eLocalUrl,omitempty" xml:"eLocalUrl,omitempty"`
	MapixCode           MapixCode           `json:"mapixCode,omitempty" xml:"mapixCode,omitempty"`
	ApiKey              string              `json:"apiKey,omitempty" xml:"apiKey,omitempty"`
	ShowWebsite         bool                `json:"showWebsite,omitempty" xml:"showWebsite,omitempty"`
	SalesRepresentative SalesRepresentative `json:"salesRepresentative,omitempty" xml:"salesRepresentative,omitempty"`
}

type Customers []Customer

type Scanner interface {
	Scan(...interface{}) error
}

type SalesRepresentative struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	Code string `json:"code,omitempty" xml:"code,omitempty"`
}

type CustomerLocation struct {
	Id              int             `json:"id,omitempty" xml:"id,omitempty"`
	Name            string          `json:"name,omitempty" xml:"name,omitempty"`
	Email           string          `json:"email,omitempty" xml:"email,omitempty"`
	Address         string          `json:"address,omitempty" xml:"address,omitempty"`
	City            string          `json:"city,omitempty" xml:"city,omitempty"`
	PostalCode      string          `json:"postalCode,omitempty" xml:"postalCode,omitempty"`
	State           geography.State `json:"state,omitempty" xml:"state,omitempty"`
	Phone           string          `json:"phone,omitempty" xml:"phone,omitempty"`
	Fax             string          `json:"fax,omitempty" xml:"fax,omitempty"`
	Latitude        float64         `json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude       float64         `json:"longitude,omitempty" xml:"longitude,omitempty"`
	CustomerId      int             `json:"customerId,omitempty" xml:"customerId,omitempty"`
	ContactPerson   string          `json:"contactPerson,omitempty" xml:"contactPerson,omitempty"`
	IsPrimary       bool            `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	ShippingDefault bool            `json:"shippingDefault,omitempty" xml:"shippingDefault,omitempty"`
	ShowWebSite     bool            `json:"showWebsite,omitempty" xml:"showWebsite,omitempty"`
	ELocalUrl       url.URL         `json:"eLocalUrl,omitempty" xml:"eLocalUrl,omitempty"`
	Website         url.URL         `json:"website,omitempty" xml:"website,omitempty"`
}

type DealerType struct {
	Id      int     `json:"id,omitempty" xml:"id,omitempty"`
	Type    string  `json:"type,omitempty" xml:"type,omitempty"`
	Label   string  `json:"label,omitempty" xml:"label,omitempty"`
	Online  bool    `json:"online,omitempty" xml:"online,omitempty"`
	Show    bool    `json:"show,omitempty" xml:"show,omitempty"`
	MapIcon MapIcon `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
}

type DealerTier struct {
	Id   int    `json:"id,omitempty" xml:"id,omitempty"`
	Tier string `json:"tier,omitempty" xml:"tier,omitempty"`
	Sort int    `json:"sort,omitempty" xml:"sort,omitempty"`
}

type MapIcon struct {
	Id            int `json:"id,omitempty" xml:"id,omitempty"`
	TierId        int `json:"tierId,omitempty" xml:"tierId,omitempty"`
	DealerTypeId  int
	MapIcon       url.URL `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
	MapIconShadow url.URL `json:"mapIconShadow,omitempty" xml:"mapIconShadow,omitempty"`
}

type MapGraphics struct {
	DealerTier DealerTier `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	DealerType DealerType `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	MapIcon    MapIcon    `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" xml:"longitude,omitempty"`
}

type MapixCode struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Code        string `json:"code,omitempty" xml:"code,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type DealerLocation struct {
	CustomerLocation    CustomerLocation    `json:"id,omitempty" xml:"id,omitempty"`
	Distance            float64             `json:"distance,omitempty" xml:"distance,omitempty"`
	Website             url.URL             `json:"website,omitempty" xml:"website,omitempty"`
	Parent              *Customer           `json:"parent,omitempty" xml:"parent,omitempty"`
	SearchUrl           url.URL             `json:"searchUrl,omitempty" xml:"searchUrl,omitempty"`
	Logo                url.URL             `json:"logo,omitempty" xml:"logo,omitempty"`
	DealerType          DealerType          `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	DealerTier          DealerTier          `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	SalesRepresentative SalesRepresentative `json:"salesRepresentative,omitempty" xml:"salesRepresentative,omitempty"`
	MapixCode           MapixCode           `json:"mapixCode,omitempty" xml:"mapixCode,omitempty"`
}
type DealerLocations []DealerLocation

type StateRegion struct {
	Id           int          `json:"id,omitempty" xml:"id,omitempty"`
	Name         string       `json:"name,omitempty" xml:"name,omitempty"`
	Abbreviation string       `json:"abbreviation,omitempty" xml:"abbreviation,omitempty"`
	Count        int          `json:"count,omitempty" xml:"count,omitempty"`
	Polygons     []MapPolygon `json:"polygon,omitempty" xml:"polygon,omitempty"`
}

type MapPolygon struct {
	Id          int           `json:"id,omitempty" xml:"id,omitempty"`
	Coordinates []GeoLocation `json:"coordinates,omitempty" xml:"coordinates,omitempty"`
}

const (
	customerFields = ` c.cust_id, c.name, c.email, c.address,  c.city, c.stateID, c.phone, c.fax, c.contact_person, c.dealer_type,
				c.latitude,c.longitude,  c.website, c.customerID, c.isDummy, c.parentID, c.searchURL, c.eLocalURL, c.logo,c.address2,
				c.postal_code, c.mCodeID, c.salesRepID, c.APIKey, c.tier, c.showWebsite `
	stateFields            = ` IFNULL(s.state, ""), IFNULL(s.abbr, ""), IFNULL(s.countryID, "0") `
	countryFields          = ` cty.name, cty.abbr `
	dealerTypeFields       = ` IFNULL(dt.type, ""), IFNULL(dt.online, ""), IFNULL(dt.show, ""), IFNULL(dt.label, "") `
	dealerTierFields       = ` IFNULL(dtr.tier, ""), IFNULL(dtr.sort, "") `
	mapIconFields          = ` IFNULL(mi.mapicon, ""), IFNULL(mi.mapiconshadow, "") ` //joins on dealer_type usually
	mapixCodeFields        = ` IFNULL(mpx.code, ""), IFNULL(mpx.description, "") `
	salesRepFields         = ` IFNULL(sr.name, ""), IFNULL(sr.code, "") `
	customerLocationFields = ` cl.locationID, cl.name, cl.address, cl.city, cl.stateID,  cl.email, cl.phone, cl.fax,
							cl.latitude, cl.longitude, cl.cust_id, cl.contact_person, cl.isprimary, cl.postalCode, cl.ShippingDefault `
	showSiteFields = ` c.showWebsite, c.website, c.eLocalURL `

	//redis
	custPrefix = "customer:"
)

var (
	getCustomerIdFromKeyStmt = `select c.customerID from Customer as c
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`
	//Old
	findCustomerIdFromCustId = `select customerID from Customer where cust_id = ? limit 1`
	findCustIdFromCustomerId = `select cust_id from Customer where customerID = ? limit 1`
	basics                   = `select ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
			from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where c.cust_id = ? `
	//affects CL methods
	customerLocation = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `
							from CustomerLocations as cl
							join Customer as c on cl.cust_id = c.cust_id
	 						left join States as s on cl.stateID = s.stateID
	 						left join Country as cty on s.countryID = cty.countryID
							where c.cust_id = ?`

	customerUser = `select cu.id, cu.name, cu.email, cu.date_added, cu.active, cu.isSudo from CustomerUser as cu
						join Customer as c on cu.cust_ID = c.cust_id
						where c.cust_id = ?
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
	etailers = `select ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where dt.online = 1 && c.isDummy = 0 order by c.name`

	localDealers = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + ` ,` + showSiteFields + `
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

	whereToBuyDealers = `select ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
			from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where dt.dealer_type = 1 and dtr.ID = 4 and c.isDummy = false and length(c.searchURL) > 1`

	customerByLocation = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `  ,` + showSiteFields + `
								from CustomerLocations as cl
								join States as s on cl.stateID = s.stateID
								left join country as cty on cty.countryID = s.countryID
								join Customer as c on cl.cust_id = c.cust_id
								join DealerTypes as dt on c.dealer_type = dt.dealer_type
								join DealerTiers as dtr on c.tier = dtr.ID
								left join MapIcons as mi on dt.dealer_type = mi.dealer_type
								left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
								left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
								where (dt.dealer_type = 2 or dt.dealer_type = 3) and c.isDummy = false
								and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`

	searchDealerLocations = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + ` ,` + showSiteFields + `
								from CustomerLocations as cl
								join States as s on cl.stateID = s.stateID
								left join country as cty on cty.countryID = s.countryID
								join Customer as c on cl.cust_id = c.cust_id
								join DealerTypes as dt on c.dealer_type = dt.dealer_type
								join DealerTiers as dtr on c.tier = dtr.ID
								left join MapIcons as mi on dt.dealer_type = mi.dealer_type
								left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
								left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
								where (dt.dealer_type = 2 or dt.dealer_type = 3) and c.isDummy = false
								and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`

	dealerLocationsByType = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + ` ,` + showSiteFields + `
								from CustomerLocations as cl
	 							join States as s on cl.stateID = s.stateID
	 							left join Country as cty ON cty.countryID = s.countryID
	 							join Customer as c on cl.cust_id = c.cust_id
	 							join DealerTypes as dt on c.dealer_type = dt.dealer_type
	 							join DealerTiers as dtr on c.tier = dtr.ID
	 							left join MapIcons as mi on dtr.tier = mi.tier
	 							left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
	 							left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
								where dt.online = false and c.isDummy = false
								and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`

	searchDealerLocationsByLatLng = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `, ` + showSiteFields + `
									from CustomerLocations as cl
									join States as s on cl.stateID = s.stateID
									left join Country as cty ON cty.countryID = s.countryID
									join Customer as c on cl.cust_id = c.cust_id
									join DealerTypes as dt on c.dealer_type = dt.dealer_type
									join DealerTiers as dtr on c.tier = dtr.ID
									left join MapIcons as mi on dtr.tier = mi.tier
		 							left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
		 							left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
									where dt.online = false and c.isDummy = false
									and dt.show = true and
									( ? * (
										2 * ATAN2(
											SQRT((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))),
											SQRT(1 - ((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))))
										)
									) < 100.0)`

	//customer Crud
	createCustomer = `insert into Customer (name, email, address,  city, stateID, phone, fax, contact_person, dealer_type, latitude,longitude,  website, customerID, isDummy, parentID, searchURL,
					eLocalURL, logo,address2, postal_code, mCodeID, salesRepID, APIKey, tier, showWebsite) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	updateCustomer = `update Customer set name = ?, email = ?, address = ?, city = ?, stateID = ?, phone = ?, fax = ?, contact_person = ?, dealer_type = ?, latitude = ?, longitude = ?,  website = ?, customerID = ?,
					isDummy = ?, parentID = ?, searchURL = ?, eLocalURL = ?, logo = ?, address2 = ?, postal_code = ?, mCodeID = ?, salesRepID = ?, APIKey = ?, tier = ?, showWebsite = ? where cust_id = ?`
	deleteCustomer = `delete from Customer where cust_id = ?`
)

func (c *Customer) Get() error {
	var err error
	redis_key := custPrefix + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c)
		if err == nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res := stmt.QueryRow(c.Id)
	c, err = ScanCustomerTableFields(res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	typeChan := make(chan int)
	go func() {
		//check redis
		redis_key := "dealerType:" + strconv.Itoa(c.DealerType.Id)
		data, err := redis.Get(redis_key)
		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &c.DealerType)
			typeChan <- 1
			return
		}

		typeMap, err := DealerTypeMap()
		if err != nil {
			typeChan <- 1
			return
		}
		if dType, ok := typeMap[c.DealerType.Id]; ok {
			c.DealerType = dType
		}
		typeChan <- 1
	}()
	tierChan := make(chan int)
	go func() {
		//check redis
		redis_key := "dealerTier:" + strconv.Itoa(c.DealerTier.Id)
		data, err := redis.Get(redis_key)
		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &c.DealerTier)
			tierChan <- 1
			return
		}
		tierMap, err := DealerTierMap()
		if err == nil {
			if dTier, ok := tierMap[c.DealerTier.Id]; ok {
				c.DealerTier = dTier
			}
		}

		tierChan <- 1
	}()
	mapixChan := make(chan int)
	go func() {
		//check redis
		redis_key := "mapixCode:" + strconv.Itoa(c.MapixCode.ID)
		data, err := redis.Get(redis_key)
		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &c.MapixCode)
			mapixChan <- 1
			return
		}
		mapixMap, err := MapixMap()
		if err == nil {
			if mapix, ok := mapixMap[c.MapixCode.ID]; ok {
				c.MapixCode = mapix
			}
		}

		mapixChan <- 1
	}()
	repChan := make(chan int)
	go func() {
		//check redis
		redis_key := "salesRep:" + strconv.Itoa(c.SalesRepresentative.ID)
		data, err := redis.Get(redis_key)
		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &c.SalesRepresentative)
			repChan <- 1
			return
		}
		repMap, err := SalesRepMap()
		if err == nil {
			if rep, ok := repMap[c.SalesRepresentative.ID]; ok {
				c.SalesRepresentative = rep
			}
		}

		repChan <- 1
	}()

	//get geography
	redis_key_state := "state:" + strconv.Itoa(c.State.Id)
	data, err = redis.Get(redis_key_state)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.State)
		if err == nil {
			redis_key_country := "country:" + strconv.Itoa(c.State.Country.Id)
			data, err = redis.Get(redis_key_country)
			if err == nil && len(data) > 0 {
				json.Unmarshal(data, &c.State.Country)
			}
		}
	} else {
		stateMap, stateErr := geography.GetStateMap()
		countryMap, countryErr := geography.GetCountryMap()
		if stateErr == nil && countryErr == nil {
			if state, ok := stateMap[c.State.Id]; ok {
				c.State = state
				if country, ok := countryMap[c.State.Country.Id]; ok {
					*c.State.Country = country
				}
			}
		}
	}

	//TODO mapicons
	<-repChan
	<-mapixChan
	<-tierChan
	<-typeChan

	redis.Setex(redis_key, c, redis.CacheTimeout)

	return nil
}

func (c *Customer) GetCustomer(key string) (err error) {

	basicsChan := make(chan error)

	go func() {
		err := c.Basics(key)
		if err == nil {
			basicsChan <- c.GetUsers(key)
		}
		basicsChan <- err
	}()

	c.GetLocations()
	err = <-basicsChan

	if err == sql.ErrNoRows {
		err = fmt.Errorf("error: %s", "failed to retrieve")
	}

	return err
}

func (c *Customer) GetCustomerIdFromKey(key string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustomerIdFromKeyStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(key).Scan(&c.Id)
	if err != nil {
		return err
	}
	return nil
}

//redundant with Get - uses SQL joins; faster?
func (c *Customer) Basics(key string) (err error) {
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

	return c.ScanCustomer(stmt.QueryRow(c.Id), key)
}

func (c *Customer) GetLocations() (err error) {
	redis_key := "customerLocations:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.Locations)
		return err
	}
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

	rows, err := stmt.Query(c.Id)

	for rows.Next() {
		loc, err := ScanLocation(rows)
		if err != nil {
			return err
		}
		c.Locations = append(c.Locations, *loc)
	}
	defer rows.Close()

	redis.Setex(redis_key, c.Locations, redis.CacheTimeout)

	return err
}

func (c *Customer) FindCustomerIdFromCustId() (err error) { //Jesus, really?
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(findCustomerIdFromCustId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.Id).Scan(&c.CustomerId)
	if err != nil {
		return err
	}
	return err
}

func (c *Customer) FindCustIdFromCustomerId() (err error) { //Jesus, really?
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(findCustIdFromCustomerId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.CustomerId).Scan(&c.Id)
	if err != nil {
		return err
	}
	return err
}

func (c *Customer) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	parentId := 0
	if c.Parent != nil {
		parentId = c.Parent.Id
	}
	res, err := stmt.Exec(
		c.Name,
		c.Email,
		c.Address,
		c.City,
		c.State.Id,
		c.Phone,
		c.Fax,
		c.ContactPerson,
		c.DealerType.Id,
		c.Latitude,
		c.Longitude,
		c.Website.String(),
		c.CustomerId,
		c.IsDummy,
		parentId,
		c.SearchUrl.String(),
		c.ELocalUrl.String(),
		c.Logo.String(),
		c.Address2,
		c.PostalCode,
		c.MapixCode.ID,
		c.SalesRepresentative.ID,
		c.ApiKey,
		c.DealerTier.Id,
		c.ShowWebsite,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	c.Id = int(id)

	return err
}

func (c *Customer) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	parentId := 0
	if c.Parent != nil {
		parentId = c.Parent.Id
	}
	_, err = stmt.Exec(
		c.Name,
		c.Email,
		c.Address,
		c.City,
		c.State.Id,
		c.Phone,
		c.Fax,
		c.ContactPerson,
		c.DealerType.Id,
		c.Latitude,
		c.Longitude,
		c.Website.String(),
		c.CustomerId,
		c.IsDummy,
		parentId,
		c.SearchUrl.String(),
		c.ELocalUrl.String(),
		c.Logo.String(),
		c.Address2,
		c.PostalCode,
		c.MapixCode.ID,
		c.SalesRepresentative.ID,
		c.ApiKey,
		c.DealerTier.Id,
		c.ShowWebsite,
		c.Id,
	)
	if err != nil {
		return err
	}
	err = redis.Set(custPrefix+strconv.Itoa(c.Id), c)
	return nil
}

func (c *Customer) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id)
	if err != nil {
		return err
	}
	redis.Delete(custPrefix + strconv.Itoa(c.Id))

	return nil
}

func (c *Customer) GetUsers(key string) (err error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(c.Id)
	var name []byte
	iter := 0
	userChan := make(chan error)
	lowerKey := strings.ToLower(key)

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
			continue
		}
		u.Name, err = conversions.ByteToString(name)

		go func(user CustomerUser) {
			if err := user.GetKeys(); err == nil {
				for _, k := range user.Keys {
					if k.Key == lowerKey {
						user.Current = true
						break
					}
				}
			}

			user.GetLocation()

			c.Users = append(c.Users, user)
			userChan <- nil

		}(u)
		iter++
	}
	defer res.Close()

	for i := 0; i < iter; i++ {
		<-userChan
	}

	return err
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

func GetEtailers(key string) (dealers []Customer, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dealers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(etailers)
	if err != nil {
		return dealers, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return dealers, err
	}
	defer rows.Close()

	for rows.Next() {
		var cust Customer
		if err := cust.ScanCustomer(rows, key); err == nil {
			dealers = append(dealers, cust)
		}
	}

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

	for res.Next() {
		l, err := ScanDealerLocation(res)
		if err != nil {
			return dealers, err
		}
		dealers = append(dealers, *l)
	}
	defer res.Close()

	sortutil.AscByField(dealers, "Distance")
	return
}

func GetLocalRegions() (regions []StateRegion, err error) {
	redis_key := "local:regions"
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
		err = json.Unmarshal(data, &regions)
		if err == nil {
			return
		}
	}
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
	defer res.Close()
	go redis.Set(redis_key, regions)
	return
}

func GetLocalDealerTiers() (tiers []DealerTier, err error) {
	redis_key := "local:tiers"
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
		err = json.Unmarshal(data, &tiers)
		if err == nil {
			return
		}
	}

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
	defer res.Close()
	go redis.Setex(redis_key, tiers, 86400)
	return
}

func GetLocalDealerTypes() (graphics []MapGraphics, err error) {
	redis_key := "local:types"
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
		err = json.Unmarshal(data, &graphics)
		if err == nil {
			return
		}
	}
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
	defer res.Close()
	go redis.Setex(redis_key, graphics, 86400)
	return
}

func GetWhereToBuyDealers(key string) (customers []Customer, err error) {
	redis_key := "dealers:wheretobuy"
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
		err = json.Unmarshal(data, &customers)
		if err == nil {
			return
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return customers, err
	}
	defer db.Close()

	stmt, err := db.Prepare(whereToBuyDealers)
	if err != nil {
		return customers, err
	}
	res, err := stmt.Query()
	if err != nil {
		return customers, err
	}
	defer res.Close()

	for res.Next() {
		var cust Customer
		if err := cust.ScanCustomer(res, key); err != nil {
			return customers, err
		}
		customers = append(customers, cust)
	}

	go redis.Setex(redis_key, customers, 86400)

	return customers, err
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
	for res.Next() {
		loc, err := ScanDealerLocation(res)
		if err != nil {
			return locations, err
		}
		locations = append(locations, *loc)
	}
	defer res.Close()

	return locations, err
}

func SearchLocationsByType(term string) (locations DealerLocations, err error) {
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
	for res.Next() {
		loc, err := ScanDealerLocation(res)
		if err != nil {
			return locations, err
		}

		locations = append(locations, *loc)
	}
	defer res.Close()

	return locations, err
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
	for res.Next() {
		loc, err := ScanDealerLocation(res)
		if err != nil {
			return locations, err
		}
		locations = append(locations, *loc)
	}
	defer res.Close()

	return locations, err
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

//Scan Methods
func (c *Customer) ScanCustomer(res Scanner, key string) error {
	var err error
	var country, countryAbbr, dealerType, dealerTypeOnline, dealerTypeShow, dealerTypeLabel *string
	var dealerTier, dealerTierSort *string
	var logo, web, searchU, icon, shadow, parentId, eLocalUrl *[]byte
	var lat, lon *string
	var mapIconId, countryId *int

	err = res.Scan(
		&c.Id,
		&c.Name,
		&c.Email,
		&c.Address,
		&c.City,
		&c.State.Id,
		&c.Phone,
		&c.Fax,
		&c.ContactPerson,
		&c.DealerType.Id,
		&lat,
		&lon,
		&web,
		&c.CustomerId,
		&c.IsDummy,
		&parentId,
		&searchU,
		&eLocalUrl,
		&logo,
		&c.Address2,
		&c.PostalCode,
		&c.MapixCode.ID,
		&c.SalesRepresentative.ID,
		&c.ApiKey,
		&c.DealerTier.Id,
		&c.ShowWebsite,
		&c.State.State,
		&c.State.Abbreviation,
		&countryId,
		&country,
		&countryAbbr,
		&dealerType,
		&dealerTypeOnline,
		&dealerTypeShow,
		&dealerTypeLabel,
		&dealerTier,
		&dealerTierSort,
		&icon,
		&shadow,
		&c.MapixCode.Code,
		&c.MapixCode.Description,
		&c.SalesRepresentative.Name,
		&c.SalesRepresentative.Code,
	)
	if err != nil {
		return err
	}

	//get parent, if has parent
	parentChan := make(chan int)
	go func() {
		if parentId != nil {
			parentInt, err := conversions.ByteToInt(*parentId)
			if err == nil && parentInt > 0 {
				par := Customer{CustomerId: parentInt}
				if err := par.FindCustIdFromCustomerId(); err != nil {
					parentChan <- 1
					return
				}

				err = par.GetCustomer(key)
				if err != nil {
					parentChan <- 1
					return
				}
				c.Parent = &par
			}
		}
		parentChan <- 1
	}()

	var coun geography.Country
	if lat != nil && *lat != "" && lon != nil && *lon != "" {
		c.Latitude, _ = strconv.ParseFloat(*lat, 64)
		c.Longitude, _ = strconv.ParseFloat(*lon, 64)
	}
	if searchU != nil {
		c.SearchUrl, err = conversions.ByteToUrl(*searchU)
	}
	if eLocalUrl != nil {
		c.ELocalUrl, err = conversions.ByteToUrl(*eLocalUrl)
	}
	if logo != nil {
		c.Logo, err = conversions.ByteToUrl(*logo)
	}
	if web != nil {
		c.Website, err = conversions.ByteToUrl(*web)
	}
	if countryId != nil {
		coun.Id = *countryId
	}
	if country != nil {
		coun.Country = *country
	}
	if countryAbbr != nil {
		coun.Abbreviation = *countryAbbr
	}
	c.State.Country = &coun

	if dealerType != nil {
		c.DealerType.Type = *dealerType
	}
	if dealerTypeOnline != nil {
		c.DealerType.Online, _ = strconv.ParseBool(*dealerTypeOnline)
	}
	if dealerTypeShow != nil {
		c.DealerType.Show, _ = strconv.ParseBool(*dealerTypeShow)
	}
	if dealerTypeLabel != nil {
		c.DealerType.Label = *dealerTypeLabel
	}
	if dealerTier != nil {
		c.DealerTier.Tier = *dealerTier
	}
	if dealerTierSort != nil {
		c.DealerTier.Sort, _ = strconv.Atoi(*dealerTierSort)
	}

	if mapIconId != nil {
		c.DealerType.MapIcon.Id = *mapIconId
	}
	if icon != nil {
		c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(*icon)
	}
	if shadow != nil {
		c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(*shadow)
	}

	<-parentChan

	return nil
}

func ScanCustomerTableFields(res Scanner) (*Customer, error) {
	var c Customer
	var err error
	var name, email, address, address2, city, phone, fax, contactPerson, postalCode, apiKey *string
	var logo, web, searchU, parentId, eLocalUrl *[]byte
	var lat, lon *string
	var stateId, dtypeId, dtierId, custID, mapixCodeID, salesRepID *int
	var isDummy, showWebsite *bool

	err = res.Scan(
		&c.Id,
		&name,
		&email,
		&address,
		&city,
		&stateId,
		&phone,
		&fax,
		&contactPerson,
		&dtypeId,
		&lat,
		&lon,
		&web,
		&custID,
		&isDummy,
		&parentId,
		&searchU,
		&eLocalUrl,
		&logo,
		&address2,
		&postalCode,
		&mapixCodeID,
		&salesRepID,
		&apiKey,
		&dtierId,
		&showWebsite,
	)
	if err != nil {
		return &c, err
	}

	if name != nil {
		c.Name = *name
	}
	if address != nil {
		c.Address = *address
	}
	if address2 != nil {
		c.Address2 = *address2
	}
	if city != nil {
		c.City = *city
	}
	if email != nil {
		c.Email = *email
	}
	if phone != nil {
		c.Phone = *phone
	}
	if fax != nil {
		c.Fax = *fax
	}
	if contactPerson != nil {
		c.ContactPerson = *contactPerson
	}
	if lat != nil && *lat != "" && lon != nil && *lon != "" {
		c.Latitude, _ = strconv.ParseFloat(*lat, 64)
		c.Longitude, _ = strconv.ParseFloat(*lon, 64)
	}
	if searchU != nil {
		c.SearchUrl, err = conversions.ByteToUrl(*searchU)
	}
	if eLocalUrl != nil {
		c.ELocalUrl, err = conversions.ByteToUrl(*eLocalUrl)
	}
	if logo != nil {
		c.Logo, err = conversions.ByteToUrl(*logo)
	}
	if web != nil {
		c.Website, err = conversions.ByteToUrl(*web)
	}
	if custID != nil {
		c.CustomerId = *custID
	}
	if isDummy != nil {
		c.IsDummy = *isDummy
	}
	if postalCode != nil {
		c.PostalCode = *postalCode
	}
	if mapixCodeID != nil {
		c.MapixCode.ID = *mapixCodeID
	}
	if salesRepID != nil {
		c.SalesRepresentative.ID = *salesRepID
	}
	if apiKey != nil {
		c.ApiKey = *apiKey
	}
	if showWebsite != nil {
		c.ShowWebsite = *showWebsite
	}
	if stateId != nil {
		c.State.Id = *stateId
	}
	if dtypeId != nil {
		c.DealerType.Id = *dtypeId
	}
	if dtierId != nil {
		c.DealerTier.Id = *dtierId
	}
	return &c, err
}

func ScanLocation(res Scanner) (*CustomerLocation, error) {
	var l CustomerLocation
	var err error
	var name, email, address, city, phone, fax, contactPerson, state, stateAbbr, country, countryAbbr, postalCode *string
	var lat, lon *float64
	var custId, stateId, countryId *int
	var isPrimary, shippingDefault *bool
	var coun geography.Country

	err = res.Scan(
		&l.Id,
		&name,
		&address,
		&city,
		&stateId,
		&email,
		&phone,
		&fax,
		&lat,
		&lon,
		&custId,
		&contactPerson,
		&isPrimary,
		&postalCode,
		&shippingDefault,
		&state,
		&stateAbbr,
		&countryId,
		&country,
		&countryAbbr,
	)
	if err != nil {
		return &l, err
	}
	if name != nil {
		l.Name = *name
	}
	if email != nil {
		l.Email = *email
	}
	if address != nil {
		l.Address = *address
	}
	if city != nil {
		l.City = *city
	}
	if postalCode != nil {
		l.PostalCode = *postalCode
	}
	if phone != nil {
		l.Phone = *phone
	}
	if fax != nil {
		l.Fax = *fax
	}
	if lat != nil {
		l.Latitude = *lat
	}
	if lon != nil {
		l.Longitude = *lon
	}
	if custId != nil {
		l.CustomerId = *custId
	}
	if contactPerson != nil {
		l.ContactPerson = *contactPerson
	}
	if isPrimary != nil {
		l.IsPrimary = *isPrimary
	}
	if shippingDefault != nil {
		l.ShippingDefault = *shippingDefault
	}
	if stateId != nil {
		l.State.Id = *stateId
	}
	if state != nil {
		l.State.State = *state
	}
	if stateAbbr != nil {
		l.State.Abbreviation = *stateAbbr
	}
	if countryId != nil {
		coun.Id = *countryId
	}
	if country != nil {
		coun.Country = *country
	}
	if countryAbbr != nil {
		coun.Abbreviation = *countryAbbr
	}
	l.State.Country = &coun
	return &l, err
}

func ScanDealerLocation(res Scanner) (*DealerLocation, error) {
	var l DealerLocation
	var err error
	var name, email, address, city, phone, fax, contactPerson, state, stateAbbr, country, countryAbbr, postalCode, mapixCode, mapixDesc, rep, repCode, dtypeType, dtypeLabel, dtierTier *string
	var lat, lon *float64
	var icon, shadow, eLocal, web *[]byte
	var custId, stateId, countryId, dtierSort *int
	var isPrimary, shippingDefault, dtypeOnline, dtypeShow, showWebsite *bool
	var coun geography.Country

	err = res.Scan(
		&l.CustomerLocation.Id,
		&name,
		&address,
		&city,
		&stateId,
		&email,
		&phone,
		&fax,
		&lat,
		&lon,
		&custId,
		&contactPerson,
		&isPrimary,
		&postalCode,
		&shippingDefault,
		&state,
		&stateAbbr,
		&countryId,
		&country,
		&countryAbbr,
		&dtypeType,
		&dtypeOnline,
		&dtypeShow,
		&dtypeLabel,
		&dtierTier,
		&dtierSort,
		&icon,
		&shadow,
		&mapixCode,
		&mapixDesc,
		&rep,
		&repCode,
		&showWebsite,
		&eLocal,
		&web,
	)
	if err != nil {
		return &l, err
	}
	if name != nil {
		l.CustomerLocation.Name = *name
	}
	if email != nil {
		l.CustomerLocation.Email = *email
	}
	if address != nil {
		l.CustomerLocation.Address = *address
	}
	if city != nil {
		l.CustomerLocation.City = *city
	}
	if postalCode != nil {
		l.CustomerLocation.PostalCode = *postalCode
	}
	if phone != nil {
		l.CustomerLocation.Phone = *phone
	}
	if fax != nil {
		l.CustomerLocation.Fax = *fax
	}
	if lat != nil {
		l.CustomerLocation.Latitude = *lat
	}
	if lon != nil {
		l.CustomerLocation.Longitude = *lon
	}
	if custId != nil {
		l.CustomerLocation.CustomerId = *custId
	}
	if contactPerson != nil {
		l.CustomerLocation.ContactPerson = *contactPerson
	}
	if isPrimary != nil {
		l.CustomerLocation.IsPrimary = *isPrimary
	}
	if shippingDefault != nil {
		l.CustomerLocation.ShippingDefault = *shippingDefault
	}
	if stateId != nil {
		l.CustomerLocation.State.Id = *stateId
	}
	if state != nil {
		l.CustomerLocation.State.State = *state
	}
	if stateAbbr != nil {
		l.CustomerLocation.State.Abbreviation = *stateAbbr
	}
	if countryId != nil {
		coun.Id = *countryId
	}
	if country != nil {
		coun.Country = *country
	}
	if countryAbbr != nil {
		coun.Abbreviation = *countryAbbr
	}
	if dtypeType != nil {
		l.DealerType.Type = *dtypeType
	}
	if dtypeOnline != nil {
		l.DealerType.Online = *dtypeOnline
	}
	if dtypeShow != nil {
		l.DealerType.Show = *dtypeShow
	}
	if dtypeLabel != nil {
		l.DealerType.Label = *dtypeLabel
	}
	if dtierSort != nil {
		l.DealerTier.Sort = *dtierSort
	}
	if dtierTier != nil {
		l.DealerTier.Tier = *dtierTier
	}
	if icon != nil {
		l.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(*icon)
	}
	if shadow != nil {
		l.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(*shadow)
	}
	if mapixCode != nil {
		l.MapixCode.Code = *mapixCode
	}
	if mapixDesc != nil {
		l.MapixCode.Description = *mapixDesc
	}
	if rep != nil {
		l.SalesRepresentative.Name = *rep
	}
	if repCode != nil {
		l.SalesRepresentative.Code = *repCode
	}
	if showWebsite != nil {
		l.CustomerLocation.ShowWebSite = *showWebsite
	}
	if eLocal != nil {
		l.CustomerLocation.ELocalUrl, err = conversions.ByteToUrl(*eLocal)
	}
	if web != nil {
		l.CustomerLocation.Website, err = conversions.ByteToUrl(*web)
	}

	l.CustomerLocation.State.Country = &coun
	return &l, err
}
