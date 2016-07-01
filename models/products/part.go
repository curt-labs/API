package products

import (
	"sort"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/customer/content"
	"github.com/curt-labs/API/models/video"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/url"
	"strings"
	"time"
)

type Part struct {
	Identifier        bson.ObjectId        `bson:"_id" json:"-" xml:"-"`
	ID                int                  `json:"id" xml:"id,attr" bson:"id"`
	PartNumber        string               `bson:"part_number" json:"part_number" xml:"part_number,attr"`
	Brand             brand.Brand          `json:"brand" xml:"brand,attr" bson:"brand"`
	Status            int                  `json:"status" xml:"status,attr" bson:"status"`
	PriceCode         int                  `json:"price_code" xml:"price_code,attr" bson:"price_code"`
	RelatedCount      int                  `json:"related_count" xml:"related_count,attr" bson:"related_count"`
	AverageReview     float64              `json:"average_review" xml:"average_review,attr" bson:"average_review"`
	DateModified      time.Time            `json:"date_modified" xml:"date_modified,attr" bson:"date_modified"`
	DateAdded         time.Time            `json:"date_added" xml:"date_added,attr" bson:"date_added"`
	ShortDesc         string               `json:"short_description" xml:"short_description,attr" bson:"short_description"`
	InstallSheet      *url.URL             `json:"install_sheet" xml:"install_sheet" bson:"install_sheet"`
	Attributes        []Attribute          `json:"attributes" xml:"attributes" bson:"attributes"`
	AcesVehicles      []AcesVehicle        `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`
	VehicleAttributes []string             `json:"vehicle_atttributes" xml:"vehicle_attributes" bson:"vehicle_attributes"`
	Vehicles          []VehicleApplication `json:"vehicle_applications,omitempty" xml:"vehicle_applications,omitempty" bson:"vehicle_applications"`
	Content           []Content            `json:"content" xml:"content" bson:"content"`
	Pricing           []Price              `json:"pricing" xml:"pricing" bson:"pricing"`
	Reviews           []Review             `json:"reviews" xml:"reviews" bson:"reviews"`
	Images            []Image              `json:"images" xml:"images" bson:"images"`
	Related           []int                `json:"related" xml:"related" bson:"related" bson:"related"`
	Categories        []Category           `json:"categories" xml:"categories" bson:"categories"`
	Videos            []video.Video        `json:"videos" xml:"videos" bson:"videos"`
	Packages          []Package            `json:"packages" xml:"packages" bson:"packages"`
	Customer          CustomerPart         `json:"customer,omitempty" xml:"customer,omitempty" bson:"v"`
	Class             Class                `json:"class,omitempty" xml:"class,omitempty" bson:"class"`
	Featured          bool                 `json:"featured,omitempty" xml:"featured,omitempty" bson:"featured"`
	AcesPartTypeID    int                  `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty" bson:"acesPartTypeId"`
	Inventory         PartInventory        `json:"inventory,omitempty" xml:"inventory,omitempty" bson:"inventory"`
	UPC               string               `json:"upc,omitempty" xml:"upc,omitempty" bson:"upc"`
	Layer             string               `json:"iconLayer" xml:"iconLayer" bson:"iconLayer"`
}

type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
}

type PaginatedProductListing struct {
	Parts         []Part `json:"parts" xml:"parts"`
	TotalItems    int    `json:"total_items" xml:"total_items,attr"`
	ReturnedCount int    `json:"returned_count" xml:"returned_count,attr"`
	Page          int    `json:"page" xml:"page,attr"`
	PerPage       int    `json:"per_page" xml:"per_page,attr"`
	TotalPages    int    `json:"total_pages" xml:"total_pages,attr"`
}

type Class struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty" bson:"id"`
	Name  string `json:"name,omitempty" xml:"name,omitempty" bson:"name"`
	Image string `json:"image,omitempty" xml:"image,omitempty" bson:"image"`
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

func GetMany(ids, brands []int, sess *mgo.Session) ([]Part, error) {

	c := sess.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)
	var statuses = []int{800, 900}
	qry := bson.M{"id": bson.M{"$in": ids}, "status": bson.M{"$in": statuses}, "brand.id": bson.M{"$in": brands}}

	var parts []Part
	err := c.Find(qry).All(&parts)

	return parts, err
}

