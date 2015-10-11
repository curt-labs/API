package mongoData

import (
	"github.com/curt-labs/GoAPI/models/brand"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type Category struct {
	Identifier       bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	CategoryID       int           `bson:"id" json:"id" xml:"id,attr"`
	ParentIdentifier bson.ObjectId `bson:"parent_identifier,omitempty" json:"parent_identifier" xml:"parent_identifier"`
	ParentID         int           `bson:"parent_id" json:"parent_id" xml:"parent_id,attr"`
	Children         []Category    `bson:"children" json:"children" xml:"children"`

	Sort               int                      `bson:"sort" json:"sort" xml:"sort,attr"`
	DateAdded          time.Time                `bson:"date_added" json:"date_added" xml:"date_added,attr"`
	Title              string                   `bson:"title" json:"title" xml:"title,attr"`
	ShortDesc          string                   `bson:"short_description" json:"short_description" xml:"short_description"`
	LongDesc           string                   `bson:"long_description" json:"long_description" xml:"long_description"`
	ColorCode          string                   `bson:"color_code" json:"color_code" xml:"color_code,attr"`
	FontCode           string                   `bson:"font_code" json:"font_code" xml:"font_code,attr"`
	Image              *url.URL                 `bson:"image" json:"image" xml:"image"`
	Icon               *url.URL                 `bson:"icon" json:"icon" xml:"icon"`
	IsLifestyle        bool                     `bson:"lifestyle" json:"lifestyle" xml:"lifestyle,attr"`
	VehicleSpecific    bool                     `bson:"vehicle_specific" json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired    bool                     `bson:"vehicle_required" json:"vehicle_required" xml:"vehicle_required,attr"`
	MetaTitle          string                   `bson:"meta_title" json:"meta_title" xml:"meta_title"`
	MetaDescription    string                   `bson:"meta_description" json:"meta_description" xml:"meta_description"`
	MetaKeywords       string                   `bson:"meta_keywords" json:"meta_keywords" xml:"meta_keywords"`
	ProductListing     *PaginatedProductListing `json:"product_listing,omitempty" xml:"product_listing,omitempty"`
	Content            []Content                `bson:"content" json:"content" xml:"content"`
	Videos             []Video                  `bson:"videos" json:"videos" xml:"videos"`
	Brand              brand.Brand              `bson:"brand" json:"brand" xml:"brand"`
	ProductIdentifiers []int                    `bson:"part_ids" json:"part_identifiers" xml:"part_identifiers"`
}

type PaginatedProductListing struct {
	Parts         []Product `json:"parts" xml:"parts"`
	TotalItems    int       `json:"total_items" xml:"total_items,attr"`
	ReturnedCount int       `json:"returned_count" xml:"returned_count,attr"`
	Page          int       `json:"page" xml:"page,attr"`
	PerPage       int       `json:"per_page" xml:"per_page,attr"`
	TotalPages    int       `json:"total_pages" xml:"total_pages,attr"`
}

// BaseVehicle ...
type BaseVehicle struct {
	Year  int    `bson:"year" json:"year" xml:"year"`
	Make  string `bson:"make" json:"make" xml:"make"`
	Model string `bson:"model" json:"model" xml:"model"`
}

// AcesVehicle ...
type AcesVehicle struct {
	Base       BaseVehicle    `bson:"base" json:"base" xml:"base"`
	Submodel   string         `bson:"submodel" json:"submodel" xml:"submodel"`
	Attributes []ConfigOption `bson:"attributes" json:"attributes" xml:"attributes"`
}

// ConfigOption ...
type ConfigOption struct {
	Key   string `bson:"name" json:"name" xml:"name,attr"`
	Value string `bson:"value" json:"value" xml:"value,attr"`
}

// Image ...
type Image struct {
	Size   string   `bson:"size" json:"size" xml:"size"`
	Sort   string   `bson:"sort" json:"sort" xml:"sort"`
	Height int      `bson:"height" json:"height" xml:"height"`
	Width  int      `bson:"width" json:"width" xml:"width"`
	Path   *url.URL `bson:"path" json:"path" xml:"path"`
}

// Video ...
type Video struct {
	Title        string            `bson:"title" json:"title" xml:"title,attr"`
	SubjectType  string            `bson:"subject_type" json:"subject_type" xml:"subject_type"`
	Description  string            `bson:"description" json:"description" xml:"description"`
	DateAdded    time.Time         `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time         `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	Thumbnail    string            `bson:"thumb_nail" json:"thumbnail" xml:"thumbnail"`
	Channels     []Channel         `bson:"channel" json:"channel" xml:"channel"`
	Files        []CdnFile         `bson:"cdn_file" json:"cdn_file" xml:"cdn_file"`
	Categories   []CatAssociation  `bson:"category_assoc" json:"cateogry_assoc" xml:"category_assoc"`
	Parts        []PartAssociation `bson:"part_assoc" json:"part_assoc" xml:"part_assoc"`
}

// Channell ...
type Channel struct {
	Type         string    `bson:"type" json:"type" xml:"type,attr"`
	Link         string    `bson:"link" json:"link" xml:"link"`
	EmbedCode    string    `bson:"embed_code" json:"embed_code" xml:"embed_code"`
	ForiegnID    string    `bson:"foreign_id" json:"foreign_id" xml:"foreign_id"`
	DateAdded    time.Time `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	Title        string    `bson:"title" json:"title" xml:"title,attr"`
	Description  string    `bson:"description" json:"description" xml:"description"`
	Duration     string    `bson:"duration" json:"duration" xml:"duration"`
}

// CdnFile ...
type CdnFile struct {
	Type         CdnFileType `bson:"type" json:"type" xml:"type"`
	Path         string      `bson:"path" json:"path" xml:"path"`
	Bucket       string      `bson:"bucket" json:"bucket" xml:"bucket"`
	ObjectName   string      `bson:"object_name" json:"object_name" xml:"object_name"`
	FileSize     string      `bson:"file_size" json:"file_size" xml:"file_size"`
	DateAdded    time.Time   `bson:"date_added" json:"date_added" xml:"date_added"`
	DateModified time.Time   `bson:"date_modified" json:"date_modified" xml:"date_modified"`
	LastUploaded string      `bson:"date_uploaded" json:"date_uploaded" xml:"date_uploaded"`
}

// CdnFileType ...
type CdnFileType struct {
	MimeType    string `bson:"mime_type" json:"mime_type" xml:"mime_type"`
	Title       string `bson:"title" json:"title" xml:"title,attr"`
	Description string `bson:"description" json:"description" xml:"description"`
}

// VideoType ..
type VideoType struct {
	Name string `bson:"name" json:"name" xml:"name"`
	Icon string `bson:"icon" json:"icon" xml:"icon"`
}

// CatAssociation ...
type CatAssociation struct {
	CatID    int    `bson:"id" json:"id" xml:"id"`
	Title    string `bson:"title" json:"title" xml:"title"`
	ParentID int    `bson:"parent_id" json:"parent_id" xml:"parent_id"`
}

// PartAssociation ...
type PartAssociation struct {
	PartID    int    `bson:"id" json:"id" xml:"id"`
	ShortDesc string `bson:"short_desc" json:"short_desc" xml:"short_desc"`
	IsPrimary bool   `bson:"primary" json:"primary" xml:"primary"`
}

// Attribute ...
type Attribute struct {
	Key   string `bson:"name" json:"name" xml:"name,attr"`
	Value string `bson:"value" json:"value" xml:"value,attr"`
}

// Content ...
type Content struct {
	Text        string      `bson:"text" json:"text" xml:"text"`
	ContentType ContentType `json:"contentType" xml:"contentType"`
}

// ContentType ...
type ContentType struct {
	Type       string `bson:"type" json:"type" xml:"type"`
	AllowsHTML bool   `bson:"allows_html" json:"allows_html" xml:"allows_html"`
}

// Warehouse ...
type Warehouse struct {
	Name          string  `bson:"name" json:"name" xml:"name,attr"`
	Code          string  `bson:"code" json:"code" xml:"code,attr"`
	Address       string  `bson:"address" json:"address" xml:"address,attr"`
	City          string  `bson:"city" json:"city" xml:"city,attr"`
	State         State   `bson:"state" json:"state" xml:"state"`
	PostalCode    string  `bson:"postal_code" json:"postal_code" xml:"postal_code,attr"`
	TollFreePhone string  `bson:"toll_free_phone" json:"toll_free_phone" xml:"toll_free_phone,attr"`
	Latitude      float64 `bson:"latitude" json:"latitude" xml:"latitude,attr"`
	Longitude     float64 `bson:"longitude" json:"longitude" xml:"longitude,attr"`
	Fax           string  `bson:"fax" json:"fax" xml:"fax,attr"`
	LocalPhone    string  `bson:"local_phone" json:"local_phone" xml:"local_phone,attr"`
	Manager       string  `bson:"manager" json:"manager" xml:"manager"`
}

// PartInventory ...
type PartInventory struct {
	TotalAvailability int         `bson:"total_availability" json:"total_availability" xml:"total_availability,attr"`
	Warehouses        []Inventory `bson:"inventory" json:"inventory" xml:"inventory"`
}

// Inventory ...
type Inventory struct {
	Warehouse   Warehouse `bson:"warehouse" json:"warehouse" xml:"warehouse"`
	Quantity    int       `bson:"quantity" json:"quantity" xml:"quantity,attr"`
	DateUpdated time.Time `bson:"date_updated" json:"date_updated" xml:"date_update,attr"`
}

// State ...
type State struct {
	State        string  `bson:"state" json:"state" xml:"state,abbr"`
	Abbreviation string  `bson:"abbreviation" json:"abbreviation" xml:"abbreviation,attr"`
	Country      Country `bson:"country" json:"country" xml:"country"`
}

// Country ...
type Country struct {
	Name         string `bson:"name" json:"name" xml:"name,attr"`
	Abbreviation string `bson:"abbreviation" json:"abbreviation" xml:"abbreviation,attr"`
}

// Package ...
type Package struct {
	Height             float64 `bson:"height" json:"height" xml:"height"`
	Width              float64 `bson:"width" json:"width" xml:"width"`
	Length             float64 `bson:"length" json:"length" xml:"length"`
	Weight             float64 `bson:"weight" json:"weight" xml:"weight"`
	DimensionUnit      string  `bson:"dimensionUnit" json:"dimensionUnit" xml:"dimensionUnit"`
	DimensionUnitLabel string  `bson:"dimensionUnitLabel"json:"dimensionUnitLabel" xml:"dimensionUnitLabel"`
	WeightUnit         string  `bson:"weightUnit" json:"weightUnit" xml:"weightUnit"`
	WeightUnitLabel    string  `bson:"weightUnitLabel" json:"weightUnitLabel" xml:"weightUnitLabel"`
	PackageUnit        string  `bson:"packageUnit" json:"packageUnit" xml:"packageUnit"`
	PackageUnitLabel   string  `bson:"packageUnitLabel" json:"packageUnitLabel" xml:"packageUnitLabel"`
	Quantity           int     `bson:"quantity" json:"quantity" xml:"quantity"`
	PackageType        string  `bson:"name" json:"name" xml:"name"`
}

// Price ...
type Price struct {
	Type         string    `bson:"type" json:"type" xml:"type"`
	Price        float64   `bson:"price" json:"price" xml:"price"`
	Enforced     bool      `bson:"enforced" json:"enforced" xml:"enforced"`
	DateModified time.Time `bson:"date_modified"json:"date_modified" xml:"date_modified"`
}

// Product ...
type Product struct {
	Identifier     bson.ObjectId        `bson:"_id" json:"-" xml:"-"`
	ProductID      int                  `bson:"id" json:"id" xml:"id,attr"`
	PartNumber     string               `bson:"part_number" json:"part_number" xml:"part_number"`
	Brand          brand.Brand          `bson:"brand" json:"brand" xml:"brand,attr"`
	Status         int                  `bson:"status" json:"status" xml:"status,attr"`
	PriceCode      int                  `bson:"price_code" json:"price_code" xml:"price_code,attr"`
	AverageReview  float64              `bson:"average_review" json:"average_review" xml:"average_review,attr"`
	DateModified   time.Time            `bson:"date_modified" json:"date_modified" xml:"date_modified,attr"`
	DateAdded      time.Time            `bson:"date_added" json:"date_added" xml:"date_added,attr"`
	ShortDesc      string               `bson:"short_description" json:"short_description" xml:"short_description,attr"`
	InstallSheet   *url.URL             `bson:"install_sheet" json:"install_sheet" xml:"install_sheet"`
	Attributes     []Attribute          `bson:"attributes" json:"attributes" xml:"attributes"`
	AcesVehicles   []AcesVehicle        `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`
	Content        []Content            `bson:"content" json:"content" xml:"content"`
	Pricing        []Price              `bson:"pricing" json:"pricing" xml:"pricing"`
	Reviews        []Review             `bson:"reviews" json:"reviews" xml:"reviews"`
	Images         []Image              `bson:"images" json:"images" xml:"images"`
	Related        []Product            `bson:"related" json:"related" xml:"related"`
	Categories     []Category           `bson:"categories" json:"categories" xml:"categories"`
	Videos         []Video              `bson:"videos" json:"videos" xml:"videos"`
	Packages       []Package            `bson:"packages" json:"packages" xml:"packages"`
	Class          Class                `bson:"class" json:"class,omitempty" xml:"class,omitempty"`
	Featured       bool                 `bson:"featured" json:"featured" xml:"featured"`
	AcesPartTypeID int                  `bson:"acesPartTypeId,omitempty" json:"acesPartTypeId" xml:"acesPartTypeId"`
	Vehicles       []VehicleApplication `bson:"vehicle_applications"  json:"vehicle_applications" xml:"vehicle_applications"`
	Inventory      PartInventory        `bson:"inventory" json:"inventory" xml:"inventory"`
	UPC            string               `bson:"upc" json:"upc" xml:"upc"`
}

// Class ...
type Class struct {
	Name  string `bson:"name"json:"name" xml:"name,omitempty"`
	Image string `bson:"image"json:"image" xml:"image,omitempty"`
}

// Review ...
type Review struct {
	Rating      int       `bson:"rating" json:"rating" xml:"rating"`
	Subject     string    `bson:"subject" json:"subject" xml:"subject"`
	ReviewText  string    `bson:"review_text" json:"review_text" xml:"review_text"`
	Name        string    `bson:"name" json:"name" xml:"name"`
	Email       string    `bson:"email" json:"email" xml:"email"`
	CreatedDate time.Time `bson:"created_date" json:"created_date" xml:"created_date"`
}

// VehicleApplication ...
type VehicleApplication struct {
	Year        string `bson:"year" json:"year" xml:"year,attr"`
	Make        string `bson:"make" json:"make" xml:"make,attr"`
	Model       string `bson:"model" json:"model" xml:"model,attr"`
	Style       string `bson:"style" json:"style" xml:"style,attr"`
	Exposed     string `bson:"exposed" json:"exposed" xml:"exposed"`
	Drilling    string `bson:"drilling" json:"drilling" xml:"drilling"`
	InstallTime string `bson:"install_time" json:"install_time" xml:"install_time"`
}
