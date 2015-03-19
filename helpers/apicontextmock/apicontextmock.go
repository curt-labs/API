package apicontextmock

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"

	"time"
)

const (
	// time format
	timeFormat = "2006-01-02 03:04:05"
)

var (
	getRegisteredUsersId = `select cu.id from CustomerUser as cu
								where cu.email = ? && cu.password = ?
								limit 1`

	insertCustomerUser = `INSERT into CustomerUser(id, name, email, password, customerID, date_added, active, locationID, isSudo, cust_ID, NotCustomer, passwordConverted)
							VALUES(UUID(),?,?,?,?,NOW(),?,?,?,?,?,1)`

	insertAPIKey = `insert into ApiKey(user_id, type_id, api_key, date_added)
						values(?,?,UUID(),NOW())` //DB schema DOES auto increment table id
	insertAPIKeyToBrand = `insert into ApiKeyToBrand(keyID, brandID)
						values(?,?)`
	insertCustomer = `insert into Customer(name,email, address, city, stateID, phone,fax, contact_person, dealer_type, latitude, longitude, password, website, customerID, isDummy, parentID, searchURL, eLocalURL, logo, address2, postal_code, mCodeID, salesRepID, APIKey, tier,showWebsite)
			values('test','w','w','w',1,'w','w','w',1,'w','w','w','w',1,1,1,'w','w','w','w','w',1,1,'w',1,1)`
	insertCustomerToBrand = `insert into CustomerToBrand (cust_id, brandID) values (?,?)`
	deleteCustomer        = `delete from Customer where cust_id = ?`
	deleteCustomerToBrand = `delete from CustomerToBrand where cust_id = ?`

	getCustomerUserKeysWithoutAuth = `select ak.api_key, akt.type from ApiKey as ak
										join ApiKeyType as akt on ak.type_id = akt.id
										where ak.user_id = ? && UPPER(akt.type) = ?`
	getUserApiKeys = `select api_key from ApiKey where user_id = ?`

	getKeyByDateType = "SELECT id FROM ApiKeyType WHERE type = ?  AND date_added = ?"
	createApiKeyType = "INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),?,?)"

	createSite                 = `INSERT INTO Website (url, description) VALUES (?,?)`
	joinToBrand                = `insert into WebsiteToBrand (WebsiteID, brandID) values (?,?)`
	deleteApiKeyType           = `delete from ApiKeyType where id = (select type_id from ApiKey where api_key = ?)`
	deleteType                 = `delete from ApiKeyType where type = (select type_id from ApiKey where api_key =?)`
	deleteCustomerUser         = `delete from CustomerUser where id = ?`
	deleteApiKey               = `delete from ApiKey where api_key =  ?`
	deleteApiKeyToBrand        = `delete from ApiKeyToBrand where keyID = (select id from ApiKey where api_key = ?)`
	deleteSite                 = `delete from Website where id = ?`
	deleteSiteToBrand          = `delete from WebsiteToBrand where WebsiteID = ?`
	deleteUsersApiKeys         = `delete from ApiKey where user_id = ?`
	deleteUsersApiKeysToBrands = `delete from ApiKeyToBrand where keyID in (select id from ApiKey where user_id = ?)`
)

func Mock() (*apicontext.DataContext, error) {
	var dtx apicontext.DataContext
	var err error
	// Needs to create records in the db for the following because of foreign key constraints:
	// Bare Min:
	dtx.CustomerID, err = InsertCustomer()
	if err != nil {
		return &dtx, err
	}

	// CustomerUser
	if dtx.UserID, err = CreateCustomerUser(dtx.CustomerID); err != nil {
		return &dtx, err
	}

	// Brand
	dtx.BrandID = 1

	// ApiKeyType
	keyTypes := []string{"Authentication", "Private", "Public"}

	keys := make(map[string]string)

	for _, keyType := range keyTypes {
		var keyTypeID string
		if keyTypeID, err = CreateApiKeyType(keyType); err != nil {
			return &dtx, err
		}

		// ApiKey and ApiKeyToBrand -
		if _, keys[keyType], err = CreateApiKey(dtx.UserID, keyTypeID, keyType, dtx.BrandID); err != nil {
			return &dtx, err
		}
	}

	dtx.Globals = make(map[string]interface{})
	dtx.Globals["keys"] = keys

	//leaves APIKey as public one
	for t, k := range keys {
		if t == "Public" {
			dtx.APIKey = k
		}
	}
	// Website
	if dtx.WebsiteID, err = CreateWebsite("http://www.testWebsite23sdf.com", "bogus website"); err != nil {
		return &dtx, err
	}

	// WebsiteToBrand
	var BrandIDs []int
	BrandIDs = append(BrandIDs, dtx.BrandID)
	if err = CreateWebsiteToBrands(BrandIDs, dtx.WebsiteID); err != nil {
		return &dtx, err
	}

	//CustomerToBrand
	if err = InsertCustomerToBrand(dtx.CustomerID, BrandIDs); err != nil {
		return &dtx, err
	}

	//brandString and array
	err = dtx.GetBrandsArrayAndString(dtx.APIKey, dtx.BrandID)
	if err != nil {
		return &dtx, err
	}
	return &dtx, nil
}

