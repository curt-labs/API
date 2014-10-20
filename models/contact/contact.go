package contact

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactsStmt = `select contactID, first_name, last_name, email, phone, subject, message, 
                          createdDate, type, address1, address2, city, state, postalcode, country
                          from Contact limit ?, ?`
	getContactStmt = `select contactID, first_name, last_name, email, phone, subject, message, 
                      createdDate, type, address1, address2, city, state, postalcode, country from Contact where contactID = ?`
	getAllContactTypesStmt     = `select contactTypeID, name from ContactType`
	getContactTypeStmt         = `select contactTypeID, name from ContactType where contactTypeID = ?`
	getAllContactReceiversStmt = `select contactReceiverID, first_name, last_name, email from ContactReceiver`
	getContactReceiverStmt     = `select contactReceiverID, first_name, last_name, email from ContactReceiver where contactReceiverID = ?`
)

type Contacts []Contact
type Contact struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Subject    string
	Message    string
	Created    time.Time
	Type       string
	Address1   string
	Address2   string
	City       string
	State      string
	PostalCode string
	Country    string
}

type ContactTypes []ContactType
type ContactType struct {
	ID   int
	Name string
}

type ContactReceivers []ContactReceiver
type ContactReceiver struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

func GetAllContacts(page, count int) (contacts Contacts, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContactsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(page, count)
	if err != nil {
		return
	}

	var addr1, addr2, city, state, postalCode, country *string

	for rows.Next() {
		var c Contact
		err = rows.Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.Phone,
			&c.Subject,
			&c.Message,
			&c.Created,
			&c.Type,
			&addr1,
			&addr2,
			&city,
			&state,
			&postalCode,
			&country,
		)
		if err != nil {
			return
		}
		if addr1 != nil {
			c.Address1 = *addr1
		}
		if addr2 != nil {
			c.Address2 = *addr2
		}
		if city != nil {
			c.City = *city
		}
		if state != nil {
			c.State = *state
		}
		if postalCode != nil {
			c.PostalCode = *postalCode
		}
		if country != nil {
			c.Country = *country
		}
		contacts = append(contacts, c)
	}
	defer rows.Close()

	return
}

func GetAllContactTypes() (types ContactTypes, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContactTypesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var ct ContactType
		err = rows.Scan(
			&ct.ID,
			&ct.Name,
		)
		if err != nil {
			return
		}
		types = append(types, ct)
	}
	defer rows.Close()

	return
}

func GetAllContactReceivers() (receivers ContactReceivers, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContactReceiversStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var cr ContactReceiver
		err = rows.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return
		}
		receivers = append(receivers, cr)
	}
	defer rows.Close()

	return
}

func (c *Contact) Get() error {
	if c.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getContactStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		var addr1, addr2, city, state, postalCode, country *string

		err = stmt.QueryRow(c.ID).Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.Phone,
			&c.Subject,
			&c.Message,
			&c.Created,
			&c.Type,
			&addr1,
			&addr2,
			&city,
			&state,
			&postalCode,
			&country,
		)
		if err != nil {
			return err
		}
		if addr1 != nil {
			c.Address1 = *addr1
		}
		if addr2 != nil {
			c.Address2 = *addr2
		}
		if city != nil {
			c.City = *city
		}
		if state != nil {
			c.State = *state
		}
		if postalCode != nil {
			c.PostalCode = *postalCode
		}
		if country != nil {
			c.Country = *country
		}
		return err
	}
	return errors.New("Invalid Contact ID")
}

func (ct *ContactType) Get() error {
	if ct.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getContactTypeStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(ct.ID).Scan(
			&ct.ID,
			&ct.Name,
		)
		return err
	}
	return errors.New("Invalid ContactType ID")
}

func (cr *ContactReceiver) Get() error {
	if cr.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getContactReceiverStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(cr.ID).Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)

		return err
	}
	return errors.New("Invalid ContactReceiver ID")
}
