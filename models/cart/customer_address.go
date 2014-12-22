package cart

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
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
	if a1.Address1 != a2.Address2 {
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
