package brand

import (
	"database/sql"
	"errors"
	"net/url"

	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	brandFields           = `ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID`
	getAllBrandsStmt      = `select ` + brandFields + ` from Brand`
	getBrandStmt          = `select ` + brandFields + ` from Brand where ID = ?`
	insertBrandStmt       = `insert into Brand(name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values (?,?,?,?,?,?,?,?)`
	updateBrandStmt       = `update Brand set name = ?, code = ?, logo = ?, logoAlt = ?, formalName = ?, longName = ?, primaryColor = ?, autocareID = ? where ID = ?`
	deleteBrandStmt       = `delete from Brand where ID = ?`
	getCustomerUserBrands = `select b.ID, b.name, b.code, b.logo, b.logoAlt, b.formalName, b.longName, b.primaryColor, b.autocareID
								from Brand as b
								join CustomerToBrand as ctb on ctb.BrandID = b.ID
								join Customer as c on c.cust_id = ctb.cust_id
								where c.cust_id = ?
								group by b.ID`
	getAllWebsitesStmt = `select w.ID, w.description, w.url, wb.brandID from Website as w
							join WebsiteToBrand as wb on w.ID = wb.WebsiteID
							order by wb.brandID, w.ID`
	getBrandWebsitesStmt = `select w.ID, w.description, w.url, wb.brandID from Website as w
							join WebsiteToBrand as wb on w.ID = wb.WebsiteID
							where wb.brandID = ?
							order by w.ID`
)

type Brands []Brand
type Brand struct {
	ID            int       `json:"id" xml:"id,attr"`
	Name          string    `json:"name" xml:"name,attr"`
	Code          string    `json:"code" xml:"code,attr"`
	Logo          *url.URL  `json:"logo" xml:"logo,attr"`
	LogoAlternate *url.URL  `json:"logo_alternate" xml:"logo_alternate,attr"`
	FormalName    string    `json:"formal_name" xml:"formal_name,attr"`
	LongName      string    `json:"long_name" xml:"long_name,attr"`
	PrimaryColor  string    `json:"primary_color" xml:"primary_color,attr"`
	AutocareID    string    `json:"autocareId" xml:"autocareId,attr"`
	Websites      []Website `json:"websites" xml:"websites"`
}

type Website struct {
	ID          int      `json:"id" xml:"id,attr"`
	Description string   `json:"description" xml:"description"`
	URL         *url.URL `json:"url" xml:"url"`
	BrandID     int      `json:"brand_id" xml:"brand_id"`
}

type Scanner interface {
	Scan(...interface{}) error
}

func ScanBrand(res Scanner) (Brand, error) {
	var logo, logoAlt, formal, long, primary, autocare *string
	var b Brand
	err := res.Scan(&b.ID, &b.Name, &b.Code, &logo, &logoAlt, &formal, &long, &primary, &autocare)
	if err != nil {
		return b, err
	}
	if logo != nil {
		b.Logo, err = url.Parse(*logo)
		if err != nil {
			return b, err
		}
	}
	if logoAlt != nil {
		b.LogoAlternate, err = url.Parse(*logoAlt)
		if err != nil {
			return b, err
		}
	}
	if formal != nil {
		b.FormalName = *formal
	}
	if long != nil {
		b.LongName = *long
	}
	if primary != nil {
		b.PrimaryColor = *primary
	}
	if autocare != nil {
		b.AutocareID = *autocare
	}
	return b, err
}

func GetAllBrands() (Brands, error) {
	err := database.Init()
	if err != nil {
		return Brands{}, nil
	}

	stmt, err := database.DB.Prepare(getAllBrandsStmt)
	if err != nil {
		return Brands{}, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return Brands{}, nil
	}

	var brands Brands

	for rows.Next() {
		var b Brand
		b, err = ScanBrand(rows)
		if err != nil {
			return brands, err
		}
		brands = append(brands, b)
	}
	defer rows.Close()

	return brands, nil
}

func GetAllBrandIds() (ids []int, err error) {
	brands, err := GetAllBrands()
	if err != nil {
		return ids, err
	}

	for _, brand := range brands {
		ids = append(ids, brand.ID)
	}
	return ids, err
}

func (b *Brand) Get() error {
	if b.ID == 0 {
		return errors.New("Invalid Brand ID")
	}

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res := stmt.QueryRow(b.ID)
	*b, err = ScanBrand(res)
	if err != nil {
		return err
	}
	return nil
}

func (b *Brand) Create() error {
	if b.Name == "" {
		return errors.New("Brand must have a name.")
	}
	if b.Code == "" {
		return errors.New("Brand must have a code.")
	}

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(insertBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(b.Name, b.Code, b.Logo, b.LogoAlternate, b.FormalName, b.LongName, b.PrimaryColor, b.AutocareID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	b.ID = int(id)
	return err
}

func (b *Brand) Update() error {
	if b.Name == "" {
		return errors.New("Brand must have a name.")
	}
	if b.Code == "" {
		return errors.New("Brand must have a code.")
	}
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(updateBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Name, b.Code, b.Logo, b.LogoAlternate, b.FormalName, b.LongName, b.PrimaryColor, b.AutocareID, b.ID)
	return err
}

func (b *Brand) Delete() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.ID)
	return err
}

func getWebsites(brandID int) ([]Website, error) {
	sites := make([]Website, 0)
	var err error

	err = database.Init()
	if err != nil {
		return sites, err
	}

	var rows *sql.Rows

	if brandID > 0 {
		stmt, err := database.DB.Prepare(getBrandWebsitesStmt)
		if err != nil {
			return sites, err
		}
		defer stmt.Close()

		rows, err = stmt.Query(brandID)
	} else {
		stmt, err := database.DB.Prepare(getAllWebsitesStmt)
		if err != nil {
			return sites, err
		}
		defer stmt.Close()

		rows, err = stmt.Query()
	}
	defer rows.Close()

	if err != nil {
		return sites, err
	}

	for rows.Next() {
		var s Website
		var u *string
		err = rows.Scan(&s.ID, &s.Description, &u, &s.BrandID)
		if err != nil || u == nil {
			continue
		}

		s.URL, err = url.Parse(*u)
		if err != nil {
			continue
		}

		sites = append(sites, s)
	}

	return sites, nil
}

func GetUserBrands(id int) ([]Brand, error) {
	brands := make([]Brand, 0)
	var err error

	err = database.Init()
	if err != nil {
		return brands, err
	}

	stmt, err := database.DB.Prepare(getAllBrandsStmt)
	if err != nil {
		return brands, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return brands, err
	}
	defer rows.Close()

	sites, err := getWebsites(0)
	if err != nil {
		return brands, err
	}

	indexedSites := make(map[int][]Website, 0)
	for _, site := range sites {
		if _, ok := indexedSites[site.BrandID]; !ok {
			indexedSites[site.BrandID] = make([]Website, 0)
		}

		indexedSites[site.BrandID] = append(indexedSites[site.BrandID], site)
	}

	for rows.Next() {
		var b Brand
		b, err = ScanBrand(rows)
		if err != nil {
			return brands, err
		}

		b.Websites = indexedSites[b.ID]
		brands = append(brands, b)
	}

	return brands, nil
}
