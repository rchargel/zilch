package zip

import "io"
import "image"
import "image/color"
import "image/draw"
import "image/png"

func drawPoint(m *image.RGBA, lat, lng float64) {
	y := int((180 - (lat + 90)) * 2)
	x := int((lng + 180) * 2)

	p := image.Rect(x - 1, y - 1, x + 1, y + 1)
	green := color.RGBA{0,255,0,255}
	draw.Draw(m, p, &image.Uniform{green}, image.ZP, draw.Src)
}

func RenderZipCodeMap(w io.Writer, mapper *ZipCodeMapper) {
	m := image.NewRGBA(image.Rect(0, 0, 180 * 4, 90 * 4))
	blue := color.RGBA{0,0,255,255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	for _, cmap := range mapper.ZipCodeMap {
		for _, entry := range cmap {
			drawPoint(m, entry.Latitude, entry.Longitude)
		}
	}

	png.Encode(w, m)
}
