package main

import (
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp/totp"
	"image/color"
	"strings"
)

const blank = " " // U+0020
const upper = "▀" // U+2580
const lower = "▄" // U+2584
const whole = "█" // U+2588

var block []string = []string{blank, upper, lower, whole}

func isNotBlack(c color.Color) int {
	if c != color.Black && c != color.Transparent {
		return 1
	}
	return 0
}

func qrline(s string) {
	// fmt.Println("\033[97;40m" + s + "\033[0m")
	fmt.Println(s)
}

func main() {
	key, err := totp.Generate(totp.GenerateOpts{AccountName: "soeiuhfiosuhfoiushgfusehfoushfuoseuhosfph", Issuer: "soph.local"})
	if err != nil {
		panic(err)
	}
	fmt.Println(key.Secret())
	fmt.Println(key.URL())
	bar, err := qr.Encode(key.URL(), qr.L, qr.Unicode)
	if err != nil {
		panic(err)
	}
	qrline(strings.Repeat(lower, bar.Bounds().Dx()+2))
	var s string
	for y := bar.Bounds().Min.Y; y < bar.Bounds().Max.Y+(bar.Bounds().Dy()%2); y += 2 {
		s = whole
		for x := bar.Bounds().Min.X; x < bar.Bounds().Max.X; x++ {
			s += block[isNotBlack(bar.At(x, y))+2*isNotBlack(bar.At(x, y+1))]
			// s += fmt.Sprintf("\033[%d;%dm%s\033[0m", 30+67*isNotBlack(bar.At(x, y)), 40+67*isNotBlack(bar.At(x, y+1)), upperblock)
		}
		qrline(s + whole)
	}
	qrline(strings.Repeat(upper, bar.Bounds().Dx()+2))
	// s, e := totp.GenerateCode("test", time.Now().UTC())
	// if e != nil {
	//	panic(e)
	// }
	// fmt.Println(s)
}
