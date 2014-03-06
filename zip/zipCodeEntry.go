package zip

import (
	"io"
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

func (z ZipCodeEntry) WriteXml(w io.Writer) {
	enc := xml.NewEncoder(w)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
}

func (z ZipCodeEntry) WriteJson(w io.Writer) {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
}
