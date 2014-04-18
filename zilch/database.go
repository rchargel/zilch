package zilch

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"
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
	start := time.Now()
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

//func (d *Database) ExecQuery(queryParams map[string]string)
