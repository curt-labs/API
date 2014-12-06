package cart

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type Fulfillment struct {
	Id                         bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	CreatedAt                  time.Time     `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	LineItems                  []LineItem    `json:"line_items" xml:"line_items>line_item" bson:"line_items"`
	OrderId                    bson.ObjectId `json:"order_id" xml:"order_id,attr" bson:"order_id"`
	Receipt                    Receipt       `json:"receipt" xml:"receipt" bson:"receipt"`
	Status                     string        `json:"status" xml:"status,attr" bson:"status"`
	TrackingCompany            string        `json:"tracking_company" xml:"tracking_company,attr" bson:"tracking_company"`
	TrackingNumber             []string      `json:"tracking_number" xml:"tracking_numbers>tracking_number" bson:"tracking_number"`
	TrackingUrls               []url.URL     `json:"tracking_urls" xml:"tracking_urls>tracking_url" bson:"tracking_urls"`
	UpdatedAt                  time.Time     `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
	VariantInventoryManagement string        `json:"variant_inventory_management" xml:"variant_inventory_management,attr" bson:"variant_inventory_management"`
}

type Receipt struct {
	TestCase      bool   `json:"test_case" xml:"test_case,attr" bson:"test_case"`
	Authorization string `json:"authorization" xml:"authorization,attr" bson:"authorization"`
}
