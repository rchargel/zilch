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
	queryURL     string = "https://maps.googleapis.com/maps/api/geocode/json?sensor=false&region=%s&address=%s%s"
	exceedsLimit string = "OVER_QUERY_LIMIT"
	zipType      string = "postal_code"
	cityType     string = "locality"
	countyType   string = "sublocality_level_1"
	stateType    string = "administrative_area_level_1"
	countryType  string = "country"
)

type location struct {
	Lat float32
	Lng float32
}

type geometry struct {
	Location     location
	LocationType string
}

type addressComponent struct {
	LongName  string
	ShortName string
	Types     []string
}

type result struct {
	AddressComponents []addressComponent
	FormattedAddress  string
	Geometry          geometry
	Types             []string
}

type gresults struct {
	Results []result
	Status  string
}

func (r result) findAddressComponent(compType string) (addressComponent, error) {
	for _, component := range r.AddressComponents {
		for _, ctype := range component.Types {
			if ctype == compType {
				return component, nil
			}
		}
	}
	return addressComponent{}, errors.New("No address component for type: " + compType)
}

// Resolver is used to gather zip code data from the Google Maps API.
type Resolver struct {
	CountryCode string
	AppKey      string
	ZipCodes    []string
}

// FindResultList pipes a set of results to a channel.
func (r Resolver) FindResultList(resultChan chan gresults) {
	for _, zipCode := range r.ZipCodes {
		var addr string
		if len(r.AppKey) == 0 {
			addr = fmt.Sprintf(queryURL, strings.ToLower(r.CountryCode), zipCode, "")
		} else {
			addr = fmt.Sprintf(queryURL, strings.ToLower(r.CountryCode), zipCode, "&key="+r.AppKey)
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

// OutputCSV creates a CSV file that can later be imported by the Database.
func (r Resolver) OutputCSV(writer io.Writer) {
	start := time.Now()
	rchan := make(chan gresults)
	w := bufio.NewWriter(writer)

	r.outputHeader(w)
	go r.FindResultList(rchan)

	for result := range rchan {
		r.outputResults(w, result)
	}

	fmt.Printf("\nFinished CSV output in %v\n", time.Since(start))
	w.Flush()
}

func (r Resolver) getResults(addr string, attempt int) (gresults, error) {
	var results gresults
	if attempt > 3 {
		return results, errors.New("Reached maximum number of attempts for: " + addr)
	}
	resp, err := http.Get(addr)
	if err != nil {
		return results, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return results, err
	}
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
				if results.Status == exceedsLimit {
					return results, errors.New("Exceeds the maximum daily query limit")
				}
				time.Sleep(1000 * time.Millisecond)
				return r.getResults(addr, attempt+1)
			} else {
				return results, fmt.Errorf("To many results returned: %v", len(results.Results))
			}
		}
	}
}

func (r Resolver) outputHeader(w *bufio.Writer) {
	w.WriteString("country,zip,primary_city,state,state_name,county,country_name,latitude,longitude\n")
	w.Flush()
}

func (r Resolver) outputResults(w *bufio.Writer, results gresults) {
	var zip, city, state, stateName, county, country, countryName string
	var latitude, longitude float32

	var comp addressComponent
	var err error
	for _, res := range results.Results {
		if comp, err = res.findAddressComponent(zipType); err == nil {
			zip = comp.LongName
		}
		if comp, err = res.findAddressComponent(cityType); err == nil {
			city = comp.LongName
		}
		if comp, err = res.findAddressComponent(countyType); err == nil {
			county = comp.LongName
		}
		if comp, err = res.findAddressComponent(stateType); err == nil {
			state = strings.ToUpper(comp.ShortName)
			stateName = comp.LongName
		}
		if comp, err = res.findAddressComponent(countryType); err == nil {
			country = strings.ToUpper(comp.ShortName)
			countryName = comp.LongName
		}
		latitude = res.Geometry.Location.Lat
		longitude = res.Geometry.Location.Lng
	}
	w.WriteString(fmt.Sprintf("%v,\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",%v,%v\n", country, zip, city, state, stateName, county, countryName, latitude, longitude))
	w.Flush()
}
