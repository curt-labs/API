package models

import (
	"../helpers/api"
	"../helpers/database"
	"errors"
	"net/url"
	"strings"
	"time"
)

const (
	AUTH_KEY_TYPE    = "AUTHENTICATION"
	PUBLIC_KEY_TYPE  = "PUBLIC"
	PRIVATE_KEY_TYPE = "PRIVATE"
)

var (
	customerUserAuthStmt = `select * from CustomerUser
							where email = ?
							&& active = 1
							limit 1`

	customerUserKeyAuthStmt = `select cu.* from CustomerUser as cu
							join ApiKey as ak on cu.id = ak.user_id
							join ApiKeyType as akt on ak.type_id = akt.id
							where UPPER(akt.type) = ? 
							&& ak.api_key = ?
							&& cu.active = 1 && ak.date_added >= ?`

	updateCustomerUserPassStmt = `update CustomerUser set proper_password = ?
							where id = ? && active = 1`

	customerUserKeysStmt = `select ak.api_key, akt.type, ak.date_added from ApiKey as ak 
							join ApiKeyType as akt on ak.type_id = akt.id
							where user_id = ? && UPPER(akt.type) NOT IN (?)`

	userAuthenticationKeyStmt = `select ak.api_key, ak.type_id, akt.type from ApiKey as ak
							join ApiKeyType as akt on ak.type_id = akt.id
							where UPPER(akt.type) = ?
							&& ak.user_id = ?`

	// This statement will run the trigger on the
	// ApiKey table to regenerate the api_key column 
	// for the updated record
	resetUserAuthenticationStmt = `update ApiKey as ak
							set ak.date_added = ?
							where ak.type_id = ? 
							&& ak.user_id = ?`

	// This statement will renew the timer on the
	// authentication API key for the given user.
	// The disabling of the trigger is to turn off the 
	// key regeneration trigger for this table
	enableTriggerStmt           = `SET @disable_trigger = 0;`
	disableTriggerStmt          = `SET @disable_trigger = 1`
	renewUserAuthenticationStmt = `update ApiKey as ak
						join ApiKeyType as akt on ak.type_id = akt.id
						set ak.date_added = ?
						where UPPER(akt.type) = ? && ak.user_id = ?`

	userCustomerStmt = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
				c.latitude, c.longitude, c.searchURL, c.logo, c.website,
				c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
				d_types.type as dealer_type, d_tier.tier as dealer_tier, mpx.code as mapix_code, mpx.description as mapic_desc,
				sr.name as rep_name, sr.code as rep_code, c.parentID
				from Customer as c
				join CustomerUser as cu on c.cust_id = cu.cust_ID
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as d_types on c.dealer_type = d_types.dealer_type
				left join DealerTiers d_tier on c.tier = d_tier.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where cu.id = ?`

	userLocationStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
					cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
					cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
					s.state, s.abbr as state_abbr, cty.name as cty_name, cty.abbr as cty_abbr
					from CustomerLocations as cl
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					join CustomerUser as cu on cl.locationID = cu.locationID
					where cu.id = ?`

	customerIDFromKeyStmt = `select c.customerID from Customer as c
					join CustomerUser as cu on c.cust_id = cu.cust_ID
					join ApiKey as ak on cu.id = ak.user_id
					where ak.api_key = ?
					limit 1`

	customerUserFromKeyStmt = `select cu.* from CustomerUser as cu
					join ApiKey as ak on cu.id = ak.user_id
					join ApiKeyType as akt on ak.type_id = akt.id
					where akt.type = ? && ak.api_key = ?
					limit 1`
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

	qry, err := database.Db.Prepare(userCustomerStmt)
	if err != nil {
		return
	}

	row, res, err := qry.ExecFirst(u.Id)
	if database.MysqlError(err) {
		return
	}

	customerID := res.Map("customerID")
	name := res.Map("name")
	email := res.Map("email")
	address := res.Map("address")
	address2 := res.Map("address2")
	city := res.Map("city")
	phone := res.Map("phone")
	fax := res.Map("fax")
	contact := res.Map("contact_person")
	lat := res.Map("latitude")
	lon := res.Map("longitude")
	search := res.Map("searchURL")
	site := res.Map("website")
	logo := res.Map("logo")
	zip := res.Map("postal_code")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealer_type := res.Map("dealer_type")
	dealer_tier := res.Map("dealer_tier")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	sURL, _ := url.Parse(row.Str(search))
	websiteURL, _ := url.Parse(row.Str(site))
	logoURL, _ := url.Parse(row.Str(logo))

	c = Customer{
		Id:                      row.Int(customerID),
		Name:                    row.Str(name),
		Email:                   row.Str(email),
		Address:                 row.Str(address),
		Address2:                row.Str(address2),
		City:                    row.Str(city),
		PostalCode:              row.Str(zip),
		Phone:                   row.Str(phone),
		Fax:                     row.Str(fax),
		ContactPerson:           row.Str(contact),
		Latitude:                row.ForceFloat(lat),
		Longitude:               row.ForceFloat(lon),
		Website:                 websiteURL,
		SearchUrl:               sURL,
		Logo:                    logoURL,
		DealerType:              row.Str(dealer_type),
		DealerTier:              row.Str(dealer_tier),
		SalesRepresentative:     row.Str(rep_name),
		SalesRepresentativeCode: row.Str(rep_code),
		MapixCode:               row.Str(mpx_code),
		MapixDescription:        row.Str(mpx_desc),
	}

	ctry := Country{
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	c.State = &State{
		State:        row.Str(state),
		Abbreviation: row.Str(state_abbr),
		Country:      &ctry,
	}

	locationChan := make(chan int)
	go func() {
		if locErr := c.GetLocations(); locErr != nil {
			err = locErr
		}
		locationChan <- 1
	}()

	if row.Int(parentID) != 0 {
		parent := Customer{
			Id: row.Int(parentID),
		}
		if err = parent.GetCustomer(); err == nil {
			c.Parent = &parent
		}
	}

	<-locationChan

	return
}

func (u *CustomerUser) AuthenticateUser(pass string) error {

	enc_pass, err := api_helpers.Md5Encrypt(pass)
	if err != nil {
		return err
	}

	qry, err := database.Db.Prepare(customerUserAuthStmt)
	if err != nil {
		return err
	}

	row, res, err := qry.ExecFirst(u.Email)
	if database.MysqlError(err) {
		return err
	}
	pwd := res.Map("password")
	prop_pass := res.Map("proper_password")
	user_id := res.Map("id")
	name := res.Map("name")
	mail := res.Map("email")
	date := res.Map("date_added")
	active := res.Map("active")
	sudo := res.Map("isSudo")

	prop := row.Str(prop_pass)

	if err != nil {
		return err
	} else if prop == "" {
		if len(enc_pass) != len(row.Str(pwd)) {
			err = errors.New("Invalid password")
			return err
		} else {
			updateQry, _ := database.Db.Prepare(updateCustomerUserPassStmt)
			dbUserId := row.Str(user_id)
			params := struct {
				pass *string
				user *string
			}{&enc_pass, &dbUserId}

			updateQry.Bind(&params)
			if updateQry != nil {
				_, _ = updateQry.Raw.Run()
			}
		}
	} else if !strings.EqualFold(prop, enc_pass) {
		err = errors.New("Invalid password")
		return err
	}

	resetChan := make(chan int)
	go func() {
		if resetErr := u.ResetAuthentication(); resetErr != nil {
			err = resetErr
		}
		resetChan <- 1
	}()

	u.Name = row.Str(name)
	u.Email = row.Str(mail)
	u.Active = row.Int(active) == 1
	u.Sudo = row.Int(sudo) == 1
	u.Current = true
	u.Id = row.Str(user_id)

	//da, _ := time.Parse("2006-01-02 15:04:15", row.Str(date))
	u.DateAdded = row.ForceLocaltime(date)

	<-resetChan

	return nil
}

func AuthenticateUserByKey(key string) (u CustomerUser, err error) {

	qry, err := database.Db.Prepare(customerUserKeyAuthStmt)
	if err != nil {
		return
	}

	params := struct {
		key_type string
		key      string
		timer    string
	}{}
	params.key_type = AUTH_KEY_TYPE
	params.key = key

	t := time.Now()
	t1 := t.Add(time.Duration(-6) * time.Hour)
	params.timer = t1.String()

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

	qry, err := database.Db.Prepare(customerUserKeysStmt)
	if err != nil {
		return err
	}
	params := struct {
		User   string
		Except string
	}{}

	params.User = u.Id
	params.Except = strings.Join([]string{AUTH_KEY_TYPE}, ",")

	rows, res, err := qry.Exec(params)
	if database.MysqlError(err) {
		return err
	}

	key := res.Map("api_key")
	typ := res.Map("type")
	dAdded := res.Map("date_added")

	var keys []ApiCredentials
	for _, row := range rows {

		//da, _ := time.Parse("2006-01-02 15:04:15", row.Str(dAdded))

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

	qry, err := database.Db.Prepare(userLocationStmt)
	if err != nil {
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
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
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

	ctry := Country{
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	l.State = &State{
		State:        row.Str(state),
		Abbreviation: row.Str(state_abbr),
		Country:      &ctry,
	}

	u.Location = &l

	return nil
}

func (u *CustomerUser) ResetAuthentication() error {

	oldQry, err := database.Db.Prepare(userAuthenticationKeyStmt)
	if err != nil {
		return err
	}

	params := struct {
		KeyType string
		User    string
	}{}
	params.KeyType = AUTH_KEY_TYPE
	params.User = u.Id

	// Retrieve the previously declared authentication key for this user
	oldRow, oldRes, err := oldQry.ExecFirst(params)

	if err != nil { // Must be something wrong with the db, lets bail
		return err
	} else if oldRow != nil { // Update the existing with a new date added and key
		old_type_id := oldRes.Map("type_id")

		updateQry, err := database.Db.Prepare(resetUserAuthenticationStmt)
		if err != nil {
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
	qry, err := database.Db.Prepare(customerIDFromKeyStmt)
	if err != nil {
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

	qry, err := database.Db.Prepare(customerUserFromKeyStmt)
	if err != nil {
		return
	}

	params := struct {
		KeyType string
		Key     string
	}{}

	params.KeyType = PRIVATE_KEY_TYPE
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
