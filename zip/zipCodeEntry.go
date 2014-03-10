package zip

import (
	"bytes"
	"fmt"
	"strings"
	"reflect"
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

type ZipSorter []ZipCodeEntry

func (z ZipSorter) Len() int {
	return len(z)
}

func (z ZipSorter) Swap(i, j int) {
	z[i], z[j] = z[j], z[i] 
}

func (z ZipSorter) Less(i, j int) bool {
	ic := 0
	jc := 0
	if z[i].Country != "US" {
		ic = 1
	}
	if z[j].Country != "US" {
		jc = 1
	}
	if ic != jc {
		return ic < jc
	}
	if z[i].Country != z[j].Country {
		return z[i].Country < z[j].Country
	}
	return z[i].ZipCode < z[j].ZipCode
}

func MarshalEntries(entries []ZipCodeEntry, format string) (string, error) {
	switch format {
	case "XML":
		return EntriesToXml(entries), nil
	case "YAML":
		return EntriesToYaml(entries), nil
	case "JSON":
		return EntriesToJson(entries), nil
	case "JS":
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

func EntriesToYaml(entries []ZipCodeEntry) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Total: %v\n\n",len(entries)))
	buf.WriteString("ZipCodeEntries:\n")

	for _, entry := range entries {
		buf.WriteString(entry.ToYaml())
	}
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
	} else if format == "JSON" || format == "JS" {
		return z.ToJson(), nil
	} else if format == "YAML" {
		return z.ToYaml(), nil
	}
	return "", Throw(fmt.Sprintf("Invalid format: %s", format))
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

func (z ZipCodeEntry) ToYaml() string {
	buf := bytes.Buffer{}
	buf.WriteString("  - ZipCodeEntry:\n")
	
	val := reflect.ValueOf(z)
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typeField := val.Type().Field(i)
		f := valField.Interface()
		val := reflect.ValueOf(f)
		buf.WriteString(fmt.Sprintf("      %s:", typeField.Name))
		for j := 0; j < (20 - len(typeField.Name)); j++ { buf.WriteString(" ") }
		switch val.Kind() {
		case reflect.String:
			buf.WriteString(val.String())
		case reflect.Float64:
			buf.WriteString(fmt.Sprintf("%v", val.Float()))
		case reflect.Slice:
			buf.WriteString("[")
			for j := 0; j < val.Len(); j++ {
				if j != 0 { buf.WriteString(", ") }
				buf.WriteString(fmt.Sprintf("%s", val.Index(j)))
			}
			buf.WriteString("]")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
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
