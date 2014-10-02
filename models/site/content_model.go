package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)

type ContentPage struct {
	SiteContent     SiteContent
	MenuWithContent MenuWithContent
	Revision        SiteContentRevision
}

type ContentPages []ContentPage

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
	Id                int               `json:"id,omitempty" xml:"id,omitempty"`
	Name              string            `json:"name,omitempty" xml:"name,omitempty"`
	StartDate         time.Time         `json:"startDate,omitempty" xml:"v,omitempty"`
	EndDate           time.Time         `json:"endDate,omitempty" xml:"endDate,omitempty"`
	Url               url.URL           `json:"url,omitempty" xml:"url,omitempty"`
	PageContent       string            `json:"pageContent,omitempty" xml:"pageContent,omitempty"`
	LinkClasses       string            `json:"linkClasses,omitempty" xml:"linkClasses,omitempty"`
	ConversionId      string            `json:"conversionId,omitempty" xml:"conversionId,omitempty"`
	ConversionLabel   string            `json:"conversionLabel,omitempty" xml:"conversionLabel,omitempty"`
	NewWindow         bool              `json:"newWindow,omitempty" xml:"newWindow,omitempty"`
	MenuPosition      string            `json:"menuPosition,omitempty" xml:"menuPosition,omitempty"`
	LandingPageDatas  LandingPageDatas  `json:"landingPageDatas,omitempty" xml:"landingPageDatas,omitempty"`
	LandingPageImages LandingPageImages `json:"landingPageImages,omitempty" xml:"landingPageImages,omitempty"`
}

type LandingPages []LandingPage

type LandingPageData struct {
	Id        int    `json:"id,omitempty" xml:"id,omitempty"`
	DataKey   string `json:"dataKey,omitempty" xml:"dataKey,omitempty"`
	DataValue string `json:"dataValue,omitempty" xml:"dataValue,omitempty"`
}
type LandingPageDatas []LandingPageData

type LandingPageImage struct {
	Id   int     `json:"id,omitempty" xml:"id,omitempty"`
	Url  url.URL `json:"url,omitempty" xml:"url,omitempty"`
	Sort int     `json:"sort,omitempty" xml:"sort,omitempty"`
}

type LandingPageImages []LandingPageImage

const (
	siteContentColumns         = "s.contentID, s.content_type, s.page_title, s.createdDate, s.lastModified, s.meta_title, s.meta_description, s.keywords, s.isPrimary, s.published, s.active, s.slug, s.requireAuthentication, s.canonical" //as s
	siteContentRevisionColumns = "revisionID, contentID, content_text, createdOn, active"
	landingPageColumns         = "l.id, l.name, l.startDate, l.endDate, l.url, l.pageContent, l.linkClasses, l.conversionID, l.conversionLabel, l.newWindow, l.menuPosition" //as l
)

var (
	getContentPageByName              = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s WHERE slug = ? LIMIT 1`
	getSiteContentRevisionByContentID = `SELECT ` + siteContentRevisionColumns + ` FROM SiteContentRevision WHERE contentID = ? ORDER BY createdOn DESC LIMIT 1`
	getContentPageByID                = `select ` + siteContentColumns + ` from SiteContent AS s WHERE contentID = ? limit 1`
	getAllSiteContent                 = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s`
	getAllSiteContentRevisions        = `SELECT ` + siteContentRevisionColumns + ` FROM SiteContentRevision `
	getPrimaryContentPage             = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s WHERE isPrimary = 1`
	getSitemapCP                      = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s JOIN Menu_SiteContent AS msc ON msc.contentID != s.contentID`
	getLandingPageByID                = `SELECT ` + landingPageColumns + ` FROM LandingPage AS l WHERE id = ? && startDate <= NOW() && endDate >= NOW() limit 1`
)

//get content page
func (cp *ContentPage) Get() (err error) {
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
	err = stmt.QueryRow(cp.SiteContent.Id).Scan(
		&cp.SiteContent.Id,
		&cType,
		&title,
		&cp.SiteContent.CreatedDate,
		&cp.SiteContent.LastModified,
		&mTitle,
		&mDesc,
		&cp.SiteContent.Keywords,
		&cp.SiteContent.IsPrimary,
		&cp.SiteContent.Published,
		&cp.SiteContent.Active,
		&slug,
		&cp.SiteContent.RequireAuthentication,
		&canon,
	)
	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}
	}
	if cType != nil {
		cp.SiteContent.Type = *cType
	}
	if title != nil {
		cp.SiteContent.Title = *title
	}
	if mTitle != nil {
		cp.SiteContent.MetaTitle = *mTitle
	}
	if mDesc != nil {
		cp.SiteContent.MetaDescription = *mDesc
	}
	if slug != nil {
		cp.SiteContent.Slug = *slug
	}
	if canon != nil {
		cp.SiteContent.Canonical = *canon
	}
	//get Revision
	err = cp.GetRevision()

	return err

}

