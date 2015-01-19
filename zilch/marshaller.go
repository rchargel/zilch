package zilch

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Marshal marshals the a DistributionEntry.
func (d DistributionMarshaller) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	buf := bytes.Buffer{}
	switch format {
	case "JS":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&d); err != nil {
			return "", err
		}
		return buf.String(), nil
	case "JSON":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&d); err != nil {
			return "", err
		}
		return buf.String(), nil
	case "XML":
		enc := xml.NewEncoder(&buf)
		if err := enc.Encode(&d); err != nil {
			return "", err
		}
		return "<?xml version=\"1.0\" encoding=\"UTF-8\"?><DistributionList>" + buf.String() + "</DistributionList>", nil
	case "YAML":
		buf.WriteString("DistributionEntries:\n")
		for _, entry := range d {
			buf.WriteString(fmt.Sprintf("  - ZipCodes:  %v\n", entry.ZipCodes))
			buf.WriteString(fmt.Sprintf("    Latitude:  %v\n", entry.Latitude))
			buf.WriteString(fmt.Sprintf("    Longitude: %v\n\n", entry.Longitude))
		}
		return buf.String(), nil
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

// Marshal marshals the a country-zipcodes map.
func (c CountryMarshaller) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	buf := bytes.Buffer{}
	switch format {
	case "JS":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&c); err != nil {
			return "", err
		}
		return buf.String(), nil
	case "JSON":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&c); err != nil {
			return "", err
		}
		return buf.String(), nil
	case "XML":
		buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Countries>")
		for key, val := range map[string]int(c) {
			buf.WriteString(fmt.Sprintf("<%v>%v</%v>", key, val, key))
		}
		buf.WriteString("</Countries>")
		return buf.String(), nil
	case "YAML":
		for key, val := range map[string]int(c) {
			buf.WriteString(fmt.Sprintf("%v: %v\n", key, val))
		}
		return buf.String(), nil
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

