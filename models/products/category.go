package products

import (
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/video"
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

	Sort            int                      `bson:"sort" json:"sort" xml:"sort,attr"`
	DateAdded       time.Time                `bson:"date_added" json:"date_added" xml:"date_added,attr"`
	Title           string                   `bson:"title" json:"title" xml:"title,attr"`
	ShortDesc       string                   `bson:"short_description" json:"short_description" xml:"short_description"`
	LongDesc        string                   `bson:"long_description" json:"long_description" xml:"long_description"`
	ColorCode       string                   `bson:"color_code" json:"color_code" xml:"color_code,attr"`
	FontCode        string                   `bson:"font_code" json:"font_code" xml:"font_code,attr"`
	Image           *url.URL                 `bson:"image" json:"image" xml:"image"`
	Icon            *url.URL                 `bson:"icon" json:"icon" xml:"icon"`
	IsLifestyle     bool                     `bson:"lifestyle" json:"lifestyle" xml:"lifestyle,attr"`
	VehicleSpecific bool                     `bson:"vehicle_specific" json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired bool                     `bson:"vehicle_required" json:"vehicle_required" xml:"vehicle_required,attr"`
	MetaTitle       string                   `bson:"meta_title" json:"meta_title" xml:"meta_title"`
	MetaDescription string                   `bson:"meta_description" json:"meta_description" xml:"meta_description"`
	MetaKeywords    string                   `bson:"meta_keywords" json:"meta_keywords" xml:"meta_keywords"`
	Content         []Content                `bson:"content" json:"content" xml:"content"`
	Videos          []video.Video            `bson:"videos" json:"videos" xml:"videos"`
	PartIDs         []int                    `bson:"part_ids" json:"part_ids" xml:"part_ids"`
	Brand           brand.Brand              `bson:"brand" json:"brand" xml:"brand"`
	ProductListing  *PaginatedProductListing `bson:"product_listing" json:"product_listing" xml:"product_listing"`
	PDFpath         *url.URL                 `bson:"pdf_path" json:"pdf_path" xml:"pdf_path"`
	XLSpath         *url.URL                 `bson:"xls_path" json:"xls_path" xml:"xls_path"`
}
