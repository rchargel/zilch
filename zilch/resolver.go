package zilch

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	QUERY_URL     string = "https://maps.googleapis.com/maps/api/geocode/json?sensor=false&region=%s&address=%s%s"
	EXCEEDS_LIMIT string = "OVER_QUERY_LIMIT"
	ZIP_TYPE      string = "postal_code"
	CITY_TYPE     string = "locality"
	COUNTY_TYPE   string = "sublocality_level_1"
	STATE_TYPE    string = "administrative_area_level_1"
	COUNTRY_TYPE  string = "country"
)

type Location struct {
	Lat float32
	Lng float32
}

type Geometry struct {
	Location      Location
	Location_Type string
}

type AddressComponent struct {
	Long_Name  string
	Short_Name string
	Types      []string
}

type Result struct {
	Address_Components []AddressComponent
	Formatted_Address  string
	Geometry           Geometry
	Types              []string
}

type GResults struct {
	Results []Result
	Status  string
}

func (r Result) FindAddressComponent(compType string) (AddressComponent, error) {
	for _, component := range r.Address_Components {
		for _, ctype := range component.Types {
			if ctype == compType {
				return component, nil
			}
		}
	}
	return AddressComponent{}, errors.New("No address component for type: " + compType)
}

type Resolver struct {
	CountryCode string
	AppKey      string
	ZipCodes    []string
}

func (r Resolver) FindResultList(resultChan chan GResults) {
	for _, zipCode := range r.ZipCodes {
		var addr string
		if len(r.AppKey) == 0 {
			addr = fmt.Sprintf(QUERY_URL, strings.ToLower(r.CountryCode), zipCode, "")
		} else {
			addr = fmt.Sprintf(QUERY_URL, strings.ToLower(r.CountryCode), zipCode, "&key="+r.AppKey)
		}
		if results, err := r.getResults(addr, 1); err == nil {
			resultChan <- results
		} else {
			fmt.Println(err)
			if strings.Index(err.Error(), "Exceed") != -1 {
				break
			}
		}
	}
	close(resultChan)
}

func (r Resolver) OutputCSV(writer io.Writer) {
	start := time.Now()
	rchan := make(chan GResults)
	w := bufio.NewWriter(writer)

	r.outputHeader(w)
	go r.FindResultList(rchan)

	for result := range rchan {
		r.outputResults(w, result)
	}

	fmt.Printf("\nFinished CSV output in %v\n", time.Since(start))
	w.Flush()
}

func (r Resolver) getResults(addr string, attempt int) (GResults, error) {
	var results GResults
	if attempt > 5 {
		return results, errors.New("Exceeded the maximum number of attempts for: " + addr)
	}
	if resp, err := http.Get(addr); err != nil {
		return results, err
	} else {
		content, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return results, err
		} else {
			dec := json.NewDecoder(strings.NewReader(string(content)))
			for {
				if err := dec.Decode(&results); err == io.EOF {
					return results, nil
				} else if err != nil {
					return results, err
				} else {
					if len(results.Results) == 1 {
						return results, nil
					} else if len(results.Results) == 0 {
						if results.Status == EXCEEDS_LIMIT {
							return results, errors.New("Exceeds the maximum daily query limit")
						}
						time.Sleep(1000 * time.Millisecond)
						return r.getResults(addr, attempt+1)
					} else {
						return results, errors.New(fmt.Sprintf("Too many results returned: %v", len(results.Results)))
					}
				}
			}
		}
	}
}

func (r Resolver) outputHeader(w *bufio.Writer) {
	w.WriteString("country,zip,primary_city,state,state_name,county,country_name,latitude,longitude\n")
	w.Flush()
}

func (r Resolver) outputResults(w *bufio.Writer, results GResults) {
	var zip, city, state, stateName, county, country, countryName string
	var latitude, longitude float32

	var comp AddressComponent
	var err error
	for _, res := range results.Results {
		if comp, err = res.FindAddressComponent(ZIP_TYPE); err == nil {
			zip = comp.Long_Name
		}
		if comp, err = res.FindAddressComponent(CITY_TYPE); err == nil {
			city = comp.Long_Name
		}
		if comp, err = res.FindAddressComponent(COUNTY_TYPE); err == nil {
			county = comp.Long_Name
		}
		if comp, err = res.FindAddressComponent(STATE_TYPE); err == nil {
			state = strings.ToUpper(comp.Short_Name)
			stateName = comp.Long_Name
		}
		if comp, err = res.FindAddressComponent(COUNTRY_TYPE); err == nil {
			country = strings.ToUpper(comp.Short_Name)
			countryName = comp.Long_Name
		}
		latitude = res.Geometry.Location.Lat
		longitude = res.Geometry.Location.Lng
	}
	w.WriteString(fmt.Sprintf("%v,\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",%v,%v\n", country, zip, city, state, stateName, county, countryName, latitude, longitude))
	w.Flush()
}
