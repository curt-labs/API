package apicontextmock

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"

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

	getCustomerUserKeysWithoutAuth = `select ak.api_key, akt.type from ApiKey as ak
										join ApiKeyType as akt on ak.type_id = akt.id
										where ak.user_id = ? && UPPER(akt.type) = ?`

	getKeyByDateType = "SELECT id FROM ApiKeyType WHERE type = ?  AND date_added = ?"
	createApiKeyType = "INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),?,?)"

	createSite  = `INSERT INTO Website (url, description) VALUES (?,?)`
	joinToBrand = `insert into WebsiteToBrand (WebsiteID, brandID) values (?,?)`
)

func Mock() (*apicontext.DataContext, error) {
	var dtx apicontext.DataContext
	// Needs to create records in the db for the following because of foreign key constraints:
	// Bare Min:

	// CustomerUser
	var err error
	CustomerUserID := ""
	if CustomerUserID, err = CreateCustomerUser(); err != nil {
		return &dtx, err
	}
	dtx.UserID = CustomerUserID
	dtx.CustomerID = 1

	// ApiKeyType
	keyType := "SUPER"
	keyTypeID := "" // needed for when you create an API Key
	if keyTypeID, err = CreateApiKeyType(keyType); err != nil {
		return &dtx, err
	}

	// ApiKey and ApiKeyToBrand
	apiKey := ""
	BrandID := 1
	if _, apiKey, err = CreateApiKey(CustomerUserID, keyTypeID, keyType, BrandID); err != nil {
		return &dtx, err
	}
	dtx.APIKey = apiKey

	// Brand
	dtx.BrandID = BrandID
	// Website
	websiteID := 0
	if websiteID, err = CreateWebsite("http://www.testWebsite23sdf.com", "bogus website"); err != nil {
		return &dtx, err
	}
	dtx.WebsiteID = websiteID
	// WebsiteToBrand
	var BrandIDs []int
	BrandIDs = append(BrandIDs, 1)
	if err = CreateWebsiteToBrands(BrandIDs, websiteID); err != nil {
		return &dtx, err
	}
	return &dtx, nil
}

func DeMock(dtx *apicontext.DataContext) error { // in place for future democking if we decide to.
	return nil
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
