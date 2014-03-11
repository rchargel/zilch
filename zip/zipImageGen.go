package zip

import "io"
import "image"
import "image/color"
import "image/draw"
import "image/png"

func drawPoint(m *image.RGBA, lat, lng float64, scale int) {
	y := int((180 - (lat + 90)) * float64(scale))
	x := int((lng + 180) * float64(scale))

	p := image.Rect(x - 1, y - 1, x+1, y+1)
	zipColor := color.RGBA{255,255,255,255}
	draw.Draw(m, p, &image.Uniform{zipColor}, image.ZP, draw.Src)
}

func RenderZipCodeMap(w io.Writer, mapper *ZipCodeMapper, scale int, transparent bool) {
	m := image.NewRGBA(image.Rect(0, 0, 180 * scale * 2, 90 * scale * 2))
	var bgColor color.RGBA 
	if (transparent) {
		bgColor = color.RGBA{0,0,0,0}
	} else {
		bgColor = color.RGBA{0,0,87,255}
	}
	draw.Draw(m, m.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

	for _, cmap := range mapper.ZipCodeMap {
		for _, entry := range cmap {
			if entry.Latitude == 0 && entry.Longitude == 0 {
				// invalid
			} else {
				drawPoint(m, entry.Latitude, entry.Longitude, scale)
			}
		}
	}

	png.Encode(w, m)
}
