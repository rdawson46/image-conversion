package conversion_test

import (
    "bytes"
    "image"
    "image/color"
    "fmt"
    "os"
    "testing"
    "path/filepath"

    "github.com/rdawson46/pic-conversion/internal/conversion"
)

// TODO: create tests
func TestConvertImage(t *testing.T) {
    width := 550

    entries, err := os.ReadDir("./images")

    if err != nil {
        t.Error(err)
    }

    for _, file := range entries {
        if !file.IsDir() {
            b, err := os.ReadFile(filepath.Join("./images", file.Name()))
            if err != nil {
                t.Error(err)
            }

            img, _, err := image.Decode(bytes.NewReader(b))

            if err != nil {
                t.Error(err)
            }

            ansi := conversion.Convert256ColorImage(img, width)
            fmt.Println(ansi)
        }
    }
}

func TestResizeImage(t *testing.T) {
    b, err := os.ReadFile("./images/sample.png")

    if err != nil {
        t.Error(err)
    }

    img, _, err := image.Decode(bytes.NewReader(b))

    if err != nil {
        t.Error(err)
    }

    new_image := conversion.ResizeImage(img, 100)

    if new_image.Bounds().Size().X > 100 {
        t.Error("Error resizing image")
    }
}

func TestColorToGrayscale(t *testing.T) {
    c := color.RGBA{
        R: 0,
        G: 0,
        B: 0,
        A: 0,
    }

    f := conversion.ColorToGrayscale(c)

    if f != 0.0 {
        t.Error("conversion to grayscale failed")
    }
}
