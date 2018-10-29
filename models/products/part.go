package products

import (
	"fmt"
	"sort"
	"strconv"

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
	PriceCode         string               `json:"price_code" xml:"price_code,attr" bson:"price_code"`
	RelatedCount      int                  `json:"related_count" xml:"related_count,attr" bson:"related_count"`
	AverageReview     float64              `json:"average_review" xml:"average_review,attr" bson:"average_review"`
	DateModified      time.Time            `json:"date_modified" xml:"date_modified,attr" bson:"date_modified"`
	DateAdded         time.Time            `json:"date_added" xml:"date_added,attr" bson:"date_added"`
	ShortDesc         string               `json:"short_description" xml:"short_description,attr" bson:"short_description"`
	InstallSheet      *url.URL             `json:"install_sheet" xml:"install_sheet" bson:"install_sheet"`
	Attributes        []Attribute          `json:"attributes" xml:"attributes" bson:"attributes"`
	AppSpecific       bool                 `bson:"application_specific" json:"application_specific" xml:"application_specific,attr"`
	AcesVehicles      []AcesVehicle        `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`
	VehicleAttributes []string             `json:"vehicle_atttributes" xml:"vehicle_attributes" bson:"vehicle_attributes"`
	Vehicles          []VehicleApplication `json:"vehicle_applications,omitempty" xml:"vehicle_applications,omitempty" bson:"vehicle_applications"`
	LuverneVehicles   []LuverneApplication `json:"luverne_applications,omitempty" xml:"luverne_applications,omitempty" bson:"luverne_applications"`
	Content           []Content            `json:"content" xml:"content" bson:"content"`
	Pricing           []Price              `json:"pricing" xml:"pricing" bson:"pricing"`
	Reviews           []Review             `json:"reviews" xml:"reviews" bson:"reviews"`
	Images            []Image              `json:"images" xml:"images" bson:"images"`
	Related           []int                `json:"related" xml:"related" bson:"related" bson:"related"`
	ReplacedBy        int                  `bson:"replaced_by" json:"replaced_by,omitempty" xml:"replaced_by,omitempty"`
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
	MappedToVehicle   bool                 `json:"mappedToVehicle" xml:"mappedToVehicle" bson:"mappedToVehicle,omitempty"`
	WebVisibility     string               `json:"-" xml:"-" bson:"web_visibility"`
	ShowOnWebsite     bool                 `json:"showOnWebsite" xml:"showOnWebsite" bson:"showOnWebsite"`
	ShowForLoggedIn   bool                 `json:"showForLoggedIn" xml:"showForLoggedIn" bson:"showForLoggedIn"`
	Tariff            string               `json:"tariff" xml:"tariff" bson:"tariff"`
	ComplexPart       *ComplexPart         `bson:"complex_part" json:"complex_part,omitempty" xml:"complex_part,omitempty"`
}

type SkuCount struct {
	Sku   string `bson:"sku" json:"sku,omitempty" xml:"sku,omitempty"`
	Count uint32 `bson:"count" json:"count,omitempty" xml:"count,omitempty"`
}

type ComplexPart struct {
	Type     string      `bson:"type" json:"type,omitempty" xml:"type,omitempty"`
	SkuCount []*SkuCount `bson:"skuCount" json:"skuCount,omitempty" xml:"skuCount,omitempty"`
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

// LuverneApplication defines a unique vehicle fitment description in the Luverne
// data set.
type LuverneApplication struct {
	Year      string `bson:"year" json:"year" xml:"year"`
	Make      string `bson:"make" json:"make" xml:"make"`
	Model     string `bson:"model" json:"model" xml:"model"`
	Body      string `bson:"body" json:"body" xml:"body"`
	BoxLength string `bson:"boxLength" json:"boxLength" xml:"boxLength"`
	CabLength string `bson:"cabLength" json:"cabLength" xml:"cabLength"`
	FuelType  string `bson:"fuelType" json:"fuelType" xml:"fuelType"`
	WheelType string `bson:"wheelType" json:"wheelType" xml:"wheelType"`
}

const (
	PUBLIC   = "Public"
	DISABLED = "Disabled"
	LOGGEDIN = "Logged In Only"
)

func GetMany(ids, brands []int, sess *mgo.Session) ([]Part, error) {

	c := sess.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)
	var statuses = []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950}
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
	parts, err := BindCustomerToSeveralParts([]Part{*p}, dtx)
	if len(parts) > 0 {
		*p = parts[0]
	}
	if err := p.FromDatabase(brands); err != nil {
		return err
	}

	return err
}

