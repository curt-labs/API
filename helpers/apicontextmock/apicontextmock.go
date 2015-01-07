package apicontextmock

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/site"

	"database/sql"
	"log"
	"strings"
)

const (
	customerUserFields = ` cu.id, cu.name, cu.email, cu.password, cu.customerID, cu.date_added, cu.active, cu.locationID, cu.isSudo, cu.cust_ID, cu.NotCustomer, cu.passwordConverted `
	customerFields     = ` c.cust_id, c.name, c.email, c.address,  c.city, c.stateID, c.phone, c.fax, c.contact_person, c.dealer_type,
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
	userCustomer = `select ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
						from Customer as c
						join CustomerUser as cu on c.cust_id = cu.cust_ID
						left join States as s on c.stateID = s.stateID
						left join Country as cty on s.countryID = cty.countryID
						left join DealerTypes as dt on c.dealer_type = dt.dealer_type
						left join MapIcons as mi on dt.dealer_type = mi.dealer_type
						left join DealerTiers dtr on c.tier = dtr.ID
						left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
						left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
						where cu.id = ?`

	getRegisteredUsersId = `select cu.id from CustomerUser as cu
								where cu.email = ? && cu.password = ?
								limit 1`

	customerUserAuth = `select cu.id, cu.name, cu.email, cu.password, cu.customerID, cu.date_added, cu.active,cu.locationID, cu.isSudo, cu.cust_ID, cu.passwordConverted from CustomerUser as cu
							where email = ?
							&& active = 1
							limit 1`
	getUserPassword        = `SELECT password, COUNT(password) AS quantity from CustomerUser where email = ?`
	updateCustomerUserPass = `update CustomerUser set password = ?, passwordConverted = 1
								where id = ? && active = 1`
	customerUserKeyAuth = `select cu.* from CustomerUser as cu
								join ApiKey as ak on cu.id = ak.user_id
								join ApiKeyType as akt on ak.type_id = akt.id
								where UPPER(akt.type) != ?
								&& ak.api_key = UPPER(?)
								&& cu.active = 1 && ak.date_added >= ?`
	customerUserKeys = `select ak.api_key, akt.type, ak.date_added from ApiKey as ak
								join ApiKeyType as akt on ak.type_id = akt.id
								where user_id = ? && UPPER(akt.type) NOT IN (?)`
	userLocation = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
									cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
									cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
									s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
									from CustomerUser as cu
									join CustomerLocations as cl on cu.locationID = cl.locationID
									left join States as s on cl.stateID = s.stateID
									left join Country as cty on s.countryID = cty.countryID
									where cu.id = ?`

	userAuthenticationKey = `select ak.api_key, akt.type, akt.id, CAST(ak.date_added as char(255)) as date_added from ApiKey as ak
									join ApiKeyType as akt on ak.type_id = akt.id
									where UPPER(akt.type) = ?
									&& ak.user_id = ?`

	resetUserAuthentication = `update ApiKey as ak
									set ak.date_added = NOW()
									where ak.type_id = ?
									&& ak.user_id = ?`

	customerUserFromKey = `select cu.* from CustomerUser as cu
								join ApiKey as ak on cu.id = ak.user_id
								join ApiKeyType as akt on ak.type_id = akt.id
								where UPPER(akt.type) != ? && UPPER(ak.api_key) = UPPER(?)
								limit 1`

	customerUserFromId = `select cu.* from CustomerUser as cu
							join ApiKey as ak on cu.id = ak.user_id
							join ApiKeyType as akt on ak.type_id = akt.id
							where cu.id = ?
							limit 1`

	insertCustomerUser = `INSERT into CustomerUser(id, name, email, password, customerID, date_added, active, locationID, isSudo, cust_ID, NotCustomer, passwordConverted)
							VALUES(UUID(),?,?,?,?,NOW(),?,?,?,?,?,1)`

	insertAPIKey = `insert into ApiKey(user_id, type_id, api_key, date_added)
						values(?,?,UUID(),NOW())` //DB schema DOES auto increment table id
	insertAPIKeyToBrand = `insert into ApiKeyToBrand(keyID, brandID)
						values(?,?)`
	deleteAPIKeyToBrand = `delete from ApiKeyToBrand where keyID = (select id from ApiKey where user_id = ? && type_id = ?)`

	getCustomerUserKeysWithoutAuth = `select ak.api_key, akt.type from ApiKey as ak
										join ApiKeyType as akt on ak.type_id = akt.id
										where ak.user_id = ? && UPPER(akt.type) = ?`
	getAPIKeyTypeID             = `select id from ApiKeyType where UPPER(type) = UPPER(?) limit 1`
	setCustomerUserPassword     = `update CustomerUser set password = ?, passwordConverted = 1 WHERE email = ?`
	deleteCustomerUser          = `DELETE FROM CustomerUser WHERE id = ?`
	deleteAPIkey                = `DELETE FROM ApiKey WHERE user_id = ? AND type_id = ?`
	getCustomerUserKeysWithAuth = `select ak.api_key, akt.type from ApiKey as ak
										join ApiKeyType as akt on ak.type_id = akt.id
										where ak.user_id = ? && (UPPER(akt.type) = ? || UPPER(akt.type) = ? || UPPER(akt.type) = ?)`
	getCustomerUserLocation = `select cl.locationID, cl.name, cl.address, cl.city,
								s.stateID, s.state,
								s.abbr, cun.countryID, cun.name as countryName, cun.abbr as countryAbbr,
								cl.email, cl.phone, cl.fax, cl.latitude, cl.longitude,
								cl.cust_id, cl.contact_person, cl.isprimary, cl.postalCode,
								cl.ShippingDefault from CustomerLocations as cl
								join CustomerUser as cu on cl.locationID = cu.locationID
								left join States as s on cl.stateID = s.stateID
								left join Country as cun on s.countryID = cun.countryID
								where cu.id = ?
								limit 1`
	getCustomerUserBrands = `select b.ID, b.name, b.code 
							from Brand as b
							join CustomerToBrand as ctb on ctb.BrandID = b.ID
							join Customer as c on c.cust_id = ctb.cust_id
							where c.cust_id = ?`

	updateCustomerUser   = `UPDATE CustomerUser SET name = ?, email = ?, active = ?, locationID = ?, isSudo = ?, NotCustomer = ? WHERE id = ?`
	getUsersByCustomerID = `SELECT id FROM CustomerUser WHERE cust_id = ?`
)

