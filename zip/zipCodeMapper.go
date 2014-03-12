package zip

import (
	"strconv"
	"time"
	"fmt"
	"sort"
	"strings"
)

type QueryResult struct {
	entry ZipCodeEntry
	error string
} 

type DistributionResult struct {
	Latitude int16
	Longitude int16
	ZipCodes uint32
}

type ZipCodeMapper struct {
	ZipCodeMap map[string]map[string]ZipCodeEntry
	Dist map[int64]DistributionResult
}

type ZipQueryResult struct {
	Entries []ZipCodeEntry
	Results int
	Found int
	StartIndex int
	EndIndex int
}

func (d DistributionResult) incr() DistributionResult {
	return DistributionResult{d.Latitude, d.Longitude, d.ZipCodes + uint32(1)}
}

func NewZipCodeMapper() *ZipCodeMapper {
	mapper := &ZipCodeMapper{ make(map[string]map[string]ZipCodeEntry), make(map[int64]DistributionResult) }
	mapper.Init()
	return mapper
}

func (z *ZipCodeMapper) putdata(ch chan ZipCodeEntry, t time.Time) {
	count := uint32(0)
	for entry := range ch {
		distKey := (int64(entry.Latitude) * int64(360)) + int64(entry.Longitude)
		_, ok := z.ZipCodeMap[entry.Country] 
		if ok {
			if _, zipExists := z.ZipCodeMap[entry.Country][entry.ZipCode]; zipExists {
				oldEntry := z.ZipCodeMap[entry.Country][entry.ZipCode].AddCity(entry)
				z.ZipCodeMap[entry.Country][entry.ZipCode] = oldEntry
			} else {
				z.ZipCodeMap[entry.Country][entry.ZipCode] = entry
				count = count + 1
				if _, distOk := z.Dist[distKey]; !distOk {
					z.Dist[distKey] = DistributionResult{int16(entry.Latitude), int16(entry.Longitude), 1}
				} else {
					z.Dist[distKey] = z.Dist[distKey].incr()
				}
			}
		} else {
			z.ZipCodeMap[entry.Country] = make(map[string]ZipCodeEntry)
			z.ZipCodeMap[entry.Country][entry.ZipCode] = entry
			count = count + 1
			if _, distOk := z.Dist[distKey]; !distOk {
				z.Dist[distKey] = DistributionResult{int16(entry.Latitude), int16(entry.Longitude), 1}
			} else {
				z.Dist[distKey] = z.Dist[distKey].incr()
			}
		}
	}
	delete(z.Dist,int64(0))
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
	bounds := make([]float64,4)
	boundString, boundOk := params["Bounds"]
	if boundOk {
		b := strings.Split(boundString, ",")
		if len(b) != 4 {
			boundOk = false
		} else {
			for i, bString := range b {
				f, err := strconv.ParseFloat(bString, 64)	
				if err != nil {
					boundOk = false
					break
				}
				bounds[i] = f
			}
		}
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
		if boundOk {
			if entry.Latitude == 0 && entry.Longitude == 0 { continue }
			if entry.Latitude > bounds[0] || entry.Latitude < bounds[2] || entry.Longitude < bounds[1] || entry.Longitude > bounds[3] {
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

func (z *ZipCodeMapper) DistributionMap() ([]DistributionResult) {
	results := make([]DistributionResult, len(z.Dist))
	i := 0
	for _, d := range z.Dist {
		results[i] = d
		i++
	}
	return results
}

func (z *ZipCodeMapper) Query(params map[string]string) (ZipQueryResult, error) {
	ch := make(chan QueryResult)
	country, cfound := params["Country"]
	countries := 0
	entries := make([]ZipCodeEntry, 0, 100)
	if (cfound) {
		if data, ok := z.ZipCodeMap[country]; ok {
			go QueryMap(data, params, ch)
			countries = 1
		} else {
			return ZipQueryResult{Entries: entries, Results: 0, Found: 0, StartIndex: 0, EndIndex: 0}, Throw(fmt.Sprintf("No country %v in database", country))
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
	total := len(entries)
	start := 0
	end := len(entries)
	sort.Sort(ZipSorter(entries))
	max := 1000
	if len(entries) > max {
		start = 0
		end = 200
		if page, pageOk := params["page"]; pageOk {
			if p, err := strconv.ParseInt(page, 10, 32); err == nil {
				start = (int(p) - 1) * max
				end = start + max
				if start >= len(entries) {
					start = len(entries)
				}
				if end > len(entries) {
					end = len(entries)
				}
			}
		}
		entries = entries[start:end]
	}
	return ZipQueryResult{Entries: entries, Results: len(entries), Found: total, StartIndex: start+1, EndIndex: end}, nil
}
