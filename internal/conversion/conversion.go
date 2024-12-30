package conversion

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"
)

type AnsiColor struct {
    code int
    r, g, b uint8
}

var ansiColors = []AnsiColor{
    {30, 0, 0, 0}, // black
    {31, 170, 0, 0}, // red
    {32, 0, 170, 0}, // green
    {33, 170, 85, 0}, // yellow
    {34, 0, 0, 170}, // blue
    {34, 170, 0, 170}, // magenta
    {35, 0, 170, 170}, // cyan
    {37, 170, 170, 170}, // white
}

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

func ColorToGrayscale(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    return (0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 65535.0
}

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

// functions for color

func colorDistance(r1, g1, b1, r2, g2, b2 uint8) float64 {
    rDiff := float64(r1) - float64(r2)
    gDiff := float64(g1) - float64(g2)
    bDiff := float64(b1) - float64(b2)

    return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
}

func closestAnsiColor(r, g, b uint8) AnsiColor {
    var closestColor AnsiColor
    minDist := math.MaxFloat64

    for _, ac := range ansiColors {
        dist := colorDistance(r, g, b, ac.r, ac.g, ac.b)

        if dist < minDist {
            minDist = dist
            closestColor = ac
        }
    }
    return closestColor
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
