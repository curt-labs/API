package cartIntegration

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
)

//TODO - extremely untested

func Upload(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		apierror.GenerateError("Error uploading file", err, rw, r)
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

	err = cartIntegration.UploadFile(file, dtx)
	return ""
}

func Download(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	customerPrices, err := cartIntegration.GetCustomerPrices(dtx)
	if err != nil {
		apierror.GenerateError("Error getting customer prices ", err, rw, r)
		return ""
	}

	//Price map
	prices, err := cartIntegration.GetPartPrices(dtx)
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
		"CURT Part ID",
		"Customer Part ID",
		"Sale Price",
		"Sale Start Date",
		"Sale End Date",
		"Map Price",
		"List Price"})

	for _, price := range customerPrices {
		mapPrice := ""
		listPrice := ""

		mapPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":Map"], 'b', 2, 64)
		listPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":List"], 'b', 2, 64)

		wr.Write([]string{
			strconv.Itoa(price.PartID),
			strconv.Itoa(price.CustomerPartID), //TODO
			strconv.FormatFloat(price.Price, 'b', 2, 64),
			price.SaleStart.String(),
			price.SaleEnd.String(),
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
