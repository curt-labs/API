package cart

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Refund struct {
	Id              bson.ObjectId `json:"id" xml:"id,attr" bson:"_id"`
	CreatedAt       time.Time     `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Note            string        `json:"note" xml:"note,attr" bson:"note"`
	RefundLineItems []LineItem    `json:"refund_line_items" xml:"refund_line_items" bson:"refund_line_items"`
	Restock         bool          `json:"restock" xml:"restock,attr" bson:"restock"`
	Transactions    []interface{} `json:"transactions" xml:"transactions" bson:"transactions"`
	UserId          bson.ObjectId `json:"user_id" xml:"user_id,attr" bson:"user_id"`
}
