package products

import (
	"time"
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
	Warehouses        []Inventory `json:"inventory,omitempty" xml:"inventory,omitempty"`
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
