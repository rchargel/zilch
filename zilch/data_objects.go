package zilch

type DistributionEntry struct {
	Latitude  int16
	Longitude int16
	ZipCodes  uint32
}

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

type ZilchSorter []ZilchEntry

func (z ZilchEntry) GetKey() uint32 {
	lon := uint32(z.Longitude + 180)
	lat := uint32(z.Latitude + 90)
	return (lon * uint32(1000)) + lat
}

func (z ZilchSorter) Len() int           { return len(z) }
func (z ZilchSorter) Swap(i, j int)      { z[i], z[j] = z[j], z[i] }
func (z ZilchSorter) Less(i, j int) bool { return z[i].ZipCode < z[j].ZipCode }
