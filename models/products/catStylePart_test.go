package products

import (
	"testing"

	"github.com/curt-labs/API/helpers/database"
)

func TestCategoryStyleParts(t *testing.T) {
	v := NoSqlVehicle{
		Year:  "2010",
		Make:  "Chevrolet",
		Model: "Silverado 1500",
	}

	if err := database.Init(); err != nil {
		t.Error(err)
	}

	session := database.ProductMongoSession

	csp, err := CategoryStyleParts(v, []int{3}, session)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(csp))
}
