package part_ctlr

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/part"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/go-martini/martini"
	"github.com/ninnemana/analytics-go"
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

	parts, err := part.All(key, page, count)
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
	part := part.Part{
		PartId: id,
	}

	err := part.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part))
}

func GetRelated(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	p := part.Part{
		PartId: id,
	}

	err := p.GetRelated()
	var parts []part.Part
	c := make(chan int, len(p.Related))
	for _, rel := range p.Related {
		go func(partId int) {
			relPart := part.Part{PartId: partId}
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

	// p := part.Part{
	// 	PartId: partID,
	// }

	// err = part.GetWithVehicle(&vehicle, key)
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
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
	}

	err := p.GetPartPackaging()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Packages))
}

func Reviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := part.Part{
		PartId: id,
	}

	err := p.GetReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(p.Reviews))
}

func Videos(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
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
	p := part.Part{
		PartId: id,
	}

	priceChan := make(chan int)
	custChan := make(chan int)

	var err error
	go func() {
		err = p.GetPricing()
		priceChan <- 1
	}()

	go func() {
		price, custErr := customer.GetCustomerPrice(key, p.PartId)
		if custErr != nil {
			err = custErr
		}

		p.Pricing = append(p.Pricing, part.Price{"Customer", price, false})
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
