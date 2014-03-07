package zip

import (
	"fmt"
)

type ZipCodeMapper struct {
	ZipCodeMap map[string]ZipCodeEntry
	AreaCodeMap map[string][]ZipCodeEntry
}

func NewZipCodeMapper() *ZipCodeMapper {
	mapper := &ZipCodeMapper{ make(map[string]ZipCodeEntry),
		make(map[string][]ZipCodeEntry) }
	mapper.Init()
	return mapper
}

func (z *ZipCodeMapper) putdata(ch chan ZipCodeEntry) {
	for entry := range ch {
		z.ZipCodeMap[entry.ZipCode] = entry
		for _, areaCode := range entry.AreaCodes {
			list, ok := z.AreaCodeMap[areaCode]
			if ok {
				z.AreaCodeMap[areaCode] = append(list, entry)
			} else {
				z.AreaCodeMap[areaCode] = make([]ZipCodeEntry, 1)
				z.AreaCodeMap[areaCode][0] = entry
			}
		}
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

func (z *ZipCodeMapper) GetEntriesByAreaCode(areaCode string) ([]ZipCodeEntry, error) {
	if entries, ok := z.AreaCodeMap[areaCode]; ok {
		return entries, nil
	}
	return make([]ZipCodeEntry, 0), Throw(fmt.Sprintf("No entries found for area code %s", areaCode))
}
