package contact

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/email"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactsStmt = `select contactID, first_name, last_name, email, phone, subject, message,
                          createdDate, type, address1, address2, city, state, postalcode, country, brandID
                          from Contact
													where (brandID = ? or 0 = ?)
                          limit ?, ?`
	getContactStmt = `select contactID, first_name, last_name, email, phone, subject, message,
                      createdDate, type, address1, address2, city, state, postalcode, country from Contact where contactID = ?`
	addContactStmt = `insert into Contact(createdDate, first_name, last_name, email, phone, subject,
                      message, type, address1, address2, city, state, postalcode, country, brandID) values (NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	updateContactStmt = `update Contact set first_name = ?, last_name = ?, email = ?, phone = ?, subject = ?,
                         message = ?, type = ?, address1 = ?, address2 = ?, city = ?, state = ?, postalCode = ?, country = ?, brandID = ? where contactID = ?`
	deleteContactStmt = `delete from Contact where contactID = ?`
)

type Contacts []Contact
type Contact struct {
	ID         int         `json:"id,omitempty" xml:"id,omitempty"`
	FirstName  string      `json:"firstName,omitempty" xml:"firstName,omitempty"`
	LastName   string      `json:"lastName,omitempty" xml:"lastName,omitempty"`
	Email      string      `json:"email" xml:"email,omitempty"`
	Phone      string      `json:"phone,omitempty" xml:"phone,omitempty"`
	Subject    string      `json:"subject,omitempty" xml:"subject,omitempty"`
	Message    string      `json:"message,omitempty" xml:"message,omitempty"`
	Created    time.Time   `json:"created,omitempty" xml:"created,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type,omitempty"`
	Address1   string      `json:"address1,omitempty" xml:"address1,omitempty"`
	Address2   string      `json:"address2,omitempty" xml:"address2,omitempty"`
	City       string      `json:"city,omitempty" xml:"city,omitempty"`
	State      string      `json:"state,omitempty" xml:"state,omitempty"`
	PostalCode string      `json:"postalCode,omitempty" xml:"postalCode,omitempty"`
	Country    string      `json:"country,omitempty" xml:"country,omitempty"`
	Brand      brand.Brand `json:"brand,omitempty" xml:"brand,omitempty"`
}
type DealerContact struct {
	Contact
	BusinessName string
	BusinessType customer.DealerType
}

func GetAllContacts(page, count int, dtx *apicontext.DataContext) (Contacts, error) {
	err := database.Init()
	if err != nil {
		return Contacts{}, nil
	}

	stmt, err := database.DB.Prepare(getAllContactsStmt)
	if err != nil {
		return Contacts{}, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.BrandID, dtx.BrandID, page, count)
	if err != nil {
		return Contacts{}, nil
	}

	var contacts Contacts

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
			&c.Brand.ID,
		)
		if err != nil {
			return contacts, err
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

	return contacts, nil
}

func (c *Contact) Get() error {
	if c.ID > 0 {
		err := database.Init()
		if err != nil {
			return err
		}

		stmt, err := database.DB.Prepare(getContactStmt)
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
			if err == sql.ErrNoRows {
				return nil
			}
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

func (c *Contact) Add(dtx *apicontext.DataContext) error {
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(addContactStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		c.FirstName, c.LastName, c.Email, c.Phone, c.Subject, c.Message,
		c.Type, c.Address1, c.Address2, c.City, c.State, c.PostalCode, c.Country, c.Brand.ID)
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

func (c *Contact) AddButLessRestrictiveYouFieldValidatinFool() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(addContactStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		c.FirstName, c.LastName, c.Email, c.Phone, c.Subject, c.Message,
		c.Type, c.Address1, c.Address2, c.City, c.State, c.PostalCode, c.Country, c.Brand.ID)
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(updateContactStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		c.FirstName, c.LastName, c.Email, c.Phone, c.Subject, c.Message, c.Type,
		c.Address1, c.Address2, c.City, c.State, c.PostalCode, c.Country, c.Brand.ID, c.ID)

	return err
}

func (c *Contact) Delete() error {
	if c.ID > 0 {
		err := database.Init()
		if err != nil {
			return err
		}

		stmt, err := database.DB.Prepare(deleteContactStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(c.ID)

		return err
	}
	return errors.New("Invalid Contact ID")
}
