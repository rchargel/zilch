package zip

import (
	"bytes"
	"fmt"
	"encoding/json"
	"encoding/xml"
)

type ZipCodeEntry struct {
	ZipCode string
	Type string
	City string
	AcceptableCities []string
	UnacceptableCities []string
	County string
	State string
	Country string
	TimeZone string
	AreaCodes []string
	Latitude float64
	Longitude float64
	Population uint32
}

func (z ZipCodeEntry) Marshal(format string) (string, error) {
	if format == "XML" {
		return z.ToXml(), nil
	} else if format == "JSON" {
		return z.ToJson(), nil
	}
	return "", Throw(fmt.Sprintf("Invalid format: %", format))
}

func (z ZipCodeEntry) ToXml() string {
	buf := bytes.Buffer{}
	enc := xml.NewEncoder(&buf)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
	return buf.String()
}

func (z ZipCodeEntry) ToJson() string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
	return buf.String()
}
