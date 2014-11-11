package techSupport

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/contact"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type TechSupport struct {
	ID            int             `json:"id,omitempty" xml:"id,omitempty"`
	VehicleMake   string          `json:"vehicleMake,omitempty" xml:"vehicleMake,omitempty"`
	VehicleModel  string          `json:"vehicleModel,omitempty" xml:"vehicleModel,omitempty"`
	VehicleYear   int             `json:"vehicleYear,omitempty" xml:"vehicleYear,omitempty"`
	PurchaseDate  time.Time       `json:"purchaseDate,omitempty" xml:"purchaseDate,omitempty"`
	PurchasedFrom string          `json:"purchasedFrom,omitempty" xml:"purchasedFrom,omitempty"`
	DealerName    string          `json:"dealerName,omitempty" xml:"dealerName,omitempty"`
	ProductCode   string          `json:"productCode,omitempty" xml:"productCode,omitempty"`
	DateCode      string          `json:"dateCode,omitempty" xml:"dateCode,omitempty"`
	Issue         string          `json:"issue,omitempty" xml:"issue,omitempty"`
	Contact       contact.Contact `json:"contact,omitempty" xml:"contact,omitempty"`
}

const (
	fields = ` ts.vehicleMake, ts.vehicleModel, ts.vehicleYear, ts.purchaseDate, ts.purchasedFrom, ts.dealerName, ts.productCode, ts.dateCode, ts.issue, ts.contactID `
)

var (
	createTechSupport          = `insert into TechSupport (vehicleMake, vehicleModel, vehicleYear, purchaseDate, purchasedFrom, dealerName, productCode, dateCode, issue, contactID ) values (?,?,?,?,?,?,?,?,?,?)`
	deleteTechSupport          = `delete from TechSupport where id = ?`
	getTechSupport             = `select ts.id, ` + fields + ` from TechSupport as ts where ts.id = ? `
	getAllTechSupport          = `select ts.id, ` + fields + ` from TechSupport as ts `
	getAllTechSupportByContact = `select ts.id, ` + fields + ` from TechSupport as ts  where contactID = ?`
)

func (t *TechSupport) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getTechSupport)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(t.ID)

	ch := make(chan TechSupport)
	go populateTechSupport(row, ch)
	*t = <-ch
	return
}

func (t *TechSupport) GetByContact() (ts []TechSupport, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllTechSupportByContact)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(t.Contact.ID)

	ch := make(chan []TechSupport)
	go populateTechSupports(rows, ch)
	ts = <-ch
	return
}

func GetAllTechSupport() (ts []TechSupport, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllTechSupport)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()

	ch := make(chan []TechSupport)
	go populateTechSupports(rows, ch)
	ts = <-ch
	return
}

func (t *TechSupport) Create() (err error) {
	//add contact if null
	if t.Contact.ID == 0 {
		t.Contact.Type = "TechSupport"
		if t.Contact.LastName != "" && t.Contact.FirstName != "" && t.Contact.Email != "" {
			err = t.Contact.AddButLessRestrictiveYouFieldValidatinFool()
			if err != nil {
				return
			}
		} else {
			return errors.New("Contact is required.")
		}
		err = t.Contact.Get()
		if err != nil {
			return err
		}

	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(createTechSupport)
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		t.VehicleMake,
		t.VehicleModel,
		t.VehicleYear,
		t.PurchaseDate,
		t.PurchasedFrom,
		t.DealerName,
		t.ProductCode,
		t.DateCode,
		t.Issue,
		t.Contact.ID,
	)

	if err != nil {
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		return
	}
	t.ID = int(id)

	return
}

func (t *TechSupport) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteTechSupport)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID)
	if err != nil {
		return
	}
	return
}

func populateTechSupport(row *sql.Row, ch chan TechSupport) {
	var t TechSupport
	err := row.Scan(
		&t.ID,
		&t.VehicleMake,
		&t.VehicleModel,
		&t.VehicleYear,
		&t.PurchaseDate,
		&t.PurchasedFrom,
		&t.DealerName,
		&t.ProductCode,
		&t.DateCode,
		&t.Issue,
		&t.Contact.ID,
	)
	if err != nil {
		ch <- t
	}
	ch <- t
	return
}

func populateTechSupports(rows *sql.Rows, ch chan []TechSupport) {
	var t TechSupport
	var ts []TechSupport
	for rows.Next() {
		err := rows.Scan(
			&t.ID,
			&t.VehicleMake,
			&t.VehicleModel,
			&t.VehicleYear,
			&t.PurchaseDate,
			&t.PurchasedFrom,
			&t.DealerName,
			&t.ProductCode,
			&t.DateCode,
			&t.Issue,
			&t.Contact.ID,
		)
		if err != nil {
			ch <- ts
		}
		ts = append(ts, t)
	}
	defer rows.Close()
	ch <- ts
	return
}
