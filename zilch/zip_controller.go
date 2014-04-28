package zilch

import (
	"fmt"
	"github.com/hoisie/web"
)

type ZipCodeController struct {
	database *Database
}

func (c ZipCodeController) Query(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	if queryResult, err := c.database.ExecQuery(writer.GetQuery()); err == nil {
		writer.SendQueryResponse(queryResult)
	} else {
		writer.SendError(err)
	}
}

func (c ZipCodeController) GetDistribution(ctx *web.Context) {
	writer := ResponseWriter{ctx, "JSON"}
	writer.SendStandardResponse(c.database.GetDistributions())
}

func (c ZipCodeController) GetCountries(ctx *web.Context) {
	writer := ResponseWriter{ctx, "JSON"}
	length := len(c.database.CountryIndexMap)
	countries := make([]map[string]string, length, length)
	index := 0
	for _, entries := range c.database.CountryIndexMap {
		if index >= length {
			break
		}
		mp := make(map[string]string, 2)
		mp["CountryCode"] = entries.CountryCode
		mp["Count"] = fmt.Sprintf("%v", len(entries.Entries))
		countries[index] = mp
		index += 1
	}
	writer.SendStandardResponse(countries)
}
