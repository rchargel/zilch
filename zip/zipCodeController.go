package zip

import (
	"fmt"
	"strings"
	"bytes"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"github.com/hoisie/web"
)

type ZipCodeController struct {
	zipCodeMapper *ZipCodeMapper
}

type ErrorString string 

type UrlRequest string 

func (e ErrorString) Error() string {
	return string(e)
}

func (e ErrorString) ToJson() string {
	return fmt.Sprintf("{\"error\":\"%v\"}", string(e))
}

func GetCallback(ctx *web.Context) string {
	callback := ctx.Request.FormValue("callback")
	if len(callback) == 0 {
		callback = ctx.Request.FormValue("jsonp")
	}
	return callback
}

func AcceptGzip(ctx *web.Context) bool {
	h := ctx.Request.Header
	e := h.Get("Accept-encoding")
	return strings.Index(e, "gzip") >= 0
}

func WriteResponse(ctx *web.Context, resp string, format string) {
	callback := GetCallback(ctx)
	if format == "JS" { format = "JSON" }
	if len(callback) != 0 {
		format = "JS"
	}
	switch format {
	case "XML":
		ctx.ResponseWriter.Header().Set("Content-type", "text/xml; charset=utf-8")
	case "JSON":
		ctx.ResponseWriter.Header().Set("Content-type", "application/json; charset=utf-8")
	case "YAML":
		ctx.ResponseWriter.Header().Set("Content-type", "text/x-yaml; charset=utf-8")
	case "JS":
		ctx.ResponseWriter.Header().Set("Content-type", "application/javascript; charset=utf-8")
		resp = callback + "(" + resp + ");"
	} 
	if len(resp) > 10 && AcceptGzip(ctx) {
		ctx.ResponseWriter.Header().Set("Content-encoding", "gzip")
		gzw := gzip.NewWriter(ctx.ResponseWriter)
		bw := bufio.NewWriter(gzw)
		bw.WriteString(resp)
		bw.Flush()
		gzw.Flush()
		gzw.Close()
	} else {
		ctx.WriteString(resp)
	}
}

func Throw(e string) error {
	return ErrorString(e)
}

func (r UrlRequest) GetValue() string {
	s := string(r)
	if idx := strings.Index(s, "."); idx >= 0 {
		return s[0:idx]
	}
	return s
}

func(r UrlRequest) GetFormat() string {
	s := strings.ToUpper(string(r))
	if idx := strings.Index(s, "."); idx >= 0 {
		return s[idx + 1:]
	}
	return "JSON"
}

func (c ZipCodeController) ListCountries(ctx *web.Context) {
	m := make(map[string]int)
	for country, cmap := range c.zipCodeMapper.ZipCodeMap {
		m[country] = len(cmap)	
	}
	
	w := &bytes.Buffer{}
	e := json.NewEncoder(w)
	e.Encode(m)
	content := w.String()

	callback := GetCallback(ctx)
	if len(callback) == 0 {
		ctx.ResponseWriter.Header().Set("Content-type", "application/json; charset=utf-8")
	} else {
		ctx.ResponseWriter.Header().Set("Content-type", "application/javascript; charset=utf-8")
		content = callback + "(" + content + ");"
	}
	if AcceptGzip(ctx) {
 		ctx.ResponseWriter.Header().Set("Content-encoding", "gzip")
		gzw := gzip.NewWriter(ctx.ResponseWriter)
		bw := bufio.NewWriter(gzw)
		bw.WriteString(content)
		bw.Flush()
		gzw.Flush()
		gzw.Close()
	} else {
		ctx.WriteString(content)
	}
}

func (c ZipCodeController) LookupZipCode(ctx *web.Context, request string) {
	req := UrlRequest(request)
	q := make(map[string]string)
	q["ZipCode"] = req.GetValue()
	c.query(ctx, q, req.GetFormat())
}

func (c ZipCodeController) LookupAreaCode(ctx *web.Context, request string) {
	req := UrlRequest(request)
	q := make(map[string]string)
	q["AreaCode"] = req.GetValue()
	c.query(ctx, q, req.GetFormat())
}

func (c ZipCodeController) QueryReq(ctx *web.Context, format string) {
	p := make(map[string]string)
	for k, _ := range ctx.Request.Form {
		p[k] = ctx.Request.FormValue(k)
	}
	c.query(ctx, p, format)
}

func (c ZipCodeController) query(ctx *web.Context, params map[string]string, format string) {
	format = strings.ToUpper(format)
	if len(format) == 0 {
		format = "JSON"
	}
	entries, err := c.zipCodeMapper.Query(params)
	if err != nil {
		WriteResponse(ctx, ErrorString(err.Error()).ToJson(), "JSON")
		return
	}
	content, err := MarshalEntries(entries, format)
	if err != nil {
		WriteResponse(ctx, ErrorString(err.Error()).ToJson(), "JSON")
		return
	}
	WriteResponse(ctx, content, format)
} 

func (c ZipCodeController) RenderMap(ctx *web.Context) {
	ctx.ResponseWriter.Header().Set("Content-type", "image/png")
	RenderZipCodeMap(ctx.ResponseWriter, c.zipCodeMapper, false)
}

func (c ZipCodeController) RenderTransparentMap(ctx *web.Context) {
	ctx.ResponseWriter.Header().Set("Content-type", "image/png")
	RenderZipCodeMap(ctx.ResponseWriter, c.zipCodeMapper, true)
}

func NewZipCodeController() ZipCodeController {
	return ZipCodeController{NewZipCodeMapper()}
}
