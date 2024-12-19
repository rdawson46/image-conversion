package conversion

import (
	"image"
	"image/color"
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
    ratio := float64(bounds.Dy()) / float64(bounds.Dx())

    height := int(float64(width) * ratio)

    resized := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        origY := y * bounds.Dy() / height

        for x := 0; x < height; x++ {
            origX := x * bounds.Dx() / width

            c := img.At(origX, origY)

            resized.Set(x, y, c)
        }
    }

    return resized
}

func ConvertImage(img image.Image, width int) string {
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

