package zip

import (
	"fmt"
	"strings"
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

func WriteResponse(ctx *web.Context, resp string, format string) {
	if format == "XML" {
		ctx.ResponseWriter.Header().Set("Content-type", "text/xml")
	} else if format == "JSON" {
		ctx.ResponseWriter.Header().Set("Content-type", "application/json")
	}
	ctx.WriteString(resp)
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
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
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
