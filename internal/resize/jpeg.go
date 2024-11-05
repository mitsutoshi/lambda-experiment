package resize

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

func ResizeJpegRGBA(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
func ResizeJpegGrayscale(graysclae image.Image, width, height int) *image.Gray {
	dst := image.NewGray(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Bounds(), graysclae, graysclae.Bounds(), draw.Over, nil)
	return dst
}

func ResizeJpegCmyk(cmyk image.Image, width, height int) *image.RGBA {
	srcRgba := ConvertCmykToRGBA(cmyk)
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Bounds(), srcRgba, srcRgba.Bounds(), draw.Over, nil)
	return dst
}

func ConvertCmykToRGBA(cmyk image.Image) *image.RGBA {
	bounds := cmyk.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := cmyk.At(x, y)
			if cmyk, ok := c.(color.CMYK); ok {
				r, g, b := color.CMYKToRGB(cmyk.C, cmyk.M, cmyk.Y, cmyk.K)
				cRgba := color.RGBA{r, g, b, 255}
				rgba.Set(x, y, cRgba)
			}
		}
	}
	return rgba
}
