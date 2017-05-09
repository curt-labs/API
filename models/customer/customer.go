package customer

import (
	"github.com/curt-labs/API/helpers/api"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/conversions"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/geography"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

//TODO: Clean up these monstrosities of scan functions. Some of these are like this
//due to the fact that too many fiels in our DB allow NULL. If we were to fix that
//then this file could be much, much shorter, and much more straightforward.

type Coordinates struct {
	Latitude  float64 `json:"latitude" xml:"latitude"`
	Longitude float64 `json:"longitude" xml:"longitude"`
}

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
	Latitude            float64             `json:"coords,latitude,omitempty" xml:"latitude,omitempty"`
	Longitude           float64             `json:"coords,longitude,omitempty" xml:"longitude,omitempty"`
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
	BrandIDs            []int               `json:"brandIds,omitempty" xml:"brandIds,omitempty"`
	Accounts            []Account           `json:"accounts,omitempty" xml:"accounts,omitempty"`
	ShippingInfo        ShippingInfo        `json:"shippingInfo,omitempty" xml:"shippingInfo,omitempty"`
}

type Customers []Customer

type Scanner interface {
	Scan(...interface{}) error
}

type Account struct {
	ID                 int          `json:"id,omitempty" xml:"id,omitempty"`
	AccountNumber      string       `json:"accountNumber,omitempty" xml:"accountNumber,omitempty"`
	Cust_id            int          `json:"cust_id,omitempty" xml:"cust_id,omitempty"`
	TypeID             int          `json:"-" xml:"-"`
	FreightLimit       float64      `json:"freightLimit,omitempty" xml:"freightLimit,omitempty"`
	DefaultWarehouseID int          `json:"defaultWarehouseId,omitempty" xml:"defaultWarehouseId,omitempty"`
	Type               AccountType  `json:"type,omitempty" xml:"type,omitempty"`
	ShippingInfo       ShippingInfo `json:"shipping_info,omitempty" xml"shipping_info,omitempty"`
}

type AccountType struct {
	ID        int     `json:"id,omitempty" xml:"id,omitempty"`
	Title     string  `json:"title,omitempty" xml:"title,omitempty"`
	ComnetURL url.URL `json:"comnetURL,omitempty" xml:"comnetURL,omitempty"`
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
	Coordinates     Coordinates     `json:"coords,omitempty" xml:"coords,omitempty"`
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
	Id      int    `json:"id,omitempty" xml:"id,omitempty"`
	Tier    string `json:"tier,omitempty" xml:"tier,omitempty"`
	Sort    int    `json:"sort,omitempty" xml:"sort,omitempty"`
	BrandID int    `json:"brandId,omitempty" xml:"brandId,omitempty"`
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
	CustomerLocation
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

type DealersResponse struct {
	Items []DealerLocation `json:"items" xml:"items"`
	Total int              `json:"total" xml:"total"`
}

type EtailerResponse struct {
	Items []Customer `json:"items" xml:"items"`
	Total int        `json:"total" xml:"total"`
}

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
				c.latitude, c.longitude,  c.website, c.customerID, c.isDummy, c.parentID, c.searchURL, c.eLocalURL, c.logo,c.address2,
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
	showSiteFields = ` c.showWebsite, c.website, c.elocalurl `

	//redis
	custPrefix = "customer:"
)

