package zilch

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/hoisie/web"
)

// StaticController this is the controller for static content.
type StaticController struct {
}

// RenderRoot outputs the root html file.
func (c StaticController) RenderRoot(ctx *web.Context) {
	c.RenderHTML(ctx, "root.html")
}

// RenderHTML outputs a specific HTML page.
func (c StaticController) RenderHTML(ctx *web.Context, page string) {
	c.renderContent(ctx, "text/html", "./web/"+page)
}

// RenderImages outputs content from the images directory.
func (c StaticController) RenderImages(ctx *web.Context, image string) {
	ext := image[strings.Index(image, ".")+1:]
	c.renderContent(ctx, "image/"+ext, "./web/images/"+image)
}

// RenderJS outputs content for the javascript directory.
func (c StaticController) RenderJS(ctx *web.Context, jsFile string) {
	c.renderContent(ctx, "application/javascript", "./web/js/"+jsFile)
}

func (c StaticController) renderContent(ctx *web.Context, contentType, filepath string) {
	if file, err := os.Open(filepath); err != nil {
		w := ResponseWriter{ctx, "HTML"}
		w.SendError(err)
	} else {
		ctx.ContentType(contentType + "; charset=utf-8")
		reader := bufio.NewReader(file)
		io.Copy(ctx.ResponseWriter, reader)
	}
}
