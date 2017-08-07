package testimonials

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
	"time"
)

const (
	testimonialFields = ` t.testimonialID, t.rating, t.title, t.testimonial, t.dateAdded, t.approved, t.active, t.first_name, t.last_name, t.location, t.brandID `
)

var (
	getAllTestimonialsStmt = `select ` + testimonialFields + ` from Testimonial as t 
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.active = 1 && t.approved = 1 order by t.dateAdded desc`
	getTestimonialsByPageStmt = `select ` + testimonialFields + ` from Testimonial as t 
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.active = 1 && t.approved = 1 order by t.dateAdded desc limit ?,?`
	getRandomTestimonalsStmt = `SELECT ` + testimonialFields + ` FROM Testimonial AS t
																WHERE t.brandID = ? && t.active = 1 && t.approved = 1
																ORDER BY Rand()
																LIMIT ?`
	getTestimonialStmt = `select ` + testimonialFields + ` from Testimonial as t WHERE t.testimonialID = ?`
	createTestimonial = `insert into Testimonial (rating, title, testimonial, dateAdded, approved, active, first_name, last_name, location, brandID) values (?,?,?,?,?,?,?,?,?,?)`
	updateTestimonial = `update Testimonial set rating = ?, title = ?, testimonial = ?, approved = ?, active = ?, first_name = ?, last_name = ?, location = ?, brandID = ? where testimonialID = ?`
	deleteTestimonial = `delete from Testimonial where testimonialID = ?`
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
	BrandID   int       `json:"brandId,omitempty" xml:"brandId,omitempty"`
}

func GetAllTestimonials(page int, count int, randomize bool, dtx *apicontext.DataContext) (tests Testimonials, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	if page == 0 && count == 0 {
		stmt, err = db.Prepare(getAllTestimonialsStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	} else if randomize {
		stmt, err = db.Prepare(getRandomTestimonalsStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, count)
	} else {
		stmt, err = db.Prepare(getTestimonialsByPageStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, page, count)
	}

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
			&t.BrandID,
		)
		if err != nil {
			return
		}

		tests = append(tests, t)
	}
	defer rows.Close()
	return
}

func (t *Testimonial) Get(dtx *apicontext.DataContext) error {
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
	err = stmt.QueryRow(dtx.APIKey, dtx.BrandID, dtx.BrandID, t.ID).Scan(
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
		&t.BrandID,
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

	res, err := stmt.Exec(t.Rating, t.Title, t.Content, t.DateAdded, t.Approved, t.Active, t.FirstName, t.LastName, t.Location, t.BrandID)
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
	_, err = stmt.Exec(t.Rating, t.Title, t.Content, t.Approved, t.Active, t.FirstName, t.LastName, t.Location, t.BrandID, t.ID)
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
