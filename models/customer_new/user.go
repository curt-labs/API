package customer_new

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/encryption"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CustomerUser struct {
	Id                 string           `json:"id" xml:"id"`
	Name               string           `json:"name" xml:"name,attr"`
	Email              string           `json:"email" xml:"email,attr"`
	Password           string           `json:"password,omitempty" xml:"password,omitempty"`
	OldCustomerID      int              `json:"oldCustomerId,omitempty" xml:"oldCustomerId,omitempty"`
	DateAdded          time.Time        `json:"date_added" xml:"date_added,attr"`
	Active             bool             `json:"active" xml:"active,attr"`
	Location           CustomerLocation `json:"location" xml:"location"`
	Sudo               bool             `json:"sudo" xml:"sudo,attr"`
	CustomerID         int              `json:"customerId,omitempty" xml:"customerId,omitempty"`
	Current            bool             `json:"current" xml:"current,attr"`
	PasswordConversion bool             `json:"passwordConversion,omitempty" xml:"passwordConversion,omitempty"`
	Keys               []ApiCredentials `json:"keys" xml:"keys"`
}

type ApiCredentials struct {
	Key       string    `json:"key" xml:"key,attr"`
	Type      string    `json:"type" xml:"type,attr"`
	TypeId    string    `json:"typeID" xml:"typeID,attr"`
	DateAdded time.Time `json:"date_added" xml:"date_added,attr"`
}

const (
	customerUserFields = ` cu.id, cu.name, cu.email, cu.password, cu.customerID, cu.date_added, cu.active, cu.locationID, cu.isSudo, cu.cust_ID, cu.NotCustomer, cu.passwordConverted `
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
								where UPPER(akt.type) = ?
								&& ak.api_key = ?
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

	updateCustomerUser   = `UPDATE CustomerUser SET name = ?, email = ?, active = ?, locationID = ?, isSudo = ?, NotCustomer = ? WHERE id = ?`
	getUsersByCustomerID = `SELECT id FROM CustomerUser WHERE customerID = ?`
)

var (
	AuthError = errors.New("failed to authenticate")
)

const (
	AUTH_KEY_TYPE    = "AUTHENTICATION"
	PUBLIC_KEY_TYPE  = "PUBLIC"
	PRIVATE_KEY_TYPE = "PRIVATE"
)

func ScanUser(res Scanner) (*CustomerUser, error) {
	var cu CustomerUser
	var err error
	var passConversionByte []byte
	var oldId *int
	var cur *bool
	var name *string
	err = res.Scan(
		&cu.Id,
		&name,
		&cu.Email,
		&cu.Password,
		&oldId,
		&cu.DateAdded,
		&cu.Active,
		&cu.Location.Id,
		&cu.Sudo,
		&cu.CustomerID,
		&cur,
		&passConversionByte,
	)
	if err != nil {
		return &cu, err
	}
	if passConversionByte != nil {
		var errConver error
		cu.PasswordConversion, errConver = strconv.ParseBool(string(passConversionByte))
		if errConver != nil {
			cu.PasswordConversion = false
		}
	}
	if name != nil {
		cu.Name = *name
	}
	if oldId != nil {
		cu.OldCustomerID = *oldId
	}
	if cur != nil {
		cu.Current = *cur
	}
	return &cu, err
}

func AuthenticateUserByKey(key string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return u, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserKeyAuth)
	if err != nil {
		return u, err
	}

	defer stmt.Close()
	t := time.Now()
	t1 := t.Add(time.Duration(-6) * time.Hour) //6 hours ago
	Timer := t1.String()
	KeyType := api_helpers.AUTH_KEY_TYPE
	params := []interface{}{
		KeyType,
		key,
		Timer,
	}
	res := stmt.QueryRow(params...)
	user, err := ScanUser(res)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("error: %s", "user does not exist")
		}
		return u, err
	}
	u = *user

	resetChan := make(chan int)
	go func() {
		if resetErr := u.ResetAuthentication(); resetErr != nil {
			err = resetErr
		}
		resetChan <- 1
	}()

	<-resetChan
	return
}

func AuthenticateAndGetCustomer(key string) (cust Customer, err error) {
	u, err := AuthenticateUserByKey(key)
	if err != nil {
		return cust, AuthError
	}

	keyChan := make(chan int)
	locChan := make(chan int)
	go func() {
		if kErr := u.GetKeys(); kErr != nil {
			err = kErr
		}
		keyChan <- 1
	}()

	go func() {
		if lErr := u.GetLocation(); lErr != nil {
			err = lErr
		}
		locChan <- 1
	}()
	cust, err = u.GetCustomer(key)
	if err != sql.ErrNoRows {
		if err != nil {
			return cust, AuthError
		}
	}
	<-keyChan
	<-locChan

	cust.Users = append(cust.Users, u)

	return cust, nil
}

