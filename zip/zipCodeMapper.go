package zip

import (
	"fmt"
	"net/http"
	"bufio"
)

type ZipCodeMapper struct {
	ZipCodeMap map[string]ZipCodeEntry
}

func NewZipCodeMapper() *ZipCodeMapper {
	mapper := &ZipCodeMapper{ make(map[string]ZipCodeEntry) }
	mapper.Init()
	return mapper
}

func (z *ZipCodeMapper) Init() {
	r := ZipReader{"./resources/us_zip_code_database.csv"}
	ch := make(chan ZipCodeEntry)

	go r.Read(ch)
	for entry := range ch {
		z.ZipCodeMap[entry.ZipCode] = entry
	}
}

func (z *ZipCodeMapper) PublishZipCode(zipCode string, format string, response http.ResponseWriter) {
	bw := bufio.NewWriter(response)
	if entry, ok := z.ZipCodeMap[zipCode]; ok {
		if format == "XML" {
			response.Header().Set("Content-Type", "text/xml")
			entry.WriteXml(bw)
		} else if format == "JSON" {
			response.Header().Set("Content-Type", "application/json")
			entry.WriteJson(bw)
		} else {
			response.Header().Set("Content-Type", "text/plain")
			bw.WriteString(fmt.Sprintf("{\"error\":\"Invalid Format: %s\"}", format))
		}
	} else {
		response.Header().Set("Content-Type", "text/plain")
		bw.WriteString(fmt.Sprintf("{\"error\":\"No entry found for zip: %s\"}", zipCode))
	}
	bw.Flush()
}
