package zilch

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"strings"

	"github.com/hoisie/web"
)

// ResponseWriter writes output to the servers response context.
type ResponseWriter struct {
	ctx    *web.Context
	format string
}

func (writer ResponseWriter) getQuery() map[string]string {
	query := make(map[string]string)
	for key := range writer.ctx.Request.Form {
		query[key] = writer.ctx.Request.FormValue(key)
	}
	return query
}

// SendError sends the supplied error to the user via an HTTP 500 error.
func (writer ResponseWriter) SendError(err error) {
	writer.ctx.Abort(500, err.Error())
}

// SendDistributionResponse sends the response as a list of DistributionEntry
// objects.
func (writer ResponseWriter) SendDistributionResponse(d []DistributionEntry) {
	format := "JSON"
	if len(writer.format) > 0 {
		format = strings.ToUpper(writer.format)
	}
	dm := DistributionMarshaller(d)
	response, err := dm.Marshal(format)

	if err != nil {
		writer.SendError(err)
	} else {
		switch format {
		case "XML":
			writer.ctx.ContentType("text/xml; charset=utf-8")
		default:
			callback := writer.getJsonpCallback()
			if len(callback) > 0 {
				writer.ctx.ContentType("application/javascript; charset=utf-8")
				response = callback + "(" + response + ");"
			} else {
				writer.ctx.ContentType("application/json; charset=utf-8")
			}
		}
		writer.compressionFilter(response)
	}
}

// SendCountryListResponse sends the response as a list of CountryEntry objects.
func (writer ResponseWriter) SendCountryListResponse(c []CountryEntry) {
	format := "JSON"
	if len(writer.format) > 0 {
		format = strings.ToUpper(writer.format)
	}
	cm := CountryEntryMarshaller(c)
	response, err := cm.Marshal(format)

	if err != nil {
		writer.SendError(err)
	} else {
		switch format {
		case "XML":
			writer.ctx.ContentType("text/xml; charset=utf-8")
		default:
			callback := writer.getJsonpCallback()
			if len(callback) > 0 {
				writer.ctx.ContentType("application/javascript; charset=utf-8")
				response = callback + "(" + response + ");"
			} else {
				writer.ctx.ContentType("application/json; charset=utf-8")
			}
		}
		writer.compressionFilter(response)
	}
}

// SendQueryResponse sends the response to a query.
func (writer ResponseWriter) SendQueryResponse(queryResult QueryResult) {
	if response, err := writer.marshalQueryResponse(queryResult); err == nil {
		writer.compressionFilter(response)
	} else {
		writer.SendError(err)
	}
}

func (writer ResponseWriter) marshalQueryResponse(queryResult QueryResult) (string, error) {
	format := "JSON"
	if len(writer.format) > 0 {
		format = strings.ToUpper(writer.format)
	}
	response, err := queryResult.Marshal(format)

	if err != nil {
		return response, err
	}

	switch format {
	case "XML":
		writer.ctx.ContentType("text/xml; charset=utf-8")
	case "YAML":
		writer.ctx.ContentType("text/yaml; charset=utf-8")
	default:
		callback := writer.getJsonpCallback()
		if len(callback) > 0 {
			writer.ctx.ContentType("application/javascript; charset=utf-8")
			response = callback + "(" + response + ");"
		} else {
			writer.ctx.ContentType("application/json; charset=utf-8")
		}
	}
	return response, err
}

func (writer ResponseWriter) compressionFilter(response string) {
	header := writer.ctx.Request.Header
	encoding := header.Get("Accept-encoding")

	if strings.Index(encoding, "gzip") != -1 && len(response) > 10 {
		writer.ctx.ResponseWriter.Header().Add("Content-encoding", "gzip")
		gzw := gzip.NewWriter(writer.ctx.ResponseWriter)
		writer := bufio.NewWriter(gzw)
		writer.WriteString(response)
		writer.Flush()
		gzw.Flush()
		gzw.Close()
	} else {
		writer.ctx.ResponseWriter.Header().Add("Content-length", fmt.Sprintf("%v", len(response)))
		writer.ctx.WriteString(response)
	}
}

func (writer ResponseWriter) getJsonpCallback() string {
	callback := writer.ctx.Request.FormValue("callback")
	if len(callback) == 0 {
		callback = writer.ctx.Request.FormValue("jsonp")
	}
	return callback
}
