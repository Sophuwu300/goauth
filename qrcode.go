package goauth

/*
 * This file is to turn strings into unicode or html qr codes.
 * Author: sophuwu <sophie@skisiel.com>
 * Feel free to use this code in any way you want.
 * Just call QR("string", header bool, html bool) to get a qr code.
 * The header will display the string above the qr code.
 */

import (
	"bytes"
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/hqbobo/text2pic"
	"image"
	"image/color"
	"image/png"
	"os"
	"slices"
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
var block = [6]rune{blank, upper, lower, whole, lower, whole}
var bw = []color.Color{color.Black, color.White, color.White, color.White}
var bwp = color.Palette(bw)

func pad(n int, r ...rune) string {
	if n == 1 {
		n += 1
	} else if n < 1 {
		n = 1
	}
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
		b = pad(0) + pad(img.Bounds().Dx()+2, whole) + pad(0) + "\n" + pad(1) + b + "\n"
		b += pad(img.Bounds().Dx()+4, whole)
	}
	return b
}

func wrap(w int, s ...string) []string {
	var line = ""
	var lines, b []string
	for n := range s {
		b = strings.Split(s[n], " ")
		for i, v := range b {
			if len(v) > w {
				lines = append(lines, v[:w-1], v[w-1:])
				continue
			}
			if len(line)+len(v) < w {
				line += v + " "
			} else {
				lines = append(lines, strings.TrimSuffix(line, " "))
				line = ""
				continue
			}
			if i == len(b)-1 {
				lines = append(lines, strings.TrimSuffix(line, " "))
				line = ""
			}
		}
	}
	return slices.Clip(lines)
}

type QRcode struct {
	Data    string
	Headers []string
	Code    image.Image
}

func GenQR(data string, header ...string) (QRcode, error) {
	code, err := qr.Encode(data, qr.H, qr.Auto)
	if err != nil {
		return QRcode{}, err
	}
	return QRcode{
		data,
		header,
		code,
	}, nil
}

func (q QRcode) String() string {
	var output = ""
	ar := wrap(q.Code.Bounds().Dx(), q.Headers...)
	width := q.Code.Bounds().Dx()
	output += fmt.Sprintln(pad(0, whole) + pad(width+2, upper) + pad(0, whole))
	for i, v := range ar {
		v = v + pad(width-len(v)-1+2*(width%2), blank)
		v = pad(0) + pad(0, blank) + v + pad(0)
		output += fmt.Sprintf("%s", v)
		if i < len(ar) {
			output += fmt.Sprintln()
		}
	}
	output += fmt.Sprintln(pad(0, whole) + pad(width+2, lower) + pad(0, whole))

	output += fmt.Sprintln(qrstr(q.Code, false))
	return output
}

func (q QRcode) HTML() string {
	var output = ""
	for _, h := range q.Headers {
		output += fmt.Sprintln("<p>" + h + "</p>")
	}
	output += fmt.Sprintln(qrstr(q.Code, true))
	return output
}

func (q QRcode) Png() []byte {
	fontBytes, _ := os.ReadFile("Oxygen.ttf")
	ttf, _ := freetype.ParseFont(fontBytes)
	// create a new image with white background
	var pic = text2pic.NewTextPicture(text2pic.Configure{Width: 1000, BgColor: text2pic.ColorWhite})
	// add the qr code header
	for _, h := range q.Headers {
		pic.AddTextLine(h, 8, ttf, text2pic.ColorBlack, text2pic.Padding{0, 0, 0, 0, 0})
	}

	var buf bytes.Buffer
	pic.Draw(&buf, text2pic.TypePng)
	img2, _ := png.Decode(&buf)

	// scale q.Code to 1000px
	img := imaging.Resize(q.Code, 900, 900, imaging.NearestNeighbor)
	img = imaging.Paste(imaging.New(1000, 1000+img2.Bounds().Size().Y, text2pic.ColorWhite), img, image.Pt(50, img2.Bounds().Size().Y+50))

	img = imaging.Paste(img, img2, image.Pt(0, 0))

	buf.Reset()
	png.Encode(&buf, img)

	return buf.Bytes()
}
