package applicationGuide

import (
	"database/sql"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/site"
	_ "github.com/go-sql-driver/mysql"
)

type ApplicationGuide struct {
	ID       int               `json:"id,omitempty" xml:"id,omitempty"`
	Url      string            `json:"url,omitempty" xml:"url,omitempty"`
	Website  site.Website      `json:"website,omitempty" xml:"website,omitempty"`
	FileType string            `json:"fileType,omitempty" xml:"fileType,omitempty"`
	Category products.Category `json:"category,omitempty" xml:"category,omitempty"`
	Icon     string            `json:"icon,omitempty" xml:"icon,omitempty"`
}

const (
	fields = ` ag.url, ag.websiteID, ag.fileType, ag.catID, ag.icon `
)

var (
	createApplicationGuide = `insert into ApplicationGuides (url, websiteID, fileType, catID, icon, brandID) values (?,?,?,?,?,?)`
	deleteApplicationGuide = `delete from ApplicationGuides where ID = ?`
	getApplicationGuide    = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
										left join Categories as c on c.catID = ag.catID
										where ag.ID = ? `
	getApplicationGuidesBySite = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
										left join Categories as c on c.catID = ag.catID
										Join ApiKeyToBrand as akb on akb.brandID = ag.brandID
										Join ApiKey as ak on akb.keyID = ak.id
										where (ak.api_key = ? && (ag.brandID = ? OR 0=?)) && websiteID = ?`
)

func (ag *ApplicationGuide) Get(dtx *apicontext.DataContext) error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getApplicationGuide)
	if err != nil {
		return err
	}

	defer stmt.Close()
	row := stmt.QueryRow(ag.ID)

	ch := make(chan ApplicationGuide)
	go populateApplicationGuide(row, ch)
	*ag = <-ch
	return nil
}

func (ag *ApplicationGuide) GetBySite(dtx *apicontext.DataContext) ([]ApplicationGuide, error) {
	err := database.Init()
	if err != nil {
		return nil, err
	}

	stmt, err := database.DB.Prepare(getApplicationGuidesBySite)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, ag.Website.ID)

	var ags []ApplicationGuide

	ch := make(chan []ApplicationGuide)
	go populateApplicationGuides(rows, ch)
	ags = <-ch
	return ags, nil
}

func (ag *ApplicationGuide) Create(dtx *apicontext.DataContext) error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createApplicationGuide)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(ag.Url, ag.Website.ID, ag.FileType, ag.Category.CategoryID, ag.Icon, dtx.BrandID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ag.ID = int(id)
	return nil
}

func (ag *ApplicationGuide) Delete() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteApplicationGuide)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(ag.ID)
	if err != nil {
		return err
	}
	return nil
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
		ag.Category.CategoryID = *catID
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
			ch <- ags
		}
		if catID != nil {
			ag.Category.CategoryID = *catID
		}
		if catName != nil {
			ag.Category.Title = *catName
		}
		if icon != nil {
			ag.Icon = string(icon[:])
		}
		ags = append(ags, ag)
	}
	defer rows.Close()

	ch <- ags
	return
}
