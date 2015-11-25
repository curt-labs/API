package cartIntegration

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	// "github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cartIntegration"
	"github.com/curt-labs/API/models/customer"
	"github.com/go-martini/martini"
)

func setDB(r *http.Request) error {
	brandID, err := strconv.Atoi(r.URL.Query().Get("brandID"))
	if err != nil {
		return err
	}
	cartIntegration.Brand_ID = brandID
	return nil
}

func setCustomerId(r *http.Request) error {
	c := customer.Customer{}
	err := c.GetCustomerIdFromKey(r.URL.Query().Get("key"))
	if err != nil {
		return err
	}
	cartIntegration.Customer_ID = c.Id
	return nil
}

// Requires APIKEY and brandID in header
func GetPricing(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	prices, err := cartIntegration.GetCustomerPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting prices by customer ID", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(prices))
}

// Requires APIKEY and brandID in header
// Requires count and page in params
func GetPricingPaged(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}

	page, err := strconv.Atoi(params["page"])
	if page < 1 || err != nil {
		apierror.GenerateError("Trouble getting page number for paged customer pricing", err, rw, r)
		return ""
	}

	count, err := strconv.Atoi(params["count"])
	if count < 1 || err != nil {
		apierror.GenerateError("Trouble getting count for paged customer pricing", err, rw, r)
		return ""
	}

	prices, err := cartIntegration.GetPricingPaged(page, count)
	if err != nil {
		apierror.GenerateError("Trouble getting prices for paged customer pricing", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(prices))
}

//Returns int
func GetPricingCount(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	count, err := cartIntegration.GetPricingCount()
	if err != nil {
		apierror.GenerateError("Trouble getting pricing count", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(count))
}

//Returns Mfr Prices for a part
func GetPartPricesByPartID(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	partID, err := strconv.Atoi(params["part"])
	if partID < 1 || err != nil {
		apierror.GenerateError("Trouble getting part number for part pricing", err, rw, r)
		return ""
	}
	prices, err := cartIntegration.GetPartPricesByPartID(partID)
	if err != nil {
		apierror.GenerateError("Trouble getting pricing", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(prices))
}

//Returns Mfr Prices
func GetAllPartPrices(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	prices, err := cartIntegration.GetPartPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting pricing", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(prices))
}

func CreatePrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble creating pricing", err, rw, r)
		return ""
	}
	var price cartIntegration.CustomerPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		apierror.GenerateError("Trouble creating pricing", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	price.CustID = cartIntegration.Customer_ID
	err = validatePrice(price)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r)
		return ""
	}
	err = price.Create()
	if err != nil {
		apierror.GenerateError("Trouble creating pricing", err, rw, r)
		return ""
	}
	err = price.InsertCartIntegration()
	if err != nil {
		apierror.GenerateError("Trouble creating CartIntegration", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(price))
}

func UpdatePrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble creating pricing", err, rw, r)
		return ""
	}
	var price cartIntegration.CustomerPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		apierror.GenerateError("Trouble creating pricing", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	price.CustID = cartIntegration.Customer_ID
	err = validatePrice(price)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r)
		return ""
	}
	err = price.Update()
	if err != nil {
		apierror.GenerateError("Trouble updating price", err, rw, r)
		return ""
	}

	err = price.UpdateCartIntegration()
	if err != nil {
		apierror.GenerateError("Trouble updating CartIntegration", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(price))
}

//set all of a customer's prices to MAP
func ResetAllToMap(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	custPrices, err := cartIntegration.GetCustomerPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting prices by customer ID", err, rw, r)
		return ""
	}

	//create map of MAP prices
	prices, err := cartIntegration.GetMAPPartPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting part prices", err, rw, r)
		return ""
	}
	priceMap := make(map[int]cartIntegration.Price)
	for _, p := range prices {
		priceMap[p.PartID] = p
	}

	//set to MAP
	for i, _ := range custPrices {
		custPrices[i].Price = priceMap[custPrices[i].PartID].Price
		if custPrices[i].CustID == 0 {
			custPrices[i].CustID = cartIntegration.Customer_ID
		}
		if custPrices[i].ID == 0 {
			err = custPrices[i].Create()
		} else {
			err = custPrices[i].Update()
		}
		if err != nil {
			apierror.GenerateError("Trouble updating price", err, rw, r)
			return ""
		}
	}
	return encoding.Must(enc.Encode(custPrices))
}

//sets all of a customer's prices to a percentage of the price type specified in params
func Global(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	err = setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	priceType := params["type"]
	percent, err := strconv.ParseFloat(params["percentage"], 64)
	if err != nil {
		apierror.GenerateError("Trouble parsing percentage", err, rw, r)
		return ""
	}
	percent = percent / 100

	//create partPriceMap
	prices, err := cartIntegration.GetPartPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting part prices", err, rw, r)
		return ""
	}
	priceMap := make(map[string]float64)
	for _, p := range prices {
		key := strconv.Itoa(p.PartID) + p.Type
		priceMap[key] = p.Price
	}

	//get CustPrices
	custPrices, err := cartIntegration.GetCustomerPrices()
	if err != nil {
		apierror.GenerateError("Trouble getting prices by customer ID", err, rw, r)
		return ""
	}

	//set to percentage
	for i, _ := range custPrices {
		if custPrices[i].CustID == 0 {
			custPrices[i].CustID = cartIntegration.Customer_ID
		}
		custPrices[i].Price = priceMap[strconv.Itoa(custPrices[i].PartID)+priceType] * percent
		if custPrices[i].ID == 0 {
			err = custPrices[i].Create()
		} else {
			err = custPrices[i].Update()

		}
		if err != nil {
			apierror.GenerateError("Trouble updating price", err, rw, r)
			return ""
		}
	}
	return encoding.Must(enc.Encode(custPrices))
}

//Get those price types
func GetAllPriceTypes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	types, err := cartIntegration.GetAllPriceTypes()
	if err != nil {
		apierror.GenerateError("Trouble getting price types", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(types))
}

//Utility
func validatePrice(p cartIntegration.CustomerPrice) error {
	if p.CustID < 1 {
		return errors.New("Customer ID cannot be less than 1")
	}
	if p.PartID < 1 {
		return errors.New("Part ID cannot be less than 1")
	}
	if p.IsSale == 1 {
		if p.SaleStart.Before(time.Now()) {
			return errors.New("The starting date is required and cannot be set to a date prior to today.")
		}

		if p.SaleStart.After(*p.SaleEnd) {
			return errors.New("The sale starting date cannot be set to a date after the sale ending date.")
		}

		if p.SaleEnd.Before(time.Now()) || p.SaleEnd.Before(*p.SaleStart) {
			return errors.New("The ending date is required and cannot be set to a date prior to today or the sale starting date.")
		}
	}
	return nil
}
