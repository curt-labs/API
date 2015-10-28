package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/brand"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/curt-labs/GoAPI/models/video"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"net/url"
	"strings"
	"time"
)

type Part struct {
	Identifier        bson.ObjectId     `bson:"_id" json:"-" xml:"-"`
	ID                int               `json:"id" xml:"id,attr" bson:"id"`
	PartNumber        string            `bson:"part_number" json:"part_number" xml:"part_number"`
	Brand             brand.Brand       `json:"brand" xml:"brand,attr" bson:"brand"`
	Status            int               `json:"status" xml:"status,attr" bson:"status"`
	PriceCode         int               `json:"price_code" xml:"price_code,attr" bson:"price_code"`
	RelatedCount      int               `json:"related_count" xml:"related_count,attr" bson:"related_count"`
	AverageReview     float64           `json:"average_review" xml:"average_review,attr" bson:"average_review"`
	DateModified      time.Time         `json:"date_modified" xml:"date_modified,attr" bson:"date_modified"`
	DateAdded         time.Time         `json:"date_added" xml:"date_added,attr" bson:"date_added"`
	ShortDesc         string            `json:"short_description" xml:"short_description,attr" bson:"short_description"`
	InstallSheet      *url.URL          `json:"install_sheet" xml:"install_sheet" bson:"install_sheet"`
	Attributes        []Attribute       `json:"attributes" xml:"attributes" bson:"attributes"`
	AcesVehicles      []AcesVehicle     `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`
	VehicleAttributes []string          `json:"vehicle_atttributes" xml:"vehicle_attributes" bson:"vehicle_attributes"`
	Vehicles          []vehicle.Vehicle `json:"vehicles,omitempty" xml:"vehicles,omitempty" bson:"vehicles"`
	Content           []Content         `json:"content" xml:"content" bson:"content"`
	Pricing           []Price           `json:"pricing" xml:"pricing" bson:"pricing"`
	Reviews           []Review          `json:"reviews" xml:"reviews" bson:"reviews"`
	Images            []Image           `json:"images" xml:"images" bson:"images"`
	Related           []int             `json:"related" xml:"related" bson:"related" bson:"related"`
	Categories        []Category        `json:"categories" xml:"categories" bson:"categories"`
	Videos            []video.Video     `json:"videos" xml:"videos" bson:"videos"`
	Packages          []Package         `json:"packages" xml:"packages" bson:"packages"`
	Customer          CustomerPart      `json:"customer,omitempty" xml:"customer,omitempty" bson:"v"`
	Class             Class             `json:"class,omitempty" xml:"class,omitempty" bson:"class"`
	Featured          bool              `json:"featured,omitempty" xml:"featured,omitempty" bson:"featured"`
	AcesPartTypeID    int               `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty" bson:"acesPartTypeId"`
	Installations     []Installation    `json:"installation,omitempty" xml:"installation,omitempty" bson:"installation"`
	Inventory         PartInventory     `json:"inventory,omitempty" xml:"inventory,omitempty" bson:"inventory"`
	UPC               string            `json:"upc,omitempty" xml:"upc,omitempty" bson:"upc"`
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

type Installation struct { //aka VehiclePart Table
	ID          int             `json:"id,omitempty" xml:"id,omitempty" bson:"id"`
	Vehicle     vehicle.Vehicle `json:"vehicle,omitempty" xml:"vehicle,omitempty" bson:"vehicle"`
	Part        Part            `json:"part,omitempty" xml:"part,omitempty" bson:"part"`
	Drilling    string          `json:"drilling,omitempty" xml:"v,omitempty" bson:"drilling"`
	Exposed     string          `json:"exposed,omitempty" xml:"exposed,omitempty" bson:"exposed"`
	InstallTime int             `json:"installTime,omitempty" xml:"installTime,omitempty" bson:"installTime"`
}

func (p *Part) Get(dtx *apicontext.DataContext) error {
	var err error
	//get brands
	brands := getBrandsFromDTX(dtx)

	customerChan := make(chan CustomerPart)
	databaseChan := make(chan error)

	go func(api_key string) {
		customerChan <- p.BindCustomer(dtx)
	}(dtx.APIKey)

	go func(brands []int) {
		if err := p.FromDatabase(brands); err != nil {
			databaseChan <- err
			return
		}
		databaseChan <- nil
	}(brands)

	p.Customer = <-customerChan
	err = <-databaseChan
	close(customerChan)
	close(databaseChan)
	return err
}

func (p *Part) FromDatabase(brands []int) error {
	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	query := bson.M{"id": p.ID, "brand.id": bson.M{"$in": brands}}
	return session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).One(&p)
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

func Featured(count int, dtx *apicontext.DataContext) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"featured": true, "brand.id": bson.M{"$in": brands}}).Sort("id:1").Limit(count).All(&parts)
	return parts, err
}

