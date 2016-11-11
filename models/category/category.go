package category

import (
	"math"
	"net/url"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/video"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	statuses = []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950}
)

type Category struct {
	Identifier       bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	CategoryID       int           `bson:"id" json:"id" xml:"id,attr"`
	ParentIdentifier bson.ObjectId `bson:"parent_identifier,omitempty" json:"parent_identifier" xml:"parent_identifier"`
	ParentID         int           `bson:"parent_id" json:"parent_id" xml:"parent_id,attr"`
	Children         []Category    `bson:"children" json:"children" xml:"children"`

	Sort               int                               `bson:"sort" json:"sort" xml:"sort,attr"`
	DateAdded          time.Time                         `bson:"date_added" json:"date_added" xml:"date_added,attr"`
	Title              string                            `bson:"title" json:"title" xml:"title,attr"`
	ShortDesc          string                            `bson:"short_description" json:"short_description" xml:"short_description"`
	LongDesc           string                            `bson:"long_description" json:"long_description" xml:"long_description"`
	ColorCode          string                            `bson:"color_code" json:"color_code" xml:"color_code,attr"`
	FontCode           string                            `bson:"font_code" json:"font_code" xml:"font_code,attr"`
	Image              *url.URL                          `bson:"image" json:"image" xml:"image"`
	Icon               *url.URL                          `bson:"icon" json:"icon" xml:"icon"`
	IsLifestyle        bool                              `bson:"lifestyle" json:"lifestyle" xml:"lifestyle,attr"`
	VehicleSpecific    bool                              `bson:"vehicle_specific" json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired    bool                              `bson:"vehicle_required" json:"vehicle_required" xml:"vehicle_required,attr"`
	MetaTitle          string                            `bson:"meta_title" json:"meta_title" xml:"meta_title"`
	MetaDescription    string                            `bson:"meta_description" json:"meta_description" xml:"meta_description"`
	MetaKeywords       string                            `bson:"meta_keywords" json:"meta_keywords" xml:"meta_keywords"`
	ProductListing     *products.PaginatedProductListing `json:"product_listing,omitempty" xml:"product_listing,omitempty"`
	Content            []Content                         `bson:"content" json:"content" xml:"content"`
	Videos             []video.Video                     `bson:"videos" json:"videos" xml:"videos"`
	Brand              brand.Brand                       `bson:"brand" json:"brand" xml:"brand"`
	ProductIdentifiers []int                             `bson:"part_ids" json:"part_identifiers" xml:"part_identifiers"`
	IsDeleted          bool                              `bson:"isdeleted" json:"-" xml:"-"`
}

type PartResponse struct {
	Parts      []products.Part `json:"parts"`
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
}

func GetCategoryTree(dtx *apicontext.DataContext) ([]Category, error) {
	var cats []Category

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return cats, err
	}
	defer session.Close()
	query := bson.M{"parent_id": 0, "isdeleted": false, "is_lifestyle": false, "brand.id": bson.M{"$in": dtx.BrandArray}}
	err = session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(query).Sort("sort").All(&cats)
	for i, _ := range cats {
		cats[i].removeDeletedChildren()
	}
	return cats, err
}

func (c *Category) Get(page, count int) error {

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(bson.M{"id": c.CategoryID, "isdeleted": false}).One(&c)
	if err != nil {
		return err
	}
	c.removeDeletedChildren()

	c.ProductListing = &products.PaginatedProductListing{
		Page:    page,
		PerPage: count,
		Parts:   []products.Part{},
	}

	c.ProductListing.TotalItems, err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"id": bson.M{"$in": c.ProductIdentifiers}}).Count()
	if err != nil {
		c.ProductListing.TotalItems = 1
	}

	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"id": bson.M{"$in": c.ProductIdentifiers}}).Sort("id").Skip((page - 1) * count).Limit(count).All(&c.ProductListing.Parts)
	if err != nil {
		return err
	}

	c.ProductListing.ReturnedCount = len(c.ProductListing.Parts)
	c.ProductListing.TotalPages = c.ProductListing.TotalItems / c.ProductListing.PerPage

	return nil
}

func GetCategoryParts(catId, page, count int) (PartResponse, error) {
	var parts PartResponse

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()

	//get category's children
	var cat Category
	err = session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(bson.M{"id": catId}).Select(bson.M{"children": 1}).One(&cat)
	if err != nil {
		return parts, err
	}

	children := []int{catId}
	for _, child := range cat.Children {
		children = append(children, child.CategoryID)
	}
	parts.Page = page

	//get parts of category and its children
	query := bson.M{
		"categories": bson.M{
			"$elemMatch": bson.M{
				"id": bson.M{
					"$in": children,
				},
			},
		},
		"status": bson.M{
			"$in": statuses,
		},
	}

	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).Limit(count).Skip((page - 1) * count).All(&parts.Parts)
	if err != nil {
		return parts, err
	}

	//get total parts count
	total_items, err := session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).Count()
	parts.TotalPages = int(math.Ceil(float64(total_items) / float64(count)))
	return parts, err
}

func (c *Category) removeDeletedChildren() {
	var newChildren []Category
	for i, _ := range c.Children {
		c.Children[i].removeDeletedChildren()
		if !c.Children[i].IsDeleted {
			newChildren = append(newChildren, c.Children[i])
		}
	}
	c.Children = newChildren
}
