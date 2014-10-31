package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type MenuWithContent struct {
	Menu     `json:"menu,omitempty" xml:"menu,omitempty"`
	Contents []MenuItem `json:"contents,omitempty" xml:"contents,omitempty"`
}
type MenuWithContents []MenuWithContent

type Menu struct {
	Id                    int    `json:"id,omitempty" xml:"id,omitempty"`
	Name                  string `json:"name,omitempty" xml:"name,omitempty"`
	IsPrimary             bool   `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Active                bool   `json:"active,omitempty" xml:"active,omitempty"`
	DisplayName           string `json:"displayName,omitempty" xml:"displayName,omitempty"`
	RequireAuthentication bool   `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	ShowOnSiteMap         bool   `json:"showOnSiteMap,omitempty" xml:"showOnSiteMap,omitempty"`
	Sort                  int    `json:"sort,omitempty" xml:"sort,omitempty"`
	WebsiteId             int    `json:"websiteId,omitempty" xml:"websiteId,omitempty"`
}

type Menus []Menu

type MenuItem struct {
	Menu_SiteContent `json:"menuSiteContent,omitempty" xml:"menuSiteContent,omitempty"`
	Content          ContentPage `json:"content,omitempty" xml:"content,omitempty"`
}
type MenuItems []MenuItem

type Menu_SiteContent struct {
	Id         int    `json:"id,omitempty" xml:"id,omitempty"`
	MenuId     int    `json:"menuId,omitempty" xml:"menuId,omitempty"`
	ContentId  int    `json:"contentId,omitempty" xml:"contentId,omitempty"`
	Sort       int    `json:"sort,omitempty" xml:"sort,omitempty"`
	Title      string `json:"title,omitempty" xml:"title,omitempty"`
	Link       string `json:"link,omitempty" xml:"link,omitempty"`
	ParentId   int    `json:"parentId,omitempty" xml:"parentId,omitempty"`
	LinkTarget bool   `json:"linkTarget,omitempty" xml:"linkTarget,omitempty"`
}
type Menu_SiteContents []Menu_SiteContent

const (
	menuColumns            = "m.menuID, m.menu_name, m.isPrimary, m.active, m.display_name, m.requireAuthentication, m.showOnSiteMap, m.sort, m.websiteID" //menu AS m
	menuSiteContentColumns = "msc.menuContentID, msc.menuID, msc.contentID, msc.menuSort, msc.menuTitle, msc.menuLink, msc.parentID, msc.linkTarget"       //as msc
)

