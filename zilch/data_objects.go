package zilch

type DistributionEntry struct {
	Latitude  int16
	Longitude int16
	ZipCodes  uint32
}

type DistributionMarshaller []DistributionEntry

type DistributionSorter []DistributionEntry

type CountryMarshaller map[string]int

type ZilchEntry struct {
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

type StateEntry struct {
	State     string
	StateName string
	ZipCodes  uint32
}

type CountryEntry struct {
	Country     string
	CountryName string
	States      []StateEntry
}

type CountryEntryMarshaller []CountryEntry

type QueryResult struct {
	ResultsReturned int
	TotalFound      int
	StartIndex      int
	EndIndex        int
	ZipCodeEntries  []ZilchEntry
}

type ZilchSorter []ZilchEntry

type StateSorter []StateEntry

type CountrySorter []CountryEntry

func (z ZilchEntry) GetKey() uint32 {
	return GetKeyFromLatitudeLongitude(z.Latitude, z.Longitude)
}

func (z ZilchSorter) Len() int      { return len(z) }
func (z ZilchSorter) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z ZilchSorter) Less(i, j int) bool {
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

func GetKeyFromLatitudeLongitude(latitude, longitude float32) uint32 {
	lon := uint32(longitude + 180)
	lat := uint32(latitude + 90)
	return (lon * uint32(1000)) + lat
}

func GetLatitudeLongitudeFromKey(key uint32) (int16, int16) {
	longitude := int16((key / 1000) - 180)
	latitude := int16((key % 1000) - 90)

	return latitude, longitude
}
