package cart

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Transaction struct {
	Id             bson.ObjectId  `json:"id" xml:"id" bson:"_id"`
	Authorization  string         `json:"authorization" xml:"authorization,attr" bson:"authorization"`
	CreatedAt      time.Time      `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	DeviceId       string         `json:"device_id" xml:"device_id,attr" bson:"device_id"`
	Gateway        string         `json:"gateway" xml:"gateway,attr" bson:"gateway"`
	SourceName     string         `json:"source_name" xml:"source_name,attr" bson:"source_name"`
	PaymentDetails PaymentDetails `json:"payment_details" xml:"payment_details" bson:"payment_details"`
	Kind           string         `json:"kind" xml:"kind,attr" bson:"kind"`
	OrderId        int            `json:"order_id" xml:"order_id,attr" bson:"order_id"`
	Receipt        Receipt        `json:"receipt" xml:"receipt" bson:"receipt"`
	ErrorCode      string         `json:"error_code" xml:"error_code,attr" bson:"error_code"`
	Status         string         `json:"status" xml:"status,attr" bson:"status"`
	Test           bool           `json:"test" xml:"test,attr" bson:"test"`
	UserId         bson.ObjectId  `json:"user_id" xml:"user_id,attr" bson:"user_id"`
	Currency       string         `json:"currency" xml:"currency,attr" bson:"currency"`
}

type PaymentDetails struct {
	AvsResultCode     string `json:"avs_result_code" xml:"avs_result_code,attr" bson:"avs_result_code"`
	CreditCardBin     string `json:"credit_card_bin" xml:"credit_card_bin,attr" bson:"credit_card_bin"`
	CvvResultCode     string `json:"cvv_result_cde" xml:"cvv_result_cde,attr" bson:"cvv_result_cde"`
	CreditCardNumber  string `json:"credit_card_number" xml:"credit_card_number,attr" bson:"credit_card_number"`
	CreditCardCompany string `json:"credit_card_company" xml:"credit_card_company,attr" bson:"credit_card_company"`
}