// GetMulti ...
func GetMulti(dtx *apicontext.DataContext, ids []string) ([]Part, error) {
	var err error
	//get brands
	brands := getBrandsFromDTX(dtx)

	if err := database.Init(); err != nil {
		return nil, err
	}

	session := database.ProductMongoSession.Copy()
	defer session.Close()

	query := bson.M{"part_number": bson.M{"$in": ids}, "brand.id": bson.M{"$in": brands}}

	var parts []Part
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).All(&parts)
	if err != nil {
		return nil, err
	}

	return BindCustomerToSeveralParts(parts, dtx)
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
			"$in": []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950},
		},
	}

	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(qry).Distinct("part_number", &parts)
	if err != nil {
		return parts, err
	}

	sort.Strings(parts)

	return parts, nil
}

func All(page, count int, dtx *apicontext.DataContext, from time.Time, to time.Time) ([]Part, int, error) {
	var total int
	var query bson.M
	//Currently a list of all visibilities, might add or subtract later
	visibility := []string{PUBLIC, DISABLED, LOGGEDIN}
	brands := getBrandsFromDTX(dtx)
	parts := make([]Part, 0)

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, total, err
	}
	defer session.Close()

	//In this block we are determining the query that we will be sending to the DB
	//time.Time.isZero is effecively the nil check for time.Time objects
	if !time.Time.IsZero(from) || !time.Time.IsZero(to) {
		//In the case that "from" is specified
		if !time.Time.IsZero(from) {
			//In the case both "to" and "from" are specified
			if !time.Time.IsZero(to) {
				query = bson.M{"brand.id": bson.M{"$in": brands},
					"web_visibility": bson.M{"$in": visibility},
					"date_modified":  bson.M{"$lte": to, "$gte": from}}
			} else { //In the case only "from" is specified
				query = bson.M{"brand.id": bson.M{"$in": brands},
					"web_visibility": bson.M{"$in": visibility},
					"date_modified":  bson.M{"$gte": from}}
			}
		} else if !time.Time.IsZero(to) { //In the case only "to" is specified
			query = bson.M{"brand.id": bson.M{"$in": brands},
				"web_visibility": bson.M{"$in": visibility},
				"date_modified":  bson.M{"$lte": to}}
		}

	} else { //In the case neither "to" or "from" are specified
		query = bson.M{"brand.id": bson.M{"$in": brands}, "web_visibility": bson.M{"$in": visibility}}
	}

	//We get the count here so that we can return it as part of the JSON response
	total, err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).Count()
	if err != nil {
		return parts, total, err
	}

	//A Mongo index is needed to ensure that the sort doesn't consume too much memory
	//See INDEX.md in root directory
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).Sort("id").Skip(page * count).Limit(count).All(&parts)

	//Determining Web Visibility, based on flags that we are given by data team
	for ind := range parts {
		if parts[ind].WebVisibility == PUBLIC {
			parts[ind].ShowForLoggedIn = false
			if parts[ind].Status >= 700 {
				parts[ind].ShowOnWebsite = true
			} else {
				parts[ind].ShowOnWebsite = false
			}

		} else if parts[ind].WebVisibility == DISABLED {
			parts[ind].ShowOnWebsite = false
			parts[ind].ShowForLoggedIn = false
		} else if parts[ind].WebVisibility == LOGGEDIN {
			parts[ind].ShowForLoggedIn = true
			if parts[ind].Status >= 700 {
				parts[ind].ShowOnWebsite = true
			} else {
				parts[ind].ShowOnWebsite = false
			}

		}
	}
	return parts, total, err
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
	var content []Content

	refChan := make(chan int)
	contentChan := make(chan int)

	price, _ = customer.GetCustomerPrice(dtx, p.ID)

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

	<-refChan
	<-contentChan
	p.Content = append(p.Content, content...)
	p.Customer.Price = price
	p.Customer.CartReference = ref
	return
}

