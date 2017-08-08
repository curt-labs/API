package customer

import (
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

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

var (
	warehouseMap           = `select id, code from Warehouses`
	updateFreightLimit     = `update Accounts set freightLimit = ? where id = ?`
	updateDefaultWarehouse = `update Accounts set defaultWarehouseId = ? where id = ?`
)

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

	for _, a := range c.Accounts {
		a.FreightLimit = c.ShippingInfo.Threshold.FreeF // Overwrite db value with shipping info resp
	}

	return err
}

//works like GetShippingInfo, but updates records
//TODO switch  user.getCustomer getShippingInfo over to use this
func (c *Customer) GetAndCompareCustomerShippingInfo() (err error) {
	//get mysql shipping info
	//get mapics shipping info
	//if mysql is different, update
	//return current-est
	var warehouseId int
	var ok bool

	warehouseMap, err := getWarehouseCodes()

	shipChan := make(chan error)
	accountChan := make(chan error)

	go func() {
		shipChan <- c.GetShippingInfo()
	}()
	go func() {
		accountChan <- c.GetAccounts()
	}()

	err = <-shipChan
	if err != nil {
		return err
	}

	err = <-accountChan
	if err != nil {
		return err
	}

	shippingInfo := c.ShippingInfo
	accounts := c.Accounts

	for i, account := range accounts {
		//adjust account freight
		if account.FreightLimit != shippingInfo.Threshold.FreeF && shippingInfo.Threshold.FreeF > 0 {
			err = account.adjustFreight(shippingInfo.Threshold.FreeF)
			if err != nil {
				return err
			}
			c.Accounts[i].FreightLimit = shippingInfo.Threshold.FreeF
		}

		if warehouseId, ok = warehouseMap[shippingInfo.DefWH.Defwh]; !ok {
			warehouseId = 0
		}
		//adjust account defwh
		if account.DefaultWarehouseID != warehouseId && warehouseId != 0 {
			err = account.adjustWarehouse(warehouseId)
			if err != nil {
				return err
			}
			c.Accounts[i].DefaultWarehouseID = warehouseId
		}
	}
	return err
}

func (a *Account) adjustFreight(limit float64) error {
	err := database.Init()
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(updateFreightLimit, limit, a.ID)
	return err
}

func (a *Account) adjustWarehouse(id int) error {
	err := database.Init()
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(updateDefaultWarehouse, id, a.ID)
	return err
}

//makes warehouse map
func getWarehouseCodes() (map[string]int, error) {
	warehouses := make(map[string]int)
	err := database.Init()
	if err != nil {
		return warehouses, err
	}

	rows, err := database.DB.Query(warehouseMap)
	if err != nil {
		return warehouses, err
	}
	defer rows.Close()

	var id *int
	var code *string
	for rows.Next() {
		err = rows.Scan(&id, &code)
		if err != nil {
			return warehouses, err
		}
		if id != nil && code != nil {
			warehouses[*code] = *id
		}
	}
	return warehouses, nil
}
