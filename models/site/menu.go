package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	// "strconv"
	"log"
)

type Menu struct { //menu  + menuCOntent(menuItem + menu_sitecontent)
	Id                    int          `json:"id,omitempty" xml:"id,omitempty"`
	Name                  string       `json:"name,omitempty" xml:"name,omitempty"`
	IsPrimary             bool         `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Active                bool         `json:"active,omitempty" xml:"active,omitempty"`
	DisplayName           string       `json:"displayName,omitempty" xml:"displayName,omitempty"`
	RequireAuthentication bool         `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	ShowOnSiteMap         bool         `json:"showOnSiteMap,omitempty" xml:"showOnSiteMap,omitempty"`
	Sort                  int          `json:"sort,omitempty" xml:"sort,omitempty"`
	MenuContents          MenuContents `json:"menuContents,omitempty" xml:"menuContents,omitempty"`
}

type Menus []Menu

type MenuContent struct { //menuitem + menu_sitecontent
	Id          int         `json:"id,omitempty" xml:"id,omitempty"`
	MenuId      int         `json:"menuId,omitempty" xml:"menuId,omitempty"`
	ContentPage ContentPage `json:"contentPage,omitempty" xml:"contentPage,omitempty"`
	Sort        int         `json:"sort,omitempty" xml:"sort,omitempty"`
	Title       string      `json:"title,omitempty" xml:"title,omitempty"`
	Link        string      `json:"link,omitempty" xml:"link,omitempty"`
	ParentId    int         `json:"parentId,omitempty" xml:"parentId,omitempty"`
	LinkTarget  bool        `json:"linkTarget,omitempty" xml:"linkTarget,omitempty"`
}
type MenuContents []MenuContent

const (
	menuColumns            = "m.menuID, m.menu_name, m.isPrimary, m.active, m.display_name, m.requireAuthentication, m.showOnSiteMap, m.sort"        //menu AS m
	menuSiteContentColumns = "msc.menuContentID, msc.menuID, msc.contentID, msc.menuSort, msc.menuTitle, msc.menuLink, msc.parentID, msc.linkTarget" //as msc
)

