package zip

import (
	"os"
	"io"
	"fmt"
	"strings"
	"strconv"
	"encoding/csv"
)

type ZipReader struct {
	Path string
} 

func getValue(record []string, columns map[string]int, colName, defVal string) string {
	if idx, ok := columns[colName]; ok {
		return record[idx]
	} else {
		return defVal
	}
}

func (r ZipReader) Read (ch chan ZipCodeEntry) {
	fmt.Println("Reading init file:",r.Path)
	file, err := os.Open(r.Path)
	count := uint32(0)
	if err != nil {
		fmt.Println("Error:",err)
		close(ch)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	columns := make(map[string]int)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:",err)
			break
		}
		if len(columns) == 0 {
			// configure columns
			for i, col := range record {
				columns[col] = i
			}
		} else {
			if dec, decOk := columns["decommissioned"]; !decOk || (decOk && record[dec] != "1") {
				latitude, err := strconv.ParseFloat(getValue(record, columns, "latitude", ""), 64)
				if err != nil {
					latitude = 0
				}
				longitude, err := strconv.ParseFloat(getValue(record, columns, "longitude", ""), 64)
				if err != nil {
					longitude = 0
				}
				acceptableCities := make([] string, 0)
				unacceptableCities := make([] string, 0)
				areaCodes := make([] string, 0)
				if len(getValue(record, columns, "acceptable_cities", "")) > 0 {
					acceptableCities = strings.Split(getValue(record, columns, "acceptable_cities", ""), ", ")
				}
				if len(getValue(record, columns, "unacceptable_cities", "")) > 0 {
					unacceptableCities = strings.Split(getValue(record, columns, "unacceptable_cities", ""), ", ")
				}
				if len(getValue(record, columns, "area_codes", "")) > 0 {
					areaCodes = strings.Split(getValue(record, columns, "area_codes", ""), ", ")
				}
				city := getValue(record, columns, "primary_city", "")
				if strings.Index(city, " (") != -1 {
					city = strings.Replace(city, " (", ", ", -1)
					city = strings.Replace(city, ")", "", -1)
				}
				if strings.Index(city, " /") != -1 {
					city = strings.Replace(city, " /", ",", -1)
				}
				if strings.Index(city, ", ") != -1 {
					cityList := strings.Split(city, ", ")
					city = cityList[0]
					acceptableCities = cityList[1:]
				}
				ch <- ZipCodeEntry { getValue(record, columns, "zip", ""), 	// Zip Code
					getValue(record, columns, "type", "STANDARD"),		// Type
					city, 							// City
					acceptableCities,					// Acceptable Cities
					unacceptableCities,					// Unacceptable Cities
					getValue(record, columns, "county", ""),		// County
					getValue(record, columns, "state", ""),			// State
					getValue(record, columns, "country", ""),		// Country
					getValue(record, columns, "timezone", ""),		// TimeZone
					areaCodes,						// Area Codes
					latitude,						// Latitude
					longitude }						// Longitude
				count = count + 1
			}
		}
	}
	fmt.Printf("Read %v records from file %s\n", count, r.Path)
	ch <- ZipCodeEntry{}
}
