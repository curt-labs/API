package customer_ctlr_new

import (
	// "encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new/content"
	// "github.com/go-martini/martini"
	// "io/ioutil"
	"net/http"
	// "strconv"
)

// Get it all
func GetAllContent(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := custcontent.AllCustomerContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	content, err := custcontent.GetCustomerContent(id, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}
