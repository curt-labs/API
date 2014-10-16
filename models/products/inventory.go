package products

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/redis"
	redix "github.com/garyburd/redigo/redis"
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
	ID            int     `json:"-" xml:"-"`
	Name          string  `json:"name" xml:"name,attr"`
	Code          string  `json:"code" xml:"code,attr"`
	Address       string  `json:"address" xml:"address,attr"`
	City          string  `json:"city" xml:"city,attr"`
	State         State   `json:"state" xml:"state"`
	PostalCode    string  `json:"postal_code" xml:"postal_code,attr"`
	TollFreePhone string  `json:"toll_free_phone" xml:"toll_free_phone,attr"`
	Latitude      float64 `json:"latitude" xml:"latitude,attr"`
	Longitude     float64 `json:"longitude" xml:"longitude,attr"`
	Fax           string  `json:"fax" xml:"fax,attr"`
	LocalPhone    string  `json:"local_phone" xml:"local_phone,attr"`
	Manager       string  `json:"manager" xml:"manager"`
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

type FeedRecord struct {
	Warehouse  Warehouse `json:"warehouse" xml:"warehouse"`
	Part       int       `json:"part" xml:"part,attr"`
	Quantity   int       `json:"quantity" xml:"quantity,attr"`
	DateUpdate time.Time `json:"date_updated" xml:"date_updated`
}

func (p *Part) GetInventory(apiKey, warehouseCode string) error {
	redis_key := fmt.Sprintf("part:%d:inventory", p.ID)

	data, err := redis.Get(redis_key)
	if err != nil {
		if err == redix.ErrNil {
			return nil
		}
		return err
	}

	var recs []FeedRecord
	if err = json.Unmarshal(data, &recs); err != nil {
		return err
	}

	for _, rec := range recs {
		i := Inventory{
			Part:        rec.Part,
			Warehouse:   rec.Warehouse,
			Quantity:    rec.Quantity,
			DateUpdated: rec.DateUpdate,
		}
		if (warehouseCode != "" && warehouseCode == i.Warehouse.Code) || warehouseCode == "" {
			p.Inventory.Warehouses = append(p.Inventory.Warehouses, i)
		}
	}

	for _, w := range p.Inventory.Warehouses {
		p.Inventory.TotalAvailability = p.Inventory.TotalAvailability + w.Quantity
	}

	return nil
}
