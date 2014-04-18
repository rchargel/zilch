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
		return z.ToXML(), nil
	case "JS":
		return z.ToJSON(), nil
	case "JSON":
		return z.ToJSON(), nil
	case "YAML":
		return z.ToYAML(), nil
	default:
		return "", errors.New("Invalid format: " + format)
	}
}

func (z ZilchEntry) ToJSON() string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	r := regexp.MustCompile("\\s+$")
	if err := enc.Encode(&z); err != nil {
		fmt.Println("Error:", err)
	}
	return r.ReplaceAllString(buf.String(), "")
}

func (z ZilchEntry) ToYAML() string {
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
	return buf.String()
}

func (z ZilchEntry) ToXML() string {
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

	return buf.String()
}
