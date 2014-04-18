package zilch

import (
	"io/ioutil"
)

type CountryIndex struct {
	CountryCode string
	Entries     []ZilchEntry
}

type Database struct {
	CountryIndexMap map[string]CountryIndex
	DistributionMap map[uint32]DistributionEntry
}

func NewDatabase(filedir string) (*Database, error) {
	d := &Database{
		CountryIndexMap: make(map[string]CountryIndex),
		DistributionMap: make(map[uint32]DistributionEntry),
	}

	channelMap := make(map[string]chan ZilchEntry)

	// start reading data
	files, err := ioutil.ReadDir(filedir)
	if err != nil {
		return d, err
	}

	for _, file := range files {
		filepath := filedir + file.Name()
		reader := CreateReader(filepath)
		readerChan := make(chan ZilchEntry)
		channelMap[reader.CountryCode] = readerChan
		go reader.Read(readerChan)
		go d.loadCountryData(reader.CountryCode, readerChan)
	}

	return d, nil
}

func (d *Database) loadCountryData(countryCode string, channel chan ZilchEntry) {
	/*
		var cIndex CountryIndex
		cIndex = CountryIndex{
			CountryCode: countryCode,
			Entries:     make([]ZilchEntry, 0, 1000),
		}

		for entry := range channel {

		}
	*/
}
