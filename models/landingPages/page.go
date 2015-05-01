package landingPage

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)

type LandingPage struct {
	Id, WebsiteID                               int
	NewWindow                                   bool
	Name, PageContent, LinkClasses              string
	StartDate, EndDate                          time.Time
	ConversionID, ConversionLabel, MenuPosition string
	Url                                         *url.URL
	LandingPageDatas                            []LandingPageData
	LandingPageImages                           []LandingPageImage
}

type LandingPageData struct {
	Id, LandingPageID  int
	DataKey, DataValue string
}

type LandingPageImage struct {
	Id, LandingPageID, Sort int
	Url                     *url.URL
}

var (
	GetLandingPageByID = `select lp.id, lp.name, lp.startDate, lp.endDate, lp.url, lp.pageContent, lp.linkClasses, lp.conversionID, lp.conversionLabel, lp.newWindow, lp.menuPosition, lp.websiteID from LandingPage as lp
							Join WebsiteToBrand as wub on wub.WebsiteID = lp.websiteID
							Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where lp.id = ? && lp.startDate <= NOW() && lp.endDate >= NOW() && (ak.api_key = ? && (wub.brandID = ? OR 0=?))
							limit 1`
)

func (lp *LandingPage) Get(dtx *apicontext.DataContext) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(GetLandingPageByID)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(lp.Id, dtx.APIKey, dtx.BrandID, dtx.BrandID)

	err = lp.PopulateLandingPageScan(row)
	if err != nil {
		return
	}
	return
}

func (lp *LandingPage) PopulateLandingPageScan(s database.Scanner) error {
	var pageContent, linkClasses, conversionID, conversionLabel, urlstr *string

	err := s.Scan(
		&lp.Id,
		&lp.Name,
		&lp.StartDate,
		&lp.EndDate,
		&urlstr,
		&pageContent,
		&linkClasses,
		&conversionID,
		&conversionLabel,
		&lp.NewWindow,
		&lp.MenuPosition,
		&lp.WebsiteID,
	)
	if err != nil {
		return err
	}

	if pageContent != nil {
		lp.PageContent = *pageContent
	}
	if linkClasses != nil {
		lp.LinkClasses = *linkClasses
	}
	if conversionID != nil {
		lp.ConversionID = *conversionID
	}
	if conversionLabel != nil {
		lp.ConversionLabel = *conversionLabel
	}
	if urlstr != nil {
		lp.Url, _ = url.Parse(*urlstr)
	}

	return nil
}
