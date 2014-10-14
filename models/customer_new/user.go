package customer_new

import (
	"code.google.com/p/go-uuid/uuid"
	// "crypto/md5"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/encryption"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type CustomerUser struct {
	Id                    string
	Name, Email           string
	DateAdded             time.Time
	Active, Sudo, Current bool
	Location              CustomerLocation
	Keys                  []ApiCredentials
}

type ApiCredentials struct {
	Key, Type, TypeId string
	DateAdded         time.Time
}

var (
	userCustomer = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
						c.latitude, c.longitude, c.searchURL, c.logo, c.website,
						c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
						dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
						dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
						mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
						mpx.code as mapix_code, mpx.description as mapic_desc,
						sr.name as rep_name, sr.code as rep_code, c.parentID
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
	getRegisteredUsersId = `
		select cu.id from CustomerUser as cu
		where cu.email = ? && cu.password = ?
		limit 1`

	customerUserAuth = `select password, id, name, email, date_added, active, isSudo, passwordConverted from CustomerUser
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
									join CustomerLocations as cl on cu.cust_ID = cl.cust_id
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
	customerIDFromKey = `select c.customerID from Customer as c
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where ak.api_key = ?
								limit 1`
	customerUserFromKey = `select cu.* from CustomerUser as cu
								join ApiKey as ak on cu.id = ak.user_id
								join ApiKeyType as akt on ak.type_id = akt.id
								where akt.type != ? && ak.api_key = ?
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
	getAPIKeyTypeID               = `select id from ApiKeyType where UPPER(type) = UPPER(?) limit 1`
	updateCustomerUserPassByEmail = `update CustomerUser set password = ?, passwordConverted = 1 WHERE email = ? AND customerID = ?`
	setCustomerUserPassword       = `update CustomerUser set password = ?, passwordConverted = 1 WHERE email = ?`
	deleteCustomerUser            = `DELETE FROM CustomerUser WHERE id = ?`
	deleteAPIkey                  = `DELETE FROM ApiKey WHERE user_id = ? AND type_id = ?`
	getCustomerUserKeysWithAuth   = `select ak.api_key, akt.type from ApiKey as ak
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

func (u CustomerUser) UserAuthentication(password string) (cust Customer, err error) {

	err = u.AuthenticateUser(password)
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

	cust, err = u.GetCustomer()
	if err != nil {
		return cust, AuthError
	}

	<-keyChan
	<-locChan

	cust.Users = append(cust.Users, u)

	return cust, nil
}

func UserAuthenticationByKey(key string) (cust Customer, err error) {
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

	cust, err = u.GetCustomer()
	if err != nil {
		return cust, AuthError
	}

	<-keyChan
	<-locChan

	cust.Users = append(cust.Users, u)

	var m *runtime.MemStats = new(runtime.MemStats)
	runtime.ReadMemStats(m)
	log.Println("Memory", m.Alloc, m.StackInuse, m.HeapAlloc)
	return cust, nil
}

func (u CustomerUser) GetCustomer() (c Customer, err error) {
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

	var logo, web, lat, lon, url, icon, shadow, mapIconId []byte
	var stateId, state, stateAbbr, countryId, country, countryAbbr, parentId, postalCode, mapixCode, mapixDesc, rep, repCode []byte
	err = stmt.QueryRow(u.Id).Scan(
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
		return c, err
	}
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
	c.MapixCode.Code, err = conversions.ByteToString(mapixCode)
	c.MapixCode.Description, err = conversions.ByteToString(mapixDesc)
	c.SalesRepresentative.Name, err = conversions.ByteToString(rep)
	c.SalesRepresentative.Code, err = conversions.ByteToString(repCode)

	parentInt, err := conversions.ByteToInt(parentId)
	if err != nil {
		return c, err
	}
	if parentInt != 0 {
		par := Customer{Id: parentInt}
		par.GetCustomer()
		c.Parent = &par
	}
	return
}

func (u *CustomerUser) AuthenticateUser(pass string) error {

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
	var passConversion bool
	err = stmt.QueryRow(u.Email).Scan(
		&dbPass,
		&u.Id,
		&u.Name,
		&u.Email,
		&u.DateAdded,
		&u.Active,
		&u.Sudo,
		&passConversion,
	)
	if err != nil {
		err = errors.New("No user found that matches: " + u.Email)
	}

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
	return nil
}

func AuthenticateUserByKey(key string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return u, AuthError
	}
	defer db.Close()

	stmt, err := db.Prepare(customerUserKeyAuth)
	if err != nil {
		return u, AuthError
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
	var dbPass, custId, customerId string
	var passConversion []byte //bools
	err = stmt.QueryRow(params...).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&dbPass,     //Not Used
		&customerId, //Not Used
		&u.DateAdded,
		&u.Active,
		&u.Location.Id,
		&u.Sudo,
		&custId, //Not Used
		&u.Current,
		&passConversion, //Not Used
	)
	if err != nil {
		return u, AuthError
	}

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

	err = stmt.QueryRow(u.Id).Scan(
		&u.Location.Id,
		&u.Name,
		&u.Email,
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
		&u.Location.State.Id,
		&u.Location.State.State,
		&u.Location.State.Abbreviation,
		&u.Location.State.Country.Id,
		&u.Location.State.Country.Country,
		&u.Location.State.Country.Abbreviation,
	)
	if err != nil {
		return err
	}
	return nil
}

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

	var dateAdded string
	err = stmt.QueryRow(params...).Scan(&a.Key, &a.Type, &a.TypeId, &dateAdded)
	if err != nil {
		return err
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

func GetCustomerIdFromKey(key string) (id int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(customerIDFromKey)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(key).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

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

	params := []interface{}{
		api_helpers.AUTH_KEY_TYPE,
		key,
	}
	var dbPass, custId, passConversion, customerId string
	err = stmt.QueryRow(params...).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&dbPass, //Not Used
		&custId, //Not Used
		&u.DateAdded,
		&u.Active,
		&u.Location.Id,
		&u.Sudo,
		&customerId,     //Not Used
		&u.Current,      //Not Used
		&passConversion, //Not Used
	)
	if err != nil {
		err = errors.New("Invalid key")
		return
	}
	return
}

func (cu *CustomerUser) Get() error {
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

	var dbPass, custId, customerId, passConversion string
	err = stmt.QueryRow(cu.Id).Scan(
		&cu.Id,
		&cu.Name,
		&cu.Email,
		&dbPass, //Not Used
		&custId, //Not User
		&cu.DateAdded,
		&cu.Active,
		&cu.Location.Id,
		&cu.Sudo,
		&customerId, //Not Used
		&cu.Current,
		&passConversion, //Not Used
	)
	if err != nil {
		return err
	}
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
func (cu *CustomerUser) Register(pass string, customerID int, isActive bool, locationID int, isSudo bool, cust_ID int, notCustomer bool) (*CustomerUser, error) {

	encryptPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Failed to generate UUID.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertCustomerUser)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(cu.Name, cu.Email, encryptPass, customerID, isActive, locationID, isSudo, cust_ID, notCustomer)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	stmt, err = db.Prepare(getRegisteredUsersId) // needs to be set on the customer user object in order to generate the keys
	if err != nil {
		return nil, err
	}

	var userID *string
	if err = stmt.QueryRow(cu.Email, encryptPass).Scan(&userID); err != nil || userID == nil {
		return nil, err
	}

	cu.Id = *userID

	// then create API keys for the user
	pubChan := make(chan error)
	privChan := make(chan error)
	authChan := make(chan error)

	// Public key:
	go func() {
		_, err := cu.generateAPIKey(PUBLIC_KEY_TYPE)
		pubChan <- err
	}()

	// Private key:
	go func() {
		_, err := cu.generateAPIKey(PRIVATE_KEY_TYPE)
		privChan <- err
	}()

	// Auth Key:
	go func() {
		_, err := cu.generateAPIKey(AUTH_KEY_TYPE)
		authChan <- err
	}()
	if e := <-pubChan; e != nil {
		return cu, e
	}
	if e := <-privChan; e != nil {
		return cu, e
	}
	if e := <-authChan; e != nil {
		return cu, e
	}

	return cu, nil
}

func (cu *CustomerUser) generateAPIKey(keyType string) (string, error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return "", err
	}
	defer db.Close()
	tx, err := db.Begin()

	stmt, err := tx.Prepare(insertAPIKey)
	if err != nil {
		return "", err
	}

	typeID, err := getAPIKeyTypeReference(keyType)
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(cu.Id, typeID)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()

	var apiKey *string
	stmt, err = db.Prepare(getCustomerUserKeysWithoutAuth)
	if err != nil {
		return "", err
	}

	err = stmt.QueryRow(cu.Id, keyType).Scan(&apiKey, &keyType)
	if err != nil {
		return "", err
	}
	return *apiKey, nil
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

func (cu *CustomerUser) ResetPass(custID string) (string, error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return "", err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare(updateCustomerUserPassByEmail)
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

	_, err = stmt.Exec(e, cu.Email, custID)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return randPass, nil
}

func (cu *CustomerUser) ChangePass(oldPass, newPass string, custID int) (string, error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return "", err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(setCustomerUserPassword)
	encryptNewPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	err = cu.AuthenticateUser(oldPass)
	if err != nil {
		return "", errors.New("Old password is incorrect.")
	}

	_, err = stmt.Exec(encryptNewPass, cu.Email)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "success", nil
}

func (cu *CustomerUser) Delete() error {

	//delete api keys
	pubChan := make(chan int)
	privChan := make(chan int)
	authChan := make(chan int)

	// Public key:
	go func() {
		cu.deleteApiKey(PUBLIC_KEY_TYPE)
		pubChan <- 1
	}()

	// Private key:
	go func() {
		cu.deleteApiKey(PRIVATE_KEY_TYPE)
		privChan <- 1
	}()

	// Auth Key:
	go func() {
		cu.deleteApiKey(AUTH_KEY_TYPE)
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

func (cu *CustomerUser) BindApiAccess() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUserKeysWithAuth)
	if err != nil {
		return err
	}

	var keys []ApiCredentials

	res, err := stmt.Query(cu.Id, PUBLIC_KEY_TYPE, PRIVATE_KEY_TYPE, AUTH_KEY_TYPE)
	for res.Next() {
		var k ApiCredentials
		err = res.Scan(&k.Key, &k.Type)
		if err != nil {
			return err
		}
		keys = append(keys, k)
	}
	cu.Keys = keys
	return nil
}

func (cu *CustomerUser) BindLocation() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUserLocation)
	if err != nil {
		return err
	}
	err = stmt.QueryRow(cu.Id).Scan(
		&cu.Location.Id,
		&cu.Location.Name,
		&cu.Location.Address,
		&cu.Location.City,
		&cu.Location.State.Id,
		&cu.Location.State.State,
		&cu.Location.State.Abbreviation,
		&cu.Location.State.Country.Id,
		&cu.Location.State.Country.Country,
		&cu.Location.State.Country.Abbreviation,
		&cu.Location.Email,
		&cu.Location.Phone,
		&cu.Location.Fax,
		&cu.Location.Latitude,
		&cu.Location.Longitude,
		&cu.Location.CustomerId,
		&cu.Location.ContactPerson,
		&cu.Location.IsPrimary,
		&cu.Location.PostalCode,
		&cu.Location.ShippingDefault,
	)

	if err != nil {
		return err
	}
	return nil
}

type ApiRequest struct {
	User        CustomerUser
	RequestTime time.Time
	Url         *url.URL
	Query       url.Values
	Form        url.Values
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

// The disabling of the triggers is failing in this method.
//
// I'm going to disable the call to it completely and expand
// the time limit of the authentication key to 6 hours.
//
// TODO: This will need to be fixed at some point in time. **Important

// func (u *CustomerUser) RenewAuthentication() error {
// 	log.Println("renewing authentication key")
// 	t := time.Now()

// 	log.Printf(renewUserAuthenticationStmt, t.String(), AUTH_KEY_TYPE, u.Id)

// 	// Excecute the update statement
// 	_, _, err := database.Db.Query(disableTriggerStmt)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	_, _, err = database.Db.Query(renewUserAuthenticationStmt, t.String(), AUTH_KEY_TYPE, u.Id)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	_, _, err = database.Db.Query(enableTriggerStmt)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }
