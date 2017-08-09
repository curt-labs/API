package warranty

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/contact"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
	"time"
)

type Warranty struct {
	ID            int             `json:"id,omitempty" xml:"id,omitempty"`
	PartNumber    int             `json:"partNumber,omitempty" xml:"partNumber,omitempty"`
	Date          *time.Time      `json:"date,omitempty" xml:"date,omitempty"`
	SerialNumber  string          `json:"serialNumber,omitempty" xml:"serialNumber,omitempty"`
	Approved      bool            `json:"approved,omitempty" xml:"approved,omitempty"`
	Contact       contact.Contact `json:"contact,omitempty" xml:"contact,omitempty"`
	OldPartNumber string          `json:"oldPartNumber,omitempty" xml:"oldPartNumber,omitempty"`
}

const (
	fields = ` w.partNumber, w.date, w.serialNumber, w.approved, w.contactID `
)

var (
	createWarranty       = `insert into Warranty (partNumber, date, serialNumber, approved, contactID) values (?,?,?,?,?)`
	deleteWarranty       = `delete from Warranty where id = ?`
	getWarranty          = `select w.id, ` + fields + ` from Warranty as w where w.id = ?`
	getWarrantyByContact = `select w.id, ` + fields + ` from Warranty as w where w.contactID = ?`
	getAllWarranties     = `select w.id, ` + fields + ` from Warranty as w
							join Part as p on p.partID = w.partNumber
							where p.brandID = ?`
)

func (w *Warranty) Create() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createWarranty)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//add contact if null
	if w.Contact.ID == 0 {
		w.Contact.Type = "Warranty"
		if w.Contact.LastName != "" && w.Contact.FirstName != "" && w.Contact.Email != "" {
			err = w.Contact.AddButLessRestrictiveYouFieldValidatinFool()
			if err != nil {
				return
			}
		} else {
			return errors.New("Contact is required.")
		}
	}
	if w.Date == nil {
		date := time.Now()
		w.Date = &date
	}
	w.Approved = false
	res, err := stmt.Exec(
		w.PartNumber,
		w.Date,
		w.SerialNumber,
		w.Approved,
		w.Contact.ID,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	w.ID = int(id)
	return
}

func (w *Warranty) Delete() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteWarranty)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.ID)
	if err != nil {
		return err
	}
	return
}

func (w *Warranty) Get() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getWarranty)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(w.ID)
	ch := make(chan Warranty)
	go populateWarranty(row, ch)
	*w = <-ch
	return
}

func (w *Warranty) GetByContact() (ws []Warranty, err error) {
	err = database.Init()
	if err != nil {
		return
	}

	stmt, err := database.DB.Prepare(getWarrantyByContact)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(w.Contact.ID)
	if err != nil {
		return
	}

	ch := make(chan []Warranty)
	go populateWarranties(rows, ch)
	ws = <-ch
	return
}

func GetAllWarranties(dtx *apicontext.DataContext) (ws []Warranty, err error) {
	err = database.Init()
	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWarranties)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(dtx.BrandID)
	if rows.Next() == false {
		err = sql.ErrNoRows
		return ws, err
	}
	if err != nil {
		return ws, err
	}
	ch := make(chan []Warranty)
	go populateWarranties(rows, ch)
	ws = <-ch
	return
}
func populateWarranty(row *sql.Row, ch chan Warranty) {
	var w Warranty
	err := row.Scan(
		&w.ID,
		&w.PartNumber,
		&w.Date,
		&w.SerialNumber,
		&w.Approved,
		&w.Contact.ID,
	)
	if err != nil {
		ch <- w
	}
	ch <- w
}
func populateWarranties(rows *sql.Rows, ch chan []Warranty) {
	var w Warranty
	var ws []Warranty
	for rows.Next() {
		err := rows.Scan(
			&w.ID,
			&w.PartNumber,
			&w.Date,
			&w.SerialNumber,
			&w.Approved,
			&w.Contact.ID,
		)
		if err != nil {
			ch <- ws
		}
		ws = append(ws, w)
	}
	defer rows.Close()
	ch <- ws
}
