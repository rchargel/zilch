package zilch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func (z ZilchEntry) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	switch format {
	case "XML":
		return z.ToXML()
	case "JS":
		return z.ToJSON()
	case "JSON":
		return z.ToJSON()
	case "YAML":
		return z.ToYAML()
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

func (q QueryResult) Marshal(format string) (string, error) {
	format = strings.ToUpper(format)
	switch format {
	case "XML":
		return q.ToXML()
	case "JS":
		return q.ToJSON()
	case "JSON":
		return q.ToJSON()
	case "YAML":
		return q.ToYAML()
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

func (q QueryResult) ToJSON() (string, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	r := regexp.MustCompile("\\s+$")
	if err := enc.Encode(&q); err != nil {
		return "", err
	}
	return r.ReplaceAllString(buf.String(), ""), nil
}

func (q QueryResult) ToXML() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString("<QueryResult>")
	buf.WriteString(fmt.Sprintf("<ResultsReturned>%v</ResultsReturned>", q.ResultsReturned))
	buf.WriteString(fmt.Sprintf("<TotalFound>%v</TotalFound>", q.TotalFound))
	buf.WriteString(fmt.Sprintf("<StartIndex>%v</StartIndex>", q.StartIndex))
	buf.WriteString(fmt.Sprintf("<EndIndex>%v</EndIndex>", q.EndIndex))
	buf.WriteString("<ZipCodeEntries>")
	for _, entry := range q.ZipCodeEntries {
		if xml, err := entry.ToXML(); err != nil {
			return "", err
		} else {
			buf.WriteString(xml)
		}
	}
	buf.WriteString("</ZipCodeEntries>")
	buf.WriteString("</QueryResult>")
	return buf.String(), nil
}

func (q QueryResult) ToYAML() (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("ResultsReturned: %v\n", q.ResultsReturned))
	buf.WriteString(fmt.Sprintf("TotalFound: %v\n", q.TotalFound))
	buf.WriteString(fmt.Sprintf("StartIndex: %v\n", q.StartIndex))
	buf.WriteString(fmt.Sprintf("EndIndex: %v\n\n", q.EndIndex))
	buf.WriteString("ZipCodeEntries:\n")

	for _, entry := range q.ZipCodeEntries {
		if yaml, err := entry.ToYAML(); err != nil {
			return "", err
		} else {
			buf.WriteString(yaml)
		}
	}
	return buf.String(), nil
}

func (z ZilchEntry) ToJSON() (string, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	r := regexp.MustCompile("\\s+$")
	if err := enc.Encode(&z); err != nil {
		return "", err
	}
	return r.ReplaceAllString(buf.String(), ""), nil
}

func (z ZilchEntry) ToYAML() (string, error) {
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

func (z ZilchEntry) ToXML() (string, error) {
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
