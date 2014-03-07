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
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:",err)
			break
		}
		if record[0] != "zip" {
			if record[13] != "1" {
				latitude, err := strconv.ParseFloat(record[9], 64)
				if err != nil {
					latitude = 0
				}
				longitude, err := strconv.ParseFloat(record[10], 64)
				if err != nil {
					longitude = 0
				}
				acceptableCities := make([] string, 0)
				unacceptableCities := make([] string, 0)
				areaCodes := make([] string, 0)
				if len(record[3]) > 0 {
					acceptableCities = strings.Split(record[3], ", ")
				}
				if len(record[4]) > 0 {
					unacceptableCities = strings.Split(record[4], ", ")
				}
				if len(record[8]) > 0 {
					areaCodes = strings.Split(record[8], ",")
				}
				ch <- ZipCodeEntry { record[0], 	// Zip Code
					record[1],			// Type
					record[2],			// City
					acceptableCities,		// Acceptable Cities
					unacceptableCities,		// Unacceptable Cities
					record[6],			// County
					record[5],			// State
					record[12],			// Country
					record[7],			// TimeZone
					areaCodes,			// Area Codes
					latitude,			// Latitude
					longitude }			// Longitude
				count = count + 1
			}
		}
	}
	fmt.Printf("Read %v records from file\n", count)
	close(ch)
}
