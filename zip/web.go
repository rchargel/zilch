package zip

import (
	"os"
	"io"
	"bufio"
	"fmt"
	"strings"
	"github.com/hoisie/web"
)

func Root(ctx *web.Context) {
	RenderHtml(ctx, "root.html")
}

func RenderJs(ctx *web.Context, page string) {
	file, err := os.Open("./web/js/" + page)
	if err != nil {
		fmt.Println("Error:",err)
	}
	ctx.ResponseWriter.Header().Set("Content-type", "text/javascript; charset=utf-8")
	r := bufio.NewReader(file)
	io.Copy(ctx.ResponseWriter, r)
}

func RenderImage(ctx *web.Context, image string) {
	file, err := os.Open("./web/images/" + image)
	if err != nil {
		fmt.Println("Error:",err)
	}
	ext := image[strings.Index(image,".") + 1:]
	ctx.ResponseWriter.Header().Set("Content-type", "image/"+ext);
	r := bufio.NewReader(file)
	io.Copy(ctx.ResponseWriter, r)
}

func RenderHtml(ctx *web.Context, page string) {
	file, err := os.Open("./web/" + page)
	if err != nil {
		fmt.Println("Error:",err)
	}
	ctx.ResponseWriter.Header().Set("Content-type", "text/html; charset=utf-8")
	r := bufio.NewReader(file)
	io.Copy(ctx.ResponseWriter, r)
}
