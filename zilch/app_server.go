package zilch

import (
	"fmt"
	"time"

	"github.com/hoisie/web"
)

// StartServer starts the Zilch Web Server.
func StartServer(resourceDir, port string) {
	start := time.Now()
	database, _ := NewDatabase(resourceDir)

	zcc := ZipCodeController{database}
	pc := PngController{database}
	sc := StaticController{}

	web.Get("/", sc.RenderRoot)
	web.Get("/query\\.?(.*)", zcc.Query)
	web.Post("/query\\.?(.*)", zcc.Query)
	web.Get("/distribution\\.?(.*)", zcc.GetDistribution)
	web.Post("/distribution\\.?(.*)", zcc.GetDistribution)
	web.Get("/countries\\.?(.*)", zcc.GetCountries)
	web.Post("/countries\\.?(.*)", zcc.GetCountries)
	web.Get("/map_(\\d*)\\.png", pc.RenderImage)
	web.Get("/distmap_(\\d*)\\.png", pc.RenderDistributionImage)
	web.Get("/images/(.*)", sc.RenderImages)
	web.Get("/js/(.*)", sc.RenderJS)
	web.Get("/(.*)", sc.RenderHTML)

	fmt.Printf("Server started on port %v in %v\n", port, time.Since(start))
	web.Run("0.0.0.0:" + port)
}
