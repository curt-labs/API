package customer

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/geography"
	_ "github.com/go-sql-driver/mysql"
	//"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	getUserCustomerStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
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
	getUserLocationStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
                          cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
                          cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
                          s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
                          from CustomerUser as cu
                          join CustomerLocations as cl on cu.cust_ID = cl.cust_id
                          left join States as s on cl.stateID = s.stateID
                          left join Country as cty on s.countryID = cty.countryID
                          where cu.id = ?`
	getCustomerUserKeysStmt = `select ak.api_key, akt.type, ak.date_added from ApiKey as ak
                              join ApiKeyType as akt on ak.type_id = akt.id
                              where user_id = ? && UPPER(akt.type) NOT IN (?)`
	getUserAuthenticationKeyStmt = `select ak.api_key, akt.type, akt.id, CAST(ak.date_added as char(255)) as date_added from ApiKey as ak
                                    join ApiKeyType as akt on ak.type_id = akt.id
                                    where UPPER(akt.type) = ?
                                    && ak.user_id = ?`
	getCustomerIdFromKeyStmt = `select c.customerID from Customer as c
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`
	getCustomerUserFromKeyStmt = `select cu.* from CustomerUser as cu
                                 join ApiKey as ak on cu.id = ak.user_id
                                 join ApiKeyType as akt on ak.type_id = akt.id
                                 where akt.type != ? && ak.api_key = ? limit 1`
	getCustomerUserFromIdStmt = `select cu.* from CustomerUser as cu
                                join ApiKey as ak on cu.id = ak.user_id
                                join ApiKeyType as akt on ak.type_id = akt.id
                                where cu.id = ? limit 1`
	resetUserAuthenticationStmt = `update ApiKey as ak
                              set ak.date_added = NOW()
                              where ak.type_id = ? && ak.user_id = ?`
	authCustomerUserStmt = `select password, id, name, email, date_added, active, isSudo, passwordConverted
                           from CustomerUser
                           where email = ? && active = 1 limit 1`
	authCustomerUserByKeyStmt = `select cu.* from CustomerUser as cu
                                join ApiKey as ak on cu.id = ak.user_id
                                join ApiKeyType as akt on ak.type_id = akt.id
                                where UPPER(akt.type) = ? && ak.api_key = ? && cu.active = 1 && ak.date_added >= ?`
	updateCustomerUserPassStmt = `update CustomerUser set password = ?, passwordConverted = 1
                                 where id = ? && active = 1`
)

var (
	AuthError = errors.New("failed to authenticate")
)

type CustomerUser struct {
	Id                    string
	Name, Email           string
	DateAdded             time.Time
	Active, Sudo, Current bool
	Location              *CustomerLocation
	Keys                  *[]ApiCredentials
}

type ApiCredentials struct {
	Key, Type string
	DateAdded time.Time
}

func (u CustomerUser) UserAuthentication(password string) (cust Customer, err error) {

	err = u.AuthenticateUser(password)
	if err != nil {
		fmt.Println("error in authenticate user")
		return
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

	<-keyChan
	<-locChan

	cust.Users = append(cust.Users, u)

	return
}

func UserAuthenticationByKey(key string) (cust Customer, err error) {
	u, err := AuthenticateUserByKey(key)
	if err != nil {
		return
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

	<-keyChan
	<-locChan
	cust.Users = append(cust.Users, u)

	return
}

func (u CustomerUser) GetCustomer() (c Customer, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getUserCustomerStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var parentID, mapIconId *int
	var logoUrl, searchUrl, websiteUrl, mapIconUrl, mapShadowUrl *string
	var salesRep, salesRepCode *string

	var stateID, countryID *int
	var stateName, stateAbbr, countryName, countryAbbr *string
	var typeID, tierID, tierSort *int
	var dealerType, dealerLabel, tier *string
	var typeOnline, typeShow *bool
	var lat, lon string

	c.State = &geography.State{}
	c.State.Country = &geography.Country{}

	err = stmt.QueryRow(u.Id).Scan(
		&c.Id,
		&c.Name,
		&c.Email,
		&c.Address,
		&c.Address2,
		&c.City,
		&c.Phone,
		&c.Fax,
		&c.ContactPerson,
		&lat,
		&lon,
		&searchUrl,
		&logoUrl,
		&websiteUrl,
		&c.PostalCode,
		&stateID,
		&stateName,
		&stateAbbr,
		&countryID,
		&countryName,
		&countryAbbr,
		&typeID,
		&dealerType,
		&typeOnline,
		&typeShow,
		&dealerLabel,
		&tierID,
		&c.DealerTier.Tier,
		&c.DealerTier.Sort,
		&mapIconId,
		&mapIconUrl,
		&mapShadowUrl,
		&c.MapixCode,
		&c.MapixDescription,
		&salesRep,
		&salesRepCode,
		&parentID,
	)
	if err != nil {
		return
	}

	if stateID != nil && *stateID > 0 && stateName != nil && stateAbbr != nil {
		c.State = &geography.State{
			Id:           *stateID,
			State:        *stateName,
			Abbreviation: *stateAbbr,
		}
		if countryID != nil && *countryID > 0 && countryName != nil && countryAbbr != nil {
			c.State.Country = &geography.Country{
				Id:           *countryID,
				Country:      *countryName,
				Abbreviation: *countryAbbr,
			}
		}
	}

	c.DealerType = DealerType{}
	c.DealerTier = DealerTier{}
	if tierID != nil && tierSort != nil && tier != nil {
		c.DealerTier.Id = *tierID
		c.DealerTier.Tier = *tier
		c.DealerTier.Sort = *tierSort
	}
	if typeID != nil && dealerType != nil && dealerLabel != nil &&
		typeOnline != nil && typeShow != nil {
		c.DealerType.Id = *typeID
		c.DealerType.Type = *dealerType
		c.DealerType.Label = *dealerLabel
		c.DealerType.Online = *typeOnline
		c.DealerType.Show = *typeShow
	}
	if lat != "" && lon != "" {
		c.Latitude, _ = strconv.ParseFloat(lat, 64)
		c.Longitude, _ = strconv.ParseFloat(lon, 64)
	}

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

	if mapIconId != nil {
		c.DealerType.MapIcon.Id = *mapIconId
	}

	if salesRep != nil {
		c.SalesRepresentative = *salesRep
	}

	if salesRepCode != nil {
		c.SalesRepresentativeCode = *salesRepCode
	}

	if parentID != nil && *parentID != 0 {
		c.Parent = &Customer{Id: *parentID}
		c.Parent.GetCustomer()
	}

	if err = c.GetLocations(); err != nil {
		return
	}

	return
}

func (u *CustomerUser) AuthenticateUser(pass string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(authCustomerUserStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var dbPass string
	var isActive, isSudo, isPassConverted int

	row := stmt.QueryRow(u.Email)
	err = row.Scan(
		&dbPass,
		&u.Id,
		&u.Name,
		&u.Email,
		&u.DateAdded,
		&isActive,
		&isSudo,
		&isPassConverted,
	)

	if err != nil {
		return err
	}

	u.Active = isActive == 1
	u.Sudo = isSudo == 1

	// Attempt to compare bcrypt strings
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass))
	if err != nil {
		// Compare unsuccessful
		enc_pass, err := api_helpers.Md5Encrypt(pass)
		if err != nil {
			return AuthError
		}
		if len(enc_pass) != len(dbPass) || isPassConverted == 1 {
			return AuthError
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return AuthError
		}

		stmt, err = db.Prepare(updateCustomerUserPassStmt)
		if err != nil {
			return AuthError
		}
		defer stmt.Close()

		_, err = stmt.Exec(hashedPass, u.Id)
		if err != nil {
			return AuthError
		}
	}

	resetChan := make(chan int)
	go func() {
		if resetErr := u.ResetAuthentication(); resetErr != nil {
			err = resetErr
		}
		resetChan <- 1
	}()

	u.Current = true

	<-resetChan

	return nil
}

func AuthenticateUserByKey(key string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(authCustomerUserByKeyStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	t := time.Now()
	t1 := t.Add(time.Duration(-6) * time.Hour) //6 hours ago
	tstr := t1.String()

	var dbPass, custId, customerId string
	var passConversion bool

	u.Location = &CustomerLocation{}

	err = stmt.QueryRow(api_helpers.PRIVATE_KEY_TYPE, key, tstr).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&dbPass,
		&customerId,
		&u.DateAdded,
		&u.Active,
		&u.Location.Id,
		&u.Sudo,
		&custId,
		&u.Current,
		&passConversion,
	)
	if err != nil {
		err = AuthError
		return
	}

	u.Current = true

	return
}

func (u *CustomerUser) GetKeys() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUserKeysStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(u.Id, strings.Join([]string{api_helpers.AUTH_KEY_TYPE}, ","))
	if err != nil {
		return err
	}

	var keys []ApiCredentials
	for rows.Next() {
		var key ApiCredentials
		err = rows.Scan(
			&key.Key,
			&key.Type,
			&key.DateAdded,
		)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}

	u.Keys = &keys

	return nil
}

func (u *CustomerUser) GetLocation() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getUserLocationStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	u.Location = &CustomerLocation{}
	u.Location.State = &geography.State{}
	u.Location.State.Country = &geography.Country{}

	var lat string
	var lon string
	var stateId, countryId *int
	var state, stateAbbr, country, countryAbbr *string

	err = stmt.QueryRow(u.Id).Scan(
		&u.Location.Id,
		&u.Name,
		&u.Email,
		&u.Location.Address,
		&u.Location.City,
		&u.Location.PostalCode,
		&u.Location.Phone,
		&u.Location.Fax,
		&lat,
		&lon,
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

	if lat != "" && lon != "" {
		u.Location.Latitude, _ = strconv.ParseFloat(lat, 64)
		u.Location.Longitude, _ = strconv.ParseFloat(lon, 64)
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

	if err != nil {
		return err
	}

	return nil
}

func (u *CustomerUser) ResetAuthentication() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getUserAuthenticationKeyStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var key ApiCredentials
	var apiKeyTypeId string

	var datedAddedStr string
	err = stmt.QueryRow(api_helpers.AUTH_KEY_TYPE, u.Id).Scan(
		&key.Key,
		&key.Type,
		&apiKeyTypeId,
		&datedAddedStr,
	)

	key.DateAdded, err = time.Parse("2006-01-02 15:04:05", datedAddedStr)

	if err != nil {
		return err
	}
	stmt, err = db.Prepare(resetUserAuthenticationStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		apiKeyTypeId,
		u.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetCustomerIdFromKey(key string) (id int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerIdFromKeyStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(key).Scan(&id)

	if err != nil {
		err = errors.New("Invalid key")
	}

	return
}

func GetCustomerUserFromKey(key string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUserFromKeyStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var custID, customerID int
	var dbpass, passconv string

	u.Location = &CustomerLocation{}

	err = stmt.QueryRow(api_helpers.AUTH_KEY_TYPE, key).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&dbpass,
		&custID,
		&u.DateAdded,
		&u.Active,
		&u.Location.Id,
		&u.Sudo,
		&customerID,
		&u.Current,
		&passconv,
	)
	if err != nil {
		err = errors.New("Invalid key")
	}

	return
}

func GetCustomerUserById(id string) (u CustomerUser, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustomerUserFromIdStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var custID, customerID int
	var dbpass, passconv string

	u.Location = &CustomerLocation{}

	err = stmt.QueryRow(u.Id).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&dbpass,
		&custID,
		&u.DateAdded,
		&u.Active,
		&u.Location.Id,
		&u.Sudo,
		&customerID,
		&u.Current,
		&passconv,
	)

	if err != nil {
		return
	}

	u.Current = true

	return
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
//  log.Println("renewing authentication key")
//  t := time.Now()

//  log.Printf(renewUserAuthenticationStmt, t.String(), AUTH_KEY_TYPE, u.Id)

//  // Excecute the update statement
//  _, _, err := database.Db.Query(disableTriggerStmt)
//  if err != nil {
//      log.Println(err)
//      return err
//  }
//  _, _, err = database.Db.Query(renewUserAuthenticationStmt, t.String(), AUTH_KEY_TYPE, u.Id)
//  if err != nil {
//      log.Println(err)
//      return err
//  }
//  _, _, err = database.Db.Query(enableTriggerStmt)
//  if err != nil {
//      log.Println(err)
//      return err
//  }
//  return nil
// }
