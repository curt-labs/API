package part

import (
	"../../helpers/database"
	"../categories"
	"../images"
	"../packages"
	"../reviews"
	"../videos"
	"errors"
	"log"
	"strings"
	"time"
)

var (
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
				from Part as p
				left join Class as pc on p.classID = pc.classID
				where p.partID = %d limit 1`

	partAttrStmt = `select field, value from PartAttribute where partID = %d`
)

type Part struct {
	PartId, CustPartId, Status, PriceCode, RelatedCount int
	InstallTime, AverageReview                          float64
	DateModified, DateAdded                             time.Time
	ShortDesc, PartClass, Drilling, Exposed             string
	Attributes                                          []Attribute
	VehicleAttributes                                   []Attribute
	Content                                             []Content
	Pricing                                             []Pricing
	Reviews                                             []reviews.Review
	Images                                              []images.Image
	Related                                             []Part
	Categories                                          []categories.Category
	Videos                                              []videos.Video
	Packages                                            []packages.Package
	Vehicles                                            []interface{}
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

func (p *Part) Get() error {

	var errs []string

	basicChan := make(chan int)
	attrChan := make(chan int)
	go func() {
		basicErr := p.Basics()
		if basicErr != nil {
			errs = append(errs, basicErr.Error())
		}
		basicChan <- 1
	}()

	go func() {
		attrErr := p.GetAttributes()
		if attrErr != nil {
			errs = append(errs, attrErr.Error())
		}
		attrChan <- 1
	}()

	<-basicChan
	<-attrChan

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}
	return nil
}

func (p *Part) GetById(id int) {
	p.PartId = id

	p.Get()
}

func (p *Part) GetAttributes() (err error) {
	db := database.Db

	rows, _, err := db.Query(partAttrStmt, p.PartId)
	if database.MysqlError(err) {
		return err
	}

	var attrs []Attribute
	for _, row := range rows {
		attr := Attribute{
			Key:   row.Str(0),
			Value: row.Str(1),
		}
		attrs = append(attrs, attr)
	}
	log.Println(attrs)
	p.Attributes = attrs

	return
}

func (p *Part) Basics() error {
	db := database.Db

	row, res, err := db.QueryFirst(basicsStmt, p.PartId)
	if database.MysqlError(err) {
		return err
	}
	status := res.Map("status")
	dateAdded := res.Map("dateAdded")
	dateModified := res.Map("dateModified")
	shortDesc := res.Map("shortDesc")
	partID := res.Map("partID")
	priceCode := res.Map("priceCode")
	class := res.Map("class")

	p.PartId = row.Int(partID)
	date_add, _ := time.Parse("2006-01-02 15:04:01", row.Str(dateAdded))
	p.DateAdded = date_add

	date_mod, _ := time.Parse("2006-01-02 15:04:01", row.Str(dateModified))
	p.DateModified = date_mod

	p.ShortDesc = row.Str(shortDesc)
	p.PriceCode = row.Int(priceCode)
	p.PartClass = row.Str(class)
	p.Status = row.Int(status)

	return nil
}
