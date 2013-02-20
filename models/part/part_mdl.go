package part

import (
	"../../helpers/database"
	"../categories"
	"../images"
	"../packages"
	"../reviews"
	"../vehicle"
	"../videos"
	"log"
	"time"
)

var (
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
				from Part as p
				left join Class as pc on p.classID = pc.classID
				where p.partID = %d limit 1`
)

type Part struct {
	PartId, CustPartId, Status, PriceCode, RelatedCount int
	InstallTime, AverageReview                          float64
	DateModified, DateAdded                             time.Time
	ShortDesc, PartClass, Drilling, Exposed             string
	Attributes                                          []Attribute
	VehicleAttributes                                   []vehicle.Attribute
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
	err := p.Basics()
	if err != nil {
		return err
	}
	return nil
}

func (p *Part) GetById(id int) {
	p.PartId = id
	p.Get()
}

func (p *Part) Basics() error {
	db := database.Db

	log.Printf(basicsStmt, p.PartId)

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
	p.DateAdded, _ = time.Parse("2006-01-02 15:04", row.Str(dateAdded))
	p.DateModified, _ = time.Parse("2006-01-02 15:04", row.Str(dateModified))
	p.ShortDesc = row.Str(shortDesc)
	p.PriceCode = row.Int(priceCode)
	p.PartClass = row.Str(class)
	p.Status = row.Int(status)

	return nil

}
