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
	fmt.Println("Reading file:",r.Path)
	file, err := os.Open(r.Path)
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
				population, err := strconv.ParseUint(record[14], 10, 32)
				latitude, err := strconv.ParseFloat(record[9], 64)
				longitude, err := strconv.ParseFloat(record[10], 64)
				if err != nil {
					population = 0
				}
				ch <- ZipCodeEntry { record[0], 	// Zip Code
					record[1],			// Type
					record[2],			// City
					strings.Split(record[3], ", "),	// Acceptable Cities
					strings.Split(record[4], ", "),	// Unacceptable Cities
					record[6],			// County
					record[5],			// State
					record[12],			// Country
					record[7],			// TimeZone
					strings.Split(record[8], ","),	// Area Codes
					latitude,			// Latitude
					longitude,			// Longitude
					uint32(population) }		// Population
			}
		}
	}
	close(ch)
}
