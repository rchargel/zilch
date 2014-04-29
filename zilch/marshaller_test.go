package zilch

import (
	"testing"
)

func Test_Marshal_Results_JSON(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}
	results := QueryResult{
		ZipCodeEntries:  []ZilchEntry{entry},
		ResultsReturned: 1,
		TotalFound:      1,
		StartIndex:      1,
		EndIndex:        1,
	}

	text := `{"ResultsReturned":1,"TotalFound":1,"StartIndex":1,"EndIndex":1,"ZipCodeEntries":[{"ZipCode":"22151","Type":"STANDARD","City":"Springfield","AcceptableCities":["N Springfield","North Springfield"],"UnacceptableCities":["N Springfld"],"County":"Fairfax County","State":"VA","Country":"US","CountryName":"United States of America","TimeZone":"America/New_York","AreaCodes":["703","202"],"Latitude":38.78,"Longitude":-77.17}]}`

	json, err := results.Marshal("JSON")
	if err != nil {
		t.Error(err)
	}
	if json == text {
		t.Log("Correct JSON Formatting")
	} else {
		t.Errorf("Invalid JSON Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", json, text)
	}
}

func Test_Marshal_Results_XML(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}
	results := QueryResult{
		ZipCodeEntries:  []ZilchEntry{entry},
		ResultsReturned: 1,
		TotalFound:      1,
		StartIndex:      1,
		EndIndex:        1,
	}

	text := "<QueryResult><ResultsReturned>1</ResultsReturned><TotalFound>1</TotalFound>" +
		"<StartIndex>1</StartIndex><EndIndex>1</EndIndex><ZipCodeEntries><ZipCodeEntry>" +
		"<ZipCode>22151</ZipCode><Type>STANDARD</Type><City>Springfield</City><AcceptableCities>" +
		"<City>N Springfield</City><City>North Springfield</City></AcceptableCities>" +
		"<UnacceptableCities><City>N Springfld</City></UnacceptableCities><County>Fairfax County</County>" +
		"<State>VA</State><Country>US</Country><CountryName>United States of America</CountryName><TimeZone>America/New_York</TimeZone><AreaCodes>" +
		"<AreaCode>703</AreaCode><AreaCode>202</AreaCode></AreaCodes><Latitude>38.78</Latitude>" +
		"<Longitude>-77.17</Longitude></ZipCodeEntry></ZipCodeEntries></QueryResult>"

	xml, err := results.Marshal("XML")
	if err != nil {
		t.Error(err)
	}
	if xml == text {
		t.Log("Correct XML Formatting")
	} else {
		t.Errorf("Invalid XML Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", xml, text)
	}
}

func Test_Marshal_Results_YAML(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}
	results := QueryResult{
		ZipCodeEntries:  []ZilchEntry{entry},
		ResultsReturned: 1,
		TotalFound:      1,
		StartIndex:      1,
		EndIndex:        1,
	}

	text := `ResultsReturned: 1
TotalFound: 1
StartIndex: 1
EndIndex: 1

ZipCodeEntries:
  - ZipCode:             22151
    Type:                STANDARD
    City:                Springfield
    AcceptableCities:    [N Springfield, North Springfield]
    UnacceptableCities:  [N Springfld]
    County:              Fairfax County
    State:               VA
    Country:             US
    CountryName:         United States of America
    TimeZone:            America/New_York
    AreaCodes:           [703, 202]
    Latitude:            38.78
    Longitude:           -77.17

`

	yaml, err := results.Marshal("YAML")
	if err != nil {
		t.Error(err)
	}
	if yaml == text {
		t.Log("Correct YAML Formatting")
	} else {
		t.Errorf("Invalid YAML Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", yaml, text)
	}
}

