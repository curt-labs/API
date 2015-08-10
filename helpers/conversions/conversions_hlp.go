package conversions

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

//Conversion funcs
func ByteToString(input []byte) (string, error) {
	if input != nil {
		return string(input), nil
	}
	return "", fmt.Errorf("%s", "failed to parse")
}

func ByteToInt(input []byte) (int, error) {
	if input == nil {
		return 0, fmt.Errorf("%s", "failed to parse")
	}

	temp, err := ByteToString(input)
	if err != nil {
		return 0, fmt.Errorf("%s", "failed to parse")
	}
	return strconv.Atoi(temp)
}

func ByteToFloat(input []byte) (float64, error) {
	if input == nil {
		return 0.0, fmt.Errorf("%s", "failed to parse")
	}

	return strconv.ParseFloat(string(input), 64)
}

func ByteToUrl(input []byte) (url.URL, error) {
	if input == nil {
		return url.URL{}, fmt.Errorf("%s", "failed to parse")
	}

	output, err := url.Parse(string(input[:]))
	if err != nil || output == nil {
		return url.URL{}, fmt.Errorf("%s", "failed to parse")
	}

	return *output, nil
}

func ByteToTime(input []byte, timeFormat string) (time.Time, error) {
	if input == nil || len(input) == 0 {
		return time.Time{}, fmt.Errorf("%s", "failed to parse")
	}

	return time.Parse(timeFormat, string(input[:]))
}

func ParseBool(input []byte) (bool, error) {
	if input == nil {
		return false, fmt.Errorf("%s", "failed to parse")
	}

	return strconv.ParseBool(string(input))
}