func (u *CustomerUser) AuthenticateUser() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return AuthError
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserAuth)
	if err != nil {
		return AuthError
	}
	defer stmt.Close()

	var dbPass string
	var passConversionByte []byte
	var passConversion bool
	var sudo, active *bool
	var id, name, email *string
	var dateAdded *time.Time
	var oldId, locId, custId *int
	err = stmt.QueryRow(u.Email).Scan(
		&id,
		&name,
		&email,
		&dbPass,
		&oldId,
		&dateAdded,
		&active,
		&locId,
		&sudo,
		&custId,
		&passConversionByte,
	)
	if err != nil {
		return err
	}
	if id != nil {
		u.Id = *id
	}
	if name != nil {
		u.Name = *name
	}
	if email != nil {
		u.Email = *email
	}
	if oldId != nil {
		u.OldCustomerID = *oldId
	}
	if dateAdded != nil {
		u.DateAdded = *dateAdded
	}
	if active != nil {
		u.Active = *active
	}
	if locId != nil {
		u.Location.Id = *locId
	}
	if sudo != nil {
		u.Sudo = *sudo
	}
	if custId != nil {
		u.CustomerID = *custId
	}
	if passConversionByte != nil {
		passConversion, err = strconv.ParseBool(string(passConversionByte))
	}
	pass := u.Password

	// Attempt to compare bcrypt strings
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass))
	if err != nil {

		// Compare unsuccessful

		enc_pass, err := api_helpers.Md5Encrypt(pass)
		if err != nil {
			return err
		}
		if len(enc_pass) != len(dbPass) || passConversion { //bool
			return errors.New("Invalid password")
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("Failed to encode the password")
		}

		stmtPass, err := db.Prepare(updateCustomerUserPass)
		if err != nil {
			return err
		}
		_, err = stmtPass.Exec(hashedPass, u.Id)
		return errors.New("Incorrect password.")
	}

	resetChan := make(chan int)
	go func() {
		if resetErr := u.ResetAuthentication(); resetErr != nil {
			err = resetErr
		}
		resetChan <- 1
	}()

	<-resetChan
	u.Current = true

	return nil
}

//like AuthenticateUserByKey, but does not update the timestamp - seems REDUNDANT
func GetCustomerUserFromKey(key string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return u, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserFromKey)
	if err != nil {
		return u, err
	}
	defer stmt.Close()

	res := stmt.QueryRow(api_helpers.AUTH_KEY_TYPE, key)
	user, err := ScanUser(res)
	if err != nil {
		err = fmt.Errorf("error: %s", "user does not exist")
		return
	}

	u = *user
	locChan := make(chan error)
	go func() {
		locChan <- u.GetLocation()
	}()

	u.GetKeys()
	<-locChan

	return
}

//Takes UUID CustomerID; deletes all CustomerUser with that CustID and their API Keys
func DeleteCustomerUsersByCustomerID(customerID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getUsersByCustomerID)
	if err != nil {
		return err
	}
	res, err := stmt.Query(customerID)
	if err != nil {
		return err
	}
	for res.Next() {
		var tempCustUser CustomerUser
		err = res.Scan(&tempCustUser.Id)
		if err != nil {
			return err
		}
		err = tempCustUser.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (u CustomerUser) GetCustomer(key string) (c Customer, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return c, err
	}
	defer db.Close()

	stmt, err := db.Prepare(userCustomer)
	if err != nil {
		return c, err
	}
	defer stmt.Close()

	res := stmt.QueryRow(u.Id)
	if err := c.ScanCustomer(res, key); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("error: %s", "user not bound to customer")
		}
		return c, err
	}

	locChan := make(chan error)
	go func() {
		locChan <- c.GetLocations()
	}()

	if u.Sudo {
		c.GetUsers(key)
	}

	<-locChan

	return
}

func (u *CustomerUser) GetKeys() error {
	var keys []ApiCredentials
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserKeys)
	if err != nil {
		return err
	}
	defer stmt.Close()

	params := []interface{}{
		u.Id,
		strings.Join([]string{api_helpers.AUTH_KEY_TYPE}, ","),
	}
	res, err := stmt.Query(params...)
	for res.Next() {
		var a ApiCredentials
		res.Scan(&a.Key, &a.Type, &a.DateAdded)
		keys = append(keys, a)
	}
	u.Keys = keys
	return nil
}

