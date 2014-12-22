package cart

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/url"
	"time"
)

type ClientDetails struct {
	AcceptLanguage string  `json:"accept_language" xml:"accept_language,attr" bson:"accept_language"`
	BrowserHeight  float64 `json:"browser_height" xml:"browser>height,attr" bson:"browser_height"`
	BrowserWidth   float64 `json:"browser_width" xml:"browser>width,attr" bson:"browser_width"`
	BrowserIP      string  `json:"browser_ip" xml:"browser>ip,attr" bson:"browser_ip"`
	SessionHash    string  `json:"session_hash" xml:"session_hash,attr" bson:"session_hash"`
	UserAgent      string  `json:"user_agent" xml:"user_agent,attr" bson:"user_agent"`
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
	Id                     bson.ObjectId    `json:"id" xml:"id,attr" bson:"_id"`
	ShopId                 bson.ObjectId    `json:"-" xml:"-" bson:"shop_id"`
	BillingAddress         *CustomerAddress `json:"billing_address" xml:"billing_address" bson:"billing_address"`
	BrowserIP              string           `json:"browser_ip" xml:"browser_ip,attr" bson:"browser_ip"`
	BuyerAcceptsMarketing  bool             `json:"buyer_accepts_marketing" xml:"buyer_accepts_marketing" bson:"buyer_accepts_marketing"`
	CancelReason           string           `json:"cancel_reason" xml:"cancel_reason" bson:"cancel_reason"`
	CancelledAt            time.Time        `json:"cancelled_at" xml:"cancelled_at" bson:"cancelled_at"`
	CartToken              *bson.ObjectId   `json:"cart_token" xml:"cart_token" bson:"cart_token"`
	ClientDetails          ClientDetails    `json:"client_details" xml:"client_details" bson:"client_details"`
	ClosedAt               time.Time        `json:"closed_at" xml:"closed_at,attr" bson:"closed_at"`
	CreatedAt              time.Time        `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Currency               string           `json:"currency" xml:"currency,attr" bson:"currency"`
	Customer               *Customer        `json:"customer" xml:"customer" bson:"customer"`
	DiscountCodes          *[]DiscountCode  `json:"discount_codes" xml:"discount_codes" bson:"discount_codes"`
	Email                  string           `json:"email" xml:"email,attr" bson:"email"`
	FinancialStatus        string           `json:"financial_status" xml:"financial_status,attr" bson:"financial_status"`
	Fulfillments           *[]Fulfillment   `json:"fulfillments" xml:"fulfillments" bson:"fulfillments"`
	FulfillmentStatus      string           `json:"fulfillment_status" xml:"fulfillment_status,attr" bson:"fulfillment_status"`
	Tags                   []string         `json:"tags" xml:"tags" bson:"tags"`
	LandingSite            *url.URL         `json:"landing_site" xml:"landing_site" bson:"landing_site"`
	LineItems              []LineItem       `json:"line_items" xml:"line_items" bson:"line_items"`
	Name                   string           `json:"name" xml:"name,attr" bson:"name"`
	Note                   string           `json:"note" xml:"note,attr" bson:"note"`
	NoteAttributes         []NoteAttribute  `json:"note_attributes" xml:"note_attributes" bson:"note_attributes"`
	Number                 int              `json:"number" xml:"number,attr" bson:"number"`
	OrderNumber            int              `json:"order_number" xml:"order_number,attr" bson:"order_number"`
	ProcessedAt            time.Time        `json:"processed_at" xml:"processed_at,attr" bson:"processed_at"`
	ProcessingMethod       string           `json:"processing_method" xml:"processing_method,attr" bson:"processing_method"`
	ReferringSite          *url.URL         `json:"referring_site" xml:"referring_site" bson:"referring_site"`
	Refund                 *[]Refund        `json:"refund" xml:"refund" bson:"refund"`
	ShippingAddress        *CustomerAddress `json:"shipping_address" xml:"shipping_address" bson:"shipping_address"`
	ShippingLines          []ShippingLine   `json:"shipping_lines" xml:"shipping_lines" bson:"shipping_lines"`
	SourceName             string           `json:"source_name" xml:"source_name,attr" bson:"source_name"`
	SubtotalPrice          float64          `json:"subtotal_price" xml:"subtotal_price,attr" bson:"subtotal_price"`
	TaxLines               []TaxLine        `json:"tax_lines" xml:"tax_lines" bson:"tax_lines"`
	TaxesIncluded          bool             `json:"taxes_included" xml:"taxes_included,attr" bson:"taxes_inclued"`
	Token                  bson.ObjectId    `json:"token" xml:"token,attr" bson:"token"`
	TotalDiscounts         float64          `json:"total_discounts" xml:"total_discounts,attr" bson:"total_discounts"`
	TotalLineItemsPrice    float64          `json:"total_line_items_price" xml:"total_line_items_price,attr" bson:"total_line_items_price"`
	TotalPrice             float64          `json:"total_price" xml:"total_price,attr" bson:"total_price"`
	TotalTax               float64          `json:"total_tax" xml:"total_tax,attr" bson:"total_tax"`
	TotalWeight            float64          `json:"total_weight" xml:"total_weight,attr" bson:"total_weight"`
	UpdatedAt              time.Time        `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
	SendWebhooks           bool             `json:"-" xml:"-" bson:"send_webhooks"`
	SendReceipt            bool             `json:"-" xml:"-" bson:"send_receipt"`
	SendFulfillmentReceipt bool             `json:"-" xml:"-" bson:"send_fulfillment_receipt"`
	InventoryBehavior      string           `json:"-" xml:"-" bson:"inventory_behavior"`
}

