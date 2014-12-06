package cart

import (
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type ClientDetails struct {
	Id             bson.ObjectId `json:"id" xml:"id" bson:"_id"`
	AcceptLanguage string        `json:"accept_language" xml:"accept_language,attr" bson:"accept_language"`
	BrowserHeight  float64       `json:"browser_height" xml:"browser>height,attr" bson:"browser_height"`
	BrowserWidth   float64       `json:"browser_width" xml:"browser>width,attr" bson:"browser_width"`
	BrowserIP      string        `json:"browser_ip" xml:"browser>ip,attr" bson:"browser_ip"`
	SessionHash    string        `json:"session_hash" xml:"session_hash,attr" bson:"session_hash"`
	UserAgent      string        `json:"user_agent" xml:"user_agent,attr" bson:"user_agent"`
}

type DiscountCode struct {
	Id     bson.ObjectId `json:"id" xml:"id,attr" bson:"_id"`
	Amount float64       `json:"amount" xml:"amount,attr" bson:"amount"`
	Code   string        `json:"code" xml:"code,attr" bson:"code"`
	Type   string        `json:"type" xml:"type,attr" bson:"type"`
}

type NoteAttribute struct {
	Name  string `json:"name" xml:"name,attr" bson:"name"`
	Value string `json:"value" xml:"value,attr" bson:"value"`
}

type TaxLine struct {
	Price float64 `json:"price" xml:"price,attr" bson:"price"`
	Rate  float64 `json:"rate" xml:"rate,attr" bson:"rate"`
	Title string  `json:"title" xml:"title,attr" bson:"title"`
}

type ShippingLine struct {
	Code     string    `json:"code" xml:"code,attr" bson:"code"`
	Price    float64   `json:"price" xml:"price,attr" bson:"price"`
	Source   string    `json:"source" xml:"source,attr" bson:"source"`
	Title    string    `json:"title" xml:"title,attr" bson:"title"`
	TaxLines []TaxLine `json:"tax_lines" xml:"tax_lines" bson:"tax_lines"`
}

type Order struct {
	Id                    bson.ObjectId   `json:"id" xml:"id,attr" bson:"_id"`
	BillingAddress        CustomerAddress `json:"billing_address" xml:"billing_address" bson:"billing_address"`
	BrowserIP             string          `json:"browser_ip" xml:"browser_ip,attr" bson:"browser_ip"`
	BuyerAcceptsMarketing bool            `json:"buyer_accepts_marketing" xml:"buyer_accepts_marketing" bson:"buyer_accepts_marketing"`
	CancelReason          string          `json:"cancel_reason" xml:"cancel_reason" bson:"cancel_reason"`
	CancelledAt           time.Time       `json:"cancelled_at" xml:"cancelled_at" bson:"cancelled_at"`
	CartToken             bson.ObjectId   `json:"cart_token" xml:"cart_token" bson:"cart_token"`
	ClientDetails         ClientDetails   `json:"client_details" xml:"client_details" bson:"client_details"`
	ClosedAt              time.Time       `json:"closed_at" xml:"closed_at,attr" bson:"closed_at"`
	CreatedAt             time.Time       `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Currency              string          `json:"currency" xml:"currency,attr" bson:"currency"`
	Customer              *Customer       `json:"customer" xml:"customer" bson:"customer"`
	DiscountCodes         []DiscountCode  `json:"discount_codes" xml:"discount_codes" bson:"discount_codes"`
	Email                 string          `json:"email" xml:"email,attr" bson:"email"`
	FinancialStatus       string          `json:"financial_status" xml:"financial_status,attr" bson:"financial_status"`
	Fulfillments          []Fulfillment   `json:"fulfillments" xml:"fulfillments" bson:"fulfillments"`
	FulfillmentStatus     string          `json:"fulfillment_status" xml:"fulfillment_status,attr" bson:"fulfillment_status"`
	Tags                  []string        `json:"tags" xml:"tags" bson:"tags"`
	LandingSite           *url.URL        `json:"landing_site" xml:"landing_site" bson:"landing_site"`
	LineItems             []LineItem      `json:"line_items" xml:"line_items" bson:"line_items"`
	Name                  string          `json:"name" xml:"name,attr" bson:"name"`
	Note                  string          `json:"note" xml:"note,attr" bson:"note"`
	NoteAttributes        []NoteAttribute `json:"note_attributes" xml:"note_attributes" bson:"note_attributes"`
	Number                int             `json:"number" xml:"number,attr" bson:"number"`
	OrderNumber           bson.ObjectId   `json:"order_number" xml:"order_number,attr" bson:"order_number"`
	ProcessedAt           time.Time       `json:"processed_at" xml:"processed_at,attr" bson:"processed_at"`
	ProcessingMethod      string          `json:"processing_method" xml:"processing_method,attr" bson:"processing_method"`
	ReferringSite         *url.URL        `json:"referring_site" xml:"referring_site" bson:"referring_site"`
	Refund                []Refund        `json:"refund" xml:"refund" bson:"refund"`
	ShippingAddress       CustomerAddress `json:"shipping_address" xml:"shipping_address" bson:"shipping_address"`
	ShippingLines         []ShippingLine  `json:"shipping_lines" xml:"shipping_lines" bson:"shipping_lines"`
	SourceName            string          `json:"source_name" xml:"source_name,attr" bson:"source_name"`
	SubtotalPrice         float64         `json:"subtotal_price" xml:"subtotal_price,attr" bson:"subtotal_price"`
	TaxLines              []TaxLine       `json:"tax_lines" xml:"tax_lines" bson:"tax_lines"`
	TaxesIncluded         bool            `json:"taxes_included" xml:"taxes_included,attr" bson:"taxes_inclued"`
	Token                 bson.ObjectId   `json:"token" xml:"token,attr" bson:"token"`
	TotalDiscounts        float64         `json:"total_discounts" xml:"total_discounts,attr" bson:"total_discounts"`
	TotalLineItemsPrice   float64         `json:"total_line_items_price" xml:"total_line_items_price,attr" bson:"total_line_items_price"`
	TotalPrice            float64         `json:"total_price" xml:"total_price,attr" bson:"total_price"`
	TotalTax              float64         `json:"total_tax" xml:"total_tax,attr" bson:"total_tax"`
	TotalWeight           float64         `json:"total_weight" xml:"total_weight,attr" bson:"total_weight"`
	UpdatedAt             time.Time       `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
}
