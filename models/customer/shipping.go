package customer

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerInfo struct {
	CustName string
}
type AccountStatus struct {
	AccountStatus string
}
type DefWH struct {
	Defwh string
}
type Threshold struct {
	FreeF float64
}
type ShippingInfo struct {
	CustomerInfo  CustomerInfo  `json:"customer_name"`
	AccountStatus AccountStatus `json:"account_status"`
	DefWH         DefWH         `json:"def_wh"`
	Threshold     Threshold     `json:"free_shipping_threshold"`
}

func (c *Customer) GetShippingInfo() (err error) {
	// http://146.148.89.23:8080/shipping is the proxy server
	u := "http://146.148.89.23:8080/shipping?id=" + strconv.Itoa(c.CustomerId)
	client := &http.Client{}
	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var s ShippingInfo
	err = dec.Decode(&s)
	if err != nil {
		return err
	}
	c.ShippingInfo = s

	return err
}