// Marshal marshals the CountryEntry object.
func (c CountryEntryMarshaller) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	buf := bytes.Buffer{}
	switch format {
	case "JS":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&c); err != nil {
			return "", err
		}
	case "JSON":
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&c); err != nil {
			return "", err
		}
	case "XML":
		buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?><Countries>`)
		enc := xml.NewEncoder(&buf)
		if err := enc.Encode(&c); err != nil {
			return "", err
		}
		buf.WriteString("</Countries>")
	case "YAML":
		for _, ce := range c {
			buf.WriteString(fmt.Sprintf("  - Country:     %v\n", ce.Country))
			buf.WriteString(fmt.Sprintf("    CountryName: %v\n", ce.CountryName))
			buf.WriteString("    States:\n")

			for _, se := range ce.States {
				buf.WriteString(fmt.Sprintf("      - State:     %v\n", se.State))
				buf.WriteString(fmt.Sprintf("        StateName: %v\n", se.StateName))
				buf.WriteString(fmt.Sprintf("        ZipCodes:  %v\n\n", se.ZipCodes))
			}
		}
	default:
		return "", errors.New("Invalid format: " + format)
	}
	return buf.String(), nil
}

// Marshal marshals the ZipEntry object.
func (z ZipEntry) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	switch format {
	case "XML":
		return z.toXML()
	case "JS":
		return z.toJSON()
	case "JSON":
		return z.toJSON()
	case "YAML":
		return z.toYAML()
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

// Marshal marshals the QueryResult object.
func (q QueryResult) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	switch format {
	case "XML":
		return q.toXML()
	case "JS":
		return q.toJSON()
	case "JSON":
		return q.toJSON()
	case "YAML":
		return q.toYAML()
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

func (q QueryResult) toJSON() (string, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	r := regexp.MustCompile("\\s+$")
	if err := enc.Encode(&q); err != nil {
		return "", err
	}
	return r.ReplaceAllString(buf.String(), ""), nil
}

func (q QueryResult) toXML() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("<QueryResult>")
	buf.WriteString(fmt.Sprintf("<ResultsReturned>%v</ResultsReturned>", q.ResultsReturned))
	buf.WriteString(fmt.Sprintf("<TotalFound>%v</TotalFound>", q.TotalFound))
	buf.WriteString(fmt.Sprintf("<StartIndex>%v</StartIndex>", q.StartIndex))
	buf.WriteString(fmt.Sprintf("<EndIndex>%v</EndIndex>", q.EndIndex))
	buf.WriteString("<ZipCodeEntries>")
	for _, entry := range q.ZipCodeEntries {
		xml, err := entry.toXML()
		if err != nil {
			return "", err
		}
		buf.WriteString(xml)
	}
	buf.WriteString("</ZipCodeEntries>")
	buf.WriteString("</QueryResult>")
	return buf.String(), nil
}

func (q QueryResult) toYAML() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("ResultsReturned: %v\n", q.ResultsReturned))
	buf.WriteString(fmt.Sprintf("TotalFound:      %v\n", q.TotalFound))
	buf.WriteString(fmt.Sprintf("StartIndex:      %v\n", q.StartIndex))
	buf.WriteString(fmt.Sprintf("EndIndex:        %v\n\n", q.EndIndex))
	buf.WriteString("ZipCodeEntries:\n")

	for _, entry := range q.ZipCodeEntries {
		yaml, err := entry.toYAML()
		if err != nil {
			return "", err
		}
		buf.WriteString(yaml)
	}
	return buf.String(), nil
}

func (z ZipEntry) toJSON() (string, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	r := regexp.MustCompile("\\s+$")
	if err := enc.Encode(&z); err != nil {
		return "", err
	}
	return r.ReplaceAllString(buf.String(), ""), nil
}

func (z ZipEntry) toYAML() (string, error) {
	buf := bytes.Buffer{}
	zval := reflect.ValueOf(z)
	for i := 0; i < zval.NumField(); i++ {
		valField := zval.Field(i)
		typeField := zval.Type().Field(i)
		f := valField.Interface()
		val := reflect.ValueOf(f)
		if i == 0 {
			buf.WriteString("  - ")
		} else {
			buf.WriteString("    ")
		}
		buf.WriteString(fmt.Sprintf("%s:", typeField.Name))
		for j := 0; j < (20 - len(typeField.Name)); j++ {
			buf.WriteString(" ")
		}
		switch val.Kind() {
		case reflect.String:
			buf.WriteString(val.String())
		case reflect.Float32:
			buf.WriteString(strconv.FormatFloat(val.Float(), 'f', -1, 32))
		case reflect.Slice:
			buf.WriteString("[")
			for j := 0; j < val.Len(); j++ {
				if j != 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("%s", val.Index(j)))
			}
			buf.WriteString("]")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	return buf.String(), nil
}

func (z ZipEntry) toXML() (string, error) {
	buf := bytes.Buffer{}

	strtoxml := func(text string) string {
		if strings.Index(text, "&") != -1 {
			return fmt.Sprintf("<![CDATA[%v]]>", text)
		}
		return text
	}

	writetag := func(name, value string) {
		if len(value) == 0 {
			buf.WriteString(fmt.Sprintf("<%v/>", name))
		} else {
			buf.WriteString(fmt.Sprintf("<%v>%v</%v>", name, strtoxml(value), name))
		}
	}

	buf.WriteString("<ZipCodeEntry>")

	zval := reflect.ValueOf(z)
	for i := 0; i < zval.NumField(); i++ {
		valField := zval.Field(i)
		typeField := zval.Type().Field(i)
		f := valField.Interface()
		val := reflect.ValueOf(f)
		tagname := typeField.Name
		switch val.Kind() {
		case reflect.String:
			writetag(tagname, val.String())
		case reflect.Float32:
			writetag(tagname, strconv.FormatFloat(val.Float(), 'f', -1, 32))
		case reflect.Slice:
			if val.Len() == 0 {
				buf.WriteString(fmt.Sprintf("<%v/>", tagname))
			} else {
				buf.WriteString(fmt.Sprintf("<%v>", tagname))
				for j := 0; j < val.Len(); j++ {
					if strings.Index(tagname, "AreaCode") != -1 {
						writetag("AreaCode", val.Index(j).String())
					} else {
						writetag("City", val.Index(j).String())
					}
				}
				buf.WriteString(fmt.Sprintf("</%v>", tagname))
			}
		}
	}

	buf.WriteString("</ZipCodeEntry>")

	return buf.String(), nil
}
