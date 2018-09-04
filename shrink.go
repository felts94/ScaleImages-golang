package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

type imWithType struct {
	im image.Image
	t  string
}

func main() {
	var w int
	if len(os.Args) < 2 {
		fmt.Println("usage: mode(print,scale) input.png int(width[optional]/scale factor) output.png(optional, defaults to out.png)")
		return
	}
	fmt.Println(os.Args)
	d, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	imtype := readimage(d, os.Args[2])
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // This example uses png.Decode which can only decode PNG images.
	// // Consider using the general image.Decode as it can sniff and decode any registered image format.
	// found := false
	// format := ".jpg"
	// var img image.Image
	// img, err = jpeg.Decode(d)

	//fmt.Println(u)
	// if err != nil {
	// 	log.Println(img, err)
	// } else {
	// 	found = true
	// }
	// if !found {
	// 	img, err = png.Decode(d)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	if os.Args[1] == "print" || os.Args[1] == "p" {
		w = 80
		if len(os.Args) > 3 {
			w, _ = strconv.Atoi(os.Args[3])
		}
		convertToStdOut(imtype.im, w)
	} else if os.Args[1] == "scale" || os.Args[1] == "s" {
		i, _ := strconv.Atoi(os.Args[3])
		convertToFile(imtype.im, "out."+imtype.t, i)
	} else if os.Args[1] == "square" {
		CenterSquare(imtype.im, "out."+imtype.t)
	} else if os.Args[1] == "pixel" {
		i, _ := strconv.Atoi(os.Args[3])
		pixelAverage(imtype.im, i, "out."+imtype.t)
	}
}

func readimage(f *os.File, n string) imWithType {

	imtype := strings.Split(n, ".")[1]

	if imtype == "jpg" || imtype == "jpeg" || imtype == "JPEG" {
		img, err := jpeg.Decode(f)

		//fmt.Println(u)
		if err != nil {
			log.Fatal(img, err)
		}
		return imWithType{
			im: img,
			t:  "png",
		}
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return imWithType{
		im: img,
		t:  "png",
	}
}

func scaleImg(img image.Image, scale int) image.Image {
	a := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Max.X/scale, img.Bounds().Max.Y/scale))
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		if y%scale == 0 {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				if x%scale == 0 {
					a.Set(x/scale, y/scale, img.At(x, y))
				}
			}
		}
	}
	return a
}
func convertToFile(img image.Image, fname string, scale int) {

	a := scaleImg(img, scale)

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func convertToStdOut(img image.Image, i int) {

	var xmod, ymod float32

	xyratio := float32(25.0) / 35.0
	if i >= img.Bounds().Max.X {
		xmod = 1.0
		ymod = 1.0
	} else {
		xmod = float32(img.Bounds().Max.X) / float32(i)
		ymod = xmod / xyratio
	}

	//levels := []string{" ", "░", "▒", "▓", "█"}
	levels := []string{"█", "▓", "▒", "░", " "}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		if y%int(ymod) == 0 {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				if x%int(xmod) == 0 {
					c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
					level := c.Y / 51 // 51 * 5 = 255
					if level == 5 {
						level--
					}
					fmt.Print(levels[level])

				}
			}
			fmt.Print("\n")
		}

	}
}

// CenterSquare makes an image a square
func CenterSquare(img image.Image, fname string) {
	//get the smaller dimention
	max := img.Bounds().Max.X
	if img.Bounds().Max.Y < max {
		max = img.Bounds().Max.Y
	}
	a := image.NewNRGBA(image.Rect(0, 0, max, max))
	//cut off either side of the bigger dim to get a square
	//could be optimized
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(x, y))
		}
	}

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

// IncreaseContrast increases the contrast in an image
func IncreaseContrast() {

}

// SameResolutionPixelation pixelates images but keeps the resolution
func SameResolutionPixelation(img image.Image, scale int, outname string) {
	//get smaller side
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	var smallside int
	if width > height {
		smallside = height
	} else {
		smallside = width
	}

	//find the width that will devide into scale well
	for smallside%scale != 0 {
		smallside--
	}

	a := image.NewNRGBA(image.Rect(0, 0, smallside, smallside))
	//cut off either side of the bigger dim to get a square
	//could be optimized
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(x, y))
		}
	}

	//now we have an image that can be devided into a grid of pixels that are size scale*scale
	midpoint := scale / 2
	xoff, yoff := 0, 0
	for y := a.Bounds().Min.Y; y < a.Bounds().Max.Y; y++ {
		for x := a.Bounds().Min.X; x < a.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(midpoint+xoff*scale, midpoint+yoff*scale))
			if x%scale == 0 {
				xoff++
			}
			//a.Set(x, y, img.At(x, y))
		}
		xoff = 0
		if y%scale == 0 {
			yoff++
		}
	}

	f, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func pixelAverage(img image.Image, scale int, outname string) {

	a := image.NewNRGBA(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, img.Bounds().Max.Y))

	for i := 1; i < scale; i++ {
		//cut off either side of the bigger dim to get a square
		//could be optimized
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x += 2 {
				//Get 4 pix
				tlr, tlg, tlb, tla := img.At(x, y).RGBA()
				trr, trg, trb, tra := img.At(x+i, y).RGBA()
				blr, blg, blb, bla := img.At(x, y+i).RGBA()
				brr, brg, brb, bra := img.At(x+i, y+i).RGBA()

				avgs := new([4]uint32)
				avgs[0] = (tlr + trr + blr + brr) / 4
				avgs[1] = (tlg + trg + blg + brg) / 4
				avgs[2] = (tlb + trb + blb + brb) / 4
				avgs[3] = (tla + tra + bla + bra) / 4
				avgColor := color.RGBA64{
					R: uint16(avgs[0]),
					G: uint16(avgs[1]),
					B: uint16(avgs[2]),
					A: uint16(avgs[3]),
				}

				//set 4 pix
				a.Set(x, y, avgColor)
				a.Set(x+1, y, avgColor)
				a.Set(x, y+1, avgColor)
				a.Set(x+1, y+1, avgColor)
			}
		}
		img = a
	}

	f, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
