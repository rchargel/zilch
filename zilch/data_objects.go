package zilch

// DistributionEntry is used to track the number of
// zip codes located in a latitude/longitude block.
type DistributionEntry struct {
	Latitude  int16
	Longitude int16
	ZipCodes  uint32
}

// DistributionMarshaller is used to marshal
// DistributionEntry objects.
type DistributionMarshaller []DistributionEntry

// DistributionSorter is used to sort
// DistributionEntry objects.
type DistributionSorter []DistributionEntry

// CountryMarshaller is used to marshal a map of
// country codes to the number of zip codes in that country.
type CountryMarshaller map[string]int

// ZipEntry is an object which holds the details of a single
// zip code.
type ZipEntry struct {
	ZipCode            string
	Type               string
	City               string
	AcceptableCities   []string
	UnacceptableCities []string
	County             string
	State              string
	StateName          string
	Country            string
	CountryName        string
	TimeZone           string
	AreaCodes          []string
	Latitude           float32
	Longitude          float32
}

// StateEntry is an object which maps the state information, to the
// number of zip codes in that state.
type StateEntry struct {
	State     string
	StateName string
	ZipCodes  uint32
}

// CountryEntry is an object which maps a country to a set of states/provinces.
type CountryEntry struct {
	Country     string
	CountryName string
	States      []StateEntry
}

// CountryEntryMarshaller is used to marshal a CountryEntry.
type CountryEntryMarshaller []CountryEntry

// QueryResult holds the result of a zip code query.
type QueryResult struct {
	ResultsReturned int
	TotalFound      int
	StartIndex      int
	EndIndex        int
	ZipCodeEntries  []ZipEntry
}

// ZipSorter sorts the ZipEntry slice.
type ZipSorter []ZipEntry

// StateSorter sorts the StateEntry slice.
type StateSorter []StateEntry

// CountrySorter sorts the CountryEntry slice.
type CountrySorter []CountryEntry

// GetKey gets a unique identifier for a zip code's latitude and longitude.
func (z ZipEntry) GetKey() uint32 {
	return getKeyFromLatitudeLongitude(z.Latitude, z.Longitude)
}

func (z ZipSorter) Len() int      { return len(z) }
func (z ZipSorter) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z ZipSorter) Less(i, j int) bool {
	if z[i].Country != z[j].Country {
		return z[i].Country < z[j].Country
	}
	return z[i].ZipCode < z[j].ZipCode
}

func (d DistributionSorter) Len() int           { return len(d) }
func (d DistributionSorter) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DistributionSorter) Less(i, j int) bool { return d[i].ZipCodes < d[j].ZipCodes }

func (s StateSorter) Len() int           { return len(s) }
func (s StateSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StateSorter) Less(i, j int) bool { return s[i].State < s[j].State }

func (c CountrySorter) Len() int           { return len(c) }
func (c CountrySorter) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CountrySorter) Less(i, j int) bool { return c[i].CountryName < c[j].CountryName }

func getKeyFromLatitudeLongitude(latitude, longitude float32) uint32 {
	lon := uint32(longitude + 180)
	lat := uint32(latitude + 90)
	return (lon * uint32(1000)) + lat
}

func getLatitudeLongitudeFromKey(key uint32) (int16, int16) {
	longitude := int16((key / 1000) - 180)
	latitude := int16((key % 1000) - 90)

	return latitude, longitude
}
