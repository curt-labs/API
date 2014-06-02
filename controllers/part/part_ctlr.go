package part_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	. "github.com/curt-labs/GoAPI/models"
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

func Get(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	part := Part{
		PartId: id,
	}

	err := part.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	track("/part/get", params, r)

	return encoding.Must(enc.Encode(part))
}

func GetRelated(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	id, _ := strconv.Atoi(params["part"])
	key := qs.Get("key")
	part := Part{
		PartId: id,
	}

	err := part.GetRelated()
	var parts []Part
	c := make(chan int, len(part.Related))
	for _, p := range part.Related {
		go func(partId int) {
			relPart := Part{PartId: partId}
			if err = relPart.Get(key); err == nil {
				parts = append(parts, relPart)
			}
			c <- 1
		}(p)
	}

	for _, _ = range part.Related {
		<-c
	}

	return encoding.Must(enc.Encode(parts))
}

func GetWithVehicle(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	partID, err := strconv.Atoi(params["part"])
	if err != nil {
		http.Error(w, "Invalid part number", http.StatusInternalServerError)
		return ""
	}
	key := qs.Get("key")
	year, err := strconv.ParseFloat(params["year"], 64)
	if err != nil {
		http.Redirect(w, r, "/part/"+params["part"]+"?key="+key, http.StatusFound)
		return ""
	}
	vMake := params["make"]
	model := params["model"]
	submodel := params["submodel"]
	config_vals := strings.Split(strings.TrimSpace(params["config"]), "/")

	vehicle := Vehicle{
		Year:          year,
		Make:          vMake,
		Model:         model,
		Submodel:      submodel,
		Configuration: config_vals,
	}

	part := Part{
		PartId: partID,
	}

	err = part.GetWithVehicle(&vehicle, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part))
}

func Vehicles(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, err := strconv.Atoi(params["part"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}

	vehicles, err := ReverseLookup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(vehicles))
}

func Images(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Images))
}

func Attributes(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {

	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetAttributes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Attributes))
}

func GetContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Content))
}

func Packaging(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetPartPackaging()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Packages))
}

func Reviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Reviews))
}

func Videos(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, _ := strconv.Atoi(params["part"])
	part := Part{
		PartId: id,
	}

	err := part.GetVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Videos))
}

func InstallSheet(w http.ResponseWriter, r *http.Request, params martini.Params) {
	id, _ := strconv.Atoi(strings.Split(params["part"], ".")[0])
	part := Part{
		PartId: id,
	}

	data, err := part.GetInstallSheet(r)
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
	part := Part{
		PartId: id,
	}

	cats, err := part.GetPartCategories(key)
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
	part := Part{
		PartId: id,
	}

	priceChan := make(chan int)
	custChan := make(chan int)

	var err error
	go func() {
		err = part.GetPricing()
		priceChan <- 1
	}()

	go func() {
		price, custErr := GetCustomerPrice(key, part.PartId)
		if custErr != nil {
			err = custErr
		}
		part.Pricing = append(part.Pricing, Pricing{"Customer", price, false})
		custChan <- 1
	}()

	<-priceChan
	<-custChan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(part.Pricing))
}
