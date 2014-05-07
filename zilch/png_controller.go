package zilch

import (
	"errors"
	"github.com/hoisie/web"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"sort"
	"strconv"
)

type PngController struct {
	database *Database
}

func (c PngController) RenderDistributionImage(ctx *web.Context, scale string) {
	int_scale := c.convertScale(scale)

	if img, err := c.getBackgroundImage(int_scale); err != nil {
		rw := ResponseWriter{ctx, "JSON"}
		rw.SendError(err)
	} else {
		distributions := c.database.GetDistributions()
		sort.Sort(DistributionSorter(distributions))

		for _, dist := range distributions {
			c.drawDistribution(img, float32(dist.Latitude), float32(dist.Longitude), float32(int_scale)/float32(2), int(dist.ZipCodes))
		}

		ctx.ContentType("image/png")
		png.Encode(ctx.ResponseWriter, img)
	}
}

func (c PngController) RenderImage(ctx *web.Context, scale string) {
	int_scale := c.convertScale(scale)

	img, err := c.getBackgroundImage(int_scale)
	if err != nil {
		rw := ResponseWriter{ctx, "JSON"}
		rw.SendError(err)
		return
	}

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

func (c PngController) drawDistribution(img *image.RGBA, lat, lng, scale float32, zipCodes int) {
	y := int((180 - (lat + 90)) * scale)
	x := int((lng + 180) * scale)
	p := image.Rect(x, y, x+int(math.Ceil(float64(scale))), y+int(math.Ceil(float64(scale))))

	var zipColor color.RGBA
	if zipCodes > 1000 {
		zipColor = color.RGBA{255, 0, 0, 255}
	} else if zipCodes > 800 {
		zipColor = color.RGBA{255, 87, 0, 255}
	} else if zipCodes > 600 {
		zipColor = color.RGBA{255, 153, 0, 255}
	} else if zipCodes > 400 {
		zipColor = color.RGBA{255, 204, 51, 255}
	} else if zipCodes > 200 {
		zipColor = color.RGBA{255, 255, 102, 255}
	} else if zipCodes > 100 {
		zipColor = color.RGBA{153, 255, 153, 255}
	} else if zipCodes > 50 {
		zipColor = color.RGBA{255, 255, 255, 255}
	} else {
		zipColor = color.RGBA{204, 204, 204, 255}
	}
	draw.Draw(img, p, &image.Uniform{zipColor}, image.ZP, draw.Src)
}

func (c PngController) convertScale(scale string) int {
	int64_scale, _ := strconv.ParseInt(scale, 10, 32)
	return int(int64_scale)
}

func (c PngController) getBackgroundImage(scale int) (*image.RGBA, error) {
	if scale < 1 || scale > 20 {
		return nil, errors.New("Image is too large or too small")
	}
	file, err := os.Open("./web/images/world-big-2-grey.jpg")
	defer file.Close()

	if err != nil {
		return nil, err
	}
	world, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	world = resize.Resize(uint(180*scale), uint(90*scale), world, resize.Bilinear)
	img := image.NewRGBA(world.Bounds())
	draw.Draw(img, img.Bounds(), world, image.ZP, draw.Src)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 87, 87}}, image.ZP, draw.Over)
	return img, nil
}