func (o *Order) Create() error {
	o.Token = bson.NewObjectId()
	o.CreatedAt = time.Now()
	o.ProcessedAt = time.Now()
	o.UpdatedAt = time.Now()

	if err := o.validate(); err != nil {
		return err
	}

	count, err := getOrderCount(o.ShopId)
	if err != nil {
		return err
	}
	o.OrderNumber = count + 1

	o.bindCustomer()

	if o.Id.Hex() == "" {
		o.Id = bson.NewObjectId()
	}

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("order")
	_, err = col.UpsertId(o.Id, o)

	return err
}

func (o *Order) Get() error {
	return fmt.Errorf("error: %s", "not implemented")
}

func (o *Order) Update() error {
	if o.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid order identifier")
	}

	o.UpdatedAt = time.Now()

	if err := o.validate(); err != nil {
		return err
	}

	o.bindCustomer()

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	updateDoc, err := o.mapUpdate()
	if err != nil {
		return err
	}

	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": updateDoc,
		},
	}

	_, err = sess.DB("CurtCart").C("order").Find(bson.M{"_id": o.Id, "shop_id": o.ShopId}).Apply(change, o)

	return err
}

func (o *Order) mapUpdate() (*map[string]interface{}, error) {
	tmp := Order{
		Id: o.Id,
	}

	if err := tmp.Get(); err != nil {
		return nil, err
	}

	doc := make(map[string]interface{})

	if !o.BillingAddress.deepEqual(tmp.BillingAddress) {
		doc["billing_address"] = o.BillingAddress
	}
	if !o.ShippingAddress.deepEqual(tmp.ShippingAddress) {
		doc["shipping_address"] = o.ShippingAddress
	}
	if o.BrowserIP != tmp.BrowserIP {
		doc["browser_ip"] = o.BrowserIP
	}
	if o.BuyerAcceptsMarketing != tmp.BuyerAcceptsMarketing {
		doc["buyer_accepts_marketing"] = o.BuyerAcceptsMarketing
	}
	if !o.ClientDetails.equal(&tmp.ClientDetails) {
		doc["client_details"] = o.ClientDetails
	}
	if o.Currency != tmp.Currency {
		doc["currency"] = o.Currency
	}

	// TODO - finish writing deep equal validation

	return &doc, nil
}

func (o *Order) bindCustomer() error {

	if (o.Customer == nil || !o.Customer.Id.Valid()) && o.Email != "" {
		c := Customer{
			ShopId: o.ShopId,
			Email:  o.Email,
		}

		if err := c.GetByEmail(); err != nil {
			return err
		}

		if c.Id.Valid() {
			o.Customer = &c
		}

		return nil
	}

	if o.Customer != nil && o.Customer.Id.Valid() {
		o.Customer.Get()
	}

	return nil
}

func getOrderCount(shopId bson.ObjectId) (int, error) {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return 0, err
	}
	defer sess.Close()

	return sess.DB("CurtCart").C("order").Find(bson.M{"shop_id": shopId}).Count()
}

func (o *Order) validate() error {
	if len(o.LineItems) == 0 {
		return fmt.Errorf("error: %s", "must have at least one line item")
	}

	for i, item := range o.LineItems {
		if item.VariantId == 0 {
			return fmt.Errorf("error: %s", "can't have missing variant details, missing variant_id")
		}
		if !item.Id.Valid() {
			o.LineItems[i].Id = bson.NewObjectId()
		}
	}

	if o.BillingAddress != nil && o.Email == "" {
		return fmt.Errorf("error: %s", "must define email address when billing address is provided")
	}

	return nil
}

func (a *ClientDetails) equal(b *ClientDetails) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if a.AcceptLanguage != b.AcceptLanguage {
		return false
	}
	if a.BrowserHeight != b.BrowserHeight {
		return false
	}
	if a.BrowserWidth != b.BrowserWidth {
		return false
	}
	if a.BrowserIP != b.BrowserIP {
		return false
	}
	if a.SessionHash != b.SessionHash {
		return false
	}
	return a.UserAgent == b.UserAgent
}
