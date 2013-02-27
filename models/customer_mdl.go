package models

import (
	"../helpers/database"
)

var (
	customerPriceStmt = `select distinct cp.price from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CustomerPricing cp on c.customerID = cp.cust_id
					where api_key = '%s'
					and cp.partID = %d`

	customerPartStmt = `select distinct ci.custPartID from ApiKey as ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join CartIntegration ci on c.customerID = ci.custID
					where ak.api_key = '%s'
					and ci.partID = %d`
)



func GetCustomerPrice(api_key string, part_id int) (price float64, err error) {
	db := database.Db

	row, _, err := db.QueryFirst(customerPriceStmt, api_key, part_id)
	if database.MysqlError(err){
		return
	}

	price = row.Float(0)
	return
}

func GetCustomerCartReference(api_key string, part_id int) (ref int, err error) {
	db := database.Db

	row, _, err := db.QueryFirst(customerPartStmt, api_key, part_id)
	if database.MysqlError(err){
		return
	}

	ref = row.Int(0)
	return
}
