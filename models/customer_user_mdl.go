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
					where email = '%s'
					limit 1`
	updateCustomerUserPassStmt = `update CustomerUser set proper_password = '%s'
						where id = '%s' && active = 1`

	customerUserKeysStmt = `select ak.api_key, akt.type, ak.date_added from ApiKey as ak 
					join ApiKeyType as akt on ak.type_id = akt.id
					where user_id = '%s' && UPPER(akt.type) NOT IN ('%s')`

	userAuthenticationKeyStmt = `select ak.api_key, ak.type_id, akt.type from ApiKey as ak
						join ApiKeyType as akt on ak.type_id = akt.id
						where UPPER(akt.type) = '%s' 
						&& ak.user_id = '%s'`

	// This statement will run the trigger on the
	// ApiKey table to regenerate the api_key column 
	// for the updated record
	resetUserAuthenticationStmt = `update ApiKey as ak
						set ak.date_added = '%s'
						where ak.type_id = '%s' 
						&& ak.user_id = '%s'`

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
				where cu.id = '%s'`

	userLocationStmt = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
					cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
					cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
					s.state, s.abbr as state_abbr, cty.name as cty_name, cty.abbr as cty_abbr
					from CustomerLocations as cl
					left join States as s on cl.stateID = s.stateID
					left join Country as cty on s.countryID = cty.countryID
					join CustomerUser as cu on cl.locationID = cu.locationID
					where cu.id = '%s'`
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

func (u CustomerUser) GetCustomer() (c Customer, err error) {
	row, res, err := database.Db.QueryFirst(userCustomerStmt, u.Id)
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
		State:                   row.Str(state),
		StateAbbreviation:       row.Str(state_abbr),
		Country:                 row.Str(country),
		CountryAbbreviation:     row.Str(country_abbr),
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
		SalesRepresentativeCode: row.Int(rep_code),
		MapixCode:               row.Str(mpx_code),
		MapixDescription:        row.Str(mpx_desc),
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

	row, res, err := database.Db.QueryFirst(customerUserAuthStmt, u.Email)
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
			_, _, _ = database.Db.Query(updateCustomerUserPassStmt, enc_pass, row.Str(user_id))
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

	da, _ := time.Parse("2006-01-02 15:04:15", row.Str(date))
	u.DateAdded = da

	<-resetChan

	return nil
}

func (u *CustomerUser) GetKeys() error {

	rows, res, err := database.Db.Query(customerUserKeysStmt, u.Id, strings.Join([]string{AUTH_KEY_TYPE}, ","))
	if database.MysqlError(err) {
		return err
	}

	key := res.Map("api_key")
	typ := res.Map("type")
	dAdded := res.Map("date_added")

	var keys []ApiCredentials
	for _, row := range rows {

		da, _ := time.Parse("2006-01-02 15:04:15", row.Str(dAdded))

		k := ApiCredentials{
			Key:       row.Str(key),
			Type:      row.Str(typ),
			DateAdded: da,
		}
		keys = append(keys, k)
	}
	u.Keys = &keys

	return nil
}

func (u *CustomerUser) GetLocation() error {

	row, res, err := database.Db.QueryFirst(userLocationStmt, u.Id)
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
		Id:                  row.Int(locationID),
		Name:                row.Str(name),
		Email:               row.Str(email),
		Address:             row.Str(address),
		City:                row.Str(city),
		State:               row.Str(state),
		StateAbbreviation:   row.Str(state_abbr),
		Country:             row.Str(country),
		CountryAbbreviation: row.Str(country_abbr),
		PostalCode:          row.Str(zip),
		Phone:               row.Str(phone),
		Fax:                 row.Str(fax),
		ContactPerson:       row.Str(contact),
		CustomerId:          row.Int(customerID),
		Latitude:            row.ForceFloat(lat),
		Longitude:           row.ForceFloat(lon),
		IsPrimary:           row.ForceBool(isPrimary),
		ShippingDefault:     row.ForceBool(shipDefault),
	}

	u.Location = &l

	return nil
}

func (u *CustomerUser) ResetAuthentication() error {

	// Retrieve the previously declared authentication key for this user
	oldRow, oldRes, err := database.Db.QueryFirst(userAuthenticationKeyStmt, AUTH_KEY_TYPE, u.Id)

	if err != nil { // Must be something wrong with the db, lets bail
		return err
	} else if oldRow != nil { // Update the existing with a new date added and key
		old_type_id := oldRes.Map("type_id")
		t := time.Now()

		// Excecute the update statement
		_, _, err := database.Db.Query(resetUserAuthenticationStmt, t.String(), oldRow.Str(old_type_id), u.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
