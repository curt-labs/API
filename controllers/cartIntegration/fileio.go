package cartIntegration

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cartIntegration"
)

//TODO - extremely untested

func Upload(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	err := setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}
	key := r.URL.Query().Get("key")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		apierror.GenerateError("Error getting file from form", err, rw, r)
		return ""
	}

	if fileHeader != nil {
		contentType := fileHeader.Header.Get("Content-Type")

		if contentType != "text/comma-separated-values" && contentType != "text/csv" &&
			contentType != "application/csv" && contentType != "application/excel" &&
			contentType != "application/vnd.ms-excel" && contentType != "application/vnd.msexcel" {
			err = errors.New("The file you tried uploading was not a valid CSV file. Please try again using a valid CSV file.")
			apierror.GenerateError("Error uploading file", err, rw, r)
			return ""
		}
	}

	go cartIntegration.UploadFile(file, key)
	if err != nil {
		apierror.GenerateError("Error uploading file", err, rw, r)
		return ""
	}
	return ""
}

func Download(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	err := setCustomerId(r)
	if err != nil {
		apierror.GenerateError("Trouble getting customer from api key", err, rw, r)
		return ""
	}
	err = setDB(r)
	if err != nil {
		apierror.GenerateError("Trouble getting brandID from query string", err, rw, r)
		return ""
	}

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	customerPricesJson, err := cartIntegration.GetCustomerPrices(0, 0)
	if err != nil {
		apierror.GenerateError("Error getting customer prices ", err, rw, r)
		return ""
	}
	customerPrices := customerPricesJson.Items

	//Price map
	prices, err := cartIntegration.GetPartPrices()
	if err != nil {
		apierror.GenerateError("Error getting part prices ", err, rw, r)
		return ""
	}
	priceMap := make(map[string]float64)
	for _, p := range prices {
		priceMap[strconv.Itoa(p.PartID)+":"+p.Type] = p.Price
	}

	//Write
	wr.Write([]string{
		"CURT Part Number",
		"Customer Part ID",
		"Sale Price",
		"Sale Start Date",
		"Sale End Date",
		"Map Price",
		"List Price"})

	for _, price := range customerPrices {
		mapPrice := ""
		listPrice := ""

		mapPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":Map"], 'f', 2, 64)
		listPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":List"], 'f', 2, 64)

		//stringify dates
		var start, end string
		if price.SaleStart != nil && !price.SaleStart.IsZero() {
			start = price.SaleStart.Format(cartIntegration.DATE_FORMAT)
		}
		if price.SaleEnd != nil && !price.SaleStart.IsZero() {
			end = price.SaleEnd.Format(cartIntegration.DATE_FORMAT)
		}
		wr.Write([]string{
			price.PartNumber,
			strconv.Itoa(price.CustomerPartID), //TODO - get CartIntegration at the same time
			strconv.FormatFloat(price.Price, 'f', 2, 64),
			start,
			end,
			mapPrice,
			listPrice,
		})

	}

	wr.Flush()
	rw.Header().Set("Content-Type", "text/csv")
	rw.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	rw.Write(b.Bytes())

	return ""
}
