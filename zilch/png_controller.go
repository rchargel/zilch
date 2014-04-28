package zilch

import (
	"errors"
	"github.com/hoisie/web"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
)

type PngController struct {
	database *Database
}

func (c PngController) RenderImage(ctx *web.Context, scale string) {
	int64_scale, _ := strconv.ParseInt(scale, 10, 32)
	int_scale := int(int64_scale)

	if int_scale < 1 || int_scale > 20 {
		rw := ResponseWriter{ctx, "JSON"}
		rw.SendError(errors.New("Image size is either too large or too small"))
		return
	}
	img := image.NewRGBA(image.Rect(0, 0, 180*int_scale, 90*int_scale))
	bgcolor := color.RGBA{0, 0, 87, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	for _, cim := range c.database.CountryIndexMap {
		for _, entry := range cim.Entries {
			c.drawPoint(img, entry.Latitude, entry.Longitude, float32(int_scale)/float32(2))
		}
	}
	ctx.ContentType("image/png")
	png.Encode(ctx.ResponseWriter, img)
}

func (c PngController) drawPoint(img *image.RGBA, lat, lng, scale float32) {
	y := int((180 - (lat + 90)) * scale)
	x := int((lng + 180) * scale)

	p := image.Rect(x-1, y-1, x+1, y+1)
	zipColor := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, p, &image.Uniform{zipColor}, image.ZP, draw.Src)
}