func DeMock(dtx *apicontext.DataContext) error {
	var err error

	keys := make(map[string]string)
	keysInterface := dtx.Globals["keys"]
	if keysInterface != nil {
		keys = keysInterface.(map[string]string)
	}

	for kType, key := range keys {
		err = DeleteApiKey(key)
		if err != nil {
			return err
		}
		err = DeleteType(kType)
		if err != nil {
			return err
		}
	}

	err = DeleteCustomerUser(dtx.UserID)
	if err != nil {
		return err
	}

	err = DeleteWebsite(dtx.WebsiteID)
	if err != nil {
		return err
	}

	err = DeleteCustomer(dtx.CustomerID)
	if err != nil {
		return err
	}
	return nil
}

func InsertCustomer() (int, error) {
	var err error
	var i int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return i, err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCustomer)
	if err != nil {
		return i, err
	}
	defer stmt.Close()
	res, err := stmt.Exec()
	if err != nil {
		return i, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return i, err
	}
	i = int(id)
	return i, err
}
func DeleteCustomer(custId int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCustomerToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(custId)

	stmt, err = db.Prepare(deleteCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(custId)
	return err
}

func InsertCustomerToBrand(custId int, brandIds []int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCustomerToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, brandId := range brandIds {
		_, err = stmt.Exec(custId, brandId)
		if err != nil {
			return err
		}
	}
	return err
}

func DeleteType(apiKey string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(apiKey)
	if err != nil {
		return err
	}
	return err
}

func GetUserApiKeys(userID string) ([]string, error) {
	var err error
	var keys []string
	var k string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return keys, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getUserApiKeys)
	if err != nil {
		return keys, err
	}
	defer stmt.Close()
	res, err := stmt.Query(userID)
	if err != nil {
		return keys, err
	}
	for res.Next() {
		err = res.Scan(&k)
		if err != nil {
			return keys, err
		}
		keys = append(keys, k)
	}
	return keys, err
}

func CreateCustomerUser(custID int) (CustomerUserID string, err error) {
	encryptPass := "bogus" // Not encrypted
	email := "TestBogus@curtmfg.com"
	customerID := 1

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
func DeleteCustomerUser(id string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCustomerUser)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return err
}

func CreateApiKeyType(keyType string) (keyTypeID string, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return keyTypeID, err
	}
	defer db.Close()
	stmt, err := db.Prepare(createApiKeyType)
	if err != nil {
		return keyTypeID, err
	}
	defer stmt.Close()
	added := time.Now().Format(timeFormat)
	_, err = stmt.Exec(keyType, added)
	if err != nil {
		return keyTypeID, err
	}

	stmt, err = db.Prepare(getKeyByDateType)
	if err != nil {
		return keyTypeID, err
	}

	defer stmt.Close()

	var typeID *string
	err = stmt.QueryRow(keyType, added).Scan(&typeID)
	if err != nil {
		return keyTypeID, err
	}
	keyTypeID = *typeID

	return keyTypeID, nil
}

func CreateApiKey(UserID string, keyTypeID string, keyType string, brandID int) (keyID int, key string, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return keyID, key, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return keyID, key, err
	}

	stmt, err := tx.Prepare(insertAPIKey)
	if err != nil {
		return keyID, key, err
	}
	res, err := stmt.Exec(UserID, keyTypeID)
	if err != nil {
		tx.Rollback()
		return keyID, key, err
	}
	keyID64, err := res.LastInsertId()
	keyID = int(keyID64)
	if err != nil {
		return keyID, key, err
	}
	stmt, err = tx.Prepare(insertAPIKeyToBrand)
	if err != nil {
		return keyID, key, err
	}
	_, err = stmt.Exec(keyID, brandID)
	if err != nil {
		tx.Rollback()
		return keyID, key, err
	}
	tx.Commit()

	var apiKey *string
	stmt, err = db.Prepare(getCustomerUserKeysWithoutAuth)
	if err != nil {
		return keyID, key, err
	}
	rows, err := stmt.Query(UserID, keyType)
	if err != nil {
		return keyID, key, err
	}

	for rows.Next() {
		var kt string
		err = rows.Scan(&apiKey, &kt)
		if err != nil {
			return keyID, key, err
		}

	}
	defer rows.Close()
	key = *apiKey

	return keyID, key, nil
}

func DeleteApiKey(key string) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteApiKeyToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(deleteApiKey)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}

	return err
}

func CreateWebsite(url, desc string) (webID int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return webID, err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createSite)
	if err != nil {
		return webID, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(url, desc)
	if err != nil {
		tx.Rollback()
		return webID, err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	webID = int(id)
	if err != nil {
		return webID, err
	}
	return webID, nil
}

func CreateWebsiteToBrands(brandIDs []int, websiteID int) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(joinToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, brandId := range brandIDs {
		_, err = stmt.Exec(websiteID, brandId)
		if err != nil {
			return err
		}
	}
	return err
}

func DeleteWebsite(siteId int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteSiteToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(siteId)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare(deleteSite)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(siteId)
	if err != nil {
		return err
	}
	return err
}
