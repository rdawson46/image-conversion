package conversion

import (
    "image/color"
    "math"
)

// ANSI COLORS
type AnsiColor struct {
    code int
    r, g, b uint8
}

var AnsiColors = []AnsiColor{
    {30, 0, 0, 0}, // black
    {31, 170, 0, 0}, // red
    {32, 0, 170, 0}, // green
    {33, 170, 85, 0}, // yellow
    {34, 0, 0, 170}, // blue
    {34, 170, 0, 170}, // magenta
    {35, 0, 170, 170}, // cyan
    {37, 170, 170, 170}, // white
}

func ColorToGrayscale(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    return (0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 65535.0
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

    for _, ac := range AnsiColors {
        dist := colorDistance(r, g, b, ac.r, ac.g, ac.b)

        if dist < minDist {
            minDist = dist
            closestColor = ac
        }
    }
    return closestColor
}




// 256 COLORS
const (
    cubeSteps = 6
    cubeStart = 16
    grayStart = 232
)

var (
    colorLevels = []uint8{0, 95, 135, 175, 215, 255}
)

func findClosest256Color(r, g, b uint8) uint8 {
    if isGrayscale(r, g, b) {
        return findClosestGrayscale(r, g, b)
    }

    return findClosestCubeColor(r, g, b)
}

func isGrayscale(r, g, b uint8) bool {
    const threshold = 30
    return math.Abs(float64(r)-float64(g)) < threshold && 
        math.Abs(float64(g)-float64(b)) < threshold && 
        math.Abs(float64(b)-float64(r)) < threshold  
}

func findClosestGrayscale(r, g, b uint8) uint8 {
    gray := (uint16(r) + uint16(g) + uint16(b)) / 3
    level := (gray * 24) / 256
    return uint8(grayStart + level)
}

func findClosestCubeColor(r, g, b uint8) uint8 {
    rIndex := findClosestColorLevel(r)
    gIndex := findClosestColorLevel(g)
    bIndex := findClosestColorLevel(b)

    return uint8(cubeStart + (36 * rIndex) + (6 * gIndex) + bIndex)
}

func findClosestColorLevel(value uint8) int {
    minDistance := uint8(255)
    closestIndex := 0

    for i, level := range colorLevels {
        distance := uint8(math.Abs(float64(value) - float64(level)))

        if distance < minDistance {
            minDistance = distance
            closestIndex = i
        }
    }

    return closestIndex
}
