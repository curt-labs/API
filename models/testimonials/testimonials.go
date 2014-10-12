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
)

type Testimonials []Testimonial
type Testimonial struct {
	ID        int
	Rating    float64
	Title     string
	Content   string
	DateAdded time.Time
	Approved  bool
	Active    bool
	FirstName string
	LastName  string
	Location  string
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
