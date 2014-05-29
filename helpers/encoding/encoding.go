package encoding

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Encoder interface {
	Encode(v ...interface{}) (string, error)
}

func Must(data string, err error) string {
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return ""
	}

	return data
}

type JsonEncoder struct{}

// JsonEncoder is an Encoder that produces JSON-formatted responses.
func (_ JsonEncoder) Encode(v ...interface{}) (string, error) {
	var data interface{} = v
	if v == nil {
		// so that empty results produce `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}
	b, err := json.Marshal(data)
	return string(b), err
}

type XmlEncoder struct{}

// XmlEncoder is an Encoder that produces XML-formatted responses.
func (_ XmlEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	if _, err := buf.Write([]byte(xml.Header)); err != nil {
		return "", err
	}
	b, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	if _, err := buf.Write(b); err != nil {
		return "", err
	}
	return buf.String(), err
}

type TextEncoder struct{}

func (_ TextEncoder) Encode(v ...interface{}) (string, error) {
	var buf bytes.Buffer
	for _, v := range v {
		if _, err := fmt.Fprintf(&buf, "%s\n", v); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
