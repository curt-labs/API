package applicationGuide

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/site_new"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type ApplicationGuide struct {
	ID       int               `json:"id,omitempty" xml:"id,omitempty"`
	Url      string            `json:"url,omitempty" xml:"url,omitempty"`
	Website  site_new.Website  `json:"website,omitempty" xml:"website,omitempty"`
	FileType string            `json:"fileType,omitempty" xml:"fileType,omitempty"`
	Category products.Category `json:"category,omitempty" xml:"category,omitempty"`
	Icon     string            `json:"icon,omitempty" xml:"icon,omitempty"`
}

const (
	fields = ` ag.url, ag.websiteID, ag.fileType, ag.catID, ag.icon `
)

var (
	createApplicationGuide     = `insert into ApplicationGuides (url, websiteID, fileType, catID, icon) values (?,?,?,?,?)`
	deleteApplicationGuide     = `delete from ApplicationGuides where ID = ?`
	getApplicationGuide        = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag left join Categories as c on c.catID = ag.catID where ag.ID = ? `
	getApplicationGuides       = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag left join Categories as c on c.catID = ag.catID`
	getApplicationGuidesBySite = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag left join Categories as c on c.catID = ag.catID where websiteID = ?`
)

func (ag *ApplicationGuide) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getApplicationGuide)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(ag.ID)

	ch := make(chan ApplicationGuide)
	go populateApplicationGuide(row, ch)
	*ag = <-ch
	return
}
func (ag *ApplicationGuide) GetBySite() (ags []ApplicationGuide, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getApplicationGuidesBySite)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(ag.Website.ID)

	ch := make(chan []ApplicationGuide)
	go populateApplicationGuides(rows, ch)
	ags = <-ch
	return
}

func (ag *ApplicationGuide) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createApplicationGuide)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(ag.Url, ag.Website.ID, ag.FileType, ag.Category.ID, ag.Icon)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ag.ID = int(id)
	return
}
func (ag *ApplicationGuide) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteApplicationGuide)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ag.ID)
	if err != nil {
		return err
	}
	return
}

func populateApplicationGuide(row *sql.Row, ch chan ApplicationGuide) {
	var ag ApplicationGuide
	var catID *int
	var icon []byte
	var catName *string
	err := row.Scan(
		&ag.ID,
		&ag.Url,
		&ag.Website.ID,
		&ag.FileType,
		&catID,
		&icon,
		&catName,
	)
	if err != nil {
		ch <- ag
	}
	if catID != nil {
		ag.Category.ID = *catID
	}
	if catName != nil {
		ag.Category.Title = *catName
	}
	if icon != nil {
		ag.Icon = string(icon[:])
	}
	ch <- ag
	return
}

func populateApplicationGuides(rows *sql.Rows, ch chan []ApplicationGuide) {
	var ag ApplicationGuide
	var ags []ApplicationGuide
	var catID *int
	var icon []byte
	var catName *string
	for rows.Next() {
		err := rows.Scan(
			&ag.ID,
			&ag.Url,
			&ag.Website.ID,
			&ag.FileType,
			&catID,
			&icon,
			&catName,
		)
		if err != nil {
			log.Print(err)
			ch <- ags
		}
		if catID != nil {
			ag.Category.ID = *catID
		}
		if catName != nil {
			ag.Category.Title = *catName
		}
		if icon != nil {
			ag.Icon = string(icon[:])
		}
		ags = append(ags, ag)
	}
	ch <- ags
	return
}
