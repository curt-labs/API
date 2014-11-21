package cart

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAdmin/helpers/geocoding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type Shop struct {
	Id                      bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	Address1                string        `json:"address1" xml:"address1,attr" bason:"address1"`
	City                    string        `json:"city" xml:"city,attr" bson:"city"`
	Country                 string        `json:"country" xml:"country>country" bson:"country"`
	CountryCode             string        `json:"country_code" xml:"country>code" bson:"country_code"`
	CountryName             string        `json:"country_name" xml:"country>name" bson:"country_name"`
	CreatedAt               time.Time     `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Currency                string        `json:"currency" xml:"currency,attr" bson:"currency"`
	Domain                  *url.URL      `json:"domain" xml:"domain,attr" bson:"domain"`
	Email                   string        `json:"email" xml:"email" bson:"email"`
	Latitude                float64       `json:"latitude" xml:"geo>latitude" bson:"latitude"`
	Longitude               float64       `json:"longitude" xml:"geo>longitude" bson:"longitude"`
	MoneyFormat             string        `json:"money_format" xml:"money_format,attr" bson:"money_format"`
	MoneyWithCurrencyFormat string        `json:"money_with_currency_fromat" xml:"money_with_currency_format,attr" bson:"money_with_currency_format"`
	Name                    string        `json:"name" xml:"name,attr" bson:"name"`
	PasswordEnabled         bool          `json:"password_enabled" xml:"password_enabled,attr" bson:"password_enabled"`
	Phone                   string        `json:"phone" xml:"phone,attr" bson:"phone"`
	Province                string        `json:"province" xml:"province,attr" bson:"province"`
	ProvinceCode            string        `json:"province_code" xml:"province_code,attr" bson:"province_code"`
	Public                  bool          `json:"public" xml:"public,attr" bson:"public"`
	ShopOwner               string        `json:"shop_owner" xml:"shop_owner,attr" bson:"shop_owner"`
	Source                  string        `json:"source" xml:"source,attr" bson:"source"`
	TaxShipping             bool          `json:"tax_shipping" xml:"taxing>tax_shipping,attr" bson:"tax_shipping"`
	TaxesIncluded           bool          `json:"taxes_included" xml:"taxing>taxes_included,attr" bson:"taxes_included"`
	CountyTaxes             bool          `json:"county_taxes" xml:"taxing>county_taxes,attr" bson:"county_taxes"`
	Timezone                string        `json:"timezone" xml:"timezone,attr" bson:"timezone"`
	Zip                     string        `json:"zip" xml:"zip,attr" bson:"zip"`
	HasStorefront           bool          `json:"has_storefront" xml:"has_storefront,attr" bson:"has_storefront"`
}

func (sh *Shop) Get() error {

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	c := sess.DB("CurtCart").C("shop")
	err = c.Find(bson.M{"_id": sh.Id}).One(&sh)
	if err != nil {
		return err
	}

	return nil
}

// This method is used explicitly for generating test data
// DO NOT EXPOSE
func insertTestData() *bson.ObjectId {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return nil
	}

	sh := Shop{}
	sh.Id = bson.NewObjectId()
	sh.Name = "Test Shop"
	sh.Address1 = "1119 Sunset Lane"
	sh.City = "Altoona"
	sh.Province = "Wisconsin"
	sh.ProvinceCode = "WI"
	sh.Country = "US"
	sh.CountryCode = "US"
	sh.CountryName = "United States"
	sh.Zip = "54720"
	sh.CreatedAt = time.Now()
	sh.Currency = "USD"
	sh.Domain, _ = url.Parse("http://store.ninneman.org")
	sh.Email = "alex@ninneman.org"
	sh.MoneyFormat = "$"
	sh.MoneyWithCurrencyFormat = "$ USD"
	sh.PasswordEnabled = true
	sh.Phone = "7153082604"
	sh.Public = false
	sh.ShopOwner = "Alex Ninneman"
	sh.TaxShipping = true
	sh.TaxesIncluded = false
	sh.CountyTaxes = false
	sh.Timezone = "US/Central"
	sh.HasStorefront = false

	l := geocoding.Lookup{
		Address: fmt.Sprintf("%s, %s, %s %s", sh.Address1, sh.City, sh.ProvinceCode, sh.Zip),
	}
	resp, err := l.Search()
	if err == nil && len(resp.Results) > 0 {
		sh.Longitude = resp.Results[0].Geometry.Location.Longitude
		sh.Latitude = resp.Results[0].Geometry.Location.Latitude
	}

	if err := sess.DB("CurtCart").C("shop").Insert(sh); err != nil {
		return nil
	}

	return &sh.Id
}
