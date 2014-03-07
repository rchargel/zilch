package zip

import (
	"fmt"
)

type ZipCodeMapper struct {
	ZipCodeMap map[string]ZipCodeEntry
}

func NewZipCodeMapper() *ZipCodeMapper {
	mapper := &ZipCodeMapper{ make(map[string]ZipCodeEntry) }
	mapper.Init()
	return mapper
}

func (z *ZipCodeMapper) putdata(ch chan ZipCodeEntry) {
	for entry := range ch {
		z.ZipCodeMap[entry.ZipCode] = entry
	}
}

func (z *ZipCodeMapper) Init() {
	r := ZipReader{"./resources/us_zip_code_database.csv"}
	ch := make(chan ZipCodeEntry)

	go r.Read(ch)
	go z.putdata(ch)
}

func (z *ZipCodeMapper) GetEntryByZipCode(zipCode string) (ZipCodeEntry, error) {
	if entry, ok := z.ZipCodeMap[zipCode]; ok {
		return entry, nil
	}
	return ZipCodeEntry{}, Throw(fmt.Sprintf("No entry found for zip %s", zipCode))
}
