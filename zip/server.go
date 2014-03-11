package zip

import (
	"fmt"
        "github.com/hoisie/web"
)

func StartServer(port string) {
	zc := NewZipCodeController() 

        web.Get("/", Root)
        web.Get("/lookup/zip/(.*)", zc.lookupZipCode)
        web.Get("/lookup/areaCode/(.*)", zc.lookupAreaCode)
        web.Get("/query\\.?(.*)", zc.queryReq)
        web.Post("/query\\.?(.*)", zc.queryReq)
        web.Get("/map.png", zc.renderMap)
	web.Get("/js/(.*)", RenderJs)
	web.Get("/(.*)", RenderHtml)

        fmt.Println("Listening on port:", port)
        web.Run("0.0.0.0:" + port)
}