func Latest(count int, dtx *apicontext.DataContext) ([]Part, error) {
	var parts []Part
	brands := getBrandsFromDTX(dtx)

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"brand.id": bson.M{"$in": brands}}).Sort("-date_added").Limit(count).All(&parts)
	return parts, err
}

func (p *Part) GetWithVehicle(vehicle *vehicle.Vehicle, api_key string, dtx *apicontext.DataContext) error {
	var errs []string

	superChan := make(chan int)
	noteChan := make(chan int)
	go func(key string) {
		p.Get(dtx)
		superChan <- 1
	}(api_key)
	go func() {
		notes, nErr := vehicle.GetNotes(p.ID)
		if nErr != nil && notes != nil {
			errs = append(errs, nErr.Error())
			p.VehicleAttributes = []string{}
		} else {
			p.VehicleAttributes = notes
		}
		noteChan <- 1
	}()

	<-superChan
	<-noteChan

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}
	return nil
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

func (p *Part) BindCustomer(dtx *apicontext.DataContext) CustomerPart {
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

	return CustomerPart{
		Price:         price,
		CartReference: ref,
	}
}

// PartBreacrumbs
//
// Description: Builds out Category breadcrumb array for the current part object.
//
// Inherited: part Part
// Returns: error
func (p *Part) PartBreadcrumbs(dtx *apicontext.DataContext) error {
	if p.ID == 0 {
		return errors.New("Invalid Part Number")
	}

	//check redis!
	redis_key := fmt.Sprintf("part:%d:breadcrumbs:%s", p.ID, dtx.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Categories); err == nil {
			return nil
		}
	}

	// Oh alright, let's talk with our database
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	partCategoryStmt, err := db.Prepare(PartCategoryStmt)
	if err != nil {
		return err
	}
	defer partCategoryStmt.Close()

	lookupCategoriesStmt, err := db.Prepare(CategoriesByBrandStmt)
	if err != nil {
		return err
	}
	defer lookupCategoriesStmt.Close()

	// Execute SQL Query against current ID
	catRow := partCategoryStmt.QueryRow(p.ID)
	if catRow == nil {
		return errors.New("No part found for " + string(p.ID))
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch, dtx)
	initCat := <-ch
	close(ch)

	// Build thee lookup
	catLookup := make(map[int]Category)
	rows, err := lookupCategoriesStmt.Query(dtx.BrandID)
	if err != nil {
		return err
	}
	defer rows.Close()

	multiChan := make(chan []Category, 0)
	go PopulateCategoryMulti(rows, multiChan)
	cats := <-multiChan
	close(multiChan)

	for _, cat := range cats {
		catLookup[cat.CategoryID] = cat
	}

	// Okay, let's put it together!
	var categories []Category
	categories = append(categories, initCat)

	nextParentID := initCat.ParentID
	for {
		if nextParentID == 0 {
			break
		}
		if c, found := catLookup[nextParentID]; found {
			nextParentID = c.ParentID
			categories = append(categories, c)
			continue
		}
		nextParentID = 0
	}

	p.Categories = categories
	if dtx.BrandString != "" {
		go func(cats []Category) {
			redis.Setex(redis_key, cats, redis.CacheTimeout)
		}(p.Categories)
	}

	return nil
}

func (p *Part) GetPartByPartNumber() (err error) {
	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"part_number": p.PartNumber}).One(&p)
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
