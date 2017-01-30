package cartIntegration

import (
	_ "github.com/go-sql-driver/mysql"

	"encoding/csv"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02"
)

func UploadFile(file multipart.File, api_key string) error {
	defer file.Close()
	csvfile := csv.NewReader(file)

	lines, err := csvfile.ReadAll()
	if err != nil {
		return err
	}

	priceLookup, err := GetCustomerPrices()
	if err != nil {
		return err
	}
	integrationLookup, err := GetCustomerCartIntegrations(api_key)
	if err != nil {
		return err
	}
	partmap, err := getPartMap()
	if err != nil {
		return err
	}
	for _, line := range lines {
		//Curt Part ID,	Customer Part ID, Sale Price, Sale Start Date, Sale End Date
		var cp CustomerPrice
		cp.CustID = Customer_ID

		//partnumber to id
		partNumber := strings.TrimSpace(line[0])
		var id int
		var ok bool
		if id, ok = partmap[partNumber]; !ok {
			continue
		}
		cp.PartID = id

		customerPartID, err := strconv.Atoi(line[1])
		if err != nil {
			return err
		}
		cp.CustomerPartID = customerPartID

		strPrice := strings.TrimSpace(strings.Replace(line[2], "$", "", -1))
		price, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return err
		}
		cp.Price = price
		if len(line) > 3 {
			if line[3] != "" {
				startd, err := time.Parse(DATE_FORMAT, line[3])
				if err != nil {
					return err
				}
				cp.SaleStart = &startd
			}
			if line[4] != "" {
				endd, err := time.Parse(DATE_FORMAT, line[4])
				if err != nil {
					return err
				}
				cp.SaleEnd = &endd
			}
		}

		cp.priceExists(priceLookup) //determine if update or create

		if cp.ID > 0 {
			err = cp.Update()
		} else {
			err = cp.Create()
		}
		if err != nil {
			return err
		}

		custPartNum := cp.integrationExists(integrationLookup)
		if custPartNum != cp.CustomerPartID {
			if custPartNum > 0 {
				err = cp.UpdateCartIntegration()
			} else {
				err = cp.InsertCartIntegration()
			}
			if err != nil {
				return err
			}
		}
	}
	return err
}

//Checks if customerPrice exists in db
func (c *CustomerPrice) priceExists(priceLookup []CustomerPrice) {
	for _, l := range priceLookup {
		if l.CustID == c.CustID && l.PartID == c.PartID {
			c.ID = l.ID
		}
	}
	return
}

//Checks if integration exists in db
func (c *CustomerPrice) integrationExists(integrationLookup []CustomerPrice) int {
	for _, l := range integrationLookup {
		if l.CustID == c.CustID && l.PartID == c.PartID {
			return l.CustomerPartID
		}
	}
	return 0
}

//getPartMap returns a map of partnumbers to partIds
func getPartMap() (map[string]int, error) {
	partmap := make(map[string]int)
	db, err := initDB()
	if err != nil {
		return partmap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`select partId, oldPartNumber from Part`)
	if err != nil {
		return partmap, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return partmap, err
	}
	defer rows.Close()

	for rows.Next() {
		var num *string
		var id *int
		err = rows.Scan(&id, &num)
		if err != nil {
			return partmap, err
		}
		if id != nil && num != nil {
			partmap[*num] = *id
		}
	}
	return partmap, nil
}
