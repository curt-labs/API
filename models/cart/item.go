package cart

import (
	"gopkg.in/mgo.v2/bson"
)

type LineItem struct {
	Id                  bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	FulFillableQuantity int           `json:"fulfillable_quantity" xml:"fulfillable_quantity,attr" bson:"fulfillable_quantity"`
	FulFillmentService  string        `json:"fulfillment_service" xml:"fulfillment_service,attr" bson:"fulfillment_service"`
	FulFillmentStatus   string        `json:"fulfillment_status" xml:"fulfillment_status,attr" bson:"fulfillment_status"`
	Grams               int           `json:"grams" xml:"grams,attr" bson:"grams"`
	Price               float64       `json:"price" xml:"price,attr" bson:"price"`
	ProductId           string        `json:"product_id" xml:"product_id,attr" bson:"product_id"`
	Quantity            int           `json:"quantity" xml:"quantity,attr" bson:"quantity"`
	RequiresShipping    bool          `json:"requires_shipping" xml:"requires_shipping,attr" bson:"requires_shipping"`
	SKU                 string        `json:"sku" xml:"sku,attr" bson:"sku"`
	Title               string        `json:"title,attr" xml:"title,attr" bson:"title"`
	VariantId           int           `json:"variant_id" xml:"variant_id,attr" bson:"variant_id"`
	Vendor              string        `json:"vendor" xml:"vendor,attr" bson:"vendor"`
	Name                string        `json:"name" xml:"name,attr" bson:"name"`
	GiftCard            bool          `json:"gift_cart" xml:"gift_cart,attr" bson:"gift_cart"`
	Taxable             bool          `json:"taxable" xml:"taxable,attr" bson:"taxable"`
	TaxLines            []TaxLine     `json:"tax_lines" xml:"tax_lines" bson:"tax_lines"`
}