//by name (slug)
func (cp *ContentPage) GetContentPageByName(menuId int, auth bool) (err error) {
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

	err = stmt.QueryRow(cp.SiteContent.Slug).Scan(
		&cp.SiteContent.Id,
		&cType,
		&title,
		&cp.SiteContent.CreatedDate,
		&cp.SiteContent.LastModified,
		&mTitle,
		&mDesc,
		&cp.SiteContent.Keywords,
		&cp.SiteContent.IsPrimary,
		&cp.SiteContent.Published,
		&cp.SiteContent.Active,
		&slug,
		&cp.SiteContent.RequireAuthentication,
		&canon,
	)

	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}
	}
	if cType != nil {
		cp.SiteContent.Type = *cType
	}
	if title != nil {
		cp.SiteContent.Title = *title
	}
	if mTitle != nil {
		cp.SiteContent.MetaTitle = *mTitle
	}
	if mDesc != nil {
		cp.SiteContent.MetaDescription = *mDesc
	}
	if slug != nil {
		cp.SiteContent.Slug = *slug
	}
	if canon != nil {
		cp.SiteContent.Canonical = *canon
	}

	//get Revision
	err = cp.GetRevision()

	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}

	}
	//get Menu
	cp.MenuWithContent.Menu.Id = menuId

	err = cp.MenuWithContent.GetMenuByContentId(menuId, auth)

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
	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *ContentPage) GetPrimaryContentPage() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPrimaryContentPage)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string
	err = stmt.QueryRow().Scan(
		&cp.SiteContent.Id,
		&cType,
		&title,
		&cp.SiteContent.CreatedDate,
		&cp.SiteContent.LastModified,
		&mTitle,
		&mDesc,
		&cp.SiteContent.Keywords,
		&cp.SiteContent.IsPrimary,
		&cp.SiteContent.Published,
		&cp.SiteContent.Active,
		&slug,
		&cp.SiteContent.RequireAuthentication,
		&canon,
	)
	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}
	}
	if cType != nil {
		cp.SiteContent.Type = *cType
	}
	if title != nil {
		cp.SiteContent.Title = *title
	}
	if mTitle != nil {
		cp.SiteContent.MetaTitle = *mTitle
	}
	if mDesc != nil {
		cp.SiteContent.MetaDescription = *mDesc
	}
	if slug != nil {
		cp.SiteContent.Slug = *slug
	}
	if canon != nil {
		cp.SiteContent.Canonical = *canon
	}
	//get Revision
	err = cp.GetRevision()
	return err
}

func GetSitemapCP() (cps ContentPages, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cps, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSitemapCP)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string
	var cp ContentPage
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&cp.SiteContent.Id,
			&cType,
			&title,
			&cp.SiteContent.CreatedDate,
			&cp.SiteContent.LastModified,
			&mTitle,
			&mDesc,
			&cp.SiteContent.Keywords,
			&cp.SiteContent.IsPrimary,
			&cp.SiteContent.Published,
			&cp.SiteContent.Active,
			&slug,
			&cp.SiteContent.RequireAuthentication,
			&canon,
		)
		if err != sql.ErrNoRows {
			if err != nil {
				return cps, err
			}
		}
		if cType != nil {
			cp.SiteContent.Type = *cType
		}
		if title != nil {
			cp.SiteContent.Title = *title
		}
		if mTitle != nil {
			cp.SiteContent.MetaTitle = *mTitle
		}
		if mDesc != nil {
			cp.SiteContent.MetaDescription = *mDesc
		}
		if slug != nil {
			cp.SiteContent.Slug = *slug
		}
		if canon != nil {
			cp.SiteContent.Canonical = *canon
		}
		//get Revision
		// err = cp.GetRevision()
		// if err != sql.ErrNoRows {
		// 	if err != nil {
		// 		return cps, err
		// 	}
		// }
		cps = append(cps, cp)
	}

	return cps, err
}

func (l *LandingPage) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getLandingPageByID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var link, convId, convLab *string
	var u []byte
	err = stmt.QueryRow(l.Id).Scan(
		&l.Id,
		&l.Name,
		&l.StartDate,
		&l.EndDate,
		&u,
		&l.PageContent,
		&link,
		&convId,
		&convLab,
		&l.NewWindow,
		&l.MenuPosition,
	)

	if err != nil {
		return err
	}
	if link != nil {
		l.LinkClasses = *link
	}
	if convId != nil {
		l.ConversionId = *convId
	}
	if convLab != nil {
		l.ConversionLabel = *convLab
	}
	if u != nil {
		tempUrl, err := url.Parse(string(u))
		if err != nil {
			return err
		}
		l.Url = *tempUrl
	}
	return nil
}