func (u *CustomerUser) GetLocation() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(userLocation)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var stateId, countryId *int
	var state, stateAbbr, country, countryAbbr *string

	err = stmt.QueryRow(u.Id).Scan(
		&u.Location.Id,
		&u.Location.Name,
		&u.Location.Email,
		&u.Location.Address,
		&u.Location.City,
		&u.Location.PostalCode,
		&u.Location.Phone,
		&u.Location.Fax,
		&u.Location.Latitude,
		&u.Location.Longitude,
		&u.Location.CustomerId,
		&u.Location.ContactPerson,
		&u.Location.IsPrimary,
		&u.Location.ShippingDefault,
		&stateId,
		&state,
		&stateAbbr,
		&countryId,
		&country,
		&countryAbbr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	var coun geography.Country

	if stateId != nil {
		u.Location.State.Id = *stateId
	}
	if state != nil {
		u.Location.State.State = *state
	}
	if stateAbbr != nil {
		u.Location.State.Abbreviation = *stateAbbr
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
	u.Location.State.Country = &coun
	return nil
}

//updates auth key dateAdded to Now()
func (u *CustomerUser) ResetAuthentication() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(userAuthenticationKey)
	if err != nil {
		return err
	}
	defer stmt.Close()

	params := []interface{}{
		api_helpers.AUTH_KEY_TYPE,
		u.Id,
	}
	var a ApiCredentials

	a.TypeId, err = getAPIKeyTypeReference(api_helpers.AUTH_KEY_TYPE)
	if err != nil {
		return fmt.Errorf("error: %s", "failed to retrieve key type reference")
	}

	var dateAdded string
	err = stmt.QueryRow(params...).Scan(&a.Key, &a.Type, &a.TypeId, &dateAdded)
	if err != nil {
		apiCredentials, err := u.GenerateAPIKey(api_helpers.AUTH_KEY_TYPE)
		if err != nil {
			return err
		}
		u.Keys = append(u.Keys, *apiCredentials)
	} else {
		loc, _ := time.LoadLocation("US/Central")
		a.DateAdded, _ = time.ParseInLocation(time.RFC3339Nano, dateAdded, loc)
		paramsNew := []interface{}{
			a.TypeId,
			u.Id,
		}

		stmtNew, err := db.Prepare(resetUserAuthentication)
		_, err = stmtNew.Exec(paramsNew...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cu *CustomerUser) GenerateAPIKey(keyType string) (*ApiCredentials, error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertAPIKey)
	if err != nil {
		return nil, err
	}

	typeID, err := getAPIKeyTypeReference(keyType)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(cu.Id, typeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	var apiKey string
	stmt, err = db.Prepare(getCustomerUserKeysWithoutAuth)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(cu.Id, keyType)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var kt string
		err = rows.Scan(&apiKey, &kt)
		if err != nil {
			return nil, err
		}

		if strings.ToLower(kt) == strings.ToLower(keyType) {
			cred := ApiCredentials{}
			cred.Key = apiKey
			cred.Type = keyType
			cred.TypeId = typeID
			cred.DateAdded = time.Now()
			return &cred, nil
		}

	}
	defer rows.Close()

	return nil, fmt.Errorf("%s", "failed to generate new key")
}

func getAPIKeyTypeReference(keyType string) (string, error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return "", err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAPIKeyTypeID)
	var apiKeyTypeId string
	err = stmt.QueryRow(keyType).Scan(&apiKeyTypeId)
	if err != nil {
		return uuid.NIL.String(), errors.New("failed to retrieve auth type")
	}
	return apiKeyTypeId, nil
}

func (cu *CustomerUser) ResetPass() (string, error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return "", err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare(setCustomerUserPassword)
	if err != nil {
		return "", err
	}

	randPass := encryption.GeneratePassword()

	// encrypt the random password:
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(randPass), bcrypt.DefaultCost)
	e, err := conversions.ByteToString(encryptPass)
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(e, cu.Email)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return randPass, nil
}

