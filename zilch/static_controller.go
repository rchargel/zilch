package zilch

import (
	"bufio"
	"github.com/hoisie/web"
	"io"
	"os"
	"strings"
)

type StaticController struct {
}

func (c StaticController) RenderRoot(ctx *web.Context) {
	c.RenderHtml(ctx, "root.html")
}

func (c StaticController) RenderHtml(ctx *web.Context, page string) {
	c.renderContent(ctx, "text/html", "./web/"+page)
}

func (c StaticController) RenderImages(ctx *web.Context, image string) {
	ext := image[strings.Index(image, ".")+1:]
	c.renderContent(ctx, "image/"+ext, "./web/images/"+image)
}

func (c StaticController) RenderJs(ctx *web.Context, jsFile string) {
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
