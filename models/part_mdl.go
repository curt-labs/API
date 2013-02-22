package models

import (
	"../helpers/database"
	"errors"
	"strings"
	"time"
)

var (
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
				from Part as p
				left join Class as pc on p.classID = pc.classID
				where p.partID = %d limit 1`

	partAttrStmt = `select field, value from PartAttribute where partID = %d`

	partPriceStmt = `select priceType, price from Price where partID = %d`

	partReviewStmt = `select rating,subject,review_text,name,email,createdDate from Review
				where partID = %d and approved = 1 and active = 1`
)

type Part struct {
	PartId, CustPartId, Status, PriceCode, RelatedCount int
	AverageReview                                       float64
	DateModified, DateAdded                             time.Time
	ShortDesc, PartClass                                string
	Attributes                                          []Attribute
	VehicleAttributes                                   []string
	Content                                             []Content
	Pricing                                             []Pricing
	Reviews                                             []Review
	Images                                              []Image
	Related                                             []Part
	Categories                                          []Category
	Videos                                              []Video
	Packages                                            []Package
	//Vehicles                                            []Vehicle
}

type Attribute struct {
	Key, Value string
}

type Content struct {
	Key, Value string
}

type Pricing struct {
	Type     string
	Price    float64
	Enforced bool
}

func (p *Part) Get() error {

	var errs []string

	basicChan := make(chan int)
	attrChan := make(chan int)
	priceChan := make(chan int)
	reviewChan := make(chan int)
	imageChan := make(chan int)

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

	go func() {
		priceErr := p.GetPricing()
		if priceErr != nil {
			errs = append(errs, priceErr.Error())
		}
		priceChan <- 1
	}()

	go func() {
		reviewErr := p.GetReviews()
		if reviewErr != nil {
			errs = append(errs, reviewErr.Error())
		}
		reviewChan <- 1
	}()

	go func() {
		imgErr := p.GetImages()
		if imgErr != nil {
			errs = append(errs, imgErr.Error())
		}
		imageChan <- 1
	}()

	<-basicChan
	<-attrChan
	<-priceChan
	<-reviewChan
	<-imageChan

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}
	return nil
}

func (p *Part) GetWithVehicle(vehicle *Vehicle) error {

	var errs []string

	superChan := make(chan int)
	noteChan := make(chan int)
	go func() {
		p.Get()
		superChan <- 1
	}()
	go func() {
		notes, nErr := vehicle.GetNotes(p.PartId)
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

func (p *Part) GetPricing() error {
	db := database.Db

	rows, res, err := db.Query(partPriceStmt, p.PartId)
	if database.MysqlError(err) {
		return err
	}

	typ := res.Map("priceType")
	price := res.Map("price")

	var prices []Pricing
	for _, row := range rows {
		pr := Pricing{
			row.Str(typ),
			row.Float(price),
			false,
		}

		if pr.Type == "Map" {
			pr.Enforced = true
		}
		prices = append(prices, pr)
	}

	p.Pricing = prices

	return nil
}
