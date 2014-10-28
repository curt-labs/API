package pagination

import (
	"strconv"
)

type Pagination struct {
	TotalItems    int `json:"total_items" xml:"total_items"`
	ReturnedCount int `json:"returned_count" xml:"returned_count"`
	Page          int `json:"page" xml:"page"`
	PerPage       int `json:"per_page" xml:"per_page"`
	TotalPages    int `json:"total_pages" xml:"total_pages"`
}

type Objects struct {
	Objects    []interface{}
	Pagination Pagination
}

func Paginate(pageStr, resultsStr string, wholeArray []interface{}) Objects {
	var o Objects
	var partialArray []interface{}
	var results int
	page, err := strconv.Atoi(pageStr)
	page = page - 1

	if resultsStr != "" {
		results, err = strconv.Atoi(resultsStr)
		if err != nil {
			return o
		}
	}

	if page < 0 {
		page = 0
	}
	startingIndex := page * results
	endingIndex := startingIndex + results
	if page > 0 || results > 0 {
		if endingIndex > len(wholeArray) {
			endingIndex = len(wholeArray)
		}

		if startingIndex > len(wholeArray) {
			startingIndex = len(wholeArray)
		}
		partialArray = wholeArray[startingIndex:endingIndex]
	} else {
		partialArray = wholeArray
	}
	totalPages := 1
	if len(partialArray) > 0 {
		totalPages = len(wholeArray) / len(partialArray)
	}

	o.Pagination = Pagination{
		TotalItems:    len(wholeArray),
		ReturnedCount: len(partialArray),
		Page:          page + 1,
		PerPage:       len(partialArray),
		TotalPages:    totalPages,
	}
	o.Objects = partialArray
	return o
}
