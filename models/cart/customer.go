package cart

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Customer struct {
	Id               bson.ObjectId     `json:"id" xml:"id" bson:"_id"`
	AcceptsMarketing bool              `json:"accepts_marketing" xml:"accepts_marketing,attr" bson:"accepts_marketing"`
	Addresses        []CustomerAddress `json:"addresses" xml:"addresses>addres" bson:"addresses"`
	DefaultAddress   []CustomerAddress `json:"default_address" xml:"default_address" bson:"default_address"`
	CreatedAt        time.Time         `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Email            string            `json:"email" xml:"email,attr" bson:"email"`
	FirstName        string            `json:"first_name" xml:"first_name,attr" bson:"first_name"`
	LastName         string            `json:"last_name" xml:"last_name,attr" bson:"last_name"`
	MetaFields       []MetaField       `json:"meta_fields" xml:"meta_fields>meta_field" bson:"meta_fields"`
	LastOrderId      *bson.ObjectId    `json:"last_order_id" xml:"last_order_id,attr" bson:"last_order_id"`
	LastOrderName    string            `json:"last_order_name,omitempty" xml:"last_order_name,attr,omitempty" bson:"last_order_name"`
	OrdersCount      int               `json:"orders_count" xml:"orders_count" bson:"orders_count"`
	Note             string            `json:"note" xml:"note" bson:"note"`
	State            string            `json:"state" xml:"state,attr" bson:"state"` // Disabled, invited, etc.
	Tags             []string          `json:"tags" xml:"tags>tag" bson:"tags"`
	TotalSpent       float64           `json:"total_spent" xml:"total_spent,attr" bson:"total_spent"`
	UpdatedAt        time.Time         `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
	VerifiedEmail    bool              `json:"verified_email" xml:"verified_email,attr" bson:"verified_email"`
	Get              func() error
	GetAddresses     func() error
}

type MetaField struct {
	Key       string `json:"key" xml:"key,attr" bson:"key"`
	Namespace string `json:"namespace" xml:"namespace,attr" bson:"namespace"`
	ValueType string `json:"value_type" xml:"value_type,attr" bson:"value_type"`
	Value     string `json:"value" xml:"value,attr" bson:"value"`
}

func SinceId(id string) ([]Customer, error) {
	custs := []Customer{}
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return custs, err
	}
	defer sess.Close()

	c := sess.DB("CurtCart").C("customer")
	qs := bson.M{
		"_id": bson.M{
			"$gt": id,
		},
	}
	c.Find(qs)

	return custs, err
}

func (c *Customer) Get() error {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	return col.Find(bson.M{"_id": c.Id}).One(&c)
}