var (
	getMenuByContentID = `SELECT ` + menuColumns + ` FROM Menu AS m
							JOIN Menu_SiteContent AS msc ON m.menuID = msc.menuID
							  WHERE msc.contentID = ? && isPrimary = 0 && m.requireAuthentication = ? && m.websiteID = 1
							  ORDER BY msc.menuID ASC
							  LIMIT 1` //TODO - V2 changed this to menuID

	getPrimaryMenu           = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.isPrimary = 1 && m.websiteID = 1`
	getMenu                  = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.menuID = ?`
	getAllMenus              = `SELECT ` + menuColumns + ` FROM Menu AS m `
	getMenuByName            = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.menu_name = ?`
	getMenuByMenuID          = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.menuID = ? && m.websiteID = 1`
	getMenuItemsByMenuID     = `SELECT ` + menuSiteContentColumns + ` FROM Menu_SiteContent AS msc WHERE msc.menuID = ? ORDER BY menuSort`
	getFooterSitemap         = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.showOnSitemap = 1 && m.websiteID = 1 ORDER BY m.sort`
	getMenuWithContentByName = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.menu_name = ? &&m.websiteID  = 1 LIMIT 1`
	getAllMenuContents       = `SELECT ` + menuSiteContentColumns + ` FROM Menu_SiteContent AS msc`
	getMenuSitemap           = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.websiteID = 1 ORDER BY m.isPrimary DESC`
	getMenuIDByContentID     = `select menuID from Menu_SiteContent where contentID = ?`
)

//get primary menu...with Menu Contents ...primary Menu is just GetMenu, with primary
func (m *MenuWithContent) GetPrimaryMenu() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPrimaryMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&m.DisplayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
		&m.WebsiteId,
	)
	if err != nil {
		return err
	}
	m.Contents, err = m.Menu.GetMenuItemsByMenuId()
	return err
}

func (m *Menu) GetMenuItemsByMenuId() (ms MenuItems, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ms, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenuItemsByMenuID)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()
	var mi MenuItem
	var title, link *string
	var contentId, parent *int
	res, err := stmt.Query(m.Id)
	for res.Next() {
		err = res.Scan(
			&mi.Menu_SiteContent.Id,
			&mi.Menu_SiteContent.MenuId,
			&contentId,
			&mi.Menu_SiteContent.Sort,
			&title,
			&link,
			&parent,
			&mi.Menu_SiteContent.LinkTarget,
		)
		if err != nil {
			return ms, err
		}
		if contentId != nil {
			mi.Menu_SiteContent.ContentId = *contentId
			mi.Content.SiteContent.Id = *contentId
		}
		if title != nil {
			mi.Menu_SiteContent.Title = *title
		}
		if link != nil {
			mi.Menu_SiteContent.Link = *link
		}
		if parent != nil {
			mi.Menu_SiteContent.ParentId = *parent
		}
		if err != nil {
			return ms, err
		}
		//getContent by pageid
		err = mi.Content.Get()
		if err != sql.ErrNoRows {
			if err != nil {
				return ms, err
			}
		}
		ms = append(ms, mi)
	}

	return ms, err
}

func GetFooterSitemap() (ms MenuWithContents, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ms, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getFooterSitemap)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	var displayName *string
	var m MenuWithContent
	for res.Next() {
		err = res.Scan(
			&m.Menu.Id,
			&m.Menu.Name,
			&m.Menu.IsPrimary,
			&m.Menu.Active,
			&displayName,
			&m.Menu.RequireAuthentication,
			&m.Menu.ShowOnSiteMap,
			&m.Menu.Sort,
			&m.WebsiteId,
		)
		if err != nil {
			return ms, err
		}
		if displayName != nil {
			m.Menu.DisplayName = *displayName
		}

		m.Contents, err = m.Menu.GetMenuItemsByMenuId()
		if err != nil {
			return ms, err
		}
		ms = append(ms, m)
	}
	if len(ms) == 0 {
		err = sql.ErrNoRows
	}
	return ms, err

}

func GetMenuSitemap() (ms MenuWithContents, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ms, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getMenuSitemap)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()
	var displayName *string
	var m MenuWithContent
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&m.Menu.Id,
			&m.Menu.Name,
			&m.Menu.IsPrimary,
			&m.Menu.Active,
			&displayName,
			&m.Menu.RequireAuthentication,
			&m.Menu.ShowOnSiteMap,
			&m.Menu.Sort,
			&m.Menu.WebsiteId,
		)
		if err != nil {
			return ms, err
		}
		if displayName != nil {
			m.Menu.DisplayName = *displayName
		}

		m.Contents, err = m.Menu.GetMenuItemsByMenuId()
		if err != nil {
			return ms, err
		}
		//append only if has contents
		if len(m.Contents) > 0 {
			ms = append(ms, m)
		}

	}
	if len(ms) == 0 {
		err = sql.ErrNoRows
	}
	return ms, err
}

//Generally used for getting content bridge From menu_sitecontent; no actual menu
func (m *MenuWithContent) GetMenuByContentId(id int, auth bool) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getMenuByContentID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var displayName *string
	err = stmt.QueryRow(id, auth).Scan(
		&m.Menu.Id,
		&m.Menu.Name,
		&m.Menu.IsPrimary,
		&m.Menu.Active,
		&displayName,
		&m.Menu.RequireAuthentication,
		&m.Menu.ShowOnSiteMap,
		&m.Menu.Sort,
		&m.Menu.WebsiteId,
	)
	if err != nil {

		return err
	}
	if displayName != nil {
		m.DisplayName = *displayName
	}

	m.Contents, err = m.GetMenuItemsByMenuId()
	if err != nil {

		return err
	}
	return err
}

func GetMenuIdByContentId(ContentId int) (menuId int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return menuId, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getMenuIDByContentID)
	if err != nil {
		return menuId, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ContentId).Scan(&menuId)
	return menuId, err
}

func (m *MenuWithContent) GetMenuWithContentByName(name string) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getMenuWithContentByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var displayName *string
	err = stmt.QueryRow(name).Scan(
		&m.Menu.Id,
		&m.Menu.Name,
		&m.Menu.IsPrimary,
		&m.Menu.Active,
		&displayName,
		&m.Menu.RequireAuthentication,
		&m.Menu.ShowOnSiteMap,
		&m.Menu.Sort,
		&m.Menu.WebsiteId,
	)
	if err != nil {
		return err
	}
	if displayName != nil {
		m.Menu.DisplayName = *displayName
	}

	m.Contents, err = m.GetMenuItemsByMenuId()

	return err

}

//NEW, non-V2, methods
func GetAllMenus() (ms Menus, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ms, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllMenus)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()

	var displayName *string
	var m Menu
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&m.Id,
			&m.Name,
			&m.IsPrimary,
			&m.Active,
			&displayName,
			&m.RequireAuthentication,
			&m.ShowOnSiteMap,
			&m.Sort,
			&m.WebsiteId,
		)
		if err != nil {
			return ms, err
		}

		if displayName != nil {
			m.DisplayName = *displayName
		}
		ms = append(ms, m)
	}
	return ms, err
}