var (
	getCustomer              = `select ` + customerFields + ` from Customer as c where c.cust_id = ? `
	getCustomerIdFromKeyStmt = `select c.customerID from Customer as c
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`

	getCustIdFromKeyStmt = `select cu.cust_ID from CustomerUser as cu
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`
	getCustIdsFromAccountNumStmt = `select c.cust_id, c.customerID from Customer as c
										Join Accounts as a on a.cust_id = c.cust_id
										where a.accountNumber = ? limit 1`
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
	customerAccounts = `select act.id, act.accountNumber, act.cust_id, act.type_id, act.freightLimit, acty.type, acty.comnet_url, act.defaultWarehouseId from Accounts as act
							Join AccountTypes as acty on acty.id = act.type_id
							where act.cust_id = ?`

	customerUser = `select cu.id, cu.name, cu.email, cu.customerID, cu.date_added, cu.active,cu.locationID, cu.isSudo, cu.cust_ID from CustomerUser as cu
						join Customer as c on cu.cust_ID = c.cust_id
						where c.cust_id = ?
						order by cu.name`
	customerPrice = `select distinct cp.price from
						 CustomerPricing cp
						 where cp.cust_ID =  ?
						 and cp.partID = ?`

	customerPart = `select distinct ci.custPartID from ApiKey as ak
						join CustomerUser cu on ak.user_id = cu.id
						join Customer c on cu.cust_ID = c.cust_id
						join CartIntegration ci on c.cust_ID = ci.custID
						where ak.api_key = ?
						and ci.partID = ?`

	etailers = `select distinct
	      ` + customerFields + `,
	      ` + stateFields + `,
	      ` + countryFields + `,
	      ` + dealerTypeFields + `,
	      ` + dealerTierFields + `,
	      ` + mapIconFields + `,
	      ` + mapixCodeFields + `,
	      ` + salesRepFields + `
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				join CustomerToBrand as ctb on ctb.cust_id = c.cust_id
				join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
				join ApiKey as a on a.id = atb.keyID
				where dt.online = 1 && c.isDummy = 0
				&& a.api_key = ? && (ctb.brandID = ? or 0 = ?)
				order by c.name
				limit ?,?`

	etailersCount = `select count(*)
				from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				join CustomerToBrand as ctb on ctb.cust_id = c.cust_id
				join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
				join ApiKey as a on a.id = atb.keyID
				where dt.online = 1 && c.isDummy = 0
				&& a.api_key = ? && (ctb.brandID = ? or 0 = ?)`

	localDealers = `select
					` + customerLocationFields + `,
					` + stateFields + `,
					` + countryFields + `,
					` + dealerTypeFields + `,
					` + dealerTierFields + `,
					` + mapIconFields + `,
					` + mapixCodeFields + `,
					` + salesRepFields + ` ,
					` + showSiteFields + `,(
						? * acos(
							cos(
								radians(?) ) * cos( radians( cl.latitude )
							) * cos(
								radians( cl.longitude ) - radians(?)
							) + sin(
								radians(?)
							) * sin(
								radians( cl.latitude )
							)
						)
					) as distance
					from CustomerLocations as cl
					join Customer as c on cl.cust_id = c.cust_id
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					left join MapIcons as mi on dt.dealer_type = mi.dealer_type
					join DealerTiers as dtr on c.tier = dtr.ID
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
					left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
					where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier && mi.brandID = ?
					having (distance < ?) || (? = 0)
					order by cl.locationID
					limit ?,?`

	localDealersNoDistance = `select
					` + customerLocationFields + `,
					` + stateFields + `,
					` + countryFields + `,
					` + dealerTypeFields + `,
					` + dealerTierFields + `,
					` + mapIconFields + `,
					` + mapixCodeFields + `,
					` + salesRepFields + ` ,
					` + showSiteFields + `
					from CustomerLocations as cl
					join Customer as c on cl.cust_id = c.cust_id
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					left join MapIcons as mi on dt.dealer_type = mi.dealer_type
					join DealerTiers as dtr on c.tier = dtr.ID
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
					left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
					where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier && mi.brandID = ?
					order by cl.locationID
					limit ?,?`

	countDealers = `select count(*)
					from CustomerLocations as cl
					join Customer as c on cl.cust_id = c.cust_id
					join DealerTypes as dt on c.dealer_type = dt.dealer_type
					left join MapIcons as mi on dt.dealer_type = mi.dealer_type
					join DealerTiers as dtr on c.tier = dtr.ID
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
					left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
					where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier && mi.brandID = ?`

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
							join CustomerToBrand as ctb on ctb.cust_id = c.cust_id
							join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
							join ApiKey as a on a.id = atb.keyID
							where dt.online = false and dt.show = true
							&& a.api_key = ? && (ctb.brandID = ? or 0 = ?)
							order by dtr.sort`
	localDealerTypes = `select distinct m.ID as iconId, m.mapicon, m.mapiconshadow,
							dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
							dt.dealer_type as dealerTypeId, dt.type as dealerType, dt.online, dt.show, dt.label
							from MapIcons as m
							join DealerTypes as dt on m.dealer_type = dt.dealer_type
							join DealerTiers as dtr on m.tier = dtr.ID
							join Customer as c on dtr.ID = c.tier
							join CustomerToBrand as ctb on ctb.cust_id = c.cust_id
							join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
							join ApiKey as a on a.id = atb.keyID
							where dt.show = true
							&& a.api_key = ? && (atb.brandID = ? or 0 = ?)
							order by dtr.sort desc`

	whereToBuyDealers = `select distinct ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
			from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				join CustomerToBrand as ctb on ctb.cust_id = c.cust_id
				join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
				join ApiKey as a on a.id = atb.keyID
				where c.dealer_type = 1 and c.tier = 4 and c.isDummy = false and length(c.searchURL) > 1
				&&(a.api_key = ? && (atb.brandID = ? or 0 = ?))`

	customerByLocation = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `  ,` + showSiteFields + `
								from CustomerLocations as cl
								join States as s on cl.stateID = s.stateID
								left join Country as cty on cty.countryID = s.countryID
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
								left join Country as cty on cty.countryID = s.countryID
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
	deleteCustomer   = `delete from Customer where cust_id = ?`
	joinUser         = `update CustomerUser set cust_ID = ? where id = ?`
	createDealerType = `insert into DealerTypes (type, online, label) values(?,?,?)`
	deleteDealerType = `delete from DealerTypes where dealer_type= ?`
)

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
	c.GetAccounts()
	err = <-basicsChan

	if err == sql.ErrNoRows {

		err = fmt.Errorf("error: %s", "failed to retrieve")
	}
	return err
}

