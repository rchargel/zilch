package zip

import (
	"time"
	"fmt"
	"sort"
	"strings"
)

type QueryResult struct {
	entry ZipCodeEntry
	error string
} 

type ZipCodeMapper struct {
	ZipCodeMap map[string]map[string]ZipCodeEntry
}

func NewZipCodeMapper() *ZipCodeMapper {
	mapper := &ZipCodeMapper{ make(map[string]map[string]ZipCodeEntry) }
	mapper.Init()
	return mapper
}

func (z *ZipCodeMapper) putdata(ch chan ZipCodeEntry, t time.Time) {
	count := uint32(0)
	for entry := range ch {
		_, ok := z.ZipCodeMap[entry.Country] 
		if ok {
			if _, zipExists := z.ZipCodeMap[entry.Country][entry.ZipCode]; zipExists {
				oldEntry := z.ZipCodeMap[entry.Country][entry.ZipCode].AddCity(entry.City)
				z.ZipCodeMap[entry.Country][entry.ZipCode] = oldEntry
			} else {
				z.ZipCodeMap[entry.Country][entry.ZipCode] = entry
				count = count + 1
			}
		} else {
			z.ZipCodeMap[entry.Country] = make(map[string]ZipCodeEntry)
			z.ZipCodeMap[entry.Country][entry.ZipCode] = entry
			count = count + 1
		}
	}
	fmt.Printf("Stored %v records in database\n", count)
	end := time.Now()
	fmt.Printf("Data loaded in %v seconds\n",end.Sub(t))
}

func (z *ZipCodeMapper) Init() {
	start := time.Now()
	r := ZipCodeDB{"./resources/"}
	ch := make(chan ZipCodeEntry)

	go r.LoadAll(ch)
	go z.putdata(ch, start)
}

func QueryMap(data map[string]ZipCodeEntry, params map[string]string, ch chan QueryResult) {
	if len(params) == 0 {
		ch <- QueryResult{ZipCodeEntry{}, "EOL"}
		return
	}
	zipCode, zipOk := params["ZipCode"]
	city, cityOk := params["City"]
	areaCode, areaCodeOk := params["AreaCode"]
	if cityOk {
		city = strings.ToLower(city)
	}
	state, stateOk := params["State"]
	if stateOk {
		state = strings.ToUpper(state)
	}
	county, countyOk := params["County"]
	if countyOk {
		county = strings.ToLower(county)
	}
	for _, entry := range data {
		if zipOk {
			if zipCode != entry.ZipCode {
				continue
			}
		}
		if stateOk {
			if state != entry.State {
				continue
			}
		}
		if countyOk {
			if !equalsIgnoreCase(county, entry.County) {
				continue
			}
		}
		if areaCodeOk {
			fail := true
			for _, ac := range entry.AreaCodes {
				if ac == areaCode {
					fail = false
					break
				}
			}
			if fail {
				continue
			}
		}
		if cityOk {
			if !equalsIgnoreCase(city, entry.City) {
				allow := false
				for _, accCity := range entry.AcceptableCities {
					if equalsIgnoreCase(city, accCity) {
						allow = true
						break
					}
				}
				if !allow {
					for _, unaccCity := range entry.UnacceptableCities {
						if equalsIgnoreCase(city, unaccCity) {
							allow = true
							break
						}
					}
				}
				if !allow {
					continue
				}
			}
		}
		ch <- QueryResult{entry, ""}
	}
	ch <- QueryResult{ZipCodeEntry{}, "EOL"}
}

func equalsIgnoreCase(expected, test string) bool {
	return strings.Index(strings.ToLower(test),expected) >= 0
}

func (z *ZipCodeMapper) Query(params map[string]string) ([]ZipCodeEntry, error) {
	ch := make(chan QueryResult)
	country, cfound := params["Country"]
	countries := 0
	entries := make([]ZipCodeEntry, 0, 100)
	if (cfound) {
		if data, ok := z.ZipCodeMap[country]; ok {
			go QueryMap(data, params, ch)
			countries = 1
		} else {
			return entries, Throw(fmt.Sprintf("No country %v in database", country))
		}
	} else {
		for _, data := range z.ZipCodeMap {
			go QueryMap(data, params, ch)
			countries = countries + 1
		}
	}
	returns := 0
	for result := range ch {
		if result.error == "EOL" {
			returns = returns + 1
			if returns == countries {
				break
			}
		} else {
			entries = append(entries, result.entry)
		}
	}
	sort.Sort(ZipSorter(entries))
	if len(entries) > 1000 {
		entries = entries[0:1000]
	}
	return entries, nil
}
