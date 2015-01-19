package zilch

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	decommissionedCol     string = "decommissioned"
	zipCodeCol            string = "zip"
	typeCol               string = "type"
	cityCol               string = "primary_city"
	acceptableCitiesCol   string = "acceptable_cities"
	unacceptableCitiesCol string = "unacceptable_cities"
	countyCol             string = "county"
	stateCol              string = "state"
	stateNameCol          string = "state_name"
	countryCol            string = "country"
	countryNameCol        string = "country_name"
	timezoneCol           string = "timezone"
	areaCodesCol          string = "area_codes"
	latitudeCol           string = "latitude"
	longitudeCol          string = "longitude"
)

// CreateReader creates a ZipEntryReader.
func CreateReader(path string) ZipEntryReader {
	r := regexp.MustCompile("\\/[a-z]{2}_")
	cc := r.FindString(path)
	cc = cc[1:3]

	return ZipEntryReader{
		Path:        path,
		CountryCode: strings.ToUpper(cc),
	}
}

// ZipEntryReader a reader used to read zip entries out of a CSV file.
type ZipEntryReader struct {
	Path        string
	CountryCode string
}

func (r ZipEntryReader) Read(ch chan ZipEntry) {
	file, err := os.Open(r.Path)
	defer file.Close()

	if err != nil {
		fmt.Println("Error:", err)
		close(ch)
		return
	}
	reader := csv.NewReader(file)
	columns := make(map[string]int)

	getVal := func(record []string, column, defaultValue string) string {
		if colIdx, colFound := columns[column]; colFound {
			val := record[colIdx]
			if len(val) != 0 {
				return val
			}
		}
		return defaultValue
	}

	getFloatVal := func(record []string, column string) float32 {
		val := getVal(record, column, "")
		if len(val) > 0 {
			flVal, err := strconv.ParseFloat(val, 32)
			if err == nil {
				return float32(flVal)
			}
		}
		return 0
	}

	getSliceVal := func(record []string, column string) []string {
		val := getVal(record, column, "")

		if len(val) != 0 {
			return strings.Split(val, ", ")
		}
		return make([]string, 0, 0)
	}

	for {
		if record, err := reader.Read(); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err.Error())
			break
		} else {
			if len(columns) == 0 {
				// setup column headers
				for i, col := range record {
					columns[col] = i
				}
			} else {
				if decom, decomFound := columns[decommissionedCol]; !decomFound || (decomFound && record[decom] != "1") {
					// not decomissioned
					latitude := getFloatVal(record, latitudeCol)
					longitude := getFloatVal(record, longitudeCol)

					if latitude < -90 || latitude > 90 {
						latitude = 0
					}
					if longitude < -180 || longitude > 180 {
						longitude = 0
					}

					acceptableCities := getSliceVal(record, acceptableCitiesCol)
					unacceptableCities := getSliceVal(record, unacceptableCitiesCol)
					areaCodes := getSliceVal(record, areaCodesCol)

					city := getVal(record, cityCol, "")
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

					ch <- ZipEntry{
						ZipCode:            getVal(record, zipCodeCol, ""),
						Type:               getVal(record, typeCol, "STANDARD"),
						City:               city,
						AcceptableCities:   acceptableCities,
						UnacceptableCities: unacceptableCities,
						County:             getVal(record, countyCol, ""),
						State:              getVal(record, stateCol, ""),
						StateName:          getVal(record, stateNameCol, ""),
						Country:            getVal(record, countryCol, r.CountryCode),
						CountryName:        getVal(record, countryNameCol, ""),
						TimeZone:           getVal(record, timezoneCol, ""),
						AreaCodes:          areaCodes,
						Latitude:           latitude,
						Longitude:          longitude,
					}
				}
			}
		}
	}
	close(ch)
}
