package mocks

import (
	"github.com/go-martini/martini"
	"net/http"
	"github.com/curt-labs/API/helpers/apicontext"
)

// Mock for Meddler (@see controllers/middlware.go). Useful for short circuiting authentication while still providing
// The needed DataContext to the Martini Handler
func Meddler(dataContext apicontext.DataContext) martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		c.Map(&dataContext)
		c.Next()
	}
}
