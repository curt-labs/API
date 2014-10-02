package site

import (
	"database/sql"
	// "encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	// "strconv"
	"log"
	"time"
)

type ContentPage struct {
	SiteContent SiteContent
	Menu        Menu
	Revision    SiteContentRevision
}

type SiteContent struct {
	Id                    int       `json:"id,omitempty" xml:"id,omitempty"`
	Type                  string    `json:"type,omitempty" xml:"type,omitempty"`
	Title                 string    `json:"title,omitempty" xml:"title,omitempty"`
	CreatedDate           time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	LastModified          time.Time `json:"lastModified,omitempty" xml:"lastModified,omitempty"`
	MetaTitle             string    `json:"metaTitle,omitempty" xml:"v,omitempty"`
	MetaDescription       string    `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	IsPrimary             bool      `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Keywords              string    `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Published             bool      `json:"published,omitempty" xml:"published,omitempty"`
	Active                bool      `json:"active,omitempty" xml:"active,omitempty"`
	Slug                  string    `json:"slug, omitempty" xml:"slug,omitempty"`
	RequireAuthentication bool      `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	Canonical             string    `json:"canonical,omitempty" xml:"canonical,omitempty"`
}

type SiteContents []SiteContent

type SiteContentRevision struct {
	Id          int       `json:"id,omitempty" xml:"id,omitempty"`
	ContentId   int       `json:"contentId,omitempty" xml:"contentId,omitempty"`
	Text        string    `json:"name,omitempty" xml:"name,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	Active      bool      `json:"active,omitempty" xml:"active,omitempty"`
}

type SiteContentRevisions []SiteContentRevision

type LandingPage struct {
	Id                int               `json:"name,omitempty" xml:"name,omitempty"`
	Name              string            `json:"name,omitempty" xml:"name,omitempty"`
	StartDate         time.Time         `json:"name,omitempty" xml:"name,omitempty"`
	EndDate           time.Time         `json:"name,omitempty" xml:"name,omitempty"`
	Url               url.URL           `json:"name,omitempty" xml:"name,omitempty"`
	PageContent       string            `json:"name,omitempty" xml:"name,omitempty"`
	LinkClasses       string            `json:"name,omitempty" xml:"name,omitempty"`
	ConversionId      string            `json:"name,omitempty" xml:"name,omitempty"`
	ConversionLabel   string            `json:"name,omitempty" xml:"name,omitempty"`
	NewWindow         bool              `json:"name,omitempty" xml:"name,omitempty"`
	MenuPosition      string            `json:"name,omitempty" xml:"name,omitempty"`
	LandingPageDatas  LandingPageDatas  `json:"name,omitempty" xml:"name,omitempty"`
	LandingPageImages LandingPageImages `json:"name,omitempty" xml:"name,omitempty"`
}

type LandingPages []LandingPage

type LandingPageData struct {
	Id        int
	DataKey   string
	DataValue string
}
type LandingPageDatas []LandingPageData

type LandingPageImage struct {
	Id   int
	Url  url.URL
	Sort int
}

type LandingPageImages []LandingPageImage

const (
	siteContentColumns         = "contentID, content_type, page_title, createdDate, lastModified, meta_title, meta_description, keywords, isPrimary, published, active, slug, requireAuthentication, canonical"
	siteContentRevisionColumns = "revisionID, contentID, content_text, createdOn, active"
)

var (
	getContentPageByName              = `SELECT ` + siteContentColumns + ` FROM SiteContent WHERE slug = ? LIMIT 1`
	getSiteContentRevisionByContentID = `SELECT ` + siteContentRevisionColumns + ` FROM SiteContentRevision WHERE contentID = ? ORDER BY createdOn DESC LIMIT 1`
	getContentPageByID                = `select ` + siteContentColumns + ` from SiteContent where contentID = ? limit 1`
	getAllSiteContent                 = `SELECT ` + siteContentColumns + ` FROM SiteContent `
	getAllSiteContentRevisions        = `SELECT ` + siteContentRevisionColumns + ` FROM SiteContentRevision `
)

//get content page
func (cp *ContentPage) Get() (err error) {

	contents, err := GetAllSiteContents()
	revisions, err := GetAllSiteContentRevisions()

	// menus, err := GetAllMenus()

	for _, c := range contents {
		if c.Id == cp.SiteContent.Id {
			cp.SiteContent = c
		}
	}

	for _, r := range revisions {
		if r.ContentId == cp.SiteContent.Id {
			cp.Revision = r
		}
	}

	err = cp.Menu.GetMenuByContentId(cp.SiteContent.Id)
	if err != sql.ErrNoRows {
		if err != nil {
			log.Print("ERRE", err, " ", cp.SiteContent.Id)
			return err
		}
	}

	return err

}

//get all menus
func GetAllSiteContents() (sc SiteContents, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return sc, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllSiteContent)
	if err != nil {
		return sc, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	var s SiteContent
	var cType, title, mTitle, mDesc, slug, canon *string
	for res.Next() {
		err = res.Scan(
			&s.Id,
			&cType,
			&title,
			&s.CreatedDate,
			&s.LastModified,
			&mTitle,
			&mDesc,
			&s.Keywords,
			&s.IsPrimary,
			&s.Published,
			&s.Active,
			&slug,
			&s.RequireAuthentication,
			&canon,
		)
		if err != nil {
			return sc, err
		}
		if cType != nil {
			s.Type = *cType
		}
		if title != nil {
			s.Title = *title
		}
		if mTitle != nil {
			s.MetaTitle = *mTitle
		}
		if mDesc != nil {
			s.MetaDescription = *mDesc
		}
		if slug != nil {
			s.Slug = *slug
		}
		if canon != nil {
			s.Canonical = *canon
		}
		sc = append(sc, s)
	}
	return sc, err
}

//get all menus
func GetAllSiteContentRevisions() (srcs SiteContentRevisions, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return srcs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllSiteContentRevisions)
	if err != nil {
		return srcs, err
	}
	defer stmt.Close()

	res, err := stmt.Query()

	var c SiteContentRevision
	for res.Next() {
		err = res.Scan(
			&c.Id,
			&c.ContentId,
			&c.Text,
			&c.CreatedDate,
			&c.Active,
		)
		if err != nil {
			return srcs, err
		}
		srcs = append(srcs, c)
	}
	return srcs, err
}

func (s *SiteContent) Get() (err error) { //by id
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getContentPageByID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string

	err = stmt.QueryRow(s.Id).Scan(
		&s.Id,
		&cType,
		&title,
		&s.CreatedDate,
		&s.LastModified,
		&mTitle,
		&mDesc,
		&s.Keywords,
		&s.IsPrimary,
		&s.Published,
		&s.Active,
		&slug,
		&s.RequireAuthentication,
		&canon,
	)
	if err != nil {
		return err
	}
	if cType != nil {
		s.Type = *cType
	}
	if title != nil {
		s.Title = *title
	}
	if mTitle != nil {
		s.MetaTitle = *mTitle
	}
	if mDesc != nil {
		s.MetaDescription = *mDesc
	}
	if slug != nil {
		s.Slug = *slug
	}
	if canon != nil {
		s.Canonical = *canon
	}
	return err
}

func (s *SiteContent) GetBySlug() (err error) { //by slug
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getContentPageByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string

	err = stmt.QueryRow(s.Slug).Scan(
		&s.Id,
		&cType,
		&title,
		&s.CreatedDate,
		&s.LastModified,
		&mTitle,
		&mDesc,
		&s.Keywords,
		&s.IsPrimary,
		&s.Published,
		&s.Active,
		&slug,
		&s.RequireAuthentication,
		&canon,
	)
	if err != nil {
		return err
	}
	if cType != nil {
		s.Type = *cType
	}
	if title != nil {
		s.Title = *title
	}
	if mTitle != nil {
		s.MetaTitle = *mTitle
	}
	if mDesc != nil {
		s.MetaDescription = *mDesc
	}
	if slug != nil {
		s.Slug = *slug
	}
	if canon != nil {
		s.Canonical = *canon
	}
	return err
}

func (c *ContentPage) GetRevision() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSiteContentRevisionByContentID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.SiteContent.Id).Scan(
		&c.Revision.Id,
		&c.SiteContent.Id,
		&c.Revision.Text,
		&c.Revision.CreatedDate,
		&c.Revision.Active,
	)
	if err != nil {
		return err
	}

	return nil
}

//map revisions
func (ms SiteContentRevisions) ToMap() map[interface{}]SiteContentRevision {
	zeeMap := make(map[interface{}]SiteContentRevision)
	for _, v := range ms {
		zeeMap[v.Id] = v
	}
	return zeeMap
}

//map site content
func (ms SiteContents) ToMap() map[interface{}]SiteContent {
	zeeMap := make(map[interface{}]SiteContent)
	for _, v := range ms {
		zeeMap[v.Id] = v
	}
	return zeeMap
}
