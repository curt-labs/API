package products

import (
	"database/sql"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	GetInventoryForPart = `select
		i.partID, i.available, i.dateUpdated,
		wh.name, wh.code, wh.address, wh.city, wh.postalCode,
		wh.tollFreePhone, wh.fax, wh.localPhone, wh.manager,
		wh.longitude, wh.latitude,
		s.abbr as stateAbbr, s.state,
		c.name as countryName, c.abbr as countryAbbr
		from Inventory as i
		join Warehouses as wh on i.warehouseID = wh.id
		left join States as s on wh.stateID = s.stateID
		left join Country as c on s.countryID = c.countryID
		where i.partID = ?
		order by wh.code`

	GetInventoryForPartByWarehouse = `select
		i.partID, i.available, i.dateUpdated,
		wh.name, wh.code, wh.address, wh.city, wh.postalCode,
		wh.tollFreePhone, wh.fax, wh.localPhone, wh.manager,
		wh.longitude, wh.latitude,
		s.abbr as stateAbbr, s.state,
		c.name as countryName, c.abbr as countryAbbr
		from Inventory as i
		join Warehouses as wh on i.warehouseID = wh.id
		left join States as s on wh.stateID = s.stateID
		left join Country as c on s.countryID = c.countryID
		where i.partID = ? && wh.code = ?`
)

type Warehouse struct {
	ID            int    `json:"-" xml:"-"`
	Name          string `json:"name" xml:"name,attr"`
	Code          string `json:"code" xml:"code,attr"`
	Address       string `json:"address" xml:"address,attr"`
	City          string `json:"city" xml:"city,attr"`
	State         State  `json:"state" xml:"state"`
	PostalCode    string `json:"postal_code" xml:"postal_code,attr"`
	TollFreePhone string `json:"toll_free_phone" xml:"toll_free_phone,attr"`
	Latitude      string `json:"latitude" xml:"latitude,attr"`
	Longitude     string `json:"longitude" xml:"longitude,attr"`
	Fax           string `json:"fax" xml:"fax,attr"`
	LocalPhone    string `json:"local_phone" xml:"local_phone,attr"`
	Manager       string `json:"manager" xml:"manager"`
}

type PartInventory struct {
	TotalAvailability int         `json:"total_availability" xml:"total_availability,attr"`
	Warehouses        []Inventory `json:"inventory" xml:"inventory"`
}

type Inventory struct {
	Part        int       `json:"part" xml:"part,attr"`
	Warehouse   Warehouse `json:"warehouse" xml:"warehouse"`
	Quantity    int       `json:"quantity" xml:"quantity,attr"`
	DateUpdated time.Time `json:"date_updated" xml:"date_update,attr"`
}

type State struct {
	ID           int     `json:"-" xml:"-"`
	State        string  `json:"state" xml:"state,abbr"`
	Abbreviation string  `json:"abbreviation" xml:"abbreviation,attr"`
	Country      Country `json:"country" xml:"country"`
}

type Country struct {
	ID           int    `json:"-" xml:"-"`
	Name         string `json:"name" xml:"name,attr"`
	Abbreviation string `json:"abbreviation" xml:"abbreviation,attr"`
}

func (p *Part) GetInventory(apiKey, warehouseCode string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	var rows *sql.Rows
	if warehouseCode == "" {
		stmt, err := db.Prepare(GetInventoryForPart)
		if err != nil {
			return err
		}
		defer stmt.Close()

		rows, err = stmt.Query(p.ID)
		if err != nil {
			return err
		}
	} else {
		stmt, err := db.Prepare(GetInventoryForPartByWarehouse)
		if err != nil {
			return err
		}
		defer stmt.Close()

		rows, err = stmt.Query(p.ID, warehouseCode)
		if err != nil {
			return err
		}
	}

	if rows == nil {
		return fmt.Errorf("error: %s", "inventory not available for this part")
	}
	defer rows.Close()

	for rows.Next() {
		var i Inventory
		// var date string
		if err = rows.Scan(
			&i.Part,
			&i.Quantity,
			&i.DateUpdated,
			&i.Warehouse.Name,
			&i.Warehouse.Code,
			&i.Warehouse.Address,
			&i.Warehouse.City,
			&i.Warehouse.PostalCode,
			&i.Warehouse.TollFreePhone,
			&i.Warehouse.Fax,
			&i.Warehouse.LocalPhone,
			&i.Warehouse.Manager,
			&i.Warehouse.Longitude,
			&i.Warehouse.Latitude,
			&i.Warehouse.State.Abbreviation,
			&i.Warehouse.State.State,
			&i.Warehouse.State.Country.Name,
			&i.Warehouse.State.Country.Abbreviation,
		); err != nil {
			return err
		}

		if i.Part != 0 {
			p.Inventory.Warehouses = append(p.Inventory.Warehouses, i)
		}
	}

	for _, w := range p.Inventory.Warehouses {
		p.Inventory.TotalAvailability = p.Inventory.TotalAvailability + w.Quantity
	}

	return nil
}
