package conversions

import (
	"net/url"
	"strconv"
	"time"
)

//Conversion funcs
func ByteToString(input []byte) (string, error) {
	var err error
	if input != nil {
		output := string(input)
		return output, err
	}
	return "", err
}

func ByteToInt(input []byte) (int, error) {
	var err error
	if input != nil {
		temp, err := ByteToString(input)
		output, err := strconv.Atoi(temp)
		return output, err
	}
	return 0, err
}

func ByteToFloat(input []byte) (float64, error) {
	var err error
	if input != nil {
		output, err := strconv.ParseFloat(string(input), 64)
		return output, err
	}
	return 0.0, err
}

func ByteToUrl(input []byte) (url.URL, error) {
	var err error
	if input != nil {
		str := string(input[:])
		output, err := url.Parse(str)
		output2 := *output
		return output2, err
	}
	output, err := url.Parse("")
	output2 := *output
	return output2, err
}
func ByteToTime(input []byte, timeFormat string) (time.Time, error) {
	var err error
	if input != nil {
		str := string(input[:])
		output, err := time.Parse(timeFormat, str)
		return output, err
	}
	output, err := time.Parse(timeFormat, "")
	return output, err
}
