package part_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"net/http"
	"strconv"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	key := params.Get("key")
	part := Part{
		PartId: id,
	}

	err := part.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part)
}

func GetRelated(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	key := params.Get("key")
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

	plate.ServeFormatted(w, r, parts)
}

func GetWithVehicle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	partID, err := strconv.Atoi(params.Get(":part"))
	if err != nil {
		http.Error(w, "Invalid part number", http.StatusInternalServerError)
		return
	}
	key := params.Get("key")
	year, err := strconv.ParseFloat(params.Get(":year"), 64)
	if err != nil {
		http.Redirect(w, r, "/part/"+params.Get(":part")+"?key="+key, http.StatusFound)
		return
	}
	vMake := params.Get(":make")
	model := params.Get(":model")
	submodel := params.Get(":submodel")
	config_vals := strings.Split(strings.TrimSpace(params.Get(":config")), "/")

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
		return
	}
	plate.ServeFormatted(w, r, part)

}

func Vehicles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":part"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}

	vehicles, err := ReverseLookup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, vehicles)
}

func Images(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Images)
}

func Attributes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetAttributes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Attributes)
}

func GetContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Content)
}

func Packaging(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetPartPackaging()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Packages)
}

func Reviews(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Reviews)
}

func Videos(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	err := part.GetVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, part.Videos)
}

func InstallSheet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(strings.Split(params.Get(":part"), ".")[0])
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
	plate.ServePdf(w, data)
}

func Categories(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	part := Part{
		PartId: id,
	}

	cats, err := part.GetPartCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, cats)
}

func Prices(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get(":part"))
	key := params.Get("key")
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
		return
	}

	plate.ServeFormatted(w, r, part.Pricing)
}
