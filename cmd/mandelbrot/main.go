package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

const (
	width      = 3840
	height     = 2160
	xmin, ymin = -2, -2
	xmax, ymax = 2, 2
	maxIter    = 1000
)

var wg sync.WaitGroup

// mandelbrot calculates the color of a point in the Mandelbrot set.
func mandelbrot(c complex128) color.Color {
	z := c
	for i := 0; i < maxIter; i++ {
		if cmplx.Abs(z) > 2 {
			return color.Gray{uint8(255 - i%256)}
		}
		z = z*z + c
	}
	return color.Black
}

// renderRow renders a single row of the Mandelbrot set.
func renderRow(img *image.RGBA, y int) {
	defer wg.Done()
	for x := 0; x < width; x++ {
		c := complex(
			float64(x)/width*(xmax-xmin)+xmin,
			float64(y)/height*(ymax-ymin)+ymin)
		color := mandelbrot(c)
		img.Set(x, y, color)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		wg.Add(1)
		go renderRow(img, y)
	}

	wg.Wait()

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
	}
}
