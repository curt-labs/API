package part_ctlr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/rest"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/go-martini/martini"
	"github.com/ninnemana/analytics-go"
)

func track(endpoint string, params map[string]string, r *http.Request) {
	client := analytics.New("sud7rjoq3o")
	client.FlushAfter = 30 * time.Second
	client.FlushAt = 25

	js, err := json.Marshal(params)
	if err != nil {
		return
	}

	client.Track(map[string]interface{}{
		"title":    "Part Endpoint",
		"url":      r.URL.String(),
		"path":     r.URL.Path,
		"referrer": r.URL.RequestURI(),
		"params":   js,
	})
}

func All(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	page := 0
	count := 10
	qs := r.URL.Query()

	if qs.Get("page") != "" {
		if pg, err := strconv.Atoi(qs.Get("page")); err == nil {
			if pg == 0 {
				pg = 1
			}
			page = pg - 1
		}
	}
	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 500 {
				apierror.GenerateError(fmt.Sprintf("maximum request size is 500, you requested: %d", ct), err, w, r)
				return ""
			}
			count = ct
		}
	}

	parts, err := products.All(page, count, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all parts", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}

func Featured(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				apierror.GenerateError(fmt.Sprintf("maximum request size is 50, you requested: %d", ct), err, w, r)
				return ""
			}
			count = ct
		}
	}

	parts, err := products.Featured(count, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting featured parts", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}

func Latest(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				apierror.GenerateError(fmt.Sprintf("maximum request size is 50, you requested: %d", ct), err, w, r)
				return ""
			}
			count = ct
		}
	}

	parts, err := products.Latest(count, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting latest parts", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}

func Get(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {

		apierror.GenerateError("Trouble getting part", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p))
}

func GetRelated(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	parts, err := p.GetRelated(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting related parts", err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(parts))
}

func GetWithVehicle(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	err = errors.New("Not Implemented")
	apierror.GenerateError("Not Implemented", err, w, r)
	return ""
	// qs := r.URL.Query()
	// partID, err := strconv.Atoi(params["part"])
	// if err != nil {
	// 	http.Error(w, "Invalid part number", http.StatusInternalServerError)
	// 	return ""
	// }
	// key := qs.Get("key")
	// year, err := strconv.ParseFloat(params["year"], 64)
	// if err != nil {
	// 	http.Redirect(w, r, "/part/"+params["part"]+"?key="+key, http.StatusFound)
	// 	return ""
	// }
	// vMake := params["make"]
	// model := params["model"]
	// submodel := params["submodel"]
	// config_vals := strings.Split(strings.TrimSpace(params["config"]), "/")

	// vehicle := Vehicle{
	// 	Year:          year,
	// 	Make:          vMake,
	// 	Model:         model,
	// 	Submodel:      submodel,
	// 	Configuration: config_vals,
	// }

	// p := products.Part{
	// 	ID: partID,
	// }

	// err = products.GetWithVehicle(&vehicle, key)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return ""
	// }

	// return encoding.Must(enc.Encode(part))
}

func Vehicles(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	vehicles, err := vehicle.ReverseLookup(id)
	if err != nil {
		apierror.GenerateError("Trouble getting part vehicles", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vehicles))
}

//Redundant
func Images(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part images", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Images))
}

//Redundant
func Attributes(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part attributes", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Attributes))
}

//Redundant
func GetContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part", err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(p.Content))
}

//Redundant
func Packaging(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part packages", err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(p.Packages))
}

//Redundant
func ActiveApprovedReviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part reviews", err, w, r)
		return ""
	}
	var revs []products.Review
	for _, rev := range p.Reviews {
		if rev.Active == true && rev.Approved == true {
			revs = append(revs, rev)
		}
	}

	return encoding.Must(enc.Encode(revs))
}

func Videos(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	p := products.Part{
		ID: id,
	}

	err = p.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting part videos", err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(p.Videos))
}

//Sort of Redundant
func InstallSheet(w http.ResponseWriter, r *http.Request, params martini.Params, dtx *apicontext.DataContext) {
	id, err := strconv.Atoi(strings.Split(params["part"], ".")[0])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return
	}
	p := products.Part{
		ID: id,
	}

	err = p.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting part", err, w, r)
		return
	}
	var text string
	for _, content := range p.Content {
		if content.ContentType.Type == "installationSheet" {
			text = content.Text
		}
	}
	if text == "" {
		apierror.GenerateError("No Installation Sheet", err, w, r)
		return
	}

	data, err := rest.GetPDF(text, r)
	if err != nil {
		apierror.GenerateError("Error getting PDF", err, w, r)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Origin")
	w.Write(data)
}

//Redundant
func Categories(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	p := products.Part{
		ID: id,
	}

	err = p.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting part categories", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Categories))
}

func Prices(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	priceChan := make(chan int)
	custChan := make(chan int)

	go func() {
		err = p.Get(dtx)
		priceChan <- 1
	}()

	go func() {
		price, custErr := customer.GetCustomerPrice(dtx, p.ID)
		if custErr != nil {
			err = custErr
		}
		p.Pricing = append(p.Pricing, products.Price{0, 0, "Customer", price, false, time.Now()})
		custChan <- 1
	}()

	<-priceChan
	<-custChan

	if err != nil {
		apierror.GenerateError("Trouble getting part prices", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Pricing))
}

func PartNumber(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Part
	var err error

	p.PartNumber = params["part"]

	if p.PartNumber == "" {
		apierror.GenerateError("Trouble getting old part number", err, rw, r)
		return ""
	}

	if err = p.GetPartByPartNumber(); err != nil {
		apierror.GenerateError("Trouble getting part by old part number", err, rw, r)
		return ""
	}

	//TODO - remove when curt & aries vehicle application data are in sync
	if p.Brand.ID == 3 {
		mgoVehicles, err := vehicle.ReverseMongoLookup(p.ID)
		if err != nil {
			apierror.GenerateError("Trouble getting part by old part number", err, rw, r)
			return ""
		}
		for _, v := range mgoVehicles {
			vehicleApplication := products.VehicleApplication{
				Year:  v.Year,
				Make:  v.Make,
				Model: v.Model,
				Style: v.Style,
			}
			p.Vehicles = append(p.Vehicles, vehicleApplication)
		}
	} //END TODO

	return encoding.Must(enc.Encode(p))
}
