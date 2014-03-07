package zip

import (
	"fmt"
	"strings"
	"bufio"
	"compress/gzip"
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
	if len(callback) != 0 {
		format = "JSONP"
	}
	switch format {
	case "XML":
		ctx.ResponseWriter.Header().Set("Content-type", "text/xml")
	case "JSON":
		ctx.ResponseWriter.Header().Set("Content-type", "application/json")
	case "JSONP":
		ctx.ResponseWriter.Header().Set("Content-type", "application/javascript")
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

func (c ZipCodeController) lookupZipCode(ctx *web.Context, request string) {
	req := UrlRequest(request)
	entry, err := c.zipCodeMapper.GetEntryByZipCode(req.GetValue())
	if err != nil {
		WriteResponse(ctx, ErrorString(err.Error()).ToJson(), "JSON")
		return
	}
	content, err := entry.Marshal(req.GetFormat())
	if err != nil {
		WriteResponse(ctx, ErrorString(err.Error()).ToJson(), "JSON")
		return
	}
	WriteResponse(ctx, content, req.GetFormat())
}

func (c ZipCodeController) root(ctx *web.Context) {
	ctx.WriteString("Welcome to Zilch!")
}

func (c ZipCodeController) Start(port string) {
	web.Get("/", c.root)
	web.Get("/lookup/(.*)", c.lookupZipCode)

	fmt.Println("Listening on port:", port)
	web.Run("0.0.0.0:" + port)
}

func NewZipCodeController() ZipCodeController {
	return ZipCodeController{NewZipCodeMapper()}
}
