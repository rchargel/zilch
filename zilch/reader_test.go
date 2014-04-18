package zilch

import (
	"testing"
)

func Test_Create_Reader(t *testing.T) {
	path := "../resources/ca_zip_code_database.csv"
	reader := CreateReader(path)

	if reader.Path != path {
		t.Errorf("Invalid path: Found %s, Expected: %s\n", reader.Path, path)
	}
	if reader.CountryCode != "CA" {
		t.Errorf("Invalid country code: Found %s, Expected: CA\n", reader.CountryCode)
	}
}

func Test_Read(t *testing.T) {
	path := "../resources/ca_zip_code_database.csv"
	reader := CreateReader(path)

	ch := make(chan ZilchEntry, 100)
	go reader.Read(ch)

	count := 0
	for entry := range ch {
		count += 1
		if entry.Country != "CA" {
			t.Errorf("Wrong country code: %s\n", entry.Country)
		}
	}
	if count != 1640 {
		t.Errorf("Expecting 1640 records, found %v", count)
	}
	t.Log("Read test passed")
}
