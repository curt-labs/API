package cart

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CustomerAddress struct {
	Id           *bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	Address1     string         `json:"address1" xml:"address1,attr" bson:"address1"`
	Address2     string         `json:"address2" xml:"address2,attr" bson:"address2"`
	City         string         `json:"city" xml:"city,attr" bson:"city"`
	Company      string         `json:"company" xml:"company,attr" bson:"company"`
	Name         string         `json:"name" xml:"name,attr" bson:"name"`
	FirstName    string         `json:"first_name" xml:"first_name,attr" bson:"first_name"`
	LastName     string         `json:"last_name" xml:"last_name,attr" bson:"last_name"`
	Phone        string         `json:"phone" xml:"phone,attr" bson:"phone"`
	Province     string         `json:"province" xml:"geo>province>province,attr" bson:"province"`
	ProvinceCode string         `json:"province_code" xml:"geo>province>code,attr" bson:"province_code"`
	Country      string         `json:"country" xml:"geo>country>country,attr" bson:"country"`
	CountryCode  string         `json:"country_code" xml:"geo>country>code,attr" bson:"country_code"`
	CountryName  string         `json:"country_name" xml:"geo>country>name,attr" bson:"country_name"`
	Zip          string         `json:"zip" xml:"geo>zip,attr" bson:"zip"`
	CreatedAt    time.Time      `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
}

func (c *Customer) AddAddress(addr CustomerAddress) error {
	if c.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "cannot update a customer that doesn't exist")
	}
	if addr.Id == nil || !addr.Id.Valid() {
		addrId := bson.NewObjectId()
		addr.Id = &addrId
	}

	if err := addr.Validate(); err != nil {
		return err
	}

	addr.UpdatedAt = time.Now()
	addr.CreatedAt = time.Now()

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$addToSet": bson.M{
				"addresses": addr,
			},
		},
	}

	_, err = sess.DB("CurtCart").C("customer").Find(bson.M{"_id": c.Id, "shop_id": c.ShopId}).Apply(change, c)

	return err
}

func (c *Customer) SaveAddress(addr CustomerAddress) error {
	if c.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "cannot update a customer that doesn't exist")
	}

	if err := addr.Validate(); err != nil {
		return err
	}

	addr.UpdatedAt = time.Now()

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	qry := bson.M{
		"addresses._id": addr.Id,
	}
	change := bson.M{
		"$set": bson.M{
			"customer.$.address": addr,
		},
	}
	return sess.DB("CurtCart").C("customer").Update(qry, change)
}

func (c *CustomerAddress) Validate() error {
	if c.Address1 == "" {
		return fmt.Errorf("error: %s", "address must be provided")
	}
	if c.City == "" {
		return fmt.Errorf("error: %s", "city must be provided")
	}
	if c.Province == "" && c.ProvinceCode == "" {
		return fmt.Errorf("error: %s", "province information must be provided")
	}
	if c.CountryName == "" && c.Country == "" && c.CountryCode == "" {
		return fmt.Errorf("error: %s", "country information must be provided")
	}
	if c.Zip == "" {
		return fmt.Errorf("error: %s", "post code must be provided")
	}
	return nil
}

func (a1 *CustomerAddress) deepEqual(a2 *CustomerAddress) bool {
	if a1 == nil && a2 == nil {
		return true
	}

	if (a1 == nil && a2 != nil) || (a1 != nil && a2 == nil) {
		return false
	}
	if a1.Address1 != a2.Address1 {
		return false
	}
	if a1.Address2 != a2.Address2 {
		return false
	}
	if a1.City != a2.City {
		return false
	}
	if a1.Company != a2.Company {
		return false
	}
	if a1.Name != a2.Name {
		return false
	}
	if a1.FirstName != a2.FirstName {
		return false
	}
	if a1.LastName != a2.LastName {
		return false
	}
	if a1.Phone != a2.Phone {
		return false
	}
	if a1.Province != a2.Province {
		return false
	}
	if a1.ProvinceCode != a2.ProvinceCode {
		return false
	}
	if a1.Country != a2.Country {
		return false
	}
	if a1.CountryCode != a2.CountryCode {
		return false
	}
	if a1.CountryName != a2.CountryName {
		return false
	}
	if a1.Zip != a2.Zip {
		return false
	}
	return true
}
