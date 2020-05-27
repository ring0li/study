package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/disintegration/imaging"
)

func main() {
	//os.Create("/Users/liuli/work/go/code/src/xiaozhu/user.xiaozhu/app/img/headImg/12.jpg")
	img, err := imaging.Open("./img/headImg/a543b186efd275aedc267141d894fe40.png")
	if err != nil {
		log.Fatal(err)
	}

	img = makeCircleSmooth(img, 1)

	err = imaging.Save(img, "./img/headImg/circle.png")
	if err != nil {
		log.Fatal(err)
	}
}

func makeCircleSmooth(src image.Image, factor float64) image.Image {
	d := src.Bounds().Dx()
	if src.Bounds().Dy() < d {
		d = src.Bounds().Dy()
	}
	dst := imaging.CropCenter(src, d, d)
	r := float64(d) / 2
	center := r - 0.5
	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			xf := float64(x)
			yf := float64(y)
			delta := math.Sqrt((xf-center)*(xf-center)+(yf-center)*(yf-center)) + factor - r
			switch {
			case delta > factor:
				dst.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0})
			case delta > 0:
				m := 1 - delta/factor
				c := dst.NRGBAAt(x, y)
				c.A = uint8(float64(c.A) * m)
				dst.SetNRGBA(x, y, c)
			}
		}
	}
	return dst
}