package geography

import (
	"net/http"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/geography"
)

func GetAllCountriesAndStates(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	cstates, err := geography.GetAllCountriesAndStates()
	if err != nil {
		apierror.GenerateError("Trouble getting all countries and states", err, rw, req)
	}
	return encoding.Must(enc.Encode(cstates))
}

func GetAllCountries(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	countries, err := geography.GetAllCountries()
	if err != nil {
		apierror.GenerateError("Trouble getting all countries", err, rw, req)
	}
	return encoding.Must(enc.Encode(countries))
}

func GetAllStates(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	states, err := geography.GetAllStates()
	if err != nil {
		apierror.GenerateError("Trouble getting all states", err, rw, req)
	}
	return encoding.Must(enc.Encode(states))
}
