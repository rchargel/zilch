package zilch

import (
	"sort"
	"testing"
)

func Test_Sort(t *testing.T) {
	entries := make([]ZipEntry, 3)

	entries[0] = ZipEntry{ZipCode: "90210", City: "Beverly Hills", Country: "US", State: "CA"}
	entries[1] = ZipEntry{ZipCode: "19103", City: "Philadelphia", Country: "US", State: "PA"}
	entries[2] = ZipEntry{ZipCode: "12345", City: "Schenectady", Country: "US", State: "NY"}

	sort.Sort(ZipSorter(entries))

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
		entry := ZipEntry{
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

func Test_GetLatitudeAndLongitudeFromKey(t *testing.T) {
	testNumbers := func(expectedLat, expectedLon int16, key uint32) {
		lat, lon := getLatitudeLongitudeFromKey(key)

		if lat != expectedLat || lon != expectedLon {
			t.Errorf("The key %v produced lat: %v / long: %v, but should have been lat: %v / long: %v\n", key, lat, lon, expectedLat, expectedLon)
		} else {
			t.Logf("The key %v produced lat: %v / long: %v, as expected\n", key, lat, lon)
		}
	}

	testNumbers(int16(-90), int16(-180), uint32(0))
	testNumbers(int16(0), int16(-180), uint32(90))
	testNumbers(int16(90), int16(-180), uint32(180))
	testNumbers(int16(-90), int16(0), uint32(180000))
	testNumbers(int16(-90), int16(180), uint32(360000))
	testNumbers(int16(0), int16(0), uint32(180090))
	testNumbers(int16(90), int16(180), uint32(360180))
	testNumbers(int16(-1), int16(-1), uint32(179089))
	testNumbers(int16(0), int16(0), uint32(180090))
	testNumbers(int16(38), int16(-78), uint32(102128))
	testNumbers(int16(-26), int16(145), uint32(325064))
}
