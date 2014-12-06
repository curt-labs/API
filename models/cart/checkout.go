package cart

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type Checkout struct {
	Id                    bson.ObjectId   `json:"id" xml:"id,attr" bson:"_id"`
	BillingAddress        CustomerAddress `json:"billing_address" xml:"billing_address" bson:"billing_address"`
	BuyerAcceptsMarketing bool            `json:"buyer_accepts_marketing" xml:"buyer_accepts_marketing,attr" bson:"buyer_accepts_marketing"`
	CancelReason          string          `json:"cancel_reason" xml:"cancel_reason,attr" bson:"cancel_reason"`
	CartToken             string          `json:"cart_token" xml:"cart_token,attr" bson:"cart_token"`
	ClosedAt              time.Time       `json:"closed_at" xml:"closed_at,attr" bson:"closed_at"`
	CreatedAt             time.Time       `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Currency              string          `json:"currency" xml:"currency,attr" bson:"currency"`
	Customer              Customer        `json:"customer" xml:"customer" bson:"customer"`
	DiscountCodes         []DiscountCode  `json:"discount_codes" xml:"discount_codes>discount_code" bson:"discount_codes"`
	Email                 string          `json:"email" xml:"email,attr" bson:"email"`
	Gateway               string          `json:"gateway" xml:"gateway,attr" bson:"gateway"`
	LandingSite           string          `json:"landing_site" xml:"landing_site" bson:"landing_site"`
	Items                 []LineItem      `json:"line_items" xml:"line_items>line_item" bson:"line_items"`
	Note                  string          `json:"note" xml:"note,attr" bson:"note"`
	ReferringSite         url.URL         `json:"referring_site" xml:"referring_site" bson:"referring_site"`
	ShippingAddress       CustomerAddress `json:"shipping_address" xml:"shipping_address" bson:"shipping_address"`
	ShippingLines         []ShippingLine  `json:"shipping_lines" xml:"shipping_lines,shipping_line" bson:"shipping_lines"`
	SourceName            string          `json:"source_name" xml:"source_name,attr" bson:"source_name"`
	SubtotalPrice         float64         `json:"subtotal_price" xml:"subtotal_price,attr" bson:"subtotal_price"`
	TaxLines              []TaxLine       `json:"tax_lines" xml:"tax_lines>tax_line" bson:"tax_lines"`
	Token                 string          `json:"token" xml:"token,attr" bson:"token"`
	TotalDiscounts        string          `json:"total_discounts" xml:"total_discounts,attr" bson:"total_discounts"`
	TotalLineItemsPrice   string          `json:"total_line_items_price" xml:"total_line_items_price,attr" bson:"total_line_items_price"`
}
