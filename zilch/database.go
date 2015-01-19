package zilch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxEntries int = 200
)

// CountryIndex tracks the country to the list of ZipEntry objects.
type CountryIndex struct {
	CountryCode string
	Entries     []ZipEntry
}

// Database is a representation of the actual database of zip codes.
type Database struct {
	CountryIndexMap map[string]CountryIndex
	DistributionMap map[uint32]DistributionEntry
	CountryList     []CountryEntry
	FullyLoaded     bool
}

// NewDatabase creates a database from the file directory.
func NewDatabase(filedir string) (*Database, error) {
	start := time.Now()
	d := &Database{
		CountryIndexMap: make(map[string]CountryIndex),
		DistributionMap: make(map[uint32]DistributionEntry),
		CountryList:     make([]CountryEntry, 0, 0),
		FullyLoaded:     false,
	}

	channelMap := make(map[string]chan ZipEntry)

	// start reading data
	files, err := ioutil.ReadDir(filedir)
	if err != nil {
		return d, err
	}

	distributionChannel := make(chan map[uint32]int)
	channels := 0

	for _, file := range files {
		var filepath string
		if filedir[len(filedir):] != "/" {
			filepath = filedir + "/" + file.Name()
		} else {
			filepath = filedir + file.Name()
		}
		reader := CreateReader(filepath)
		readerChan := make(chan ZipEntry, 20)
		channelMap[reader.CountryCode] = readerChan
		go reader.Read(readerChan)
		go d.loadCountryData(reader.CountryCode, readerChan, distributionChannel)
		channels++
	}

	go d.finishDistributionChannels(distributionChannel, channels, start)

	return d, nil
}

func (d *Database) loadCountryData(countryCode string, channel chan ZipEntry, distChannel chan map[uint32]int) {
	entries := make([]ZipEntry, 0, 1000)
	distMap := make(map[uint32]int)

	type stateData struct {
		State     string
		StateName string
		ZipCodes  int
	}
	type countryData struct {
		Country     string
		CountryName string
		States      map[string]stateData
	}

	country := countryData{
		Country:     countryCode,
		CountryName: "",
		States:      make(map[string]stateData),
	}

	for entry := range channel {
		entries = append(entries, entry)
		distMap[entry.GetKey()]++

		if len(country.CountryName) == 0 {
			country.CountryName = entry.CountryName
		}

		if len(entry.State) > 0 {
			if state, found := country.States[entry.State]; found {
				country.States[entry.State] = stateData{
					State:     state.State,
					StateName: state.StateName,
					ZipCodes:  state.ZipCodes + 1,
				}
			} else {
				country.States[entry.State] = stateData{
					State:     entry.State,
					StateName: entry.StateName,
					ZipCodes:  1,
				}
			}
		}
	}

	countryEntry := CountryEntry{
		Country:     country.Country,
		CountryName: country.CountryName,
		States:      make([]StateEntry, len(country.States)),
	}
	var idx int
	for _, st := range country.States {
		countryEntry.States[idx] = StateEntry{
			State:     st.State,
			StateName: st.StateName,
			ZipCodes:  uint32(st.ZipCodes),
		}
		idx++
	}

	sort.Sort(StateSorter(countryEntry.States))

	d.CountryIndexMap[countryCode] = CountryIndex{
		CountryCode: countryCode,
		Entries:     entries,
	}
	d.CountryList = append(d.CountryList, countryEntry)

	distChannel <- distMap
}

// IsFullyLoaded determines whether the database has finished being
// initialized out of the filesystem.
func (d *Database) IsFullyLoaded() bool {
	return d.FullyLoaded
}

func (d *Database) finishDistributionChannels(distChannel chan map[uint32]int, totalChannels int, startTime time.Time) {
	channels := 0

	for distMap := range distChannel {
		channels++

		for key, total := range distMap {
			// key == 180090 = where equater meets prime meridian, not a real place
			if key != 180090 {
				lat, lon := getLatitudeLongitudeFromKey(key)
				if _, found := d.DistributionMap[key]; !found {
					d.DistributionMap[key] = DistributionEntry{
						Latitude:  lat,
						Longitude: lon,
						ZipCodes:  uint32(total),
					}
				} else {
					zipCodes := d.DistributionMap[key].ZipCodes + uint32(total)

					d.DistributionMap[key] = DistributionEntry{
						Latitude:  lat,
						Longitude: lon,
						ZipCodes:  uint32(zipCodes),
					}
				}
			}
		}

		if channels == totalChannels {
			break
		}
	}

	sort.Sort(CountrySorter(d.CountryList))

	d.FullyLoaded = true

	ellapsedTime := time.Since(startTime)
	fmt.Printf("Finished reading database in %s.\n", ellapsedTime)
}

// GetDistributions gets the list of DistributionEntry objects.
func (d *Database) GetDistributions() []DistributionEntry {
	entries := make([]DistributionEntry, len(d.DistributionMap))
	idx := 0

	for _, entry := range d.DistributionMap {
		entries[idx] = entry
		idx++
	}

	sort.Sort(DistributionSorter(entries))
	return entries
}