var (
	getMenuByContentID = `SELECT ` + menuColumns + ` FROM Menu AS m
							JOIN Menu_SiteContent AS msc ON m.menuID = msc.menuID
							  WHERE msc.contentID = ? && isPrimary = 0
							  ORDER BY msc.menuID ASC
							  LIMIT 1`

	getPrimaryMenu           = `SELECT ` + menuColumns + ` FROM Menu as m WHERE m.isPrimary = 1`
	getMenu                  = `SELECT ` + menuColumns + ` FROM Menu as m WHERE m.menuID = ?`
	getAllMenus              = `SELECT ` + menuColumns + ` FROM Menu as m `
	getMenuByName            = `SELECT ` + menuColumns + ` FROM Menu as m WHERE m.menu_name = ?`
	getMenuByMenuID          = `SELECT ` + menuColumns + ` FROM Menu as m WHERE m.menuID = ?`
	getMenuItemsByMenuID     = `SELECT ` + menuSiteContentColumns + ` FROM Menu_SiteContent as msc WHERE menuID = ? ORDER BY menuSort`
	getFooterSitemap         = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.showOnSitemap = 1 ORDER BY m.sort`
	getMenuWithContentByName = `SELECT ` + menuColumns + ` FROM Menu AS m WHERE m.menu_name = ? LIMIT 1`
	getAllMenuContents       = `SELECT ` + menuSiteContentColumns + ` FROM Menu_SiteContent AS msc`
)

//get primary menu...with Menu Contents ...primary Menu is just GetMenu, with primary
func (m *Menu) GetPrimaryMenu() (err error) {
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

	mc, err := GetMenuContents()
	if err != nil {
		return err
	}
	contentMap := mc.ToMap()
	err = stmt.QueryRow().Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&m.DisplayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
	)
	if err != nil {
		return err
	}

	if cmc, found := contentMap[m.Id]; found {
		m.MenuContents = append(m.MenuContents, cmc)
	}

	return err
}

//get all menus
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

	mc, err := GetMenuContents()
	if err != nil {
		return ms, err
	}
	contentMap := mc.ToMap()
	res, err := stmt.Query()
	//m.menuID, m.menu_name, m.isPrimary, m.active, m.display_name, m.requireAuthentication, m.showOnSiteMap, m.sort
	var m Menu
	for res.Next() {
		err = res.Scan(
			&m.Id,
			&m.Name,
			&m.IsPrimary,
			&m.Active,
			&m.DisplayName,
			&m.RequireAuthentication,
			&m.ShowOnSiteMap,
			&m.Sort,
		)
		if err != nil {
			log.Print("ERR HERE", err)
			return ms, err
		}

		for _, c := range contentMap {
			if c.MenuId == m.Id {
				m.MenuContents = append(m.MenuContents, c)
			}
		}

	}
	ms = append(ms, m)
	return ms, err
}

//get Menu with menuContents...by id
func (m *Menu) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	mc, err := GetMenuContents()
	if err != nil {
		return err
	}
	contentMap := mc.ToMap()
	err = stmt.QueryRow(m.Id).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&m.DisplayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
	)
	if err != nil {
		return err
	}
	if cmc, found := contentMap[m.Id]; found {
		m.MenuContents = append(m.MenuContents, cmc)
	}
	return err
}

//get Menu with menuContents...by name; data appears jacked - multiple menus with same name, we grab first
//I hate this method, but similar functionality exists in v2...do we need it?
func (m *Menu) GetByName() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenuByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	mc, err := GetMenuContents()
	if err != nil {
		return err
	}
	contentMap := mc.ToMap()
	err = stmt.QueryRow(m.Name).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&m.DisplayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
	)
	if err != nil {
		return err
	}
	if cmc, found := contentMap[m.Id]; found {
		m.MenuContents = append(m.MenuContents, cmc)
	}
	return err
}

func (m *Menu) GetMenuContents() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenuItemsByMenuID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(m.Id)
	var mm MenuContent
	var title, link *string
	var contentId, parent *int
	for res.Next() {
		err = res.Scan(
			&mm.Id,
			&mm.MenuId,
			&contentId,
			&mm.Sort,
			&title,
			&link,
			&parent,
			&mm.LinkTarget,
		)
		if err != nil {
			return err
		}
		if contentId != nil {
			mm.ContentPage.SiteContent.Id = *contentId
		}
		if title != nil {
			mm.Title = *title
		}
		if link != nil {
			mm.Link = *link
		}
		if parent != nil {
			mm.ParentId = *parent
		}
		if err != nil {
			return err
		}
		err = mm.ContentPage.Get() //get content page
		if err != nil {
			return err
		}
		m.MenuContents = append(m.MenuContents, mm)
	}
	return err
}

//Generally used for getting content bridge From menu_sitecontent; no actual menu
func (m *Menu) GetMenuByContentId(id int) (err error) {
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
	err = stmt.QueryRow(id).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&displayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
	)
	if err != nil {
		return err
	}
	if displayName != nil {
		m.DisplayName = *displayName
	}
	// err = m.GetMenuContents()
	if err != nil {
		return err
	}
	return err

}

func (m *Menu) GetMenuWithContentByName(slug string) (err error) {
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
	err = stmt.QueryRow(slug).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&displayName,
		&m.RequireAuthentication,
		&m.ShowOnSiteMap,
		&m.Sort,
	)
	if err != nil {
		return err
	}
	if displayName != nil {
		m.DisplayName = *displayName
	}
	err = m.GetMenuContents()
	return err

}

func GetFooterSitemap() (ms Menus, err error) {
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
	var m Menu
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
		)
		if err != nil {
			return ms, err
		}
		if displayName != nil {
			m.DisplayName = *displayName
		}
		err = m.GetMenuContents()
		if err != nil {
			return ms, err
		}
		ms = append(ms, m)
	}
	return ms, err

}

//new
func GetMenuContents() (mc MenuContents, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return mc, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllMenuContents)
	if err != nil {
		return mc, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	var mm MenuContent
	var title, link *string
	var contentId, parent *int
	for res.Next() {
		err = res.Scan(
			&mm.Id,
			&mm.MenuId,
			&contentId,
			&mm.Sort,
			&title,
			&link,
			&parent,
			&mm.LinkTarget,
		)
		if err != nil {
			return mc, err
		}
		if contentId != nil {
			mm.ContentPage.SiteContent.Id = *contentId
		}
		if title != nil {
			mm.Title = *title
		}
		if link != nil {
			mm.Link = *link
		}
		if parent != nil {
			mm.ParentId = *parent
		}
		err = mm.ContentPage.Get()
		mc = append(mc, mm)
	}
	return mc, err
}

//map contents
func (cs MenuContents) ToMap() map[interface{}]MenuContent {
	zeeMap := make(map[interface{}]MenuContent)
	for _, v := range cs {
		zeeMap[v.Id] = v
	}
	return zeeMap
}

//map menus
func (ms Menus) ToMap() map[interface{}]Menu {
	zeeMap := make(map[interface{}]Menu)
	for _, v := range ms {
		zeeMap[v.Id] = v
	}
	return zeeMap
}