func BindCustomerToSeveralParts(parts []Part, dtx *apicontext.DataContext) ([]Part, error) {
	if len(parts) < 1 {
		return parts, nil
	}
	var partIDs string
	var err error
	for i, p := range parts {
		if i > 0 {
			partIDs += ","
		}
		partIDs += strconv.Itoa(p.ID)
	}

	err = database.Init()
	if err != nil {
		return parts, err
	}

	statement := fmt.Sprintf(`select distinct ci.custPartID, cp.price, cp.partID from ApiKey as ak
						join CustomerUser cu on ak.user_id = cu.id
						join Customer c on cu.cust_ID = c.cust_id
						left join CustomerPricing cp on cp.cust_ID = cu.cust_ID
						left join CartIntegration ci on c.cust_ID = ci.custID && cp.partID = ci.partID
						where ak.api_key = '%s'
						and cp.partID in (%s)`, dtx.APIKey, partIDs)

	stmt, err := database.DB.Prepare(statement)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return parts, err
	}
	var custPartID, partID *int
	var price *float64
	custPartMap := make(map[int]int)
	custPriceMap := make(map[int]float64)

	for res.Next() {
		err = res.Scan(
			&custPartID,
			&price,
			&partID,
		)
		if err != nil {
			return parts, err
		}
		if custPartID != nil && partID != nil {
			custPartMap[*partID] = *custPartID
		}
		if price != nil && partID != nil {
			custPriceMap[*partID] = *price
		}
	}

	custContentMap := make(map[int][]Content)
	allPartContent, err := custcontent.GetAllPartContent(dtx.APIKey)
	if err != nil {
		return parts, err
	}
	for _, c := range allPartContent {
		for _, content := range c.Content {
			custContentMap[c.PartId] = append(custContentMap[c.PartId], Content{Text: content.Text, ContentType: ContentType{Type: content.ContentType.Type, AllowsHTML: content.ContentType.AllowHtml}})
		}
	}

	for i, part := range parts {
		var ok bool
		if _, ok = custPriceMap[part.ID]; ok {
			parts[i].Customer.Price = custPriceMap[part.ID]
		}
		if _, ok = custPartMap[part.ID]; ok {
			parts[i].Customer.CartReference = custPartMap[part.ID]
		}
		if _, ok = custContentMap[part.ID]; ok {
			parts[i].Content = custContentMap[part.ID]
		}
	}
	return parts, err
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

	parts, err := BindCustomerToSeveralParts([]Part{*p}, dtx)
	if len(parts) > 0 {
		*p = parts[0]
	}

	//Determining Web Visibility, based on flags that we are given by data team

	if p.WebVisibility == PUBLIC {
		p.ShowForLoggedIn = false
		if p.Status >= 700 {
			p.ShowOnWebsite = true
		} else {
			p.ShowOnWebsite = false
		}

	} else if p.WebVisibility == DISABLED {
		p.ShowOnWebsite = false
		p.ShowForLoggedIn = false
	} else if p.WebVisibility == LOGGEDIN {
		p.ShowForLoggedIn = true
		if p.Status >= 700 {
			p.ShowOnWebsite = true
		} else {
			p.ShowOnWebsite = false
		}

	}

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
