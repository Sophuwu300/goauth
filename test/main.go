package main

import (
	"fmt"
	"github.com/boombuler/barcode/qr"
	"image"
	"image/color"
	_ "image/draw"
	"os"
	"strings"
)

type utf8r rune

func (u utf8r) String(html bool) string {
	f := "%c"
	if html {
		f = "<td>&#x%x;</td>"
	}
	return fmt.Sprintf(f, u)
}

const blank rune = ' ' // U+0020
const upper rune = '▀' // U+2580
const lower rune = '▄' // U+2584
const whole rune = '█' // U+2588
var block = [4]rune{blank, upper, lower, whole}
var bw = []color.Color{color.Black, color.White}
var bwp = color.Palette(bw)

func getblock(img image.Image, x int, y int) utf8r {
	return utf8r(block[bwp.Index(img.At(x, y))+2*bwp.Index(img.At(x, y+1))])
}

func qrstr(img image.Image, html bool) string {
	var b string = ""
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y+(img.Bounds().Dy()%2); y += 2 {
		// b += block[3].String()
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			b += getblock(img, x, y).String(html)
		}
		b += "\n"
	}
	if html {
		b = "<tr>" + strings.ReplaceAll(b, "\n", "</tr><tr>")
		b = strings.TrimSuffix(b, "<tr>")
		b = `<table style="letter-spacing: -2px;line-height: 100%;border-spacing: 0;color: white;background-color: black;">` + b + "</table>"
	}
	return b
}

func main() {
	s := "Hello, 世界"
	if len(os.Args) > 1 {
		s = os.Args[1]
	}
	bar, err := qr.Encode(s, qr.H, qr.Unicode)
	if err != nil {
		panic(err)
	}
	fmt.Println(qrstr(bar, false))
	fmt.Println(qrstr(bar, true))
}
