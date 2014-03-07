package zip

import (
	"bytes"
	"fmt"
	"strings"
	"encoding/json"
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

func MarshalEntries(entries []ZipCodeEntry, format string) (string, error) {
	switch format {
	case "XML":
		return EntriesToXml(entries), nil
	case "JSON":
		return EntriesToJson(entries), nil
	case "JSONP":
		return EntriesToJson(entries), nil
	}
	return "", Throw(fmt.Sprintf("Invalid format: %s", format))
}

func EntriesToXml(entries []ZipCodeEntry) string {
	buf := bytes.Buffer{}
	buf.WriteString("<Response>")
	buf.WriteString("<Total>"+fmt.Sprintf("%v", len(entries))+"</Total>")
	buf.WriteString("<ZipCodeEntries>")
	for _, entry := range entries {
		buf.WriteString(entry.ToXml())
	}
	buf.WriteString("</ZipCodeEntries>")
	buf.WriteString("</Response>")
	return buf.String()
}

func EntriesToJson(entries []ZipCodeEntry) string {
	buf := bytes.Buffer{}
	buf.WriteString("{\"Total\":")
	buf.WriteString(fmt.Sprintf("%v",len(entries)))
	buf.WriteString(",\"ZipCodeEntries\":")
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&entries); err != nil {
		fmt.Println("Error:", err)
	}
	buf.WriteString("}")
	return buf.String()
}

func (z ZipCodeEntry) Marshal(format string) (string, error) {
	if format == "XML" {
		return z.ToXml(), nil
	} else if format == "JSON" || format == "JSONP" {
		return z.ToJson(), nil
	}
	return "", Throw(fmt.Sprintf("Invalid format: %s", format))
}

func (z ZipCodeEntry) ToXml() string {
	buf := bytes.Buffer{}
	buf.WriteString("<ZipCodeEntry>")
	buf.WriteString("<ZipCode>"+z.ZipCode+"</ZipCode>")
	buf.WriteString("<Type>"+z.Type+"</Type>")
	buf.WriteString("<City>"+strtoxml(z.City)+"</City>")
	if len(z.AcceptableCities) > 0 {
		buf.WriteString("<AcceptableCities>")
		for _, city := range z.AcceptableCities {
			buf.WriteString("<City>"+strtoxml(city)+"</City>")
		}
		buf.WriteString("</AcceptableCities>")
	} else {
		buf.WriteString("<AcceptableCities/>")
	}
	if len(z.UnacceptableCities) > 0 {
		buf.WriteString("<UnacceptableCities>")
		for _, city := range z.UnacceptableCities {
			buf.WriteString("<City>"+strtoxml(city)+"</City>")
		}
		buf.WriteString("</UnacceptableCities>")
	} else {
		buf.WriteString("<UnacceptableCities/>")
	}
	buf.WriteString("<County>"+strtoxml(z.County)+"</County>")
	buf.WriteString("<State>"+z.State+"</State>")
	buf.WriteString("<Country>"+strtoxml(z.Country)+"</Country>")
	buf.WriteString("<TimeZone>"+strtoxml(z.TimeZone)+"</TimeZone>")
	if len(z.AreaCodes) > 0 {
		buf.WriteString("<AreaCodes>")
		for _, areaCode := range z.AreaCodes {
			buf.WriteString("<AreaCode>"+areaCode+"</AreaCode>")
		}
		buf.WriteString("</AreaCodes>")
	} else {
		buf.WriteString("<AreaCodes/>")
	}
	buf.WriteString("<Latitude>"+fmt.Sprintf("%v",z.Latitude)+"</Latitude>")
	buf.WriteString("<Longitude>"+fmt.Sprintf("%v",z.Longitude)+"</Longitude>")
	buf.WriteString("<Population>"+fmt.Sprintf("%v",z.Population)+"</Population>")
	buf.WriteString("</ZipCodeEntry>")
	return buf.String()
}

func strtoxml(text string) string {
	if strings.Index(text, "&") >= 0 {
		return "<![CDATA[" + text + "]]>"
	}
	return text
}

func (z ZipCodeEntry) ToJson() string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
	return buf.String()
}
