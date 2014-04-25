package zilch

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"strings"
)

type ResponseWriter struct {
	ctx    *web.Context
	format string
}

func (writer ResponseWriter) GetQuery() map[string]string {
	query := make(map[string]string)
	for key, _ := range writer.ctx.Request.Form {
		query[key] = writer.ctx.Request.FormValue(k)
	}
	return query
}

func (writer ResponseWriter) SendError(err error) {
	writer.ctx.Abort(500, err.Error())
}

func (writer ResponseWriter) SendDistributionResponse(distribution []DistributionEntry) {
	if response, err := writer.marshalDistributionResponse(distribution); err == nil {
		writer.compressionFilter(response)
	} else {
		writer.SendError(err)
	}
}

func (writer ResponseWriter) SendQueryResponse(queryResult QueryResult) {
	if response, err := writer.marshalQueryResponse(queryResult); err == nil {
		writer.compressionFilter(response)
	} else {
		writer.SendError(err)
	}
}

func (writer ResponseWriter) marshalDistributionResponse(distribution []DistributionEntry) (string, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(&distribution); err != nil {
		return "", err
	} else {
		response := buf.String()
		callback := writer.getJsonpCallback()
		if len(callback) > 0 {
			writer.ctx.ContentType("application/javascript; charset=utf-8")
			return callback + "(" + response + ");", nil
		} else {
			writer.ctx.ContentType("application/json; charset=utf-8")
			return response, nil
		}
	}
}

func (writer ResponseWriter) marshalQueryResponse(queryResult QueryResult) (string, error) {
	format := strings.ToUpper(writer.format)
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
