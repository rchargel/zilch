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

func (z ZipCodeEntry) AddCity(city string) ZipCodeEntry {
	acceptableCities := append(z.AcceptableCities, city)
	return ZipCodeEntry { z.ZipCode,
		z.Type,
		z.City,
		acceptableCities,
		z.UnacceptableCities,
		z.County,
		z.State,
		z.Country,
		z.TimeZone,
		z.AreaCodes,
		z.Latitude,
		z.Longitude }
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
	buf.WriteString(tag("ZipCode", z.ZipCode))
	buf.WriteString(tag("Type", z.Type))
	buf.WriteString(tag("City", z.City))
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
	buf.WriteString(tag("County", z.County))
	buf.WriteString(tag("State", z.State))
	buf.WriteString(tag("Country", z.Country))
	buf.WriteString(tag("TimeZone", z.TimeZone))
	if len(z.AreaCodes) > 0 {
		buf.WriteString("<AreaCodes>")
		for _, areaCode := range z.AreaCodes {
			buf.WriteString("<AreaCode>"+areaCode+"</AreaCode>")
		}
		buf.WriteString("</AreaCodes>")
	} else {
		buf.WriteString("<AreaCodes/>")
	}
	buf.WriteString(tag("Latitude", fmt.Sprintf("%v", z.Latitude)))
	buf.WriteString(tag("Longitude", fmt.Sprintf("%v", z.Longitude)))
	buf.WriteString("</ZipCodeEntry>")
	return buf.String()
}

func strtoxml(text string) string {
	if strings.Index(text, "&") >= 0 {
		return "<![CDATA[" + text + "]]>"
	}
	return text
}

func tag(tagname, text string) string {
	if len(text) == 0 {
		return fmt.Sprintf("<%s/>", tagname)
	}
	return fmt.Sprintf("<%s>%s</%s>", tagname, strtoxml(text), tagname)
}

func (z ZipCodeEntry) ToJson() string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:",err)
	}
	return buf.String()
}
