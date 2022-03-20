package ch3

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
)

const (
	xMin, yMin, xMax, yMax = -2, -2, +2, +2
	mWidth, mHeight        = 1024, 1024
)

// Mandelbrot Mandelbrot图像
func Mandelbrot(writer io.Writer) {
	img := image.NewRGBA(image.Rect(0, 0, mWidth, mHeight))
	for py := 0; py < mHeight; py++ {
		y := float64(py)/mHeight*(yMax-yMin) + yMin
		for px := 0; px < mWidth; px++ {
			x := float64(px)/mWidth*(xMax-xMin) + xMin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}
	err := png.Encode(writer, img)
	if err != nil {
		return
	} // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{
				R: 255 - contrast*n,
				G: contrast * n,
				B: 255 - contrast*n,
				A: 0xff,
			}
		}
	}
	return color.Black
}
