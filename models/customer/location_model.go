package customer

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/conversions"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type CustomerLocations []CustomerLocation

var (
	getLocation  = "SELECT locationID, name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault FROM CustomerLocations WHERE locationID= ? "
	getLocations = `SELECT cl.locationID, cl.name, cl.address, cl.city, cl.stateID, cl.email,cl.phone, cl.fax, cl.latitude, cl.longitude, cl.cust_id, cl.contact_person, cl.isprimary, cl.postalCode, cl.ShippingDefault 
			FROM CustomerLocations as cl
			join CustomerToBrand as ctb on ctb.cust_id = cl.cust_id
			join apiKeyToBrand as akb on akb.brandID = ctb.brandID
			join apiKey as ak on ak.id = akb.keyID	
			where ak.api_key = ? && (ctb.BrandID = ? or 0 = ?)`
	createLocation = "INSERT INTO CustomerLocations (name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	updateLocation = "UPDATE CustomerLocations SET name = ?, address = ?,  city = ?,  stateID = ?, email = ?,  phone = ?,  fax = ?,  latitude = ?,  longitude = ?,  cust_id = ?, contact_person = ?,  isprimary = ?, postalCode = ?, ShippingDefault = ? WHERE locationID = ?"
	deleteLocation = "DELETE FROM CustomerLocations WHERE locationID = ? "
)

func (l *CustomerLocation) Get() error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getLocation)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var name, address, city, email, phone, fax, contactPerson, postal []byte
	err = stmt.QueryRow(l.Id).Scan(
		&l.Id,
		&name,
		&address,
		&city,
		&l.State.Id,
		&email,
		&phone,
		&fax,
		&l.Latitude,
		&l.Longitude,
		&l.CustomerId,
		&contactPerson,
		&l.IsPrimary,
		&postal,
		&l.ShippingDefault,
	)
	if err != nil {
		return err
	}
	l.Name, err = conversions.ByteToString(name)
	l.Address, err = conversions.ByteToString(address)
	l.City, err = conversions.ByteToString(city)
	l.Email, err = conversions.ByteToString(email)
	l.Phone, err = conversions.ByteToString(phone)
	l.Fax, err = conversions.ByteToString(fax)
	l.ContactPerson, err = conversions.ByteToString(contactPerson)
	l.PostalCode, err = conversions.ByteToString(postal)
	if err != nil {
		return err
	}

	return err
}

func GetAllLocations(apikey string, brandID int) (CustomerLocations, error) {
	var ls CustomerLocations
	var err error
	redis_key := "customers:locations"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ls)
		return ls, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ls, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getLocations)
	if err != nil {
		return ls, err
	}
	defer stmt.Close()
	res, err := stmt.Query(apikey, brandID, brandID)
	var name, address, city, email, phone, fax, contactPerson, postal []byte
	for res.Next() {
		var l CustomerLocation
		err = res.Scan(
			&l.Id,
			&name,
			&address,
			&city,
			&l.State.Id,
			&email,
			&phone,
			&fax,
			&l.Latitude,
			&l.Longitude,
			&l.CustomerId,
			&contactPerson,
			&l.IsPrimary,
			&postal,
			&l.ShippingDefault,
		)
		if err != nil {
			return ls, err
		}
		l.Name, err = conversions.ByteToString(name)
		l.Address, err = conversions.ByteToString(address)
		l.City, err = conversions.ByteToString(city)
		l.Email, err = conversions.ByteToString(email)
		l.Phone, err = conversions.ByteToString(phone)
		l.Fax, err = conversions.ByteToString(fax)
		l.ContactPerson, err = conversions.ByteToString(contactPerson)
		l.PostalCode, err = conversions.ByteToString(postal)
		if err != nil {
			return ls, err
		}
		ls = append(ls, l)
	}
	defer res.Close()
	go redis.Setex(redis_key, ls, 86400)
	return ls, err
}

func (l *CustomerLocation) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createLocation)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		l.Name,
		l.Address,
		l.City,
		l.State.Id,
		l.Email,
		l.Phone,
		l.Fax,
		l.Latitude,
		l.Longitude,
		l.CustomerId,
		l.ContactPerson,
		l.IsPrimary,
		l.PostalCode,
		l.ShippingDefault,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	l.Id = int(id)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
func (l *CustomerLocation) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateLocation)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		l.Name,
		l.Address,
		l.City,
		l.State.Id,
		l.Email,
		l.Phone,
		l.Fax,
		l.Latitude,
		l.Longitude,
		l.CustomerId,
		l.ContactPerson,
		l.IsPrimary,
		l.PostalCode,
		l.ShippingDefault,
		l.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (l *CustomerLocation) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteLocation)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(l.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
