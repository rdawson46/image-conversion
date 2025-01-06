package conversion

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"
)


var asciiChars = []rune{
    ' ',
    '.',
    ':',
    '-',
    '~',
    '+',
    '=',
    '*',
    '#',
    '%',
    '@',
}

/*
var asciiChars = []rune {
    'I', 'l', 'i', 'j', 't', 'r', 'f',
    'c', 'v','u','n','x','z','J','L',
    'C','Y','U','O','Q','D','P','X',
    'M','W','H','K','B','R','N',
}
*/

func MapToAscii(grayValue float64) rune {
    i := int(math.Floor(grayValue * float64(len(asciiChars)-1)))
    return asciiChars[i]
}

func ResizeImage(img image.Image, width int) image.Image {
    bounds := img.Bounds()
    dy := bounds.Dy()
    dx := bounds.Dx()

    ratio := float64(dy) / float64(dx)
    height := int(float64(width) * ratio)
    resized := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        origY := y * dy / height

        for x := 0; x < width; x++ {
            origX := x * dx / width

            c := img.At(origX, origY)

            resized.Set(x, y, c)
        }
    }

    return resized
}

func OldConvertImage(img image.Image, width int) string {
    var ans string

    resizedImage := ResizeImage(img, width)
    bounds := resizedImage.Bounds()
    height := bounds.Max.Y

    for y := 0; y < height; y++ {
        l := ""

        for x := 0; x < width; x++ {
            pixelColor := resizedImage.At(x, y)
            grayValue := ColorToGrayscale(pixelColor)
            asciiChar := MapToAscii(grayValue)

            l += string(asciiChar)
        }

        if len(strings.TrimSpace(l)) != 0 {
            ans += l + "\n"
        }
    }

    return ans
}

func ConvertColorImage(img image.Image, width int) string {
    var ans string

    resizedImage := ResizeImage(img, width)
    bounds := resizedImage.Bounds()
    height := bounds.Max.Y
    width = bounds.Max.X

    for y := 0; y < height; y++ {
        l := ""

        for x := 0; x < width; x++ {
            pixelColor := resizedImage.At(x, y)
            r, g, b, _ := pixelColor.RGBA()

            r8 := uint8(r >> 8)
            g8 := uint8(g >> 8)
            b8 := uint8(b >> 8)

            ansiColor := closestAnsiColor(r8, g8, b8)

            grayValue := ColorToGrayscale(pixelColor)
            asciiChar := MapToAscii(grayValue)

            var asciiArt string
            if asciiChar != ' '{
                asciiArt = fmt.Sprintf("\033[%dm%c", ansiColor.code, asciiChar)
            } else {
                asciiArt = string(asciiChar)
            }

            l += string(asciiArt)
       }

        if len(strings.TrimSpace(l)) != 0 {
            ans += l + "\033[0m\n"
        }
    }

    return ans
}

func Convert256ColorImage(img image.Image, width int) string {
    var ans string

    resizedImage := ResizeImage(img, width)
    bounds := resizedImage.Bounds()
    height := bounds.Max.Y
    width = bounds.Max.X

    for y := 0; y < height; y++ {
        l := ""

        for x := 0; x < width; x++ {
            pixelColor := resizedImage.At(x, y)
            r, g, b, _ := pixelColor.RGBA()

            r8 := uint8(r >> 8)
            g8 := uint8(g >> 8)
            b8 := uint8(b >> 8)

            colorCode := findClosest256Color(r8, g8, b8)

            grayValue := ColorToGrayscale(pixelColor)
            asciiChar := MapToAscii(grayValue)

            var asciiArt string
            if asciiChar != ' '{
                asciiArt = fmt.Sprintf("\033[38;5;%dm%c", colorCode, asciiChar)
            } else {
                asciiArt = string(asciiChar)
            }
            asciiArt += "\033[0m"

            l += string(asciiArt)
       }

        if len(strings.TrimSpace(l)) != 0 {
            ans += l + "\033[0m\n"
        }
    }

    return ans
}

func Compress(ans string) (string, error) {
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)

    _, err := gz.Write([]byte(ans))

    if err != nil {
        return ans, errors.New("Couldn't compress string")
    }

    if err := gz.Close(); err != nil {
        return ans, errors.New("Error closing compressor")
    }

    return b.String(), nil
}