func (cu *CustomerUser) ChangePass(oldPass, newPass string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(setCustomerUserPassword)
	encryptNewPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	err = cu.AuthenticateUser()
	if err != nil {
		return errors.New("Old password is incorrect.")
	}

	_, err = stmt.Exec(encryptNewPass, cu.Email)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (cu *CustomerUser) Get(key string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserFromId)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var dbPass, passConversion string
	err = stmt.QueryRow(cu.Id).Scan(
		&cu.Id,
		&cu.Name,
		&cu.Email,
		&dbPass, //Not Used
		&cu.OldCustomerID,
		&cu.DateAdded,
		&cu.Active,
		&cu.Location.Id,
		&cu.Sudo,
		&cu.CustomerID,
		&cu.Current,
		&passConversion, //Not Used
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("error: %s", "user does not exist")
		}
		return err
	}

	keyChan := make(chan error)
	go func() {
		err := cu.GetKeys()
		if err == nil {
			for _, k := range cu.Keys {
				if strings.ToLower(k.Key) == strings.ToLower(key) {
					cu.Current = true
					break
				}
			}
		}
		keyChan <- err
	}()

	cu.GetLocation()

	<-keyChan

	return nil
}

//Update customerUser
func (cu *CustomerUser) UpdateCustomerUser() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(updateCustomerUser)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		cu.Name,
		cu.Email,
		cu.Active,
		cu.Location.Id,
		cu.Sudo,
		cu.Current,
		cu.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//Create CustomerUser
func (cu *CustomerUser) Create() error {
	var err error
	encryptPass, err := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to generate encrypted password.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertCustomerUser)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(cu.Name, cu.Email, encryptPass, cu.OldCustomerID, cu.Active, cu.Location.Id, cu.Sudo, cu.CustomerID, cu.Current)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	stmt, err = db.Prepare(getRegisteredUsersId) // needs to be set on the customer user object in order to generate the keys
	if err != nil {
		return err
	}

	var userID *string
	if err = stmt.QueryRow(cu.Email, encryptPass).Scan(&userID); err != nil || userID == nil {
		return err
	}

	cu.Id = *userID

	// then create API keys for the user
	pubChan := make(chan error)
	privChan := make(chan error)
	authChan := make(chan error)

	// Public key:
	go func() {
		pub, err := cu.GenerateAPIKey(PUBLIC_KEY_TYPE)
		if pub != nil {
			cu.Keys = append(cu.Keys, *pub)
		}
		pubChan <- err
	}()

	// Private key:
	go func() {
		pri, err := cu.GenerateAPIKey(PRIVATE_KEY_TYPE)
		if pri != nil {
			cu.Keys = append(cu.Keys, *pri)
		}
		privChan <- err
	}()

	// Auth Key:
	go func() {
		auth, err := cu.GenerateAPIKey(AUTH_KEY_TYPE)
		if auth != nil {
			cu.Keys = append(cu.Keys, *auth)
		}
		authChan <- err

	}()

	if e := <-pubChan; e != nil {
		return e
	}
	if e := <-privChan; e != nil {
		return e
	}
	if e := <-authChan; e != nil {
		return e
	}

	return nil
}

func (cu *CustomerUser) Delete() error {

	//delete api keys
	pubChan := make(chan int)
	privChan := make(chan int)
	authChan := make(chan int)
	var errStr string

	// Public key:
	go func() {
		err := cu.deleteApiKey(PUBLIC_KEY_TYPE)
		if err != nil {
			errStr += err.Error()
		}
		pubChan <- 1
	}()

	// Private key:
	go func() {
		err := cu.deleteApiKey(PRIVATE_KEY_TYPE)
		if err != nil {
			errStr += err.Error()
		}
		privChan <- 1
	}()

	// Auth Key:
	go func() {
		err := cu.deleteApiKey(AUTH_KEY_TYPE)
		if err != nil {
			errStr += err.Error()
		}
		authChan <- 1
	}()
	<-pubChan
	<-privChan
	<-authChan

	//delete CustomerUser
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteCustomerUser)
	_, err = stmt.Exec(cu.Id)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (cu *CustomerUser) deleteApiKey(keyType string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(deleteAPIkey)
	if err != nil {
		return err
	}

	typeID, err := getAPIKeyTypeReference(keyType)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(cu.Id, typeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

type ApiRequest struct {
	User        CustomerUser `json:"user" xml:"user"`
	RequestTime time.Time    `json:"request_time" xml:"request_time,attr"`
	Url         *url.URL     `json:"url" xml:"url"`
	Query       url.Values   `json:"query" xml:"query"`
	Form        url.Values   `json:"form" xml:"form"`
}

func (u *CustomerUser) LogApiRequest(r *http.Request) {
	var ar ApiRequest
	ar.User = *u
	ar.RequestTime = time.Now()
	ar.Url = r.URL
	ar.Query = r.URL.Query()
	ar.Form = r.Form

	redis.Lpush(fmt.Sprintf("log:%s", u.Id), ar)
}
