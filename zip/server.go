package zip

import (
	"fmt"
        "github.com/hoisie/web"
)

func StartServer(port string) {
	zc := NewZipCodeController() 

        web.Get("/", Root)
        web.Get("/lookup/zip/(.*)", zc.LookupZipCode)
        web.Get("/lookup/areaCode/(.*)", zc.LookupAreaCode)
        web.Get("/query\\.?(.*)", zc.QueryReq)
        web.Post("/query\\.?(.*)", zc.QueryReq)
	web.Get("/distribution.json", zc.DistributionMap)
	web.Get("/distribution.js", zc.DistributionMap)
	web.Get("/countries.json", zc.ListCountries)
	web.Get("/countries.js", zc.ListCountries)
        web.Get("/map.png", zc.RenderMap)
        web.Get("/map_small.png", zc.RenderSmallMap)
        web.Get("/map_transparent.png", zc.RenderTransparentMap)
	web.Get("/images/(.*)", RenderImage)
	web.Get("/js/(.*)", RenderJs)
	web.Get("/(.*)", RenderHtml)

        fmt.Println("Listening on port:", port)
        web.Run("0.0.0.0:" + port)
}
