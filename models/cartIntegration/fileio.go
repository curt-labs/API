package cartIntegration

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
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

func UploadFile(file multipart.File, dtx *apicontext.DataContext) error {
	defer file.Close()
	csvfile := csv.NewReader(file)
	lines, err := csvfile.ReadAll()

	if err != nil {
		return err
	}

	priceLookup, err := GetCustomerPrices(dtx)
	if err != nil {
		return err
	}
	integrationLookup, err := GetCustomerCartIntegrations(dtx)
	if err != nil {
		return err
	}

	for i, line := range lines {
		//Curt Part ID,	Customer Part ID, Sale Price, Sale Start Date, Sale End Date
		var cp CustomerPrice
		cp.CustID = dtx.CustomerID

		partID, err := strconv.Atoi(line[0])
		if err != nil && i == 0 {
			continue //checks for line headers
		}
		if err != nil {
			return err
		}
		cp.PartID = partID

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
				*cp.SaleStart = startd
			}
			if line[4] != "" {
				endd, err := time.Parse(DATE_FORMAT, line[4])
				if err != nil {
					return err
				}
				*cp.SaleEnd = endd
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
