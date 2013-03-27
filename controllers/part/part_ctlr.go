package part_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"net/http"
	"strconv"
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
