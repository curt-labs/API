package apifilter

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/products"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ExcludedPartAttributes = []string{"UPC"}
)

func PartFilter(parts []products.Part, specs []interface{}) ([]Options, error) {

	var filtered FilteredOptions

	attrChan := make(chan error)
	priceChan := make(chan error)
	catChan := make(chan error)
	classChan := make(chan error)
	go func() {
		filtered = append(filtered, filtered.partAttributes(parts)...)
		attrChan <- nil
	}()
	go func() {
		filtered = append(filtered, filtered.partPrices(parts))
		priceChan <- nil
	}()
	go func() {
		filtered = append(filtered, filtered.partCategory(parts))
		catChan <- nil
	}()
	go func() {
		filtered = append(filtered, filtered.partClass(parts))
		classChan <- nil
	}()

	select {
	case err := <-attrChan:
		if err != nil {
			log.Printf("filter error: %s\n", err.Error())
		}
	case err := <-priceChan:
		if err != nil {
			log.Printf("filter error: %s\n", err.Error())
		}
	case err := <-catChan:
		if err != nil {
			log.Printf("filter error: %s\n", err.Error())
		}
	case err := <-classChan:
		if err != nil {
			log.Printf("filter error: %s\n", err.Error())
		}
	case <-time.After(1 * time.Second):
		log.Println("filter generation timed out")
	}

	sortutil.AscByField(filtered, "Key")

	return filtered, nil
}

func (filtered FilteredOptions) partAttributes(parts []products.Part) FilteredOptions {
	attributeDefinitions := make(map[string]Options, 0)
	for _, part := range parts {
		for _, attr := range part.Attributes {

			// Check Excluded attributes
			exclude := false
			for _, ex := range ExcludedPartAttributes {
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
					break
				}
			}

			if !exists {
				newOption := Option{
					Value:    attr.Value,
					Selected: false,
					Products: []int{part.PartId},
				}
				vals.Options = append(vals.Options, newOption)
				attributeDefinitions[attr.Key] = vals
			} else {
				for i, opt := range attributeDefinitions[attr.Key].Options {
					if opt.Value == attr.Value {
						prods := attributeDefinitions[attr.Key].Options[i].Products
						prods = append(prods, part.PartId)
						sort.Ints(prods)
						attributeDefinitions[attr.Key].Options[i].Products = prods
					}
				}
			}
		}
	}

	var f FilteredOptions
	for _, vals := range attributeDefinitions {
		if len(vals.Options) > 1 {
			sortutil.AscByField(vals.Options, "Value")
			f = append(f, vals)
		}
	}
	return f
}

func (filtered FilteredOptions) partPrices(parts []products.Part) Options {
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
		if _, ok := existing[low]; ok {
			continue
		}
		val := fmt.Sprintf("$%d - $%d", low, low+50)
		opt := Option{
			Value:    val,
			Selected: false,
		}

		for _, p := range parts {
			for _, pr := range p.Pricing {
				if pr.Type == "List" && (int(pr.Price) >= low && int(pr.Price) <= (low+50)) {
					opt.Products = append(opt.Products, p.PartId)
					break
				}
			}
		}

		priceDefinitions.Options = append(priceDefinitions.Options, opt)
		existing[low] = low
	}

	return priceDefinitions
}

func (filtered FilteredOptions) partCategory(parts []products.Part) Options {
	var opt Options

	existing := make(map[string]string, 0)
	for _, p := range parts {
		if len(p.Categories) > 0 {
			opt.Key = "Category"
			cat := p.Categories[0]

			if _, ok := existing[cat.Title]; !ok {
				newOption := Option{
					Value:    cat.Title,
					Products: []int{p.PartId},
				}
				opt.Options = append(opt.Options, newOption)
				existing[cat.Title] = cat.Title
				continue
			}

			for i, o := range opt.Options {
				if o.Value == cat.Title {
					prods := opt.Options[i].Products
					prods = append(prods, p.PartId)
					sort.Ints(prods)
					opt.Options[i].Products = prods
				}
			}
		}
	}

	sortutil.AscByField(opt.Options, "Value")

	return opt
}

func (filtered FilteredOptions) partClass(parts []products.Part) Options {
	opt := Options{
		Key: "Class",
	}

	existing := make(map[string]string, 0)
	for _, p := range parts {
		if p.PartClass == "" {
			p.PartClass = "Other"
		}

		if _, ok := existing[p.PartClass]; !ok {
			newOption := Option{
				Value: p.PartClass,
			}
			opt.Options = append(opt.Options, newOption)
			existing[p.PartClass] = p.PartClass
		}

		for i, o := range opt.Options {
			if p.PartClass == o.Value {
				o.Products = append(o.Products, p.PartId)
				opt.Options[i] = o
			}
		}
	}

	sortutil.AscByField(opt.Options, "Value")

	return opt
}
