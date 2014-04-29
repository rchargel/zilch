package zilch

import (
	"fmt"
	"github.com/hoisie/web"
	"time"
)

func StartServer(resourceDir, port string) {
	start := time.Now()
	database, _ := NewDatabase(resourceDir)

	zcc := ZipCodeController{database}
	pc := PngController{database}
	sc := StaticController{}

	web.Get("/", sc.RenderRoot)
	web.Get("/query\\.?(.*)", zcc.Query)
	web.Post("/query\\.?(.*)", zcc.Query)
	web.Get("/distribution.jso?n?", zcc.GetDistribution)
	web.Post("/distribution.jso?n?", zcc.GetDistribution)
	web.Get("/countries.jso?n?", zcc.GetCountries)
	web.Post("/countries.jso?n?", zcc.GetCountries)
	web.Get("/map_(\\d*)\\.png", pc.RenderImage)
	web.Get("/images/(.*)", sc.RenderImages)
	web.Get("/js/(.*)", sc.RenderJs)
	web.Get("/(.*)", sc.RenderHtml)

	fmt.Printf("Server started on port %v in %v\n", port, time.Since(start))
	web.Run("0.0.0.0:" + port)
}
