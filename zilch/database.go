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
	max_entries int = 200
)

type CountryIndex struct {
	CountryCode string
	Entries     []ZilchEntry
}

type Database struct {
	CountryIndexMap map[string]CountryIndex
	DistributionMap map[uint32]DistributionEntry
	CountryList     []CountryEntry
	FullyLoaded     bool
}

func NewDatabase(filedir string) (*Database, error) {
	start := time.Now()
	d := &Database{
		CountryIndexMap: make(map[string]CountryIndex),
		DistributionMap: make(map[uint32]DistributionEntry),
		CountryList:     make([]CountryEntry, 0, 0),
		FullyLoaded:     false,
	}

	channelMap := make(map[string]chan ZilchEntry)

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
		readerChan := make(chan ZilchEntry, 20)
		channelMap[reader.CountryCode] = readerChan
		go reader.Read(readerChan)
		go d.loadCountryData(reader.CountryCode, readerChan, distributionChannel)
		channels += 1
	}

	go d.finishDistributionChannels(distributionChannel, channels, start)

	return d, nil
}

func (d *Database) loadCountryData(countryCode string, channel chan ZilchEntry, distChannel chan map[uint32]int) {
	entries := make([]ZilchEntry, 0, 1000)
	distMap := make(map[uint32]int)

	type state_data struct {
		State     string
		StateName string
		ZipCodes  int
	}
	type country_data struct {
		Country     string
		CountryName string
		States      map[string]state_data
	}

	country := country_data{
		Country:     countryCode,
		CountryName: "",
		States:      make(map[string]state_data),
	}

	for entry := range channel {
		entries = append(entries, entry)
		distMap[entry.GetKey()] += 1

		if len(country.CountryName) == 0 {
			country.CountryName = entry.CountryName
		}

		if len(entry.State) > 0 {
			if state, found := country.States[entry.State]; found {
				country.States[entry.State] = state_data{
					State:     state.State,
					StateName: state.StateName,
					ZipCodes:  state.ZipCodes + 1,
				}
			} else {
				country.States[entry.State] = state_data{
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

func (d *Database) IsFullyLoaded() bool {
	return d.FullyLoaded
}

func (d *Database) finishDistributionChannels(distChannel chan map[uint32]int, totalChannels int, startTime time.Time) {
	channels := 0

	for distMap := range distChannel {
		channels += 1

		for key, total := range distMap {
			// key == 180090 = where equater meets prime meridian, not a real place
			if key != 180090 {
				lat, lon := GetLatitudeLongitudeFromKey(key)
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

func (d *Database) GetDistributions() []DistributionEntry {
	entries := make([]DistributionEntry, len(d.DistributionMap))
	idx := 0

	for _, entry := range d.DistributionMap {
		entries[idx] = entry
		idx += 1
	}

	sort.Sort(DistributionSorter(entries))
	return entries
}

func (d *Database) ExecQuery(queryParams map[string]string) (QueryResult, error) {
	if len(queryParams) == 0 {
		return QueryResult{}, errors.New("There are no query parameters")
	}
	var entries []ZilchEntry
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
	sort.Sort(ZilchSorter(entries))
	if total > max_entries {
		start = 0
		end = max_entries
		if page, page_found := queryParams["page"]; page_found {
			if p, perr := strconv.ParseUint(page, 10, 32); perr != nil {
				return QueryResult{}, perr
			} else {
				start = (int(p) - 1) * max_entries
				end = start + max_entries
				if start >= total {
					start = total
				}
				if end >= total {
					end = total
				}
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

func (d *Database) querySingleCountry(country string, queryParams map[string]string) ([]ZilchEntry, error) {
	entries := make([]ZilchEntry, 0, 20)
	if countryIndex, found := d.CountryIndexMap[country]; found {
		ch := make(chan ZilchEntry)
		go countryIndex.QueryIndex(queryParams, ch)
		for entry := range ch {
			entries = append(entries, entry)
		}
		return entries, nil
	} else {
		return entries, errors.New(fmt.Sprintf("No country %s found", country))
	}
}

func (d *Database) queryAllCountries(queryParams map[string]string) ([]ZilchEntry, error) {
	entries := make([]ZilchEntry, 0, 40)
	totalCountries := len(d.CountryIndexMap)
	completed := 0
	ch := make(chan ZilchEntry)
	for _, countryIndex := range d.CountryIndexMap {
		go countryIndex.QueryIndexNoClose(queryParams, ch)
	}
	for entry := range ch {
		if entry.Type == "EOL" {
			completed += 1
			if completed >= totalCountries {
				break
			}
		} else {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func (c CountryIndex) QueryIndex(queryParams map[string]string, ch chan ZilchEntry) {
	nonalphanum := regexp.MustCompile("[^0-9A-Za-z]")
	string_data := func(paramName string, params map[string]string) (string, bool) {
		if value, valExists := params[paramName]; valExists {
			return strings.ToLower(value), true
		} else {
			return "", false
		}
	}
	bounds_data := func(params map[string]string) ([]float32, bool) {
		b := make([]float32, 4)
		if val, valExists := params["Bounds"]; valExists {
			ba := strings.Split(val, ",")
			if len(ba) != 4 {
				return b, false
			} else {
				for i, bastring := range ba {
					if f, err := strconv.ParseFloat(bastring, 32); err == nil {
						b[i] = float32(f)
					} else {
						return b, false
					}
				}
				return b, true
			}
		} else {
			return b, false
		}
	}
	starts_with := func(expected, actual string) bool {
		return strings.Index(actual, expected) == 0
	}
	contains := func(expected, actual string) bool {
		return strings.Index(actual, expected) != -1
	}
	in_array := func(expected string, actual []string) bool {
		for _, val := range actual {
			if strings.Index(strings.ToLower(val), expected) != -1 {
				return true
			}
		}
		return false
	}
	in_bounds := func(bounds []float32, latitude, longitude float32) bool {
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
	bounds, boundsTest := bounds_data(queryParams)
	zipCode, zipCodeTest := string_data("ZipCode", queryParams)
	city, cityTest := string_data("City", queryParams)
	areaCode, areaCodeTest := string_data("AreaCode", queryParams)
	state, stateTest := string_data("State", queryParams)
	county, countyTest := string_data("County", queryParams)

	for _, entry := range c.Entries {
		if zipCodeTest {
			if !starts_with(nonalphanum.ReplaceAllString(zipCode, ""), strings.ToLower(nonalphanum.ReplaceAllString(entry.ZipCode, ""))) {
				continue
			}
		}
		if cityTest {
			if !contains(city, strings.ToLower(entry.City)) {
				valid := false
				if len(entry.AcceptableCities) > 0 && in_array(city, entry.AcceptableCities) {
					valid = true
				}
				if len(entry.UnacceptableCities) > 0 && in_array(city, entry.UnacceptableCities) {
					valid = true
				}
				if !valid {
					continue
				}
			}
		}
		if areaCodeTest {
			if !in_array(areaCode, entry.AreaCodes) {
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
			if !in_bounds(bounds, entry.Latitude, entry.Longitude) {
				continue
			}
		}
		ch <- entry
	}
	close(ch)
}

func (c CountryIndex) QueryIndexNoClose(queryParams map[string]string, ch chan ZilchEntry) {
	ch2 := make(chan ZilchEntry)
	go c.QueryIndex(queryParams, ch2)
	for entry := range ch2 {
		ch <- entry
	}
	ch <- ZilchEntry{Type: "EOL"}
}
