package main

import (
	"fmt"
	"github.com/boombuler/barcode/qr"
	"image"
	"image/color"
	"os"
	"strings"
	"time"
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
var block = [6]rune{blank, upper, lower, whole, lower, whole}
var bw = []color.Color{color.Black, color.White, color.White, color.White}
var bwp = color.Palette(bw)

func pad(n int, r ...rune) string {
	if len(r) == 0 {
		return strings.Repeat(utf8r(whole).String(false), n)
	}
	return strings.Repeat(utf8r(r[0]).String(false), n)
}

func getblock(img image.Image, x int, y int, yn int) utf8r {
	// fmt.Printf("%v ", image.Pt(x, y+1).In(img.Bounds()))
	return utf8r(block[((bwp.Index(img.At(x, y)))+2*bwp.Index(img.At(x, y+1)))+2*(yn%2)])
}

/*.Intersect(image.Rect(x, y, x+1, y+2))*/
func xys(img image.Image, y *int) any {
	un := img.Bounds().Intersect(image.Rect(img.Bounds().Min.X, *y, img.Bounds().Max.X, *y+2)).Size().Y
	*y = *y + un
	return un
	// return image.Rect(x, y, x+1, y+2).Intersect(img.Bounds().Union(image.Rect(0, 0, 1, 2)))
}
func qrstr(img image.Image, html bool) string {
	var b string = ""
	var y, i int = 0, 0
	yn := 2
	for 1 <= yn {
		yn = xys(img, &y).(int)
		for i = 0; i < img.Bounds().Size().X; i++ {
			b += getblock(img, i, y-yn, yn).String(html)
		}
		if yn <= 1 {
			break
		}
		b += "\n"

	}
	if html {
		b = "<tr>" + strings.ReplaceAll(b, "\n", "</tr><tr>")
		b = strings.TrimSuffix(b, "<tr>")
		b = `<table style="letter-spacing: -2px;line-height: 100%;border-spacing: 0;color: white;background-color: black;">` + b + "</table>"
	} else {
		b = strings.ReplaceAll(b, "\n", pad(1)+"\n"+pad(1))
		b = strings.TrimPrefix(b, pad(1)+"\n") + pad(1)
		b = pad(img.Bounds().Dx()+2, lower) + "\n" + pad(1) + b
	}
	return b
}

func main() {
	s := time.Now().String()
	if len(os.Args) > 1 {
		s = os.Args[1]
	}
	bar, err := qr.Encode(s, qr.M, qr.Auto)
	if err != nil {
		panic(err)
	}
	fmt.Println(qrstr(bar, false))
	fmt.Println(s)
}