//gets cust_id, not customerId
func (c *Customer) GetCustomerIdFromKey(key string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustIdFromKeyStmt)
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

func (c *Customer) GetCustomerIdsFromAccountNumber(accountNum string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustIdsFromAccountNumStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var cust_id, customerID *int
	err = stmt.QueryRow(accountNum).Scan(&cust_id, &customerID)
	if err != nil {
		return err
	}
	if cust_id != nil {
		c.Id = *cust_id
	}
	if customerID != nil {
		c.CustomerId = *customerID
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

func (c *Customer) GetAccounts() (err error) {

	redis_key := "CustAccount:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.Accounts)
		return err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerAccounts)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(c.Id)

	var accts []Account
	for rows.Next() {
		acc, err := ScanAccount(rows)
		if err != nil {
			return err
		}

		accts = append(accts, *acc)
	}
	defer rows.Close()

	c.Accounts = accts

	redis.Setex(redis_key, c.Accounts, redis.CacheTimeout)

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
	if err != nil {
		return err
	}
	c.Id = int(id)

	for _, brandID := range c.BrandIDs {
		err = c.CreateCustomerBrand(brandID)
		if err != nil {
			return err
		}
	}
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
	err = c.DeleteAllCustomerBrands()
	if err != nil {
		return err
	}
	for _, brandID := range c.BrandIDs {
		err = c.CreateCustomerBrand(brandID)
	}
	go redis.Set(custPrefix+strconv.Itoa(c.Id), c)
	go redis.Delete("customerLocations:" + strconv.Itoa(c.Id))
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
	err = c.DeleteAllCustomerBrands()
	if err != nil {
		return err
	}
	go redis.Delete(custPrefix + strconv.Itoa(c.Id))
	go redis.Delete("customerLocations:" + strconv.Itoa(c.Id))
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
	if err != nil {
		return err
	}
	iter := 0
	userChan := make(chan error)
	lowerKey := strings.ToLower(key)

	for res.Next() {
		var u CustomerUser
		err = res.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.CustomerID,
			&u.DateAdded,
			&u.Active,
			&u.Location.Id,
			&u.Sudo,
			&u.CustID,
		)
		if err != nil {
			continue
		}
		go func(user CustomerUser) {
			if err := user.GetKeys(); err == nil {
				for _, k := range user.Keys {
					if k.Key == lowerKey {
						user.Current = true
						break
					}
				}
			}

			user.Brands, err = brand.GetUserBrands(c.Id)
			if err != nil {
				userChan <- nil
				return
			}

			user.GetComnetAccounts()
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

func (c *Customer) JoinUser(u CustomerUser) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(joinUser)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id, u.Id)
	if err != nil {
		return err
	}
	return err
}

func GetCustomerPrice(dtx *apicontext.DataContext, part_id int) (price float64, err error) {
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

	err = stmt.QueryRow(dtx.CustomerID, part_id).Scan(&price)
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
	defer stmt.Close()

	err = stmt.QueryRow(api_key, part_id).Scan(&ref)
	return ref, err
}

func GetEtailers(dtx *apicontext.DataContext, count int, page int) (EtailerResponse, error) {
	redis_key := "dealers:etailer:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	var dealers []Customer
	var etailResp EtailerResponse
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &dealers)
		if err != nil {
			return etailResp, err
		}
		return etailResp, err
	}

	skip := (page - 1) * count

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return etailResp, err
	}
	defer db.Close()

	var total int

	row := db.QueryRow(etailersCount, dtx.APIKey, dtx.BrandID, dtx.BrandID)
	err = row.Scan(&total)
	if err != nil {
		return etailResp, err
	}

	rows, err := db.Query(etailers, dtx.APIKey, dtx.BrandID, dtx.BrandID, skip, count)
	if err != nil {
		return etailResp, err
	}
	defer rows.Close()

	for rows.Next() {
		var cust Customer
		if err := cust.ScanCustomer(rows, dtx.APIKey); err == nil {
			dealers = append(dealers, cust)
		}
	}
	redis.Setex(redis_key, dealers, 86400)

	etailResp = EtailerResponse{Items: dealers, Total: total}

	return etailResp, err
}