// ExecQuery executes a query against the database.
func (d *Database) ExecQuery(queryParams map[string]string) (QueryResult, error) {
	if len(queryParams) == 0 {
		return QueryResult{}, errors.New("There are no query parameters")
	}
	var entries []ZipEntry
	var err error
	if country, found := queryParams["Country"]; found {
		entries, err = d.querySingleCountry(country, queryParams)
	} else {
		entries, err = d.queryAllCountries(queryParams)
	}
	if err != nil {
		return QueryResult{}, err
	}

	total := len(entries)
	start := 0
	end := len(entries)
	sort.Sort(ZipSorter(entries))
	if total > maxEntries {
		start = 0
		end = maxEntries
		if page, pageFound := queryParams["page"]; pageFound {
			p, perr := strconv.ParseUint(page, 10, 32)
			if perr != nil {
				return QueryResult{}, perr
			}
			start = (int(p) - 1) * maxEntries
			end = start + maxEntries
			if start >= total {
				start = total
			}
			if end >= total {
				end = total
			}
		}
		entries = entries[start:end]
	}
	return QueryResult{
		ResultsReturned: len(entries),
		TotalFound:      total,
		StartIndex:      start + 1,
		EndIndex:        end,
		ZipCodeEntries:  entries,
	}, nil
}

func (d *Database) querySingleCountry(country string, queryParams map[string]string) ([]ZipEntry, error) {
	entries := make([]ZipEntry, 0, 20)
	if countryIndex, found := d.CountryIndexMap[country]; found {
		ch := make(chan ZipEntry)
		go countryIndex.QueryIndex(queryParams, ch)
		for entry := range ch {
			entries = append(entries, entry)
		}
		return entries, nil
	}
	return entries, fmt.Errorf("No country %s found", country)
}

func (d *Database) queryAllCountries(queryParams map[string]string) ([]ZipEntry, error) {
	entries := make([]ZipEntry, 0, 40)
	totalCountries := len(d.CountryIndexMap)
	completed := 0
	ch := make(chan ZipEntry)
	for _, countryIndex := range d.CountryIndexMap {
		go countryIndex.queryIndexNoClose(queryParams, ch)
	}
	for entry := range ch {
		if entry.Type == "EOL" {
			completed++
			if completed >= totalCountries {
				break
			}
		} else {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

// QueryIndex executes a query against the CountryIndex.
func (c CountryIndex) QueryIndex(queryParams map[string]string, ch chan ZipEntry) {
	nonalphanum := regexp.MustCompile("[^0-9A-Za-z]")
	stringData := func(paramName string, params map[string]string) (string, bool) {
		if value, valExists := params[paramName]; valExists {
			return strings.ToLower(value), true
		}
		return "", false
	}
	boundsData := func(params map[string]string) ([]float32, bool) {
		b := make([]float32, 4)
		if val, valExists := params["Bounds"]; valExists {
			ba := strings.Split(val, ",")
			if len(ba) != 4 {
				return b, false
			}
			for i, bastring := range ba {
				if f, err := strconv.ParseFloat(bastring, 32); err == nil {
					b[i] = float32(f)
				} else {
					return b, false
				}
			}
			return b, true
		}
		return b, false
	}
	startsWith := func(expected, actual string) bool {
		return strings.Index(actual, expected) == 0
	}
	contains := func(expected, actual string) bool {
		return strings.Index(actual, expected) != -1
	}
	inArray := func(expected string, actual []string) bool {
		for _, val := range actual {
			if strings.Index(strings.ToLower(val), expected) != -1 {
				return true
			}
		}
		return false
	}
	inBounds := func(bounds []float32, latitude, longitude float32) bool {
		if latitude == 0 && longitude == 0 {
			return false
		}
		if latitude > bounds[0] || latitude < bounds[2] {
			return false
		}
		if longitude < bounds[1] || longitude > bounds[3] {
			return false
		}
		return true
	}
	bounds, boundsTest := boundsData(queryParams)
	zipCode, zipCodeTest := stringData("ZipCode", queryParams)
	city, cityTest := stringData("City", queryParams)
	areaCode, areaCodeTest := stringData("AreaCode", queryParams)
	state, stateTest := stringData("State", queryParams)
	county, countyTest := stringData("County", queryParams)

	for _, entry := range c.Entries {
		if zipCodeTest {
			if !startsWith(nonalphanum.ReplaceAllString(zipCode, ""), strings.ToLower(nonalphanum.ReplaceAllString(entry.ZipCode, ""))) {
				continue
			}
		}
		if cityTest {
			if !contains(city, strings.ToLower(entry.City)) {
				valid := false
				if len(entry.AcceptableCities) > 0 && inArray(city, entry.AcceptableCities) {
					valid = true
				}
				if len(entry.UnacceptableCities) > 0 && inArray(city, entry.UnacceptableCities) {
					valid = true
				}
				if !valid {
					continue
				}
			}
		}
		if areaCodeTest {
			if !inArray(areaCode, entry.AreaCodes) {
				continue
			}
		}
		if stateTest {
			if state != strings.ToLower(entry.State) {
				if len(state) == 2 || !contains(state, strings.ToLower(entry.StateName)) {
					continue
				}
			}
		}
		if countyTest {
			if !contains(county, entry.County) {
				continue
			}
		}
		if boundsTest {
			if !inBounds(bounds, entry.Latitude, entry.Longitude) {
				continue
			}
		}
		ch <- entry
	}
	close(ch)
}

func (c CountryIndex) queryIndexNoClose(queryParams map[string]string, ch chan ZipEntry) {
	ch2 := make(chan ZipEntry)
	go c.QueryIndex(queryParams, ch2)
	for entry := range ch2 {
		ch <- entry
	}
	ch <- ZipEntry{Type: "EOL"}
}
