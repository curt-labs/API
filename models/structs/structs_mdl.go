package structs

import (
	"net/url"
	"time"
)

type ConfigResponse struct {
	ConfigOption ConfigOption
	Matched      *ProductMatch
}

type ConfigOption struct {
	Type    string
	Options []string
}

type ProductMatch struct {
	Parts  []Part
	Groups []int
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
	Parts                 []interface{}
	Groups                []interface{}
}

type Part struct {
	PartId, CustPartId, Status, PriceCode, RelatedCount int
	InstallTime, AverageReview                          float64
	DateModified, DateAdded                             time.Time
	ShortDesc, PartClass, Drilling, Exposed             string
	Attributes                                          []Attribute
	VehicleAttributes                                   []Attribute
	Content                                             []Content
	Pricing                                             []Pricing
	Reviews                                             []Review
	Images                                              []Image
	Related                                             []Part
	Categories                                          []Category
	Videos                                              []Video
	Packages                                            []Package
	Vehicles                                            []Vehicle
}

type Attribute struct {
	Key, Value string
}

type Content struct {
	Key, Value string
}

type Pricing struct {
	Type  string
	Price float64
}

type Category struct {
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image                        url.URL
	IsLifestyle, VehicleSpecific bool
}

type Image struct {
	Size, Sort            string
	PartId, Height, Width int
	Path                  url.URL
}

type Package struct {
	Height, Width, Length, Quantity   int
	Weight                            float64
	DimensionUnit, DimensionUnitLabel string
	WeightUnit, WeightUnitLabel       string
	PackageUnti, PackageUnitLabel     string
}

type Review struct {
	PartId, Rating                   int
	Subject, ReviewText, Name, Email string
	CreatedDate                      time.Time
}

type Video struct {
	YouTubeVideoId, Type string
	IsPrimary            bool
	TypeId               int
	TypeIcon             url.URL
}
