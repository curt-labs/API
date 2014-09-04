package apifilter

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"log"
	"sort"

	// "math"
	"strconv"
	"strings"
	"time"

	// "github.com/curt-labs/GoAPI/models/category"
	"github.com/curt-labs/GoAPI/models/part"
)

var (
	ExcludedAttributes = []string{"UPC"}
)

type FilteredOptions []Options

type Options struct {
	Key     string
	Options []Option
}

type Option struct {
	Value    string
	Selected bool
}

type Decision struct {
	Field  string
	Values map[string]interface{}
}

func PartFilter(parts []part.Part, specs []interface{}) ([]Options, error) {

	var filtered FilteredOptions

	attrChan := make(chan error)
	priceChan := make(chan error)
	go func() {
		filtered = append(filtered, filtered.attributes(parts)...)
		attrChan <- nil
	}()
	go func() {
		filtered = append(filtered, filtered.prices(parts))
		priceChan <- nil
	}()

	select {
	case err := <-attrChan:
		if err != nil {
			log.Println("filter error: %s", err.Error())
		}
	case err := <-priceChan:
		if err != nil {
			log.Println("filter error: %s", err.Error())
		}
	case <-time.After(1 * time.Second):
		log.Println("filter attributes timed out")
	}

	sortutil.AscByField(filtered, "Key")

	return filtered, nil
}

func (filtered FilteredOptions) attributes(parts []part.Part) FilteredOptions {
	attributeDefinitions := make(map[string]Options, 0)
	for _, part := range parts {
		for _, attr := range part.Attributes {

			// Check Excluded attributes
			exclude := false
			for _, ex := range ExcludedAttributes {
				if ex == attr.Key {
					exclude = true
				}
			}
			if exclude {
				continue
			}

			vals, ok := attributeDefinitions[attr.Key]
			if !ok {
				vals = Options{
					Key: attr.Key,
				}
			}

			exists := false
			for _, val := range vals.Options {
				if vals.Key == attr.Key && val.Value == attr.Value {
					exists = true
				}
			}

			if !exists {
				newOption := Option{
					Value:    attr.Value,
					Selected: false,
				}
				vals.Options = append(vals.Options, newOption)
				attributeDefinitions[attr.Key] = vals
			}
		}
	}

	var f FilteredOptions
	for _, vals := range attributeDefinitions {
		if len(vals.Options) > 1 {
			f = append(f, vals)
		}
	}
	return f
}

func (filtered FilteredOptions) prices(parts []part.Part) Options {
	priceDefinitions := Options{
		Key:     "Price",
		Options: make([]Option, 0),
	}

	lows := make([]int, 0)
	for _, p := range parts {

		// get list price
		var list float64
		for _, pr := range p.Pricing {
			if pr.Type == "List" {
				list = pr.Price
				break
			}
		}

		exists := false
		for _, def := range priceDefinitions.Options {
			val := strings.Replace(def.Value, "$", "", -1)
			segs := strings.Split(val, " - ")
			if len(segs) < 2 {
				continue
			}

			low, err := strconv.ParseFloat(segs[0], 64)
			if err != nil {
				continue
			}
			high, err := strconv.ParseFloat(segs[1], 64)
			if err != nil {
				continue
			}

			if list >= low && list <= high {
				exists = true
			}
		}

		if !exists {
			lows = append(lows, (int(list)/50)*50)
		}
	}

	sort.Ints(lows)
	existing := make(map[int]int, 0)
	for _, low := range lows {
		if _, ok := existing[low]; !ok {
			val := fmt.Sprintf("$%d - $%d", low, low+50)
			opt := Option{
				Value:    val,
				Selected: false,
			}
			priceDefinitions.Options = append(priceDefinitions.Options, opt)
			existing[low] = low
		}
	}
	// sortutil.AscByField(priceDefinitions.Options, "Value")

	return priceDefinitions
}
