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
	if format == "JS" { format = "JSON" }
	if len(callback) != 0 {
		format = "JS"
	}
	switch format {
	case "XML":
		ctx.ResponseWriter.Header().Set("Content-type", "text/xml")
	case "JSON":
		ctx.ResponseWriter.Header().Set("Content-type", "application/json")
	case "YAML":
		ctx.ResponseWriter.Header().Set("Content-type", "text/x-yaml")
	case "JS":
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
	q := make(map[string]string)
	q["ZipCode"] = req.GetValue()
	c.query(ctx, q, req.GetFormat())
}

func (c ZipCodeController) lookupAreaCode(ctx *web.Context, request string) {
	req := UrlRequest(request)
	q := make(map[string]string)
	q["AreaCode"] = req.GetValue()
	c.query(ctx, q, req.GetFormat())
}

func (c ZipCodeController) queryReq(ctx *web.Context, format string) {
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

func (c ZipCodeController) root(ctx *web.Context) {
	ctx.Header().Set("Content-type","text/html")
	content := `
<!doctype HTML>
<html>
	<head>
		<title>ZiLCh - Zipcode Lookup Channel</title>
		<meta charset="UTF-8"/>
		<style type="text/css">
		html,body { font-family: Arial, sans-serif; background: #ddd; }
		h1,h2,h3 { text-align: center; }
		h4,h5,h6 { text-align: left; }
		table { width: 80%; border-spacing: 0; border-collapse: separate; }
		td,th { border: 1px solid #000; text-align: left; }
		</style>
	</head>
	<body>
		<h1>Welcome to ZiLCh</h1>
		<h3>Zipcode Lookup Channel</h3>
		<p>
			ZiLCh is a simple REST-like API used to lookup information about zip codes. 
			It's very simple to use, as it has only one URL to retrieve from (/query). Here's an
			<a href="/query.yaml?ZipCode=12345">example</a>: 
		</p>
		<h4>How much does it cost?</h4>
		<p>Zilch! It's free to use.</p>
		<h4>What query parameters are supported?</h4>
		<table>
			<thead>
				<tr><th>Parameter</th><th>Description</th></tr>
			</thead>
			<tbody>
				<tr> <td>ZipCode</td> <td>The zip code you're looking for.</td> </tr>
				<tr> <td>City</td> <td>Any part of the city's name (eg: Phila will find Philadelphia). Not case sensitive.</td> </tr>
				<tr> <td>Country</td> <td>The 2-letter country code. Uppercase only.</td> </tr>
				<tr> <td>State</td> <td>The state abbreviation. Uppercase only.</td> </tr>
				<tr> <td>County</td> <td>Any part of the county's name. Not all countries have county data. Not case sensitive.</td> </tr>
				<tr> <td>AreaCode</td> <td>The 3-digit area code for a phone number. United States Only.</td> </tr>
			</tbody>
		</table>
		<p><strong>Note:</strong> There is a hard limit of 1,000 records returned in any query.</p>
		<h4>What response formats are supported?</h4>
		<p>The response format is selected by changing the file extension of "/query" (eg: /query.xml)
			<ul>
				<li>JSON: /query.json or /query.js (this is the default format if no extension is added)</li>
				<li>XML: /query.xml</li>
				<li>YAML: /query.yaml</li>
			</ul>
		</p>
		<h4>What about JSONP support?</h4>
		<p>
			If you hit this service directly from a browser, you'll probably want to use JSONP in order to
			avoid cross-domain scripting issues.  Simply add your callback to your GET request via the "callback"
			or "jsonp" query parameters.  This is supported out-of-the-box by Javascript APIs like 
			<a href="http://jquery.org">jQuery</a>.  The only response format for JSONP is JSON.
			Here's an example: <a href="/query.js?callback=mycallback&ZipCode=12345">/query.js?<b>callback=mycallback</b>&amp;ZipCode=12345</a>
		</p>
		<h4>What countries are supported?</h4>
		<p>Currently the following countries are supported
			<ul>
				<li>United States (data from US Postal Service)</li>
				<li>Canada (data from GeoNames.org)</li>
				<li>Mexico (data from GeoNames.org)</li>
				<li>Great Britain (data from GeoNames.org)</li>
				<li>Brazil (data from GeoNames.org)</li>
			</ul>
		</p>
		<h4>Where can I get the source?</h4>
		<p>
			The source is available on <a href="https://github.com/rchargel/zilch">GitHub</a>. You'll need <a href="http://golang.org/">Go</a> to
			run or install this on your own server.
		</p>
	</body>
</html>
`
	ctx.WriteString(content)
}

func (c ZipCodeController) Start(port string) {
	web.Get("/", c.root)
	web.Get("/lookup/zip/(.*)", c.lookupZipCode)
	web.Get("/lookup/areaCode/(.*)", c.lookupAreaCode)
	web.Get("/query\\.?(.*)", c.queryReq)
	web.Post("/query\\.?(.*)", c.queryReq)

	fmt.Println("Listening on port:", port)
	web.Run("0.0.0.0:" + port)
}

func NewZipCodeController() ZipCodeController {
	return ZipCodeController{NewZipCodeMapper()}
}