// Get ...
func (p *Part) Get(dtx *apicontext.DataContext) error {
	var err error
	//get brands
	brands := getBrandsFromDTX(dtx)

	customerChan := make(chan bool)

	go func(api_key string) {
		// p.BindCustomer(dtx)
		customerChan <- true
	}(dtx.APIKey)

	if err := p.FromDatabase(brands); err != nil {
		return err
	}

	<-customerChan

	return err
}

func (p *Part) GetNoCust(dtx *apicontext.DataContext, sess *mgo.Session) error {
	var err error
	//get brands
	brands := getBrandsFromDTX(dtx)
	if err := p.FromMongoDatabase(brands, sess); err != nil {
		return err
	}
	return err
}

func (p *Part) FromMongoDatabase(brands []int, session *mgo.Session) error {
	query := bson.M{"id": p.ID, "brand.id": bson.M{"$in": brands}}
	return session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).One(&p)
}

// FromDatabase ...
func (p *Part) FromDatabase(brands []int) error {
	if err := database.Init(); err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	query := bson.M{"id": p.ID, "brand.id": bson.M{"$in": brands}}

	return session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).One(&p)
}

// Identifiers ...
func Identifiers(brand int, dtx *apicontext.DataContext) ([]string, error) {
	var parts []string
	brands := []int{brand}

	if brand == 0 {
		brands = getBrandsFromDTX(dtx)
	}

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()

	qry := bson.M{
		"brand.id": bson.M{
			"$in": brands,
		},
		"status": bson.M{
			"$in": []int{800, 900},
		},
	}

	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(qry).Distinct("part_number", &parts)
	if err != nil {
		return parts, err
	}

	sort.Strings(parts)

	return parts, nil
}

func All(page, count int, dtx *apicontext.DataContext) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"brand.id": bson.M{"$in": brands}}).Sort("id:1").Skip(page * count).Limit(count).All(&parts)
	return parts, err
}

func Featured(count int, dtx *apicontext.DataContext, brand int) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)
	if brand > 0 {
		brands = []int{brand}
	}

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"featured": true, "brand.id": bson.M{"$in": brands}}).Sort("id:1").Limit(count).All(&parts)
	return parts, err
}

func Latest(count int, dtx *apicontext.DataContext, brand int) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)
	if brand > 0 {
		brands = []int{brand}
	}

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"brand.id": bson.M{"$in": brands}}).Sort("-date_added").Limit(count).All(&parts)
	return parts, err
}

func (p *Part) GetRelated(dtx *apicontext.DataContext) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()

	query := bson.M{
		"id": bson.M{
			"$in": p.Related,
		},
		"brand.id": bson.M{
			"$in": brands,
		},
	}
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).Sort("id:1").All(&parts)
	return parts, err
}

func (p *Part) BindCustomer(dtx *apicontext.DataContext) {
	var price float64
	var ref int

	priceChan := make(chan int)
	refChan := make(chan int)
	contentChan := make(chan int)

	go func() {
		price, _ = customer.GetCustomerPrice(dtx, p.ID)
		priceChan <- 1
	}()

	go func() {
		ref, _ = customer.GetCustomerCartReference(dtx.APIKey, p.ID)
		refChan <- 1
	}()

	go func() {
		content, _ := custcontent.GetPartContent(p.ID, dtx.APIKey)
		for _, con := range content {

			strArr := strings.Split(con.ContentType.Type, ":")
			cType := con.ContentType.Type
			if len(strArr) > 1 {
				cType = strArr[1]
			}
			var c Content
			c.ContentType.Type = cType
			c.Text = con.Text
			p.Content = append(p.Content, c)
		}
		contentChan <- 1
	}()

	<-priceChan
	<-refChan
	<-contentChan

	p.Customer.Price = price
	p.Customer.CartReference = ref
	return
}

func (p *Part) GetPartByPartNumber(dtx *apicontext.DataContext) (err error) {
	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	pattern := bson.RegEx{
		Pattern: "^" + p.PartNumber + "$",
		Options: "i",
	}
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"part_number": pattern}).One(&p)
	if err != nil {
		return err
	}
	// p.BindCustomer(dtx)
	return err
}

func getBrandsFromDTX(dtx *apicontext.DataContext) []int {
	var brands []int
	if dtx.BrandID == 0 {
		for _, b := range dtx.BrandArray {
			brands = append(brands, b)
		}
	} else {
		brands = append(brands, dtx.BrandID)
	}
	return brands
}