func GetLocalDealers(latlng string, distance int, skip int, count int, brandID int) (DealersResponse, error) {
	var err error
	var dealers []DealerLocation
	var dealerResp DealersResponse

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dealerResp, err
	}
	defer db.Close()

	var total int

	row := db.QueryRow(countDealers, brandID)
	row.Scan(&total)

	var latitude string
	var longitude string
	var res *sql.Rows

	// Get the boundary points
	if latlng != "" {
		latlngs := strings.Split(latlng, ",")
		if len(latlngs) != 2 {
			err = fmt.Errorf("%s", "failed to parse the latitude and longitude")
			return dealerResp, err
		}
		latitude = latlngs[0]
		longitude = latlngs[1]
	}

	if latlng == "" {
		res, err = db.Query(localDealersNoDistance, brandID, skip, count)
		if err != nil {
			return dealerResp, err
		}
	} else {
		res, err = db.Query(localDealers, api_helpers.EARTH, latitude, longitude, latitude, distance, distance, brandID, skip, count)
		if err != nil {
			return dealerResp, err
		}
	}
	defer res.Close()

	for res.Next() {
		cols, err := res.Columns()
		if err != nil {
			return dealerResp, err
		}
		var l *DealerLocation
		if latlng != "" {
			l, err = ScanDealerLocation(res, len(cols), true)
		} else {
			l, err = ScanDealerLocation(res, len(cols), false)
		}

		if err != nil {
			return dealerResp, err
		}

		dealers = append(dealers, *l)
	}

	if latlng != "" {
		sortutil.AscByField(dealers, "Distance")
	}

	dealerResp = DealersResponse{Items: dealers, Total: total}
	return dealerResp, err
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
	defer stmtPolygon.Close()
	stmtCoordinates, err := db.Prepare(MapPolygonCoordinatesForState)
	if err != nil {
		return regions, err
	}
	defer stmtCoordinates.Close()
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

func GetLocalDealerTiers(dtx *apicontext.DataContext) (tiers []DealerTier, err error) {
	redis_key := "local:tiers:" + dtx.BrandString
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
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	var brandID *int
	for res.Next() {
		var t DealerTier
		err = res.Scan(&t.Id, &t.Tier, &t.Sort, &brandID)
		if err != nil {
			return tiers, err
		}
		tiers = append(tiers, t)
	}
	defer res.Close()
	go redis.Setex(redis_key, tiers, 86400)
	return
}

