package testimonials

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllTestimonialsStmt = `select * from Testimonial`
	getTestimonialStmt     = `select * from Testimonial where testimonialID = ?`
	createTestimonial      = `insert into Testimonial (rating, title, testimonial, dateAdded, approved, active, first_name, last_name, location) values (?,?,?,?,?,?,?,?,?)`
	updateTestimonial      = `update Testimonial set rating = ?, title = ?, testimonial = ?, approved = ?, active = ?, first_name = ?, last_name = ?, location = ? where testimonialID = ?`
	deleteTestimonial      = `delete from Testimonial where testimonialID = ?`
)

type Testimonials []Testimonial
type Testimonial struct {
	ID        int       `json:"id,omitempty" xml:"id,omitempty"`
	Rating    float64   `json:"rating,omitempty" xml:"rating,omitempty"`
	Title     string    `json:"title,omitempty" xml:"title,omitempty"`
	Content   string    `json:"content,omitempty" xml:"content,omitempty"`
	DateAdded time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	Approved  bool      `json:"approved,omitempty" xml:"approved,omitempty"`
	Active    bool      `json:"active,omitempty" xml:"active,omitempty"`
	FirstName string    `json:"firstName,omitempty" xml:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty" xml:"lastName,omitempty"`
	Location  string    `json:"location,omitempty" xml:"location,omitempty"`
}

func GetAllTestimonials() (tests Testimonials, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllTestimonialsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var t Testimonial
		err = rows.Scan(
			&t.ID,
			&t.Rating,
			&t.Title,
			&t.Content,
			&t.DateAdded,
			&t.Approved,
			&t.Active,
			&t.FirstName,
			&t.LastName,
			&t.Location,
		)
		if err != nil {
			return
		}

		tests = append(tests, t)
	}
	defer rows.Close()

	return
}

func (t *Testimonial) Get() error {
	if t.ID == 0 {
		return errors.New("Invalid testimonial ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getTestimonialStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(t.ID).Scan(
		&t.ID,
		&t.Rating,
		&t.Title,
		&t.Content,
		&t.DateAdded,
		&t.Approved,
		&t.Active,
		&t.FirstName,
		&t.LastName,
		&t.Location,
	)

	return err
}

func (t *Testimonial) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()
	t.DateAdded = time.Now()
	res, err := stmt.Exec(t.Rating, t.Title, t.Content, t.DateAdded, t.Approved, t.Active, t.FirstName, t.LastName, t.Location)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	t.ID = int(id)
	return nil
}

func (t *Testimonial) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()
	t.DateAdded = time.Now()
	_, err = stmt.Exec(t.Rating, t.Title, t.Content, t.Approved, t.Active, t.FirstName, t.LastName, t.Location, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *Testimonial) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID)
	if err != nil {
		return err
	}
	return nil
}
