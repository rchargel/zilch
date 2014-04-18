package zilch

import (
	"fmt"
	"io/ioutil"
)

type CountryIndex struct {
	CountryCode string
	Entries     []ZilchEntry
}

type Database struct {
	CountryIndexMap map[string]CountryIndex
	DistributionMap map[uint32]DistributionEntry
	FullyLoaded     bool
}

func NewDatabase(filedir string) (*Database, error) {
	d := &Database{
		CountryIndexMap: make(map[string]CountryIndex),
		DistributionMap: make(map[uint32]DistributionEntry),
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
		filepath := filedir + file.Name()
		reader := CreateReader(filepath)
		readerChan := make(chan ZilchEntry, 20)
		channelMap[reader.CountryCode] = readerChan
		go reader.Read(readerChan)
		go d.loadCountryData(reader.CountryCode, readerChan, distributionChannel)
		channels += 1
	}

	go d.finishDistributionChannels(distributionChannel, channels)

	return d, nil
}

func (d *Database) loadCountryData(countryCode string, channel chan ZilchEntry, distChannel chan map[uint32]int) {
	entries := make([]ZilchEntry, 0, 1000)
	distMap := make(map[uint32]int)

	for entry := range channel {
		entries = append(entries, entry)
		distMap[entry.GetKey()] += 1
	}

	d.CountryIndexMap[countryCode] = CountryIndex{
		CountryCode: countryCode,
		Entries:     entries,
	}

	distChannel <- distMap
}

func (d *Database) finishDistributionChannels(distChannel chan map[uint32]int, totalChannels int) {
	channels := 0

	for distMap := range distChannel {
		channels += 1

		for key, total := range distMap {
			if _, found := d.DistributionMap[key]; !found {
				lat, lon := GetLatitudeLongitudeFromKey(key)
				d.DistributionMap[key] = DistributionEntry{
					Latitude:  lat,
					Longitude: lon,
					ZipCodes:  uint32(total),
				}
			}
		}

		if channels == totalChannels {
			break
		}
	}

	d.FullyLoaded = true
	fmt.Println("Finished reading database")
}
