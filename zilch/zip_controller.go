package zilch

import (
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

func (c ZipCodeController) GetDistribution(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	writer.SendDistributionResponse(c.database.GetDistributions())
}

func (c ZipCodeController) GetCountries(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	countries := make(map[string]int)
	for _, entries := range c.database.CountryIndexMap {
		countries[entries.CountryCode] = len(entries.Entries)
	}
	writer.SendCountryListResponse(countries)
}
