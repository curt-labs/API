package cart

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CustomerAddress struct {
	Id           bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	Address1     string        `json:"address1" xml:"address1,attr" bson:"address1"`
	Address2     string        `json:"address2" xml:"address2,attr" bson:"address2"`
	City         string        `json:"city" xml:"city,attr" bson:"city"`
	Company      string        `json:"company" xml:"company,attr" bson:"company"`
	Name         string        `json:"name" xml:"name,attr" bson:"name"`
	FirstName    string        `json:"first_name" xml:"first_name,attr" bson:"first_name"`
	LastName     string        `json:"last_name" xml:"last_name,attr" bson:"last_name"`
	Phone        string        `json:"phone" xml:"phone,attr" bson:"phone"`
	Province     string        `json:"province" xml:"geo>province>province,attr" bson:"province"`
	ProvinceCode string        `json:"province_code" xml:"geo>province>code,attr" bson:"province_code"`
	Country      string        `json:"country" xml:"geo>country>country,attr" bson:"country"`
	CountryCode  string        `json:"country_code" xml:"geo>country>code,attr" bson:"country_code"`
	CountryName  string        `json:"country_name" xml:"geo>country>name,attr" bson:"country_name"`
	Zip          string        `json:"zip" xml:"geo>zip,attr" bson:"zip"`
}
