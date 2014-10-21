package contact

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/email"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactsStmt = `select contactID, first_name, last_name, email, phone, subject, message, 
                          createdDate, type, address1, address2, city, state, postalcode, country
                          from Contact limit ?, ?`
	getContactStmt = `select contactID, first_name, last_name, email, phone, subject, message, 
                      createdDate, type, address1, address2, city, state, postalcode, country from Contact where contactID = ?`
	addContactStmt = `insert into Contact(createdDate, first_name, last_name, email, phone, subject, 
                      message, type, address1, address2, city, state, postalcode, country) values (NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	updateContactStmt = `update Contact set first_name = ?, last_name = ?, email = ?, phone = ?, subject = ?, 
                         message = ?, type = ?, address1 = ?, address2 = ?, city = ?, state = ?, postalCode = ?, country = ? where contactID = ?`
	deleteContactStmt = `delete from Contact where contactID = ?`
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

func (c *Contact) Add() error {
	if strings.TrimSpace(c.FirstName) == "" {
		return errors.New("First name is required")
	}
	if strings.TrimSpace(c.LastName) == "" {
		return errors.New("Last name is required")
	}
	if !email.IsEmail(c.Email) {
		return errors.New("Empty or invalid email address")
	}
	if strings.TrimSpace(c.Type) == "" {
		return errors.New("Type can't be empty")
	}
	if strings.TrimSpace(c.Subject) == "" {
		return errors.New("Subject can't be empty")
	}
	if strings.TrimSpace(c.Message) == "" {
		return errors.New("Message can't be empty")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addContactStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		c.FirstName, c.LastName, c.Email, c.Phone, c.Subject, c.Message,
		c.Type, c.Address1, c.Address2, c.City, c.State, c.PostalCode, c.Country)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		c.ID = int(id)
	}

	return nil
}

func (c *Contact) Update() error {
	if c.ID == 0 {
		return errors.New("Invalid Contact ID")
	}
	if strings.TrimSpace(c.FirstName) == "" {
		return errors.New("First name is required")
	}
	if strings.TrimSpace(c.LastName) == "" {
		return errors.New("Last name is required")
	}
	if !email.IsEmail(c.Email) {
		return errors.New("Empty or invalid email address")
	}
	if strings.TrimSpace(c.Type) == "" {
		return errors.New("Type can't be empty")
	}
	if strings.TrimSpace(c.Subject) == "" {
		return errors.New("Subject can't be empty")
	}
	if strings.TrimSpace(c.Message) == "" {
		return errors.New("Message can't be empty")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateContactStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		c.FirstName, c.LastName, c.Email, c.Phone, c.Subject, c.Message, c.Type,
		c.Address1, c.Address2, c.City, c.State, c.PostalCode, c.Country, c.ID)

	return err
}

func (c *Contact) Delete() error {
	if c.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(deleteContactStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(c.ID)

		return err
	}
	return errors.New("Invalid Contact ID")
}
