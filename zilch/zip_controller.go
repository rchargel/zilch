package zilch

import (
	"github.com/hoisie/web"
)

type ZipCodeController struct {
	database *Database
}

func (c ZipCodeController) Query(ctx *web.Context, format string) {
	writer := ResponseWriter{ctx, format}
	if queryResult, err := database.ExecQuery(writer.GetQuery()); err == nil {
		writer.SendQueryResponse(queryResult)
	} else {
		writer.SendError(err)
	}
}

func (c ZipCodeController) Query(ctx *web.Context) {
	writer := ResponseWriter(ctx, "JSON")
	writer.SendDistributionResponse(database.GetDistributions())
}
