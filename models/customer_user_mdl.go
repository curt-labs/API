package models

import (
	"../helpers/api"
	"../helpers/database"
	"errors"
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

	qry, err := database.GetStatement("UserCustomerStmt")
	if database.MysqlError(err) {
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
	stateID := res.Map("stateID")
	state := res.Map("state")
	state_abbr := res.Map("state_abbr")
	countryID := res.Map("countryID")
	country := res.Map("country_name")
	country_abbr := res.Map("country_abbr")
	dealerTypeId := res.Map("typeID")
	dealerType := res.Map("dealerType")
	typeOnline := res.Map("typeOnline")
	typeShow := res.Map("typeShow")
	typeLabel := res.Map("typeLabel")
	tierID := res.Map("tierID")
	tier := res.Map("tier")
	tierSort := res.Map("tierSort")
	mpx_code := res.Map("mapix_code")
	mpx_desc := res.Map("mapic_desc")
	rep_name := res.Map("rep_name")
	rep_code := res.Map("rep_code")
	parentID := res.Map("parentID")

	sURL, _ := url.Parse(row.Str(search))
	websiteURL, _ := url.Parse(row.Str(site))
	logoURL, _ := url.Parse(row.Str(logo))

	c = Customer{
		Id:            row.Int(customerID),
		Name:          row.Str(name),
		Email:         row.Str(email),
		Address:       row.Str(address),
		Address2:      row.Str(address2),
		City:          row.Str(city),
		PostalCode:    row.Str(zip),
		Phone:         row.Str(phone),
		Fax:           row.Str(fax),
		ContactPerson: row.Str(contact),
		Latitude:      row.ForceFloat(lat),
		Longitude:     row.ForceFloat(lon),
		Website:       websiteURL,
		SearchUrl:     sURL,
		Logo:          logoURL,
		DealerType: DealerType{
			Id:     row.Int(dealerTypeId),
			Type:   row.Str(dealerType),
			Label:  row.Str(typeLabel),
			Online: row.ForceBool(typeOnline),
			Show:   row.ForceBool(typeShow),
		},
		DealerTier: DealerTier{
			Id:   row.Int(tierID),
			Tier: row.Str(tier),
			Sort: row.Int(tierSort),
		},
		SalesRepresentative:     row.Str(rep_name),
		SalesRepresentativeCode: row.Str(rep_code),
		MapixCode:               row.Str(mpx_code),
		MapixDescription:        row.Str(mpx_desc),
	}

	ctry := Country{
		Id:           row.Int(countryID),
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	c.State = &State{
		Id:           row.Int(stateID),
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

	qry, err := database.GetStatement("CustomerUserAuthStmt")
	if database.MysqlError(err) {
		return err
	}

	row, res, err := qry.ExecFirst(u.Email)

	// Check for error while executing query
	if database.MysqlError(err) {
		return err
	}

	// Make sure we have a record for this email
	if row == nil {
		return errors.New("No user found that matches: " + u.Email)
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
			updateQry, _ := database.GetStatement("UpdateCustomerUserPassStmt")
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

	ctry := Country{
		Id:           row.Int(countryID),
		Country:      row.Str(country),
		Abbreviation: row.Str(country_abbr),
	}

	l.State = &State{
		Id:           row.Int(stateID),
		State:        row.Str(state),
		Abbreviation: row.Str(state_abbr),
		Country:      &ctry,
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

	params.KeyType = api_helpers.PRIVATE_KEY_TYPE
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
