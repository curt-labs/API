package showcase

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"errors"
	"net/url"
	"time"
)

const (
	showcaseFields      = ` s.showcaseID, s.rating, s.title, s.text, s.dateAdded, s.approved, s.active, s.first_name, s.last_name, s.location, s.brandID `
	showcaseImageFields = ` si.showcaseImageID, si.path `
)

var (
	getAllShowcasesStmt = `select ` + showcaseFields + ` from Showcase as s
																	Join ApiKeyToBrand as akb on akb.brandID = s.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (s.brandID = ? OR 0=?)) && s.active = 1 && s.approved = 1 order by s.dateAdded desc`
	getShowcaseByPageStmt = `select ` + showcaseFields + ` from Showcase as s
																	Join ApiKeyToBrand as akb on akb.brandID = s.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (s.brandID = ? OR 0=?)) && s.active = 1 && s.approved = 1 order by s.dateAdded desc limit ?,?`
	getRandomShowcasesStmt = `select ` + showcaseFields + ` from Showcase as s
																	Join ApiKeyToBrand as akb on akb.brandID = s.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (s.brandID = ? OR 0=?)) && s.active = 1 && s.approved = 1 order by Rand() limit ?`
	getShowcaseStmt = `select ` + showcaseFields + ` from Showcase as s
																	Join ApiKeyToBrand as akb on akb.brandID = s.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (s.brandID = ? OR 0=?)) && s.showcaseID = ?`
	getShowcaseImages = `select ` + showcaseImageFields + ` from ShowcaseImage si
		join ShowcaseToShowcaseImage sti on sti.showcaseImageID = si.showcaseImageID
		where sti.showcaseID = ?`
	createShowcase = `insert into Showcase (rating, title, text, dateAdded, approved, active, first_name, last_name, location, brandID) values (?,?,?,?,?,?,?,?,?,?)`
	updateShowcase = `update Showcase set rating = ?, title = ?, text = ?, approved = ?, active = ?, first_name = ?, last_name = ?, location = ?, brandID = ? where showcaseID = ?`
	deleteShowcase = `delete from Showcase where showcaseID = ?`

	createImage            = `insert into ShowcaseImage (path) values (?)`
	deleteImage            = `delete from ShowcaseImage where showcaseImageID = ?`
	createImageJoin        = `insert into ShowcaseToShowcaseImage (showcaseID, showcaseImageID) values (?, ?)`
	deleteImageJoin        = `delete from ShowcaseToShowcaseImage where showcaseID = ? && showcaseImageID = ?`
	deleteImageJoinByImage = `delete from ShowcaseToShowcaseImage where showcaseImageID = ?`
	updateImage            = `update ShowcaseImage set path = ? where showcaseImageID = ?`
)

type Showcase struct {
	ID        int        `json:"id,omitempty" xml:"id,omitempty"`
	Rating    float64    `json:"rating,omitempty" xml:"rating,omitempty"`
	Title     string     `json:"title,omitempty" xml:"title,omitempty"`
	Text      string     `json:"text,omitempty" xml:"text,omitempty"`
	DateAdded *time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	Approved  bool       `json:"approved,omitempty" xml:"approved,omitempty"`
	Active    bool       `json:"active,omitempty" xml:"active,omitempty"`
	FirstName string     `json:"firstName,omitempty" xml:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty" xml:"lastName,omitempty"`
	Location  string     `json:"location,omitempty" xml:"location,omitempty"`
	BrandID   int        `json:"brandId,omitempty" xml:"brandId,omitempty"`
	Images    []Image    `json:"images,omitempty" xml:"images,omitempty"`
}

type Image struct {
	ID   int      `json:"id,omitempty" xml:"id,omitempty"`
	Path *url.URL `json:"path,omitempty" xml:"path,omitempty"`
}

type ShowScanner interface {
	Scan(...interface{}) error
}

func GetAllShowcases(page int, count int, randomize bool, dtx *apicontext.DataContext) (shows []Showcase, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	if page == 0 && count == 0 {
		stmt, err = db.Prepare(getAllShowcasesStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	} else if randomize {
		stmt, err = db.Prepare(getRandomShowcasesStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, count)
	} else {
		stmt, err = db.Prepare(getShowcaseByPageStmt)
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
		var s Showcase
		err = s.Scan(rows)
		if err != nil {
			return shows, err
		}
		err = s.GetImages()
		if err != nil {
			return shows, err
		}
		shows = append(shows, s)
	}
	defer rows.Close()
	return
}

func (s *Showcase) Scan(rows ShowScanner) error {
	var title, text, first, last, location *string
	err := rows.Scan(
		&s.ID,
		&s.Rating,
		&title,
		&text,
		&s.DateAdded,
		&s.Approved,
		&s.Active,
		&first,
		&last,
		&location,
		&s.BrandID,
	)
	if err != nil {
		return err
	}
	if title != nil {
		s.Title = *title
	}
	if text != nil {
		s.Text = *text
	}
	if first != nil {
		s.FirstName = *first
	}
	if last != nil {
		s.LastName = *last
	}
	if location != nil {
		s.Location = *location
	}
	return nil
}

func (s *Showcase) Get(dtx *apicontext.DataContext) error {
	if s.ID == 0 {
		return errors.New("Invalid showcase ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getShowcaseStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(dtx.APIKey, dtx.BrandID, dtx.BrandID, s.ID)
	err = s.Scan(row)
	if err != nil {
		return err
	}
	return s.GetImages()
}

func (s *Showcase) GetImages() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getShowcaseImages)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(s.ID)
	if err != nil {
		return err
	}
	var i Image
	var path *string
	for res.Next() {
		err = res.Scan(&i.ID, &path)
		if err != nil {
			return err
		}
		if path != nil {
			i.Path, err = url.Parse(*path)
			if err != nil {
				return err
			}
		}
		s.Images = append(s.Images, i)
	}
	return nil
}

func (s *Showcase) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createShowcase)
	if err != nil {
		return err
	}
	defer stmt.Close()
	now := time.Now()
	s.DateAdded = &now

	res, err := stmt.Exec(s.Rating, s.Title, s.Text, s.DateAdded, s.Approved, s.Active, s.FirstName, s.LastName, s.Location, s.BrandID)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	s.ID = int(id)
	stmt, err = tx.Prepare(createImage)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	joinStmt, err := tx.Prepare(createImageJoin)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer joinStmt.Close()

	for i, _ := range s.Images {
		res, err = stmt.Exec(s.Images[i].Path.String())
		if err != nil {
			tx.Rollback()
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		s.Images[i].ID = int(id)

		_, err = joinStmt.Exec(s.ID, s.Images[i].ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Showcase) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(updateShowcase)
	if err != nil {
		return err
	}
	defer stmt.Close()
	now := time.Now()
	s.DateAdded = &now

	_, err = stmt.Exec(s.Rating, s.Title, s.Text, s.Approved, s.Active, s.FirstName, s.LastName, s.Location, s.BrandID, s.ID)
	if err != nil {
		return err
	}
	for _, i := range s.Images {
		err = i.Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Showcase) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	for _, i := range s.Images {
		err = i.Delete()
		if err != nil {
			return err
		}
	}

	stmt, err := db.Prepare(deleteShowcase)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(s.ID)
	return err
}

func (i *Image) Create(showcaseID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createImage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(i.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	i.ID = int(id)
	stmt, err = tx.Prepare(createImageJoin)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(showcaseID, i.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (i *Image) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteImageJoinByImage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(i.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare(deleteImage)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (i *Image) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateImage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(i.ID, i.Path.String())
	return err
}
