package zilch

import (
	"github.com/hoisie/web"
)

// ZipCodeController main controller for this application.
type ZipCodeController struct {
	database *Database
}

// Query controller method to respond to a query.
func (c ZipCodeController) Query(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	if queryResult, err := c.database.ExecQuery(writer.getQuery()); err == nil {
		writer.SendQueryResponse(queryResult)
	} else {
		writer.SendError(err)
	}
}

// GetDistribution controller method to get the distribution response.
func (c ZipCodeController) GetDistribution(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	writer.SendDistributionResponse(c.database.GetDistributions())
}

// GetCountries controller method to get the list of countries and country
// details.
func (c ZipCodeController) GetCountries(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	writer.SendCountryListResponse(c.database.CountryList)
}