func Mock() (*apicontext.DataContext, error) {
	var dtx apicontext.DataContext
	var c customer.Customer
	var cu customer.CustomerUser
	var w site.Website
	c.Name = "test cust"
	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		//setup apiKeyTypes
		pub.Type = "Public"
		pri.Type = "Private"
		auth.Type = "Authentication"
		pub.Create()
		pri.Create()
		auth.Create()
	}
	c.Create()
	cu.CustomerID = c.Id
	cu.Name = "test cust content user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var err error
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	err = w.Create()
	w.BrandIDs = append(w.BrandIDs, 1)
	err = w.JoinToBrand()
	if err != nil {
		return &dtx, err
	}

	dtx.WebsiteID = w.ID
	dtx.BrandID = 1
	dtx.APIKey = apiKey
	dtx.CustomerID = c.Id
	dtx.UserID = cu.Id
	return &dtx, err
}

func DeMock(dtx *apicontext.DataContext) error {
	var err error
	var cust customer.Customer
	var user customer.CustomerUser

	cust.Id = dtx.CustomerID
	user.Id = dtx.UserID

	var w site.Website

	w.ID = dtx.WebsiteID
	w.Get()

	err = w.Delete()

	err = cust.Delete()

	err = user.Delete()

	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		for _, key := range user.Keys {
			if strings.ToLower(key.Type) == "public" {
				pub.ID = key.TypeId
			}
			if strings.ToLower(key.Type) == "private" {
				pri.ID = key.TypeId
			}
			if strings.ToLower(key.Type) == "authentication" {
				auth.ID = key.TypeId
			}
		}
		err = pub.Delete()
		err = pri.Delete()
		err = auth.Delete()

	}

	return err
}

func Mock2() (*apicontext.DataContext, error) {
	var dtx apicontext.DataContext
	// Needs to create records in the db for the following because of foreign key constraints:
	// Bare Min:
	// CustomerUser
	var err error
	CustomerUserID := ""
	if CustomerUserID, err = CreateCustomerUser(); err != nil {
		return &dtx, err
	}
	log.Println(CustomerUserID)
	dtx.UserID = CustomerUserID
	dtx.CustomerID = 1

	// ApiKeyType
	// ApiKey
	// ApiKeyToBrand
	// Brand
	return &dtx, nil
}

func CreateCustomerUser() (CustomerUserID string, err error) {
	encryptPass := "bogus" // Not encrypted
	email := "TestBogus@curtmfg.com"
	customerID := 1
	custID := 1

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return CustomerUserID, err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertCustomerUser)
	if err != nil {
		return CustomerUserID, err
	}
	_, err = stmt.Exec("Test", "TestBogus@curtmfg.com", encryptPass, customerID, true, 0, true, custID, true)
	if err != nil {
		tx.Rollback()
		return CustomerUserID, err
	}
	if err = tx.Commit(); err != nil {
		return CustomerUserID, err
	}

	stmt, err = db.Prepare(getRegisteredUsersId) // needs to be set on the customer user object in order to generate the keys
	if err != nil {
		return CustomerUserID, err
	}

	var userID *string
	if err = stmt.QueryRow(email, encryptPass).Scan(&userID); err != nil || userID == nil {
		return CustomerUserID, err
	}

	CustomerUserID = *userID

	return CustomerUserID, nil
}