func Test_Marshal_JSON(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}

	text := `{"ZipCode":"22151","Type":"STANDARD","City":"Springfield","AcceptableCities":["N Springfield","North Springfield"],"UnacceptableCities":["N Springfld"],"County":"Fairfax County","State":"VA","Country":"US","CountryName":"United States of America","TimeZone":"America/New_York","AreaCodes":["703","202"],"Latitude":38.78,"Longitude":-77.17}`

	json, err := entry.Marshal("JSON")
	if err != nil {
		t.Error(err)
	}
	if json == text {
		t.Log("Correct JSON Formatting")
	} else {
		t.Errorf("Invalid JSON Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", json, text)
	}
}

func Test_Marshal_JS(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}

	text := `{"ZipCode":"22151","Type":"STANDARD","City":"Springfield","AcceptableCities":["N Springfield","North Springfield"],"UnacceptableCities":["N Springfld"],"County":"Fairfax County","State":"VA","Country":"US","CountryName":"United States of America","TimeZone":"America/New_York","AreaCodes":["703","202"],"Latitude":38.78,"Longitude":-77.17}`

	json, err := entry.Marshal("JS")
	if err != nil {
		t.Error(err.Error())
	}
	if json == text {
		t.Log("Correct JS Formatting")
	} else {
		t.Errorf("Invalid JS Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", json, text)
	}
}

func Test_Marshal_YAML(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}

	text := `  - ZipCode:             22151
    Type:                STANDARD
    City:                Springfield
    AcceptableCities:    [N Springfield, North Springfield]
    UnacceptableCities:  [N Springfld]
    County:              Fairfax County
    State:               VA
    Country:             US
    CountryName:         United States of America
    TimeZone:            America/New_York
    AreaCodes:           [703, 202]
    Latitude:            38.78
    Longitude:           -77.17

`

	if yaml, err := entry.Marshal("YAML"); err == nil {
		if yaml == text {
			t.Log("Correct YAML Formatting")
		} else {
			t.Errorf("Invalid YAML Formatting\nFound:\n'%s'\n\nExpecting:\n'%s'\n", yaml, text)
		}
	} else {
		t.Error(err.Error())
	}
}

func Test_Marshal_XML(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}

	text := `<ZipCodeEntry><ZipCode>22151</ZipCode><Type>STANDARD</Type><City>Springfield</City><AcceptableCities><City>N Springfield</City><City>North Springfield</City></AcceptableCities><UnacceptableCities><City>N Springfld</City></UnacceptableCities><County>Fairfax County</County><State>VA</State><Country>US</Country><CountryName>United States of America</CountryName><TimeZone>America/New_York</TimeZone><AreaCodes><AreaCode>703</AreaCode><AreaCode>202</AreaCode></AreaCodes><Latitude>38.78</Latitude><Longitude>-77.17</Longitude></ZipCodeEntry>`

	if xml, err := entry.Marshal("XML"); err == nil {
		if xml == text {
			t.Log("Correct XML Formatting")
		} else {
			t.Errorf("Invalid XML Formatting\nFound: %s\n\nExpecting: %s\n", xml, text)
		}
	} else {
		t.Error(err.Error())
	}
}

func Test_Marshal_Invalid(t *testing.T) {
	entry := ZilchEntry{
		ZipCode:            "22151",
		Type:               "STANDARD",
		City:               "Springfield",
		AcceptableCities:   []string{"N Springfield", "North Springfield"},
		UnacceptableCities: []string{"N Springfld"},
		County:             "Fairfax County",
		Country:            "US",
		CountryName:        "United States of America",
		State:              "VA",
		TimeZone:           "America/New_York",
		AreaCodes:          []string{"703", "202"},
		Latitude:           float32(38.78),
		Longitude:          float32(-77.17),
	}

	if _, err := entry.Marshal("blah"); err == nil {
		t.Error("No error thrown")
	} else if err.Error() == "Invalid format: BLAH" {
		t.Log("Invalid format thrown")
	} else {
		t.Error("Wrong exception:", err.Error())
	}
}