func GetLocalDealerTypes(dtx *apicontext.DataContext) (graphics []MapGraphics, err error) {
	redis_key := "local:types:" + dtx.BrandString
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
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
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

func GetWhereToBuyDealers(dtx *apicontext.DataContext) (customers []Customer, err error) {
	redis_key := "dealers:wheretobuy:" + dtx.BrandString
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
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if err != nil {
		return customers, err
	}
	defer res.Close()

	for res.Next() {
		var cust Customer
		if err := cust.ScanCustomer(res, dtx.APIKey); err != nil {
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
	defer stmt.Close()
	term = "%" + term + "%"

	res, err := stmt.Query(term, term)
	if err != nil {
		return locations, err
	}

	cols, err := res.Columns()
	if err != nil {
		return locations, err
	}

	for res.Next() {
		loc, err := ScanDealerLocation(res, len(cols), true)
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
	defer stmt.Close()
	term = "%" + term + "%"

	res, err := stmt.Query(term, term)
	if err != nil {
		return locations, err
	}

	cols, err := res.Columns()
	if err != nil {
		return locations, err
	}

	for res.Next() {
		loc, err := ScanDealerLocation(res, len(cols), true)
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
	defer stmt.Close()
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

	cols, err := res.Columns()
	if err != nil {
		return locations, err
	}

	for res.Next() {
		loc, err := ScanDealerLocation(res, len(cols), true)
		if err != nil {
			return locations, err
		}
		locations = append(locations, *loc)
	}
	defer res.Close()

	return locations, err
}

//Dealer Types
func (d *DealerType) Create() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createDealerType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(d.Type, d.Online, d.Label)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = int(id)
	return err
}

func (d *DealerType) Delete(dtx *apicontext.DataContext) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteDealerType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(d.Id)
	go redis.Delete("local:types:" + dtx.BrandString)
	go redis.Delete("dealers:etailer:" + dtx.BrandString)
	return err
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
	var email, address, city, phone, fax, contact, address2, postal, api, state, stateAb *string
	var stateId, dType, mapixId, salesRepId, dTier, customerid *int
	var show, isdummy *bool

	err = res.Scan(
		&c.Id,
		&c.Name,
		&email,
		&address,
		&city,
		&stateId,
		&phone,
		&fax,
		&contact,
		&dType,
		&lat,
		&lon,
		&web,
		&customerid,
		&isdummy,
		&parentId,
		&searchU,
		&eLocalUrl,
		&logo,
		&address2,
		&postal,
		&mapixId,
		&salesRepId,
		&api,
		&dTier,
		&show,
		&state,
		&stateAb,
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
	<-parentChan
	if city != nil {
		c.City = *city
	}
	if address != nil {
		c.Address = *address
	}
	if stateId != nil {
		c.State.Id = *stateId
	}
	if phone != nil {
		c.Phone = *phone
	}
	if fax != nil {
		c.Fax = *fax
	}
	if email != nil {
		c.Email = *email
	}
	if contact != nil {
		c.ContactPerson = *contact
	}
	if dType != nil {
		c.DealerType.Id = *dType
	}
	if customerid != nil {
		c.CustomerId = *customerid
	}
	if isdummy != nil {
		c.IsDummy = *isdummy
	}
	if address2 != nil {
		c.Address2 = *address2
	}
	if postal != nil {
		c.PostalCode = *postal
	}
	if api != nil {
		c.ApiKey = *api
	}
	if mapixId != nil {
		c.MapixCode.ID = *mapixId
	}
	if salesRepId != nil {
		c.SalesRepresentative.ID = *salesRepId
	}
	if dTier != nil {
		c.DealerTier.Id = *dTier
	}
	if state != nil {
		c.State.State = *state
	}
	if stateAb != nil {
		c.State.Abbreviation = *stateAb
	}

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

func ScanAccount(res Scanner) (*Account, error) {
	var a Account
	var err error

	var accID *int
	var accountNumber *string
	var cust_id *int
	var typeID *int
	var typeText *string
	var comnetURL *[]byte
	var freightLimit *float64
	var defaultWare *int

	err = res.Scan(
		&accID,
		&accountNumber,
		&cust_id,
		&typeID,
		&freightLimit,
		&typeText,
		&comnetURL,
		&defaultWare,
	)
	if err != nil {
		return &a, err
	}
	if accID != nil {
		a.ID = *accID
	}
	if accountNumber != nil {
		a.AccountNumber = *accountNumber
	}
	if cust_id != nil {
		a.Cust_id = *cust_id
	}
	if typeID != nil {
		a.TypeID = *typeID
		a.Type.ID = a.TypeID
	}
	if freightLimit != nil {
		a.FreightLimit = *freightLimit
	}
	if typeText != nil {
		a.Type.Title = *typeText
	}
	if comnetURL != nil {
		a.Type.ComnetURL, err = conversions.ByteToUrl(*comnetURL)
	}
	if defaultWare != nil {
		a.DefaultWarehouseID = *defaultWare
	}

	return &a, err
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
		l.Coordinates.Latitude = *lat
	}
	if lon != nil {
		l.Coordinates.Longitude = *lon
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

func ScanDealerLocation(res *sql.Rows, count int, isDistance bool) (*DealerLocation, error) {
	var l DealerLocation
	var err error
	var mapIconString string
	var mapIconShadowString string
	var websiteString string
	var elocalString string

	if l.State.Country == nil {
		l.State.Country = &geography.Country{}
	}

	if isDistance == true {
		err = res.Scan(&l.CustomerLocation.Id, &l.Name, &l.Address, &l.City, &l.State.Id,
			&l.Email, &l.Phone, &l.Fax, &l.Coordinates.Latitude, &l.Coordinates.Longitude,
			&l.CustomerId, &l.ContactPerson, &l.IsPrimary, &l.PostalCode, &l.ShippingDefault,
			&l.State.State, &l.State.Abbreviation, &l.State.Country.Id, &l.State.Country.Country,
			&l.State.Country.Abbreviation, &l.DealerType.Type, &l.DealerType.Online,
			&l.DealerType.Show, &l.DealerType.Label, &l.DealerTier.Tier, &l.DealerTier.Sort,
			&mapIconString, &mapIconShadowString, &l.MapixCode.Code, &l.MapixCode.Description,
			&l.SalesRepresentative.Name, &l.SalesRepresentative.Code, &l.ShowWebSite,
			&websiteString, &elocalString, &l.Distance)

		if err != nil {
			return &l, err
		}
	}

	if isDistance == false {
		err = res.Scan(&l.CustomerLocation.Id, &l.Name, &l.Address, &l.City, &l.State.Id,
			&l.Email, &l.Phone, &l.Fax, &l.Coordinates.Latitude, &l.Coordinates.Longitude,
			&l.CustomerId, &l.ContactPerson, &l.IsPrimary, &l.PostalCode, &l.ShippingDefault,
			&l.State.State, &l.State.Abbreviation, &l.State.Country.Id, &l.State.Country.Country,
			&l.State.Country.Abbreviation, &l.DealerType.Type, &l.DealerType.Online,
			&l.DealerType.Show, &l.DealerType.Label, &l.DealerTier.Tier, &l.DealerTier.Sort,
			&mapIconString, &mapIconShadowString, &l.MapixCode.Code, &l.MapixCode.Description,
			&l.SalesRepresentative.Name, &l.SalesRepresentative.Code, &l.ShowWebSite,
			&websiteString, &elocalString)

		if err != nil {
			return &l, err
		}
	}

	mapIconURL, err := url.Parse(mapIconString)
	mapIconShadowURL, err := url.Parse(mapIconShadowString)
	websiteURL, err := url.Parse(websiteString)
	eLocalURL, err := url.Parse(elocalString)
	if err != nil {
		return &l, err
	}

	l.DealerType.MapIcon.MapIcon = *mapIconURL
	l.DealerType.MapIcon.MapIconShadow = *mapIconShadowURL
	l.Website = *websiteURL
	l.ELocalUrl = *eLocalURL

	return &l, err
}
