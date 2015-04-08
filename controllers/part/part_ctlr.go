package part_ctlr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/curt-labs/GoAPI/models/video"
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
			if ct > 50 {
				apierror.GenerateError(fmt.Sprintf("maximum request size is 50, you requested: %d", ct), err, w, r)
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

	vehicleChan := make(chan error)
	go func() {
		vs, err := vehicle.ReverseLookup(p.ID)
		if err != nil {
			vehicleChan <- err
		} else {
			p.Vehicles = vs

			vehicleChan <- nil
		}
	}()

	if err = p.Get(dtx); err != nil {

		apierror.GenerateError("Trouble getting part", err, w, r)
		return ""
	}

	<-vehicleChan

	close(vehicleChan)

	sortutil.AscByField(p.Vehicles, "ID")

	return encoding.Must(enc.Encode(p))
}

func GetRelated(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetRelated(dtx)
	var parts []products.Part
	c := make(chan int, len(p.Related))
	for _, rel := range p.Related {
		go func(partId int) {
			relPart := products.Part{ID: partId}
			if err = relPart.Get(dtx); err == nil {
				parts = append(parts, relPart)
			}
			c <- 1
		}(rel)
	}

	for _, _ = range p.Related {
		<-c
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

func Images(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.GetImages(dtx); err != nil {
		apierror.GenerateError("Trouble getting part images", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Images))
}

func Attributes(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.GetAttributes(dtx); err != nil {
		apierror.GenerateError("Trouble getting part attributes", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Attributes))
}

func GetContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.GetContent(dtx); err != nil {
		apierror.GenerateError("Trouble getting part content", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Content))
}

func Packaging(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.GetPartPackaging(dtx); err != nil {
		apierror.GenerateError("Trouble getting part packages", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Packages))
}

func ActiveApprovedReviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}
	p := products.Part{
		ID: id,
	}

	if err = p.GetActiveApprovedReviews(dtx); err != nil {
		apierror.GenerateError("Trouble getting part reviews", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(p.Reviews))
}

func Videos(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	p := products.Part{
		ID: id,
	}

	var vs video.Videos
	if vs, err = video.GetPartVideos(p.ID); err != nil {
		apierror.GenerateError("Trouble getting part videos", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(vs))
}

func InstallSheet(w http.ResponseWriter, r *http.Request, params martini.Params, dtx *apicontext.DataContext) {
	id, err := strconv.Atoi(strings.Split(params["part"], ".")[0])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return
	}
	p := products.Part{
		ID: id,
	}

	data, err := p.GetInstallSheet(r, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting part installsheets", err, w, r)
		return
	} else if len(data) == 0 {
		apierror.GenerateError("No Installation Sheet found", err, w, r)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Origin")
	w.Write(data)
}

func Categories(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, w, r)
		return ""
	}

	p := products.Part{
		ID: id,
	}

	cats, err := p.GetPartCategories(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting part categories", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
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
		err = p.GetPricing(dtx)
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

func SavePrice(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Price
	var err error
	idStr := params["id"]
	if idStr != "" {
		p.Id, err = strconv.Atoi(idStr)
		err = p.Get()
		if err != nil {
			apierror.GenerateError("Trouble getting part ID", err, rw, r)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while saving part price", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &p)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling josn request body while saving part price", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if p.Id > 0 {
		err = p.Update(dtx)
	} else {
		err = p.Create(dtx)
	}
	if err != nil {
		msg := "Trouble while creating part price"
		if p.Id > 0 {
			msg = "Trouble while updating part price"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func DeletePrice(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Price
	var err error
	idStr := params["id"]

	if p.Id, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting price ID", err, rw, r)
		return ""
	}

	if err = p.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting price", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func GetPrice(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	var p products.Price
	var err error
	idStr := params["id"]

	if p.Id, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting price ID", err, rw, r)
		return ""
	}

	if err = p.Get(); err != nil {
		apierror.GenerateError("Trouble getting price", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func OldPartNumber(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var pa products.Part
	var err error

	pa.OldPartNumber = params["part"]

	if pa.OldPartNumber == "" {
		apierror.GenerateError("Trouble getting old part number", err, rw, r)
		return ""
	}

	if err = pa.GetPartByOldPartNumber(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting part by old part number", err, rw, r)
		return ""
	}

	// have to create a new object otherwise the properties of the old object(status) will intefere with getting the full part
	p := products.Part{
		ID: pa.ID,
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part by old part number", err, rw, r)
		return ""
	}

	vehicleChan := make(chan error)
	go func() {
		vs, err := vehicle.ReverseLookup(p.ID)
		if err != nil {
			vehicleChan <- err
		} else {
			p.Vehicles = vs
			vehicleChan <- nil
		}
	}()

	<-vehicleChan

	return encoding.Must(enc.Encode(p))
}

func CreatePart(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Part
	var err error

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while creating part", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &p); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while creating part", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = p.Create(dtx); err != nil {
		apierror.GenerateError("Trouble creating part", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func UpdatePart(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Part
	var err error

	idStr := params["id"]
	if idStr == "" {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	if p.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	if err = p.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting part", err, rw, r)
		return ""
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while updating part", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = json.Unmarshal(requestBody, &p); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while updating part", err, rw, r)
		return encoding.Must(enc.Encode(false))
	}

	if err = p.Update(dtx); err != nil {
		apierror.GenerateError("Trouble updating part", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func DeletePart(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var p products.Part
	var err error
	idStr := params["id"]

	if p.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	if err = p.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting part", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}
