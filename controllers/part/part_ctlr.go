package part_ctlr

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/go-martini/martini"
	"github.com/ninnemana/analytics-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func track(endpoint string, params map[string]string, r *http.Request) {
	client := analytics.New("sud7rjoq3o")
	client.FlushAfter = 30 * time.Second
	client.FlushAt = 25

	js, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
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

func All(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {

	page := 0
	count := 10
	qs := r.URL.Query()
	key := qs.Get("key")
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
				http.Error(w, fmt.Sprintf("maximum request size is 50, you requested: %d", ct), http.StatusInternalServerError)
				return ""
			}
			count = ct
		}
	}

	parts, err := products.All(key, page, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}

func Get(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
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

	err := p.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	<-vehicleChan

	return encoding.Must(enc.Encode(p))
}

func GetRelated(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	p := products.Part{
		ID: id,
	}

	err := p.GetRelated()
	var parts []products.Part
	c := make(chan int, len(p.Related))
	for _, rel := range p.Related {
		go func(partId int) {
			relPart := products.Part{ID: partId}
			if err = relPart.Get(key); err == nil {
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
	http.Error(w, "NotImplemented", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusFound)
	}

	vehicles, err := vehicle.ReverseLookup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(vehicles))
}

func Images(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Images))
}

func Attributes(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {

	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetAttributes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Attributes))
}

func GetContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Content))
}

func Packaging(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetPartPackaging()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Packages))
}

func ActiveApprovedReviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetActiveApprovedReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Reviews))
}

func Videos(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := products.Part{
		ID: id,
	}

	err := p.GetVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Videos))
}

func InstallSheet(w http.ResponseWriter, r *http.Request, params martini.Params) {
	id, _ := strconv.Atoi(strings.Split(params["part"], ".")[0])
	p := products.Part{
		ID: id,
	}

	data, err := p.GetInstallSheet(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(data) == 0 {
		http.Error(w, "No Installation Sheet found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Origin")
	w.Write(data)
}

func Categories(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	p := products.Part{
		ID: id,
	}

	cats, err := p.GetPartCategories(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
}

func Prices(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	p := products.Part{
		ID: id,
	}

	priceChan := make(chan int)
	custChan := make(chan int)

	var err error
	go func() {
		err = p.GetPricing()
		priceChan <- 1
	}()

	go func() {
		price, custErr := customer.GetCustomerPrice(key, p.ID)
		if custErr != nil {
			err = custErr
		}

		p.Pricing = append(p.Pricing, products.Price{0, 0, "Customer", price, false, time.Now()})
		custChan <- 1
	}()

	<-priceChan
	<-custChan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Pricing))
}

func SavePrice(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var p products.Price
	var err error
	idStr := params["id"]
	if idStr != "" {
		p.Id, err = strconv.Atoi(idStr)
		err = p.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	log.Print(p)
	//create or update
	if p.Id > 0 {
		err = p.Update()
	} else {
		err = p.Create()
	}
	log.Print(p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func DeletePrice(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var p products.Price
	var err error
	idStr := params["id"]

	p.Id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = p.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}

func GetPrice(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var p products.Price
	var err error
	idStr := params["id"]

	p.Id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = p.Get()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(p))
}
