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
	DECOMISSIONED_COL       string = "decommissioned"
	ZIP_CODE_COL            string = "zip"
	TYPE_COL                string = "type"
	CITY_COL                string = "primary_city"
	ACCEPTABLE_CITIES_COL   string = "acceptable_cities"
	UNACCEPTABLE_CITIES_COL string = "unacceptable_cities"
	COUNTY_COL              string = "county"
	STATE_COL               string = "state"
	COUNTRY_COL             string = "country"
	COUNTRY_NAME_COL        string = "country_name"
	TIMEZONE_COL            string = "timezone"
	AREA_CODES_COL          string = "area_codes"
	LATITUDE_COL            string = "latitude"
	LONGITUDE_COL           string = "longitude"
)

func CreateReader(path string) ZilchEntryReader {
	r := regexp.MustCompile("\\/[a-z]{2}_")
	cc := r.FindString(path)
	cc = cc[1:3]

	return ZilchEntryReader{
		Path:        path,
		CountryCode: strings.ToUpper(cc),
	}
}

type ZilchEntryReader struct {
	Path        string
	CountryCode string
}

func (r ZilchEntryReader) Read(ch chan ZilchEntry) {
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
				if decom, decomFound := columns[DECOMISSIONED_COL]; !decomFound || (decomFound && record[decom] != "1") {
					// not decomissioned
					latitude := getFloatVal(record, LATITUDE_COL)
					longitude := getFloatVal(record, LONGITUDE_COL)

					if latitude < -90 || latitude > 90 {
						latitude = 0
					}
					if longitude < -180 || longitude > 180 {
						longitude = 0
					}

					acceptableCities := getSliceVal(record, ACCEPTABLE_CITIES_COL)
					unacceptableCities := getSliceVal(record, UNACCEPTABLE_CITIES_COL)
					areaCodes := getSliceVal(record, AREA_CODES_COL)

					city := getVal(record, CITY_COL, "")
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

					ch <- ZilchEntry{
						ZipCode:            getVal(record, ZIP_CODE_COL, ""),
						Type:               getVal(record, TYPE_COL, "STANDARD"),
						City:               city,
						AcceptableCities:   acceptableCities,
						UnacceptableCities: unacceptableCities,
						County:             getVal(record, COUNTY_COL, ""),
						State:              getVal(record, STATE_COL, ""),
						Country:            getVal(record, COUNTRY_COL, r.CountryCode),
						CountryName:        getVal(record, COUNTRY_NAME_COL, ""),
						TimeZone:           getVal(record, TIMEZONE_COL, ""),
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
