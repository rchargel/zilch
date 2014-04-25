package zilch

type DistributionEntry struct {
	Latitude  int16
	Longitude int16
	ZipCodes  uint32
}

type DistributionSorter []DistributionEntry

type ZilchEntry struct {
	ZipCode            string
	Type               string
	City               string
	AcceptableCities   []string
	UnacceptableCities []string
	County             string
	State              string
	Country            string
	TimeZone           string
	AreaCodes          []string
	Latitude           float32
	Longitude          float32
}

type QueryResult struct {
	ResultsReturned int
	TotalFound      int
	StartIndex      int
	EndIndex        int
	ZipCodeEntries  []ZilchEntry
}

type ZilchSorter []ZilchEntry

func (z ZilchEntry) GetKey() uint32 {
	return GetKeyFromLatitudeLongitude(z.Latitude, z.Longitude)
}

func (z ZilchSorter) Len() int           { return len(z) }
func (z ZilchSorter) Swap(i, j int)      { z[i], z[j] = z[j], z[i] }
func (z ZilchSorter) Less(i, j int) bool { return z[i].ZipCode < z[j].ZipCode }

func (d DistributionSorter) Len() int           { return len(d) }
func (d DistributionSorter) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DistributionSorter) Less(i, j int) bool { return d[i].ZipCodes > d[j].ZipCodes }

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