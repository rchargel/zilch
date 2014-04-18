package zilch

import (
	"sort"
	"testing"
)

func Test_Sort(t *testing.T) {
	entries := make([]ZilchEntry, 3)

	entries[0] = ZilchEntry{ZipCode: "90210", City: "Beverly Hills", Country: "US", State: "CA"}
	entries[1] = ZilchEntry{ZipCode: "19103", City: "Philadelphia", Country: "US", State: "PA"}
	entries[2] = ZilchEntry{ZipCode: "12345", City: "Schenectady", Country: "US", State: "NY"}

	sort.Sort(ZilchSorter(entries))

	if entries[0].ZipCode != "12345" {
		t.Errorf("Sort failed, was %v, should have been 12345", entries[0].ZipCode)
	} else if entries[1].ZipCode != "19103" {
		t.Errorf("Sort failed, was %v, should have been 19103", entries[1].ZipCode)
	} else if entries[2].ZipCode != "90210" {
		t.Errorf("Sort failed, was %v, should have been 90210", entries[2].ZipCode)
	} else {
		t.Log("Sort test passed")
	}
}

func Test_GetKey(t *testing.T) {

	testNumbers := func(lat, lon float32, expected uint32) {
		entry := ZilchEntry{
			Latitude:  lat,
			Longitude: lon,
		}

		if entry.GetKey() != expected {
			t.Errorf("The key for %v / %v should be %v but was %v\n", entry.Latitude, entry.Longitude, expected, entry.GetKey())
		} else {
			t.Logf("The key for %v / %v = %v, as expected\n", entry.Latitude, entry.Longitude, entry.GetKey())
		}
	}

	testNumbers(float32(-90), float32(-180), uint32(0))
	testNumbers(float32(0), float32(-180), uint32(90))
	testNumbers(float32(90), float32(-180), uint32(180))
	testNumbers(float32(-90), float32(0), uint32(180000))
	testNumbers(float32(-90), float32(180), uint32(360000))
	testNumbers(float32(0), float32(0), uint32(180090))
	testNumbers(float32(90), float32(180), uint32(360180))
	testNumbers(float32(-0.01), float32(-0.01), uint32(179089))
	testNumbers(float32(0.01), float32(0.01), uint32(180090))
	testNumbers(float32(38.78), float32(-77.17), uint32(102128))
	testNumbers(float32(-25.34), float32(145.89), uint32(325064))
	testNumbers(float32(-25.99), float32(145.01), uint32(325064))
}
