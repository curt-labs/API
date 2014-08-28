package customer_new

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/geography"
	// "log"
	"net/http"
	"net/url"
	"strings"
	"time"
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

	customerUserAuth = `select password, id, name, email, date_added, active, isSudo, passwordConverted from CustomerUser
							where email = ?
							&& active = 1
							limit 1`
	updateCustomerUserPass = `update CustomerUser set password = ?, passwordConverted = 1
								where id = ? && active = 1`
)

func (u CustomerUser) UserAuthentication(password string) (cust Customer, err error) {

	err = u.AuthenticateUser(password)
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
	err = stmt.QueryRow(c.Id).Scan(
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
	c.Latitude, err = byteToFloat(lat)
	c.Longitude, err = byteToFloat(lon)
	c.SearchUrl, err = byteToUrl(url)
	c.Logo, err = byteToUrl(logo)
	c.Website, err = byteToUrl(web)
	c.DealerType.MapIcon.MapIcon, err = byteToUrl(icon)
	c.DealerType.MapIcon.MapIconShadow, err = byteToUrl(shadow)
	c.PostalCode, err = byteToString(postalCode)
	c.State.Id, err = byteToInt(stateId)
	c.State.State, err = byteToString(state)
	c.State.Abbreviation, err = byteToString(stateAbbr)
	c.State.Country.Id, err = byteToInt(countryId)
	c.State.Country.Country, err = byteToString(country)
	c.State.Country.Abbreviation, err = byteToString(countryAbbr)
	c.DealerType.MapIcon.Id, err = byteToInt(mapIconId)
	c.DealerType.MapIcon.MapIcon, err = byteToUrl(icon)
	c.DealerType.MapIcon.MapIconShadow, err = byteToUrl(shadow)
	c.MapixCode, err = byteToString(mapixCode)
	c.MapixDescription, err = byteToString(mapixDesc)
	c.SalesRepresentative, err = byteToString(rep)
	c.SalesRepresentativeCode, err = byteToString(repCode)

	parentInt, err := byteToInt(parentId)
	if err != nil {
		return c, err
	}
	if parentInt != 0 {
		par := Customer{Id: parentInt}
		par.GetCustomer_New()
		c.Parent = &par
	}
	return
}

func (u *CustomerUser) AuthenticateUser(pass string) error {

	// password, id, name, email, date_added, active, isSudo, passwordConverted from CustomerUser
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(userCustomer)
	if err != nil {
		return err
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
		&u.Active,
		&u.Sudo,
		&passConversion,
	)
	if err == nil {
		err = errors.New("No user found that matches: " + u.Email)
	}

	// Attempt to compare bcrypt strings
	if bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) != nil {
		// Compare unsuccessful
		enc_pass, err := api_helpers.Md5Encrypt(pass)
		if err != nil {
			return err
		}
		if len(enc_pass) != len(dbPass) || passConversion { //bool
			return errors.New("Invalid password")
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("Failed to encode the password")
		}

		stmtPass, err := db.Prepare(updateCustomerUserPass)
		if err != nil {
			return err
		}
		_, err = stmtPass.Exec(hashedPass, u.Id)
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

	qry, err := database.GetStatement("CustomerUserKeyAuthStmt")
	if database.MysqlError(err) {
		return
	}

	params := struct {
		KeyType string
		Key     string
		Timer   string
	}{}
	params.KeyType = api_helpers.AUTH_KEY_TYPE
	params.Key = key

	t := time.Now()
	t1 := t.Add(time.Duration(-6) * time.Hour)
	params.Timer = t1.String()

	row, res, err := qry.ExecFirst(params)
	if database.MysqlError(err) {
		return
	}
	user_id := res.Map("id")
	name := res.Map("name")
	mail := res.Map("email")
	date := res.Map("date_added")
	active := res.Map("active")
	sudo := res.Map("isSudo")

	if err != nil {
		return
	} else if row == nil {
		err = errors.New("Invalid password")
		return
	}

	//
	// DISABLED: See RenewAuthentication() below
	//
	// resetChan := make(chan int)
	// go func() {
	// 	if resetErr := u.RenewAuthentication(); resetErr != nil {
	// 		err = resetErr
	// 	}
	// 	resetChan <- 1
	// }()

	u.Name = row.Str(name)
	u.Email = row.Str(mail)
	u.Active = row.Int(active) == 1
	u.Sudo = row.Int(sudo) == 1
	u.Current = true
	u.Id = row.Str(user_id)

	//da, _ := time.Parse("2006-01-02 15:04:15", row.Str(date))
	u.DateAdded = row.ForceLocaltime(date)

	//<-resetChan

	return
}

func (u *CustomerUser) GetKeys() error {

	qry, err := database.GetStatement("CustomerUserKeysStmt")
	if database.MysqlError(err) {
		return err
	}
	params := struct {
		User   string
		Except string
	}{}

	params.User = u.Id
	params.Except = strings.Join([]string{api_helpers.AUTH_KEY_TYPE}, ",")

	rows, res, err := qry.Exec(params)
	if database.MysqlError(err) {
		return err
	}

	key := res.Map("api_key")
	typ := res.Map("type")
	dAdded := res.Map("date_added")

	var keys []ApiCredentials
	for _, row := range rows {
		k := ApiCredentials{
			Key:       row.Str(key),
			Type:      row.Str(typ),
			DateAdded: row.ForceLocaltime(dAdded),
		}
		keys = append(keys, k)
	}
	u.Keys = &keys

	return nil
}

func (u *CustomerUser) GetLocation() error {

	qry, err := database.GetStatement("UserLocationStmt")
	if database.MysqlError(err) {
		return err
	}

	row, res, err := qry.ExecFirst(u.Id)
	if database.MysqlError(err) {
		return err
	} else if row == nil {
		return nil
	}

	locationID := res.Map("locationID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	zip := res.Map("postalCode")
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("cty_name")
	country_abbr := res.Map("cty_abbr")
	customerID := res.Map("cust_id")
	isPrimary := res.Map("isprimary")
	shipDefault := res.Map("ShippingDefault")

	l := CustomerLocation{
		Id:              row.Int(locationID),
		Name:            row.Str(name),
		Email:           row.Str(email),
		Address:         row.Str(address),
		City:            row.Str(city),
		PostalCode:      row.Str(zip),
		Phone:           row.Str(phone),
		Fax:             row.Str(fax),
		ContactPerson:   row.Str(contact),
		CustomerId:      row.Int(customerID),
		Latitude:        row.ForceFloat(lat),
		Longitude:       row.ForceFloat(lon),
		IsPrimary:       row.ForceBool(isPrimary),
		ShippingDefault: row.ForceBool(shipDefault),
	}

	ctry := geography.Country{
		Id:           row.Int(countryID),
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	l.State = geography.State_New{
		Id:           row.Int(stateID),
		State:        row.Str(state),
		Abbreviation: row.Str(state_abbr),
		Country:      ctry,
	}

	u.Location = &l

	return nil
}

func (u *CustomerUser) ResetAuthentication() error {

	oldQry, err := database.GetStatement("UserAuthenticationKeyStmt")
	if database.MysqlError(err) {
		return err
	}

	params := struct {
		KeyType string
		User    string
	}{}
	params.KeyType = api_helpers.AUTH_KEY_TYPE
	params.User = u.Id

	// Retrieve the previously declared authentication key for this user
	oldRow, oldRes, err := oldQry.ExecFirst(params)

	if err != nil { // Must be something wrong with the db, lets bail
		return err
	} else if oldRow != nil { // Update the existing with a new date added and key
		old_type_id := oldRes.Map("type_id")

		updateQry, err := database.GetStatement("ResetUserAuthenticationStmt")
		if database.MysqlError(err) {
			return err
		}

		params := struct {
			Now     string
			OldType string
			User    string
		}{}
		updateQry.Bind(&params)

		params.Now = time.Now().String()
		params.OldType = oldRow.Str(old_type_id)
		params.User = u.Id

		// Excecute the update statement
		_, err = updateQry.Raw.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCustomerIdFromKey(key string) (id int, err error) {
	qry, err := database.GetStatement("CustomerIDFromKeyStmt")
	if database.MysqlError(err) {
		return
	}

	row, _, err := qry.ExecFirst(key)
	if database.MysqlError(err) {
		return 0, err
	} else if row == nil {
		return 0, errors.New("Invalid API Key")
	}

	return row.Int(0), nil
}

func GetCustomerUserFromKey(key string) (u CustomerUser, err error) {

	qry, err := database.GetStatement("CustomerUserFromKeyStmt")
	if database.MysqlError(err) {
		return
	}

	params := struct {
		KeyType string
		Key     string
	}{}

	params.KeyType = api_helpers.AUTH_KEY_TYPE
	params.Key = key

	row, res, err := qry.ExecFirst(params)
	if database.MysqlError(err) {
		return
	}
	user_id := res.Map("id")
	name := res.Map("name")
	mail := res.Map("email")
	date := res.Map("date_added")
	active := res.Map("active")
	sudo := res.Map("isSudo")

	if err != nil {
		return
	} else if row == nil {
		err = errors.New("Invalid key")
		return
	}

	u.Name = row.Str(name)
	u.Email = row.Str(mail)
	u.Active = row.Int(active) == 1
	u.Sudo = row.Int(sudo) == 1
	u.Current = true
	u.Id = row.Str(user_id)
	u.DateAdded = row.ForceLocaltime(date)

	return
}

func GetCustomerUserById(id string) (u CustomerUser, err error) {
	qry, err := database.GetStatement("CustomerUserFromId")
	if database.MysqlError(err) {
		return
	}

	row, res, err := qry.ExecFirst(id)
	if database.MysqlError(err) {
		return
	}
	user_id := res.Map("id")
	name := res.Map("name")
	mail := res.Map("email")
	date := res.Map("date_added")
	active := res.Map("active")
	sudo := res.Map("isSudo")

	if err != nil {
		return
	} else if row == nil {
		err = errors.New("Invalid key")
		return
	}

	u.Name = row.Str(name)
	u.Email = row.Str(mail)
	u.Active = row.Int(active) == 1
	u.Sudo = row.Int(sudo) == 1
	u.Current = true
	u.Id = row.Str(user_id)
	u.DateAdded = row.ForceLocaltime(date)

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
