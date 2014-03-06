package zip

import (
	"fmt"
	"github.com/hoisie/web"
)

type ZipCodeController struct {
	zipCodeMapper *ZipCodeMapper
}

func (c ZipCodeController) lookupZipCode(ctx *web.Context, path string) {
	c.zipCodeMapper.PublishZipCode(path, "JSON", ctx.ResponseWriter)
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
