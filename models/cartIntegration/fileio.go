package cartIntegration

import (
	"github.com/curt-labs/API/helpers/database"
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
	// only add if it hasnt already been added

	// clean up excel file
	added := make(map[string]bool)
	var finalLines [][]string
	for _, line := range lines {
		partNumber := strings.TrimSpace(line[0])
		if !added[partNumber] { // if it has not been added to final lines, then add it
			added[partNumber] = true              // mark that it has been added
			finalLines = append(finalLines, line) // actually add it to the list of finalLines
		}
	}

	priceLookupJson, err := GetCustomerPrices(0, 0)
	if err != nil {
		return err
	}
	priceLookup := priceLookupJson.Items
	integrationLookup, err := GetCustomerCartIntegrations(api_key)
	if err != nil {
		return err
	}
	partmap, err := getPartMap()
	if err != nil {
		return err
	}

	for _, line := range finalLines {
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
		// value passed in could be 0, and there could be a 0 in the DB, so that would return an existing integration
		// value passed in could be 0, and no value could exist in the DB, thats why there is a bool and a int passed back.
		custPartNum, integExists := cp.integrationExists(integrationLookup)

		if custPartNum != cp.CustomerPartID { // if value from csv file does not match the found integration value from the DB
			if integExists { // if there is a found integration, then update the existing one
				err = cp.UpdateCartIntegration()
			} else { // if there is no curent integration then insert one
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
func (c *CustomerPrice) integrationExists(integrationLookup []CustomerPrice) (int, bool) {
	for _, l := range integrationLookup {
		if l.CustID == c.CustID && l.PartID == c.PartID {
			return l.CustomerPartID, true
		}
	}
	return 0, false
}

//getPartMap returns a map of partnumbers to partIds
func getPartMap() (map[string]int, error) {
	partmap := make(map[string]int)
	err := database.Init()
	if err != nil {
		return partmap, err
	}

	stmt, err := database.DB.Prepare(`select partId, oldPartNumber from Part`)
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
