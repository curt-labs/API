package customer_new

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/goacesapi/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// type Location struct {
// 	ID            int
// 	Name          string
// 	Email         string
// 	Address       string
// 	City          string
// 	StateID       int
// 	Phone         string
// 	Fax           string
// 	ContactPerson string
// 	GeoLocation
// 	CustomerID      int
// 	PostalCode      string
// 	IsPrimary       bool
// 	ShippingDefault bool
// }
// type Locations []Location

// type GeoLocation struct {
// 	Latitude, Longitude, Distance float64
// }

type CustomerLocations []CustomerLocation

var (
	getLocation    = "SELECT locationID, name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault FROM CustomerLocations WHERE locationID= ? "
	getLocations   = "SELECT locationID, name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault FROM CustomerLocations"
	createLocation = "INSERT INTO CustomerLocations (name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	updateLocation = "UPDATE CustomerLocations SET name = ?, address = ?,  city = ?,  stateID = ?, email = ?,  phone = ?,  fax = ?,  latitude = ?,  longitude = ?,  cust_id = ?, contact_person = ?,  isprimary = ?, postalCode = ?, ShippingDefault = ? WHERE locationID = ?"
	deleteLocation = "DELETE FROM CustomerLocations WHERE locationID = ?"
)

func (l *CustomerLocation) Get() error {
	redis_key := "goapi:customers:location:" + strconv.Itoa(l.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &l)
		return err
	}

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
	err = stmt.QueryRow(l.Id).Scan(&l.Id, &l.Name, &l.Address, &l.City, &l.State, &l.Email, &l.Phone, &l.Fax, &l.Latitude, &l.Longitude, &l.CustomerId, &l.ContactPerson, &l.IsPrimary, &l.PostalCode, &l.ShippingDefault)
	if err != nil {
		return err
	}

	go redis.Setex(redis_key, l, 86400)
	return err

}

func GetAllLocations() (CustomerLocations, error) {
	var ls CustomerLocations
	var err error
	redis_key := "goapi:customers:locations"
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
	res, err := stmt.Query()
	for res.Next() {
		var l CustomerLocation
		err = res.Scan(&l.Id, &l.Name, &l.Address, &l.City, &l.State, &l.Email, &l.Phone, &l.Fax, &l.Latitude, &l.Longitude, &l.CustomerId, &l.ContactPerson, &l.IsPrimary, &l.PostalCode, &l.ShippingDefault)
		if err != nil {
			return ls, err
		}
		ls = append(ls, l)
	}
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
	res, err := stmt.Exec(l.Name, l.Address, l.City, l.State, l.Email, l.Phone, l.Fax, l.Latitude, l.Longitude, l.CustomerId, l.ContactPerson, l.IsPrimary, l.PostalCode, l.ShippingDefault)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	l.Id = int(id)
	if err != nil {
		tx.Rollback()
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

	_, err = stmt.Exec(l.Name, l.Address, l.City, l.State, l.Email, l.Phone, l.Fax, l.Latitude, l.Longitude, l.CustomerId, l.ContactPerson, l.IsPrimary, l.PostalCode, l.ShippingDefault, l.Id)
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
	_, err = stmt.Exec(l.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
